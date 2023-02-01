package main

import (
	utils "AuthService/utils"
	"net/http"

	"fmt"
	"net"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func signUp(w http.ResponseWriter, r *http.Request) {

	utils.InsertUser(w, r)
}

func signIn() { fmt.Println("signin") }

func getUser() { fmt.Println("user data") }

func signOut() {}

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
}

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		println("failed!")
	}

	grpcServer := grpc.NewServer()

	if err := grpcServer.Serve(lis); err != nil {
		println("failed!")
	}

}
