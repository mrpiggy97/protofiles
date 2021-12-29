package user_test

import (
	"context"
	"fmt"
	"net"
	"sync"
	"testing"

	"github.com/mrpiggy97/shared-protofiles/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var listener *bufconn.Listener = bufconn.Listen(1024 * 1024)

func bufDialer(cxt context.Context, str string) (net.Conn, error) {
	return listener.Dial()
}

type serverTestResponse struct {
	response      *user.User
	responseError error
}

// server is a testing server meant to be run concurrently
// and consumed by a client.
func runServer(waiter *sync.WaitGroup) {
	var userServer *user.Server = &user.Server{}
	var grpcServer *grpc.Server = grpc.NewServer()
	user.RegisterUserServiceServer(grpcServer, userServer)
	grpcServer.Serve(listener)
	waiter.Wait()
	grpcServer.GracefulStop()
	defer listener.Close()
}

// client is a testing client server meant to run concurrently
// and consumed by server above.
func runClient(testChannel chan<- serverTestResponse, waiter *sync.WaitGroup) {
	var cxt context.Context = context.Background()
	// connection for test
	connection, connError := grpc.DialContext(
		cxt,
		"bufnet",
		grpc.WithInsecure(),
		grpc.WithContextDialer(bufDialer),
	)

	if connError != nil {
		panic(connError)
	}

	defer connection.Close()
	defer listener.Close()
	// consume server
	var client user.UserServiceClient = user.NewUserServiceClient(connection)
	var request *user.UserRequest = &user.UserRequest{
		UserId: "123234sdasd",
	}
	response, responseErr := client.GetUser(cxt, request)
	var testResponse serverTestResponse = serverTestResponse{
		response:      response,
		responseError: responseErr,
	}
	testChannel <- testResponse
	// tell server we are done using it and that it now can stop
	waiter.Done()
}

// GetUser will run server and client concurrently, it will
// test if we can connect,send a request and recieve a response
// from user.UserServer.go.
func GetUser(testCase *testing.T) {
	var serverChannel chan serverTestResponse = make(chan serverTestResponse, 2)
	var waiter *sync.WaitGroup = new(sync.WaitGroup)
	waiter.Add(1)
	go runServer(waiter)
	go runClient(serverChannel, waiter)
	var serverResponse serverTestResponse = <-serverChannel
	fmt.Println(serverResponse.response.String())

	if serverResponse.responseError != nil {
		testCase.Error(serverResponse.responseError)
	}
}

func TestServer(testCase *testing.T) {
	testCase.Run("Action=get-user", GetUser)
}
