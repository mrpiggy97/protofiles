package tests

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/mrpiggy97/shared-protofiles/stringMethods"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type formattingClient struct {
	client stringMethods.StringFormattingClient
	conn   *grpc.ClientConn
}

var formattingServerListener *bufconn.Listener = bufconn.Listen(BufSize)

func formattingBufDialer(cxt context.Context, str string) (net.Conn, error) {
	return formattingServerListener.Dial()
}

func runFormattingServer(closeServer <-chan bool) {
	var formattingServer *stringMethods.FormattingServer = new(stringMethods.FormattingServer)
	var grpcServer *grpc.Server = grpc.NewServer()
	stringMethods.RegisterStringFormattingServer(grpcServer, formattingServer)
	var err error = grpcServer.Serve(formattingServerListener)
	fmt.Println(err)
	<-closeServer
	grpcServer.GracefulStop()
	defer formattingServerListener.Close()
}

func runFormattingClient(clientChannel chan<- formattingClient) {
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

	var client stringMethods.StringFormattingClient = stringMethods.NewStringFormattingClient(
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
	go runFormattingServer(closeServer)
	go runFormattingClient(getClient)

	//recieve client
	var client formattingClient = <-getClient

	//make request to server
	var request *stringMethods.FormattingRequest = &stringMethods.FormattingRequest{
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
	go runFormattingServer(closeServer)
	go runFormattingClient(getClient)

	//get client
	var client formattingClient = <-getClient

	//run request
	var request *stringMethods.FormattingRequest = &stringMethods.FormattingRequest{
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
	go runFormattingServer(stopServer)
	go runFormattingClient(getClient)

	//get client
	var formatClient formattingClient = <-getClient

	//make request
	var request *stringMethods.FormattingRequest = &stringMethods.FormattingRequest{
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
func TestFormattingServer(testCase *testing.T) {
	testCase.Run("Action=to-camel-case", toCamelCase)
	testCase.Run("Action=to-lower-case", toLowerCase)
	testCase.Run("Action=to-upper-case", toUpperCase)
}
