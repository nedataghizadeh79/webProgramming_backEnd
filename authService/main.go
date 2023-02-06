package main

import (
	"AuthService/models"
	utils "AuthService/utils"
	"log"

	"net/http"

	"github.com/golang-jwt/jwt"
)

func signUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	utils.InsertUser(w, r)
}

func signIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	utils.SignInUser(w, r)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
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

func signOut(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	var token models.AuthToken = models.AuthToken{r.Header["Token"][0]}

	log_err := utils.AddExpiredToken(token)

	if log_err != nil {

	}

}

func main() {
	utils.ConnectToDb()
	utils.GetUserData("user:1", "data")

	mux := http.NewServeMux()

	mux.HandleFunc("/signUp", signUp)
	mux.HandleFunc("/signIn", signIn)
	mux.Handle("/user", utils.ValidateJWT(getUser))
	mux.HandleFunc("/logOut", signOut)

	err := http.ListenAndServeTLS(":9000", "./certs/server.crt", "./certs/server.key", mux)

	if err != nil {
		log.Fatal(err)
	}

}
