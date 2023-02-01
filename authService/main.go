package main

import (
	utils "AuthService/utils"
	pb "AuthService/proto"
	// "net/http"

	"fmt"
	"net"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)


// func (s *routeGuideService) signUp(ctx context.Context, user *pb.SignUpUser) (*pb.Token, error) {

// 	return utils.InsertUser(user)
// }

// func signIn() { fmt.Println("signin") }

// func getUser() { fmt.Println("user data") }

// func signOut() {}

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

	s := utils.RouteGuideServer{}
	grpcServer := grpc.NewServer()
	
	pb.RegisterRouteGuideServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		println("failed!")
	}

}
