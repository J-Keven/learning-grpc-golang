package main

import (
	"log"
	"net"

	pb "github.com/j-keven/learning-grpc-golang/pb"
	"github.com/j-keven/learning-grpc-golang/services"
	"google.golang.org/grpc"
)

func main() {

	listern, err := net.Listen("tcp", ":3333")

	if err != nil {
		log.Fatalf("Cauld not connetc: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterUserServiceServer(grpcServer, services.NewUserService())

	if err := grpcServer.Serve(listern); err != nil {
		log.Fatalf("could not server: %v", err)
	}
}
