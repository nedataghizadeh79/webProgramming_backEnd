package main

import (
	utils "AuthService/utils"
	"errors"
	"fmt"
	"net/http"
	"os"

    "github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func signUp(w http.ResponseWriter, r *http.Request) { fmt.Println("signup") }

func signIn(w http.ResponseWriter, r *http.Request) { fmt.Println("signin") }

func getUser(w http.ResponseWriter, r *http.Request) { fmt.Println("user data") }

func signOut(w http.ResponseWriter, r *http.Request) {}

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
}

func main() {
	utils.ConnectToDb()
	utils.GetUserData("user:1", "data")

    r := gin.Default

	r.POST("/signUp", signUp)
	r.POST("/signIn", signIn)
	r.POST("/signOut", signOut)
	r.GET("/userInfo", getUser)

	r.run()
}
