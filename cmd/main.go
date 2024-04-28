package main

import (
	"log"
	"net"
	"os"

	"github.com/akshaybt001/DatingApp_UserService/db"
	"github.com/akshaybt001/DatingApp_UserService/initializer"
	"github.com/akshaybt001/DatingApp_proto_files/pb"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf(err.Error())
	}
	addr := os.Getenv("DB_KEY")
	DB, err := db.InitDB(addr)
	if err != nil {
		log.Fatal(err.Error())
	}
	services := initializer.Initializer(DB)
	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, services)
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("failed to listen on port 8081 %v", err)
	}
	log.Printf("user service listening on port 8081")
	if err = server.Serve(listener); err != nil {
		log.Fatalf("failed to listen on port 8081 %v", err)
	}
}
