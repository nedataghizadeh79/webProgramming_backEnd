package main
 
import (
    "net/http"
    "database/sql"
    "fmt"
    "github.com/go-redis/redis"
    _ "github.com/lib/pq"
)

func signUp(w http.ResponseWriter, r *http.Request) {}

func signIn(w http.ResponseWriter, r *http.Request) {}

func getUser(w http.ResponseWriter, r *http.Request) {}

func signOut(w http.ResponseWriter, r *http.Request) {}



func main() {
     

    http.handleFunc("/signUp", signUp)
    http.handleFunc("/signIn", signIn)
    http.handleFunc("/signOut", signOut)
    http.handleFunc("/userInfo", getUser)
}
