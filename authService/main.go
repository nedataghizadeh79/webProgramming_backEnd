package main

import (
	"AuthService/models"
	utils "AuthService/utils"
	"log"

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

	utils.ValidateSignUpData(user)

	utils.InsertUser(w, r)
}

func signIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	utils.SignInUser(w, r)
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

	if r.Header.Get("Token") == "" {
		http.Error(w, "Token not set", http.StatusUnauthorized)
	}

	var token models.AuthToken = models.AuthToken{Token: r.Header["Token"][0]}

	log_err := utils.AddExpiredToken(token)

	if log_err != nil {
		http.Error(w, log_err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)

}

func main() {

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
