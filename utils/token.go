package utils

import (
	"fmt"
	"net/http"
	"time"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

var key = []byte(GetEnv("SECRET_JWT_KEY"))

func CreateToken(userId string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	claims["sub"] = payload

	tokenStr, err := token.SignedString(GetEnv("PRIVATE_JWT_TOKEN"))

	if err != nil {
		return "", err
	}

	return tokenStr, nil

}

func ValidateJWT(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(t *jwt.Token) (interface{}, error) {
				_, ok := t.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("not authorized"))
				}
				return GetEnv("PRIVATE_JWT_KEY"), nil
			})

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("not authorized: " + err.Error()))
			}

			if token.Valid {
				next(w, r)
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("not authorized"))
		}
	})
}

func GetJwt(w http.ResponseWriter, r *http.Request) {
	if r.Header["Access"] != nil {
		token, err := CreateToken()
		if err != nil {
			return
		}
		fmt.Fprint(w, token)
	}
}
