package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/mrpiggy97/shared-protofiles/stream"
	"google.golang.org/grpc/metadata"
)

type serverStreamMock struct {
}

func (mockServer serverStreamMock) Send(response *stream.RandomNumberResponse) error {
	return nil
}

func (mockServer serverStreamMock) Context() context.Context {
	return context.Background()
}

func (mockServer serverStreamMock) RecvMsg(obj interface{}) error {
	return nil
}

func (mockServer serverStreamMock) SendHeader(data metadata.MD) error {
	return nil
}

func (mockserver serverStreamMock) SendMsg(obj interface{}) error {
	return nil
}

func (mockServer serverStreamMock) SetHeader(data metadata.MD) error {
	return nil
}

func (mockServer serverStreamMock) SetTrailer(data metadata.MD) {
	fmt.Println("set trailer")
}

func StreamCheck(testCase *testing.T) {
	var testServer *stream.RandomNumberServer = &stream.RandomNumberServer{}
	var streamMock serverStreamMock = serverStreamMock{}
	var testRequest *stream.RandomNumberRequest = &stream.RandomNumberRequest{
		Number: 2000,
	}
	var resError = testServer.AddRandomNumber(testRequest, streamMock)
	if resError != nil {
		var message string = fmt.Sprintf("exepcted resError to be nil instead got %v", resError)
		testCase.Error(message)
	}
}

func TestRandomNumberServer(testCase *testing.T) {
	testCase.Run("Action=stream-check", StreamCheck)
}
