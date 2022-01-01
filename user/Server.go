package user

import (
	"context"
	"fmt"
	"io"
)

type Server struct {
	UnimplementedUserServiceServer
}

// GetUser will just send back a user
func (server *Server) GetUser(cxt context.Context, req *UserRequest) (*User, error) {
	fmt.Println("server recieved:", req.String())
	var newUser *User = &User{UserId: "John", Email: "example@email.com"}
	return newUser, nil
}

func (server *Server) RegisterUsers(stream UserService_RegisterUsersServer) error {
	for {
		currentRequest, requestError := stream.Recv()
		if requestError != nil && requestError != io.EOF {
			return requestError
		}
		fmt.Println("server recieved: ", currentRequest)
		if requestError == io.EOF {
			break
		}

		var newUser *User = &User{
			UserId: currentRequest.Username,
			Email:  "",
		}

		var response *RegisterUserResponse = &RegisterUserResponse{
			User: newUser,
		}

		stream.Send(response)
	}
	return nil
}
