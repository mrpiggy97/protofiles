package calculation_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"reflect"
	"sync"
	"testing"

	"github.com/mrpiggy97/sharedProtofiles/calculation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var listener *bufconn.Listener = bufconn.Listen(1024 * 1024)

type responseWrap struct {
	client calculation.CalculationServiceClient
	conn   *grpc.ClientConn
}

func bufDialer(cxt context.Context, str string) (net.Conn, error) {
	return listener.Dial()
}

// runServer will be running a testing server that will recieve
// requests and send responses.
func runServer(closeServer *sync.WaitGroup) {
	var grpcServer *grpc.Server = grpc.NewServer()
	var server *calculation.Server = new(calculation.Server)
	calculation.RegisterCalculationServiceServer(grpcServer, server)
	grpcServer.Serve(listener)
	closeServer.Wait()
	grpcServer.GracefulStop()
	defer listener.Close()
}

// runClient will run a testing server for the client to
// make requests and recieve responses.
func runClient(sendClient chan<- responseWrap) {
	connection, connectionError := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithInsecure(),
		grpc.WithContextDialer(bufDialer),
	)

	if connectionError != nil {
		panic("failed to establish connection between testing servers")
	}

	var client calculation.CalculationServiceClient = calculation.NewCalculationServiceClient(
		connection,
	)

	sendClient <- responseWrap{
		client: client,
		conn:   connection,
	}
}

// TestSumStream will test Server.SumStream, it will be sending
// a stream of requests, and at the same time a stream of responses.
func TestSumStream(testCase *testing.T) {
	//run test servers
	var getClient chan responseWrap = make(chan responseWrap, 1)
	var closeServer *sync.WaitGroup = new(sync.WaitGroup)
	closeServer.Add(1)
	go runServer(closeServer)
	go runClient(getClient)

	//get client
	var client responseWrap = <-getClient
	stream, streamErr := client.client.SumStream(
		context.Background(),
	)
	if streamErr != nil {
		testCase.Error("streamErr expected to be nil,instead got ", streamErr)
	}
	//make requests
	for i := 0; i < 10; i++ {
		var request *calculation.SumStreamRequest = &calculation.SumStreamRequest{
			A: int32(i),
			B: int32(i + 5),
		}
		//this line sends request
		requestError := stream.Send(request)
		if requestError != nil {
			testCase.Error("resError expected to be nil,instead got ", requestError)
		}

	}
	stream.CloseSend()

	//consume responses
	for {
		res, resError := stream.Recv()
		if resError != nil && resError != io.EOF {
			testCase.Error(resError)
		}
		if resError == io.EOF {
			break
		}
		fmt.Println("client recieved ", res.String())
		var expectedType string = "*calculation.SumStreamResponse"
		if reflect.TypeOf(res).String() != expectedType {
			testCase.Error("expected type of res to be ", expectedType)
		}
	}

	defer closeServer.Done()
	defer client.conn.Close()
}
