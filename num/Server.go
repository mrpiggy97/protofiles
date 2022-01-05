package num

import (
	"errors"
	"fmt"
	"io"
	"time"
)

type Server struct {
	UnimplementedNumServiceServer
}

// counter will recieve 3 values from request object From, To and
// Number, Number will be the final possible number count
// To will be the absoulute limit for loop
// From will set the starting point
func (serverInstance *Server) Rnd(request *NumRequest, stream NumService_RndServer) error {
	fmt.Println(request.String())
	if request.Number <= 0 {
		return errors.New("number cannot be less or equal than 0")
	}

	if request.To <= request.From {
		return errors.New("req.to cannot be less or eqaul than req.From")
	}
	var response *NumResponse = new(NumResponse)
	for i := int(request.From); i <= int(request.To); i++ {
		response = &NumResponse{
			CurrentNumber: int64(i),
			Remaining:     int64(int(request.Number) - i),
		}
		time.Sleep(time.Millisecond * 50)
		var err error = stream.Send(response)
		if err != nil {
			return err
		}
	}
	return nil
}

func (serverInstance *Server) Sum(stream NumService_SumServer) error {
	var total int64 = 0
	for {
		req, reqError := stream.Recv()
		if reqError != nil && reqError != io.EOF {
			return reqError
		}
		if reqError == io.EOF {
			break
		}
		total = total + req.GetNumber()
		fmt.Println("server recieved ", req.String())
	}
	var response *SumResponse = &SumResponse{
		Total: total,
	}
	stream.SendAndClose(response)
	fmt.Println("server sent ", response.String())
	return nil
}
