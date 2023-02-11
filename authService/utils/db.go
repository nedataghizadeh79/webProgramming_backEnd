package utils

import (
	models "AuthService/models"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	_ "github.com/lib/pq"
)

var (
	host     = GetEnv("POSTGRES_HOST")
	port     = GetEnv("POSTGRES_PORT")
	user     = GetEnv("POSTGRES_USER")
	password = GetEnv("POSTGRES_PASSWORD")
	dbname   = GetEnv("POSTGRES_DB")
)

func ConnectToDb(w http.ResponseWriter) *sql.DB {

	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err, w)

	err = db.Ping()
	CheckError(err, w)

	return db
}

func AddToRedis(key string, value string, expiry time.Duration) error {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	defer client.Close()

	err := client.Set(ctx, key, value, expiry).Err()

	return err

}

func GetFromRedis(key string) (string, error) {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	defer client.Close()

	value, err := client.Get(ctx, key).Result()

	return value, err

}

func CheckError(err error, w http.ResponseWriter) {
	if err != nil {
		if w != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		panic(err)
	}
}

func GetToken(key string) (tokenJson []byte, err error) {
	token, token_err := CreateToken(key)

	if token_err != nil {
		return nil, token_err
	}

	tokenString, j_err := json.Marshal(&models.AuthToken{Token: token})

	if j_err != nil {
		return nil, j_err
	}

	return tokenString, nil
}

func InsertUser(w http.ResponseWriter, user models.UserAccount) {

	db := ConnectToDb(w)
	defer db.Close()

	hashedPassword := HashPassword(user.Password)

	sqlStatement := "INSERT INTO user_account (email, phone_number, gender, first_name, last_name, password_hash) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := db.Exec(sqlStatement, user.Email, user.PhoneNumber, user.Gender, user.FirstName, user.LastName, hashedPassword)
	CheckError(err, w)

	if err != nil {
		fmt.Println("db error")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	tokenString, token_err := GetToken(user.Email)

	if token_err != nil {
		fmt.Println("token error")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(token_err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(tokenString)

	var token models.AuthToken
	json.Unmarshal(tokenString, &token)

	redis_err := AddToRedis(token.Token, "valid", time.Hour)

	CheckError(redis_err, nil)

}

func SignInUser(w http.ResponseWriter, user models.SignInInput) {

	db := ConnectToDb(w)
	defer db.Close()

	sqlStatement := "SELECT password_hash FROM user_account WHERE email =$1"
	row := db.QueryRow(sqlStatement, user.Email)

	var password_hash string

	switch err := row.Scan(&password_hash); err {
	case sql.ErrNoRows:
		w.WriteHeader(http.StatusNotFound)
	case nil:
		tokenString, token_err := GetToken(user.Email)
		if token_err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(token_err.Error()))
			return
		}

		passErr := ComparePassword(password_hash, user.Password)
		if passErr != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("password incorrect"))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(tokenString)

		var token models.AuthToken
		json.Unmarshal(tokenString, &token)

		redis_err := AddToRedis(token.Token, "valid", time.Hour)

		CheckError(redis_err, nil)
	default:
		panic(err)
	}
}

func FindUser(w http.ResponseWriter, r *http.Request) {

	hmacSecret := []byte(GetEnv("SECRET_JWT_KEY"))
	token, _ := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
		return hmacSecret, nil
	})

	var email string

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		email = claims["sub"].(string)
	}

	db := ConnectToDb(w)
	defer db.Close()

	sqlStatement := "SELECT first_name, last_name, gender, email, phone_number FROM user_account WHERE email = $1"
	row := db.QueryRow(sqlStatement, email)

	var user models.UserResponse = models.UserResponse{}

	switch err := row.Scan(&user.FirstName, &user.LastName, &user.Gender, &user.Email, &user.PhoneNumber); err {
	case sql.ErrNoRows:
		w.WriteHeader(http.StatusNotFound)
	case nil:
		userJson, json_err := json.Marshal(&user)

		if json_err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(userJson)
	default:
		panic(err)
	}

}

func AddExpiredToken(token models.AuthToken) error {

	hmacSecret := []byte(GetEnv("SECRET_JWT_KEY"))
	jwtToken, _ := jwt.Parse(token.Token, func(token *jwt.Token) (interface{}, error) {
		return hmacSecret, nil
	})

	var email string

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok {
		email = claims["sub"].(string)
	}

	db := ConnectToDb(nil)
	defer db.Close()

	sqlStatement := "SELECT user_id FROM user_account WHERE email = $1"
	row := db.QueryRow(sqlStatement, email)

	var user_id int

	if err := row.Scan(&user_id); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user not found")
		} else {
			return errors.New("database crash")
		}
	}

	timestamp := time.Now()
	now := timestamp.Format("2006-01-02 15:04:05")

	_, err := db.Exec("INSERT INTO unauthorized_token VALUES($1, $2, $3)", user_id, token.Token, now)

	return err
}
