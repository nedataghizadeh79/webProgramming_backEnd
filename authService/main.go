package main

import (
	"AuthService/models"
	utils "AuthService/utils"
	"log"
	"time"

	"encoding/json"
	"net/http"
)

func signUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var user models.UserAccount
	error := json.NewDecoder(r.Body).Decode(&user)
	if error != nil {
		http.Error(w, error.Error(), http.StatusBadRequest)
		return
	}

	err := utils.ValidateSignUpData(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.InsertUser(w, user)
}

func signIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user models.SignInInput
	error := json.NewDecoder(r.Body).Decode(&user)
	if error != nil {
		http.Error(w, error.Error(), http.StatusBadRequest)
		return
	}

	err := utils.ValidateSignInData(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SignInUser(w, user)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	utils.FindUser(w, r)
}

func signOut(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var token models.AuthToken = models.AuthToken{Token: r.Header["Token"][0]}

	log_err := utils.AddExpiredToken(token)

	if log_err != nil {
		http.Error(w, log_err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)

	redis_err := utils.AddToRedis(token.Token, "invalid", time.Hour)

	utils.CheckError(redis_err, nil)

}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/signUp", signUp)
	mux.HandleFunc("/signIn", signIn)
	mux.Handle("/user", utils.ValidateJWT(getUser))
	mux.Handle("/logOut", utils.ValidateJWT(signOut))

	err := http.ListenAndServeTLS(":9000", "./certs/server.crt", "./certs/server.key", mux)

	if err != nil {
		log.Fatal(err)
	}

}
