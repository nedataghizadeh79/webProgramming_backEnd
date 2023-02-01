package utils

import (
	models "AuthService/models"
	pb "AuthService/proto"
	"context"
	// "database/sql"
	// "encoding/json"
	"fmt"
	// "net/http"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

type RouteGuideServer struct {
	pb.UnimplementedRouteGuideServer
}


func GetUserData(username string, password string) {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	val, err := client.Get(ctx, username).Result()

	CheckError(err)
	if val != "" {
		fmt.Println(val)
	} else {
		fmt.Println("key not found")
	}
}

func (s *RouteGuideServer) SignUp(ctx context.Context, user *pb.SignUpUser) (*pb.Token, error) {

	db := ConnectToDb()
	defer db.Close()

	hashedPassword := HashPassword(user.Password)

	sqlStatement := "INSERT INTO user_account (email, phone_number, gender, first_name, last_name, password_hash) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := db.Exec(sqlStatement, user.Email, user.PhoneNumber, user.Gender, user.FirstName, user.LastName, hashedPassword)
	CheckError(err)
	token, _ := CreateToken(user.Email);

	return &pb.Token{Token: token, Error: "nothing"}, nil

}

func FindUser(user models.SignInInput) {
	db := ConnectToDb()
	defer db.Close()

	sqlStatement := "SELECT * FROM user_account WHERE email = $1"
	_, err := db.Query(sqlStatement, user.Email)
	CheckError(err)
}