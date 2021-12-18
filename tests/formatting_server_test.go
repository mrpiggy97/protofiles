package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/mrpiggy97/shared-protofiles/stringMethods"
)

func ToCamelCase(testCase *testing.T) {
	var emptyContext context.Context = context.Background()
	var testRequest *stringMethods.FormattingRequest = &stringMethods.FormattingRequest{
		StringToConvert: "FABIAN-jEsuS-rivas",
	}
	var server *stringMethods.FormattingServer = &stringMethods.FormattingServer{}
	_, resError := server.ToCamelCase(emptyContext, testRequest)
	if resError != nil {
		var message string = fmt.Sprintf("expected resError to be nil,got %v instead", resError)
		testCase.Error(message)
	}
}

func TestFormattingServer(testCase *testing.T) {
	testCase.Run("Action=to-camel-case", ToCamelCase)
}
