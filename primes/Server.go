package primes

import (
	"fmt"
	"io"
)

type Server struct {
	UnimplementedPrimeServiceServer
}

func (serverInstance *Server) IsPrime(number int64) bool {
	if number == 2 || number == 3 || number == 5 || number == 7 || number == 9 {
		return true
	}
	if number%2 == 0 {
		return false
	}
	if number%3 == 0 {
		return false
	}
	if number%5 == 0 {
		return false
	}
	if number%7 == 0 {
		return false
	}
	if number%9 == 0 {
		return false
	}
	return true
}

func (serverInstance *Server) PrimeCount(number int64) int64 {
	var count int64 = 0
	for i := 0; i <= int(number); i++ {
		var isPrime bool = serverInstance.IsPrime(int64(i))
		if isPrime {
			count = count + 1
		}
	}
	return count
}

func (serverInstance *Server) GetCount(stream PrimeService_GetCountServer) error {
	for {
		//consume stream
		currentRequest, requestErr := stream.Recv()
		if requestErr != nil && requestErr != io.EOF {
			return requestErr
		}
		if requestErr == io.EOF {
			fmt.Println("finished consuming prime stream")
			return nil
		}
		var primeCount int64 = serverInstance.PrimeCount(currentRequest.Number)
		var newResponse *PrimeResponse = &PrimeResponse{
			Count: primeCount,
		}
		message := fmt.Sprintf(
			"current count for number %v is %v",
			currentRequest.Number,
			primeCount,
		)
		fmt.Println(message)
		var err error = stream.Send(newResponse)
		if err != nil {
			return err
		}
	}
}
