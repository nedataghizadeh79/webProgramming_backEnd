package main
 
import (
    "net/http"
    "os"
    "fmt"
    "errors"
    utils "AuthService/utils"
)

func signUp(w http.ResponseWriter, r *http.Request) {fmt.Println("signup")}

func signIn(w http.ResponseWriter, r *http.Request) {fmt.Println("signin")}

func getUser(w http.ResponseWriter, r *http.Request) {fmt.Println("user data")}

func signOut(w http.ResponseWriter, r *http.Request) {}



func main() {
    utils.ConnectToDb()
    utils.GetUserData("user:1", "data")

    http.HandleFunc("/signUp", signUp)
    http.HandleFunc("/signIn", signIn)
    http.HandleFunc("/signOut", signOut)
    http.HandleFunc("/userInfo", getUser)

    err := http.ListenAndServe(":3333", nil)

    if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
