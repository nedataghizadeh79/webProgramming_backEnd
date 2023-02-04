package utils

import (
	models "AuthService/models"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

var (
	host     = GetEnv("POSTGRES_HOST")
	port     = GetEnv("POSTGRES_PORT")
	user     = GetEnv("POSTGRES_USER")
	password = GetEnv("POSTGRES_PASSWORD")
	dbname   = GetEnv("POSTGRES_DB")
)

func ConnectToDb() *sql.DB {

	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	err = db.Ping()
	CheckError(err)

	return db
}

func CheckError(err error) {
	if err != nil {
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

func GetUserData(username string, password string) {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	val, err := client.Get(ctx, username).Result()

	CheckError(err)
	if val != "" {
		fmt.Println(val)
	} else {
		fmt.Println("key not found")
	}
}

func InsertUser(w http.ResponseWriter, r *http.Request) {
	var user models.UserAccount
	en_err := json.NewDecoder(r.Body).Decode(&user)
	if en_err != nil {
		http.Error(w, en_err.Error(), http.StatusBadRequest)
		return
	}

	db := ConnectToDb()
	defer db.Close()

	hashedPassword := HashPassword(user.Password)

	sqlStatement := "INSERT INTO user_account (email, phone_number, gender, first_name, last_name, password_hash) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := db.Exec(sqlStatement, user.Email, user.PhoneNumber, user.Gender, user.FirstName, user.LastName, hashedPassword)
	CheckError(err)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	tokenString, token_err := GetToken(user.Email)

	if token_err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(token_err.Error()))
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(tokenString)

}

func SignInUser(w http.ResponseWriter, r *http.Request) {
	var user models.SignInInput

	en_err := json.NewDecoder(r.Body).Decode(&user)

	if en_err != nil {
		http.Error(w, en_err.Error(), http.StatusBadRequest)
		return
	}

	db := ConnectToDb()
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
	default:
		panic(err)
	}
}

func FindUser(email string) {
	db := ConnectToDb()
	defer db.Close()

	sqlStatement := "SELECT * FROM user_account WHERE email = $1"
	data, err := db.Query(sqlStatement, email)
	CheckError(err)

	fmt.Println(data)
}
