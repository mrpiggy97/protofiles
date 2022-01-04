package randomNumber

import (
	"math/rand"
)

type Server struct {
	UnimplementedRandomServiceServer
}

func (server *Server) AddRandomNumber(request *RandomNumberRequest, stream RandomService_AddRandomNumberServer) error {
	var randomInt int64 = rand.Int63()
	for i := 0; i < 100; i++ {
		var response *RandomNumberResponse = &RandomNumberResponse{
			OriginalNumber: request.Number,
			TotalNumber:    randomInt + int64(request.Number),
		}

		var streamErr error = stream.Send(response)
		if streamErr != nil {
			return streamErr
		}
	}

	return nil
}

func (server *Server) SubstractRandomNumber(request *RandomNumberRequest, stream RandomService_SubstractRandomNumberServer) error {
	var randomInt int64 = rand.Int63()
	for i := 0; i <= 100; i++ {
		var response *RandomNumberResponse = &RandomNumberResponse{
			OriginalNumber: request.Number,
			TotalNumber:    (randomInt - int64(i)) - request.Number,
		}

		var streamErr error = stream.Send(response)
		if streamErr != nil {
			return streamErr
		}
	}

	return nil
}
