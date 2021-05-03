package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	uuid "github.com/google/uuid"
	pb "github.com/j-keven/learning-grpc-golang/pb"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func NewUserService() *UserService {
	return &UserService{}
}

// Request -> Response
func (*UserService) AddUser(ctx context.Context, req *pb.User) (*pb.User, error) {

	// insert into database
	fmt.Println(req)

	return &pb.User{
		Id:    "1",
		Name:  req.GetName(),
		Email: req.GetEmail(),
	}, nil
}

//Request -> Response(send partial response using data streams)
func (*UserService) AddUserVerbose(req *pb.User, stream pb.UserService_AddUserVerboseServer) error {
	stream.Send(&pb.UserResultStream{
		Status: "init",
		User:   &pb.User{},
	})

	time.Sleep(time.Second * 3)

	stream.Send(&pb.UserResultStream{
		Status: "inserting in to database",
		User: &pb.User{
			Id:    req.GetId(),
			Name:  req.GetName(),
			Email: req.GetEmail(),
		},
	})

	time.Sleep(time.Second * 3)

	stream.Send(&pb.UserResultStream{
		Status: "inserting in to database finished",
		User: &pb.User{
			Id:    "123",
			Name:  req.GetName(),
			Email: req.GetEmail(),
		},
	})

	time.Sleep(time.Second * 3)

	stream.Send(&pb.UserResultStream{
		Status: "Completed",
		User: &pb.User{
			Id:    "123",
			Name:  req.GetName(),
			Email: req.GetEmail(),
		},
	})

	time.Sleep(time.Second * 3)
	// response finished
	return nil
}

// Request(data stream) -> Response
func (*UserService) AddUsers(stream pb.UserService_AddUsersServer) error {
	users := []*pb.User{}

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&pb.Users{
				Users: users,
			})
		}

		if err != nil {
			log.Fatalf("Error from data stream: %v", err)
		}

		u1 := uuid.NewString()

		newUser := &pb.User{
			Id:    u1,
			Name:  req.GetName(),
			Email: req.GetEmail(),
		}

		users = append(users, newUser)
		fmt.Println("inset a new users:", newUser.Name)
	}
}

func (*UserService) AddUserStreamBoth(stream pb.UserService_AddUserStreamBothServer) error {
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error receiving stream from the client: %v", err)
		}

		u1 := uuid.NewString()
		err = stream.Send(&pb.UserResultStream{
			Status: "Added",
			User: &pb.User{
				Id:    u1,
				Name:  req.GetName(),
				Email: req.GetEmail(),
			},
		})

		if err != nil {
			log.Fatalf("Error sendng stream from the client: %v", err)
		}
	}
}
