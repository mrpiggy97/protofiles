package stringMethods

import (
	"context"
	"fmt"
	"strings"
)

type FormattingServer struct {
	UnimplementedStringFormattingServer
}

func (server *FormattingServer) ToCamelCase(cxt context.Context, request *FormattingRequest) (*FormattingResponse, error) {
	var currentString string = request.GetStringToConvert()
	var stringSlice []string = strings.Split(currentString, "-")
	var formattedString string = ""
	for _, word := range stringSlice {
		for index, letter := range word {
			if index == 0 {
				formattedString = formattedString + strings.ToUpper(string(letter))
			} else {
				formattedString = formattedString + strings.ToLower(string(letter))
			}
		}
	}

	var response *FormattingResponse = &FormattingResponse{ConvertedString: formattedString}
	return response, nil
}

func (server *FormattingServer) ToLowerCase(cxt context.Context, req *FormattingRequest) (*FormattingResponse, error) {
	var stringToConvert string = req.GetStringToConvert()
	var stringSlice []string = strings.Split(stringToConvert, "-")
	var newString string = ""
	for _, word := range stringSlice {
		newString = newString + strings.ToLower(word)
	}

	var response *FormattingResponse = &FormattingResponse{
		ConvertedString: newString,
	}
	fmt.Println("server recieved ", req)
	fmt.Println("server sent ", response)
	return response, nil
}

func (server *FormattingServer) ToUpperCase(cxt context.Context, req *FormattingRequest) (*FormattingResponse, error) {
	var stringToConvert string = req.GetStringToConvert()
	var stringSlice []string = strings.Split(stringToConvert, "-")
	var newString string = ""
	for _, word := range stringSlice {
		newString = newString + strings.ToUpper(word)
	}

	var response *FormattingResponse = &FormattingResponse{
		ConvertedString: newString,
	}

	fmt.Println("server recieved ", req)
	fmt.Println("server sent ", response)

	return response, nil
}
