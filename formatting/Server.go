package formatting

import (
	"context"
	"fmt"
	"strings"
)

type Server struct {
	UnimplementedStringFormattingServer
}

func (server *Server) ToCamelCase(cxt context.Context, request *FormattingRequest) (*FormattingResponse, error) {
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

func (server *Server) ToLowerCase(cxt context.Context, req *FormattingRequest) (*FormattingResponse, error) {
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

func (server *Server) ToUpperCase(cxt context.Context, req *FormattingRequest) (*FormattingResponse, error) {
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
