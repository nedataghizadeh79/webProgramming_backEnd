package utils

import (
	// models "AuthService/models"
	// pb "AuthService/proto"
	// "context"
	"database/sql"
	// "encoding/json"
	"fmt"
	// "net/http"

	// "github.com/go-redis/redis/v8"
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
