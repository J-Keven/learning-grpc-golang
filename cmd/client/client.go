package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/j-keven/learning-grpc-golang/pb"
	"google.golang.org/grpc"
)

func main() {
	connectio, err := grpc.Dial("localhost:3333", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect in gRPC server: %v", err)
	}

	client := pb.NewUserServiceClient(connectio)
	// AddUser(client)
	// AddUserVerbose(client)
	// AddUsers(client)
	AddUserStreamBoth(client)
}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Jhon",
		Email: "jhon@jhon.com",
	}

	res, err := client.AddUser(context.Background(), req)

	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	fmt.Println(res)
}

// Recive msg in data stream
func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Jhon",
		Email: "jhon@jhon.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)

	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not recive the msg: %v", err)
		}

		fmt.Println("Status:", stream.Status, "- User:", stream.GetUser())
	}
}

// Send users by data stream and receive all inserted users
func AddUsers(client pb.UserServiceClient) {

	users := []*pb.User{
		&pb.User{
			Id:    "1",
			Name:  "Jhon 01",
			Email: "jhon1@ex.com",
		},
		&pb.User{
			Id:    "2",
			Name:  "Jhon 02",
			Email: "jhon2@ex.com",
		},
		&pb.User{
			Id:    "3",
			Name:  "Jhon 03",
			Email: "jhon3@ex.com",
		},
		&pb.User{
			Id:    "4",
			Name:  "Jhon 04",
			Email: "jhon4@ex.com",
		},
	}
	stream, err := client.AddUsers(context.Background())

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, user := range users {
		stream.Send(user)

		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient) {
	users := []*pb.User{
		&pb.User{
			Id:    "1",
			Name:  "Jhon 01",
			Email: "jhon1@ex.com",
		},
		&pb.User{
			Id:    "2",
			Name:  "Jhon 02",
			Email: "jhon2@ex.com",
		},
		&pb.User{
			Id:    "3",
			Name:  "Jhon 03",
			Email: "jhon3@ex.com",
		},
		&pb.User{
			Id:    "4",
			Name:  "Jhon 04",
			Email: "jhon4@ex.com",
		},
	}
	stream, err := client.AddUserStreamBoth(context.Background())

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	await := make(chan int)
	go func() {
		for _, user := range users {
			stream.Send(user)
			fmt.Println("Send: ", user.Name)

			time.Sleep(time.Second * 3)
		}

		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}

			fmt.Println("Received user:", res.User.GetName())
		}

		close(await)

	}()

	<-await
}
