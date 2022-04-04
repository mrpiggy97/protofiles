package primes_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"testing"

	"github.com/mrpiggy97/sharedProtofiles/primes"
	"google.golang.org/grpc"
)

type primesClient struct {
	conn   *grpc.ClientConn
	client primes.PrimeServiceClient
}

func runServer(waiter *sync.WaitGroup, listener net.Listener) {
	waiter.Add(1)
	var grpcServer *grpc.Server = grpc.NewServer()
	var primeServer *primes.Server = &primes.Server{}
	primes.RegisterPrimeServiceServer(grpcServer, primeServer)
	grpcServer.Serve(listener)
	waiter.Wait()
	grpcServer.GracefulStop()
}

func getClient(sendClientChannel chan<- primesClient) {
	conn, _ := grpc.Dial("localhost:8000", grpc.WithInsecure())
	var serviceClient primes.PrimeServiceClient = primes.NewPrimeServiceClient(conn)
	var mainClient primesClient = primesClient{
		conn:   conn,
		client: serviceClient,
	}
	sendClientChannel <- mainClient
}

func TestPrimesServer(testCase *testing.T) {
	//run testing servers
	var waiter *sync.WaitGroup = new(sync.WaitGroup)
	waiter.Add(1)
	listener, listenErr := net.Listen("tcp", "localhost:8000")
	if listenErr != nil {
		testCase.Error("failed to establish a listener")
	}
	go runServer(waiter, listener)
	//get primesClient
	var getClientChannel chan primesClient = make(chan primesClient, 1)
	getClient(getClientChannel)
	var testClient primesClient = <-getClientChannel
	stream, _ := testClient.client.GetCount(context.Background())
	for i := 100; i <= 1000; i = i + 100 {
		var newRequest *primes.PrimeRequest = &primes.PrimeRequest{
			Number: int64(i),
		}
		var sendingError error = stream.Send(newRequest)
		if sendingError != nil {
			testCase.Error(sendingError)
		}
	}
	stream.CloseSend()
	for {
		res, resError := stream.Recv()
		if resError != nil && resError != io.EOF {
			testCase.Error(resError)
		}
		if resError == io.EOF {
			fmt.Println("prime count finished")
			break
		}
		fmt.Println("count is ", res.Count)
	}
	defer testClient.conn.Close()
	defer waiter.Done()
}
