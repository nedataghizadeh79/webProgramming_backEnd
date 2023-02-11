package utils

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
)

var key = []byte(GetEnv("SECRET_JWT_KEY"))

func verifyJwtExpiration(token string) (string, error) {
	sqlStatement := "SELECT token FROM unauthorized_token WHERE token = $1"
	db := ConnectToDb(nil)
	defer db.Close()

	row := db.QueryRow(sqlStatement, token)

	var foundToken string

	switch err := row.Scan(&foundToken); err {
	case sql.ErrNoRows:
		return "valid", nil
	case nil:
		return "invalid", nil
	default:
		return "", err

	}

}

func CreateToken(email string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	claims["sub"] = email

	tokenStr, err := token.SignedString(key)

	if err != nil {
		return "", err
	}

	return tokenStr, nil

}

func ValidateJWT(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenStr := r.Header["Token"][0]
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				_, ok := t.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("not authorized"))
				}
				return key, nil
			})

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("not authorized: " + err.Error()))
			}

			switch val, redErr := GetFromRedis(tokenStr); redErr {
			case redis.Nil:
				result, checkErr := verifyJwtExpiration(tokenStr)
				if checkErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("error in server"))
					return
				} else {
					AddToRedis(tokenStr, result, time.Hour)
					if result == "invalid" {
						w.WriteHeader(http.StatusUnauthorized)
						w.Write([]byte("not authorized"))
						return
					}
				}
			case nil:
				if val == "invalid" {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("not authorized"))
					return
				}
			default:
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("error in server"))
				return
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
