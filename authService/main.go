package main

import (
	utils "AuthService/utils"
	"log"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func signUp(w http.ResponseWriter, r *http.Request) {

	utils.InsertUser(w, r)
}

func signIn(w http.ResponseWriter, r *http.Request) {
	utils.SignInUser(w, r)
}

func getUser(w http.ResponseWriter, r *http.Request) {
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

func signOut(c *gin.Context) {}

func main() {
	utils.ConnectToDb()
	utils.GetUserData("user:1", "data")

	mux := http.NewServeMux()

	mux.HandleFunc("/signUp", signUp)
	mux.HandleFunc("/signIn", signIn)
	mux.Handle("/user", utils.ValidateJWT(getUser))

	err := http.ListenAndServeTLS(":9000", "./certs/server.crt", "./certs/server.key", mux)

	if err != nil {
		log.Fatal(err)
	}

}
