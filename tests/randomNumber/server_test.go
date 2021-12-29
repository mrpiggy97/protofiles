package randomNumber_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"reflect"
	"testing"

	"github.com/mrpiggy97/shared-protofiles/randomNumber"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var listener *bufconn.Listener = bufconn.Listen(1024 * 1024)

type randomNumberClient struct {
	client randomNumber.RandomServiceClient
	conn   *grpc.ClientConn
}

type responseWrapper struct {
	response      *randomNumber.RandomNumberResponse
	responseError error
}

func consumeStream(sendResponse chan<- responseWrapper, stream randomNumber.RandomService_AddRandomNumberClient) {
	for {
		res, resErr := stream.Recv()
		if resErr != nil && resErr != io.EOF {
			sendResponse <- responseWrapper{
				response:      res,
				responseError: resErr,
			}
			close(sendResponse)
			break
		}

		if resErr == io.EOF {
			close(sendResponse)
			break
		}

		sendResponse <- responseWrapper{
			response:      res,
			responseError: nil,
		}
	}
}

func bufDialer(ctx context.Context, str string) (net.Conn, error) {
	return listener.Dial()
}

func runServer(stopServer <-chan bool) {
	var grpcServer *grpc.Server = grpc.NewServer()
	var server *randomNumber.Server = new(randomNumber.Server)
	randomNumber.RegisterRandomServiceServer(grpcServer, server)
	grpcServer.Serve(listener)
	<-stopServer
	grpcServer.GracefulStop()
}

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

func TestAddRandomNumber(testCase *testing.T) {
	//run server and client
	var closeServer chan bool = make(chan bool, 1)
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
	var resChannel chan responseWrapper = make(chan responseWrapper, 1)
	var expectedType string = "*randomNumber.RandomNumberResponse"
	go consumeStream(resChannel, stream)
	for {
		resWrapper, channelAvailable := <-resChannel
		if channelAvailable {
			if resWrapper.responseError != nil {
				message := fmt.Sprintf("resWrapper.responseError should be nil, got %v instead",
					resWrapper.responseError,
				)
				testCase.Error(message)
			}
			if reflect.TypeOf(resWrapper.response).String() != expectedType {
				message := fmt.Sprintf("expected resWrapper.response to have type %v, instead",
					expectedType)
				message = fmt.Sprintf("%v, it has %v", message, reflect.TypeOf(resWrapper.response).String())
				testCase.Error(message)
			}
		} else {
			break
		}
	}
	//close servers
	closeServer <- true
	defer client.conn.Close()
}
