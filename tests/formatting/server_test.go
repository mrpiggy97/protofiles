package formatting_test

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/mrpiggy97/shared-protofiles/formatting"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type formattingClient struct {
	client formatting.StringFormattingClient
	conn   *grpc.ClientConn
}

var listener *bufconn.Listener = bufconn.Listen(1024 * 1024)

func formattingBufDialer(cxt context.Context, str string) (net.Conn, error) {
	return listener.Dial()
}

func runServer(closeServer <-chan bool) {
	var formattingServer *formatting.Server = new(formatting.Server)
	var grpcServer *grpc.Server = grpc.NewServer()
	formatting.RegisterStringFormattingServer(grpcServer, formattingServer)
	var err error = grpcServer.Serve(listener)
	fmt.Println(err)
	<-closeServer
	grpcServer.GracefulStop()
	defer listener.Close()
}

func runClient(clientChannel chan<- formattingClient) {
	var cxt context.Context = context.Background()
	connection, connError := grpc.DialContext(
		cxt,
		"bufnet",
		grpc.WithInsecure(),
		grpc.WithContextDialer(formattingBufDialer),
	)

	if connError != nil {
		panic("failed to establish connection between testing server and client")
	}

	var client formatting.StringFormattingClient = formatting.NewStringFormattingClient(
		connection,
	)
	var formatClient formattingClient = formattingClient{
		client: client,
		conn:   connection,
	}
	clientChannel <- formatClient
}

func toCamelCase(testCase *testing.T) {

	//prepare server and client
	var closeServer chan bool = make(chan bool, 1)
	var getClient chan formattingClient = make(chan formattingClient, 1)
	go runServer(closeServer)
	go runClient(getClient)

	//recieve client
	var client formattingClient = <-getClient

	//make request to server
	var request *formatting.FormattingRequest = &formatting.FormattingRequest{
		StringToConvert: "FABIAN-jEsus-rivas",
	}
	response, resError := client.client.ToCamelCase(context.Background(), request)

	//make tests
	var expectedResponse string = "FabianJesusRivas"
	if resError != nil {
		message := fmt.Sprintf("expected resError to be nil,instead got %v", resError)
		testCase.Error(message)
	}

	if response.GetConvertedString() != expectedResponse {
		message := fmt.Sprintf("expected response.ConvertedString to be %v, instead got %v",
			expectedResponse, response.ConvertedString,
		)

		testCase.Error(message)
	}

	//close servers
	closeServer <- true
	defer client.conn.Close()
}

func toLowerCase(testCase *testing.T) {
	//set up servers
	var closeServer chan bool = make(chan bool, 1)
	var getClient chan formattingClient = make(chan formattingClient, 1)
	go runServer(closeServer)
	go runClient(getClient)

	//get client
	var client formattingClient = <-getClient

	//run request
	var request *formatting.FormattingRequest = &formatting.FormattingRequest{
		StringToConvert: "FABIAN-jeSus-rivas",
	}

	response, resErr := client.client.ToLowerCase(context.Background(), request)

	//run tests
	var expectedResponse string = "fabianjesusrivas"

	if resErr != nil {
		message := fmt.Sprintf("expected resError to be nil, instead got %v", resErr)
		testCase.Error(message)
	}

	if response.GetConvertedString() != expectedResponse {
		message := fmt.Sprintf("expected response.ConvertedString to be %v,instead got %v",
			response.GetConvertedString(), expectedResponse)
		testCase.Error(message)
	}

	//close servers
	closeServer <- true
	defer client.conn.Close()
}

func toUpperCase(testCase *testing.T) {

	//run servers
	var stopServer chan bool = make(chan bool, 1)
	var getClient chan formattingClient = make(chan formattingClient, 1)
	go runServer(stopServer)
	go runClient(getClient)

	//get client
	var formatClient formattingClient = <-getClient

	//make request
	var request *formatting.FormattingRequest = &formatting.FormattingRequest{
		StringToConvert: "fabian-JeSus-RIVAS",
	}

	response, resError := formatClient.client.ToUpperCase(context.Background(), request)

	//run tests
	var expectedResponse string = "FABIANJESUSRIVAS"

	if resError != nil {
		message := fmt.Sprintf("expected resError to be nil, instead got %v", expectedResponse)
		testCase.Error(message)
	}

	if response.GetConvertedString() != expectedResponse {
		message := fmt.Sprintf("expected response.ConvertedString to be %v, instead got %v",
			expectedResponse, response.GetConvertedString())
		testCase.Error(message)
	}

	//close servers
	stopServer <- true
	defer formatClient.conn.Close()
}
func TestServer(testCase *testing.T) {
	testCase.Run("Action=to-camel-case", toCamelCase)
	testCase.Run("Action=to-lower-case", toLowerCase)
	testCase.Run("Action=to-upper-case", toUpperCase)
}
