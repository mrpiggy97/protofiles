package stringMethods

import (
	"context"
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
