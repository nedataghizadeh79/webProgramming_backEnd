package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var secret = []byte("secret-auth-token")

func CreateToken() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour).Unix()

	tokenStr, err := token.SignedString(secret)

	if err != nil {
		return "", err
	}

	return tokenStr, nil

}
