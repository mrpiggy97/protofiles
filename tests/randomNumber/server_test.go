package randomNumber_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"reflect"
	"sync"
	"testing"

	"github.com/mrpiggy97/sharedProtofiles/randomNumber"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var listener *bufconn.Listener = bufconn.Listen(1024 * 1024)

type randomNumberClient struct {
	client randomNumber.RandomServiceClient
	conn   *grpc.ClientConn
}

func bufDialer(ctx context.Context, str string) (net.Conn, error) {
	return listener.Dial()
}

// runServer will run a test server to be consumed by tests.
func runServer(waiter *sync.WaitGroup) {
	var grpcServer *grpc.Server = grpc.NewServer()
	var server *randomNumber.Server = new(randomNumber.Server)
	randomNumber.RegisterRandomServiceServer(grpcServer, server)
	grpcServer.Serve(listener)
	waiter.Wait()
	grpcServer.GracefulStop()
	defer listener.Close()
}

// runClient will run a test server to be consumed by tests.
func runClient(getClient chan<- randomNumberClient) {
	connection, connError := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithInsecure(),
		grpc.WithContextDialer(bufDialer),
	)

	if connError != nil {
		panic("failed to establish connection between server and client")
	}

	var client randomNumber.RandomServiceClient = randomNumber.NewRandomServiceClient(
		connection,
	)

	getClient <- randomNumberClient{
		client: client,
		conn:   connection,
	}
}

// TestAddRandomNumer will test randomNumber.Server.AddRandomNumber method
// Methods that return or have a stream as a request cannot be run effectively
// by testCase.run(),so it's best to run every test for every method of randomNumber.Server
// on their own rather then group them in a single test.
func TestAddRandomNumber(testCase *testing.T) {
	//run server and client
	var closeServer *sync.WaitGroup = new(sync.WaitGroup)
	closeServer.Add(1)
	var getClient chan randomNumberClient = make(chan randomNumberClient, 1)
	go runServer(closeServer)
	go runClient(getClient)

	//get client
	var client randomNumberClient = <-getClient

	//make request
	var request *randomNumber.RandomNumberRequest = &randomNumber.RandomNumberRequest{
		Number: 1000,
	}

	stream, streamError := client.client.AddRandomNumber(
		context.Background(),
		request,
	)

	if streamError != nil {
		panic("failed to establish connection with test server")
	}

	//consume stream
	var expectedType string = "*randomNumber.RandomNumberResponse"
	for {
		response, responseError := stream.Recv()
		if responseError != nil && responseError != io.EOF {
			testCase.Error("responseError can only be io.EOF")
		}
		if responseError == io.EOF {
			break
		}
		fmt.Println(response.String())
		if reflect.TypeOf(response).String() != expectedType {
			message := fmt.Sprintf("expected response to be of type %v,instead it is %v",
				expectedType, reflect.TypeOf(response).String())
			testCase.Error(message)
			break
		}

	}
	//close servers
	defer closeServer.Done()
	defer client.conn.Close()
}

// TestSubstractRandomNumber will test randomNumber.Server.SubstractRandomNumber method
// Methods that return or have a stream as a request cannot be run effectively
// by testCase.run(),so it's best to run every test for every method of randomNumber.Server
// on their own rather then group them in a single test.
func TestSubstractRandomNumber(testCase *testing.T) {
	//run test servers
	var getClient chan randomNumberClient = make(chan randomNumberClient, 1)
	var closeServer *sync.WaitGroup = new(sync.WaitGroup)
	closeServer.Add(1)
	go runServer(closeServer)
	go runClient(getClient)

	//get client
	var client randomNumberClient = <-getClient

	//make request
	var request *randomNumber.RandomNumberRequest = &randomNumber.RandomNumberRequest{
		Number: 20,
	}
	stream, streamError := client.client.SubstractRandomNumber(
		context.Background(),
		request,
	)

	if streamError != nil {
		panic("failed to establish connection between test servers")
	}

	//make tests
	var expectedType string = "*randomNumber.RandomNumberResponse"
	for {
		response, responseError := stream.Recv()
		if responseError != nil && responseError != io.EOF {
			testCase.Error("responseError should only be nil or io.EOF,instead got ", responseError)
		}
		if responseError == io.EOF {
			break
		}
		fmt.Println(response.String())
		if reflect.TypeOf(response).String() != expectedType {
			message := fmt.Sprintf("expected type of response to be %v, instea got %v",
				expectedType, reflect.TypeOf(response))
			testCase.Error(message)
		}
	}

	defer closeServer.Done()
	defer client.conn.Close()
}
