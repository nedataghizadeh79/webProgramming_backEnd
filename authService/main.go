package main

import (
	utils "AuthService/utils"
	"log"

	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func signUp(w http.ResponseWriter, r *http.Request) {

	utils.InsertUser(w, r)
}

func check(w http.ResponseWriter, r *http.Request) {
	hmacSecret := []byte(utils.GetEnv("SECRET_JWT_KEY"))
	token, _ := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return hmacSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		w.Write([]byte(claims["sub"].(string)))
	}
}

func signIn(c *gin.Context) { fmt.Println("signin") }

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
	mux.Handle("/checkAuth", utils.ValidateJWT(check))

	err := http.ListenAndServe(":3333", mux)

	if err != nil {
		log.Fatal(err)
	}

}
