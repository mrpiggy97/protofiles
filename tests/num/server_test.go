package num_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"reflect"
	"sync"
	"testing"

	"github.com/mrpiggy97/sharedProtofiles/num"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var listener *bufconn.Listener = bufconn.Listen(1024 * 1024)

type responseWrapper struct {
	client num.NumServiceClient
	conn   *grpc.ClientConn
}

func bufDialer(cxt context.Context, str string) (net.Conn, error) {
	return listener.Dial()
}

func runServer(waiter *sync.WaitGroup) {
	var grpcServer *grpc.Server = grpc.NewServer()
	var server *num.Server = new(num.Server)
	num.RegisterNumServiceServer(grpcServer, server)
	grpcServer.Serve(listener)
	waiter.Wait()
	grpcServer.GracefulStop()
	defer listener.Close()
}

func runClient(getClient chan<- responseWrapper) {
	connection, connectionError := grpc.DialContext(
		context.Background(),
		"buffnet",
		grpc.WithInsecure(),
		grpc.WithContextDialer(bufDialer),
	)

	if connectionError != nil {
		panic("failure to establish connection between testing servers")
	}

	var client num.NumServiceClient = num.NewNumServiceClient(connection)
	var resWrapper responseWrapper = responseWrapper{
		client: client,
		conn:   connection,
	}

	getClient <- resWrapper
}

func TestRnd(testCase *testing.T) {
	//run test servers
	var waiter *sync.WaitGroup = new(sync.WaitGroup)
	waiter.Add(1)
	var getClient chan responseWrapper = make(chan responseWrapper, 1)
	go runServer(waiter)
	go runClient(getClient)

	//get client
	var client responseWrapper = <-getClient

	//make request
	var request *num.NumRequest = &num.NumRequest{
		From:   0,
		To:     40,
		Number: 100,
	}

	stream, streamError := client.client.Rnd(context.Background(), request)
	if streamError != nil {
		message := fmt.Sprintf("expected streamError to be nil, got %v", streamError)
		testCase.Error(message)
	}
	stream.CloseSend()

	var expectedType string = "*num.NumResponse"

	for {
		response, resError := stream.Recv()
		if resError != nil && resError != io.EOF {
			testCase.Error("resError should only be io.EOF")
		}

		if resError == io.EOF {
			break
		}

		if reflect.TypeOf(response).String() != expectedType {
			message := fmt.Sprintf("response type expected to be %v, instead got %v",
				expectedType, reflect.TypeOf(response).String())
			testCase.Error(message)
		}
		fmt.Println(response.String())
	}

	defer waiter.Done()
	defer client.conn.Close()
}
