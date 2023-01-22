package main
 
import (
    "net/http"
    "fmt"
    "utils"
)

func signUp(w http.ResponseWriter, r *http.Request) {}

func signIn(w http.ResponseWriter, r *http.Request) {}

func getUser(w http.ResponseWriter, r *http.Request) {}

func signOut(w http.ResponseWriter, r *http.Request) {}



func main() {
     
    utils.connectToDb()

    http.HandleFunc("/signUp", signUp)
    http.HandleFunc("/signIn", signIn)
    http.HandleFunc("/signOut", signOut)
    http.HandleFunc("/userInfo", getUser)
}
