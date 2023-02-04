package utils

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
}

func GetEnv(key string) string {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return ""
}
