package utils

import (
	models "AuthService/models"
	"context"
	"database/sql"
	"fmt"

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

	defer db.Close()

	err = db.Ping()
	CheckError(err)

	return db
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
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

func InserUser(user models.UserAccount) {
	db := ConnectToDb()
	defer db.Close()

	sqlStatement := "INSERT INTO user_account (email, phone_number, gender, first_name, last_name, password_hash) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := db.Exec(sqlStatement, user.Email, user.PhoneNumber, user.Gender, user.FirstName, user.LastName, user.PasswordHash)
	CheckError(err)
}

func FindUser(user models.SignInInput) {
	db := ConnectToDb()
	defer db.Close()

	sqlStatement := "SELECT * FROM user_account WHERE email = $1"
	_, error := db.Query(sqlStatement, user.Email)
	CheckError(error)
}
