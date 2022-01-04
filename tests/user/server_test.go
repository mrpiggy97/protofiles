package user_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"reflect"
	"sync"
	"testing"

	"github.com/mrpiggy97/sharedProtofiles/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var listener *bufconn.Listener = bufconn.Listen(1024 * 1024)

func bufDialer(cxt context.Context, str string) (net.Conn, error) {
	return listener.Dial()
}

//wrapper for client and client connection
type userClient struct {
	client user.UserServiceClient
	conn   *grpc.ClientConn
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
// and consume the server above
func runClient(testChannel chan<- userClient) {
	var cxt context.Context = context.Background()
	// connection for test
	connection, connError := grpc.DialContext(
		cxt,
		"bufnet",
		grpc.WithInsecure(),
		grpc.WithContextDialer(bufDialer),
	)

	if connError != nil {
		panic("failed to establish connectioj between testing servers")
	}
	// consume server
	var client user.UserServiceClient = user.NewUserServiceClient(connection)
	testChannel <- userClient{
		client: client,
		conn:   connection,
	}
}

// TestGetUser will run server and client concurrently, it will
// test if we can connect,send a request and recieve a response
// from user.Server.GetUser

func TestGetUser(testCase *testing.T) {
	//run test servers
	var stopServer *sync.WaitGroup = new(sync.WaitGroup)
	stopServer.Add(1)
	var getClient chan userClient = make(chan userClient, 1)
	go runServer(stopServer)
	go runClient(getClient)

	//get client
	var client userClient = <-getClient

	//make request and recieve response
	var request *user.UserRequest = &user.UserRequest{
		UserId: "john fitz",
	}
	response, responseError := client.client.GetUser(
		context.Background(),
		request,
	)

	//make tests
	if responseError != nil {
		message := fmt.Sprintf("expected responseError to be nil, got %v instead",
			responseError)
		testCase.Error(message)
	}

	var expectedType string = "*user.User"
	if reflect.TypeOf(response).String() != expectedType {
		message := fmt.Sprintf("expected response to be of type %v, instead it is of type %v",
			expectedType, reflect.TypeOf(response).String())
		testCase.Error(message)
	}

	defer stopServer.Done()
	defer client.conn.Close()
}

func TestRegisterUsers(testCase *testing.T) {
	//run test servers
	var stopServer *sync.WaitGroup = new(sync.WaitGroup)
	stopServer.Add(1)
	var getClient chan userClient = make(chan userClient, 1)
	go runServer(stopServer)
	go runClient(getClient)

	//get client
	var client userClient = <-getClient

	//make stream of requests and test
	stream, streamErr := client.client.RegisterUsers(
		context.Background(),
	)

	if streamErr != nil {
		testCase.Error(streamErr)
	}

	for i := 0; i <= 10; i++ {
		var request *user.RegisterUserRequest = &user.RegisterUserRequest{
			Username: fmt.Sprintf("%v username", i),
		}
		sendingErr := stream.Send(request)

		if sendingErr != nil {
			testCase.Error(sendingErr)
		}
	}

	stream.CloseSend()
	var expectedType string = "*user.RegisterUserResponse"
	for {
		response, responseErr := stream.Recv()
		fmt.Println(responseErr, "error")
		if responseErr != nil && responseErr != io.EOF {
			testCase.Error("responseErr should only be io.EOF")
		}
		if responseErr == io.EOF {
			break
		}
		if reflect.TypeOf(response).String() != expectedType {
			message := fmt.Sprintf("expected response to be of type %v, instead it is of type %v",
				expectedType, reflect.TypeOf(response).String())
			testCase.Error(message)
		}
		fmt.Println(response.String())
	}
}
