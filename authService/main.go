package main

import (
	"AuthService/models"
	utils "AuthService/utils"
	"log"

	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func signUp(w http.ResponseWriter, r *http.Request) {
	var user models.UserAccount
	error := json.NewDecoder(r.Body).Decode(&user)
	if (error != nil) {
		http.Error(w, error.Error(), http.StatusBadRequest)
		return
	}

	utils.ValidateSignUpData(user)

	utils.InsertUser(w, r)
}

func signIn(w http.ResponseWriter, r *http.Request) {
	hmacSecret := []byte(utils.GetEnv("SECRET_JWT_KEY"))
	token, _ := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
		return hmacSecret, nil
	})

	var email string

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		email = claims["sub"].(string)
	}

	utils.FindUser(email)
}

func getUser(c *gin.Context) { fmt.Println("user data") }

func signOut(c *gin.Context) {}

func main() {
	utils.ConnectToDb()
	utils.GetUserData("user:1", "data")

	// r := gin.Default()

	// r.POST("/signUp", signUp)
	// r.POST("/signIn", signIn)
	// r.POST("/signOut", signOut)
	// r.GET("/userInfo", getUser)

	// r.Run()

	mux := http.NewServeMux()

	mux.HandleFunc("/signUp", signUp)
	mux.Handle("/signIn", utils.ValidateJWT(signIn))

	err := http.ListenAndServe(":3333", mux)

	if err != nil {
		log.Fatal(err)
	}

}
