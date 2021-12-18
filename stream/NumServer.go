package stream

import (
	"errors"
	"fmt"
	"time"
)

// counter will recieve 3 values from request object From, To and
// Number, Number will be the final possible number count
// To will be the absoulute limit for loop
// From will set the starting point
func counter(request *NumRequest, streamInstance NumService_RndServer, channel chan<- bool) {
	var response *NumResponse = new(NumResponse)
	for i := int(request.From); i <= int(request.To); i++ {
		response = &NumResponse{
			CurrentNumber: int64(i),
			Remaining:     int64(int(request.Number) - i),
		}
		time.Sleep(time.Millisecond * 50)
		var err error = streamInstance.Send(response)
		if err != nil {
			panic(err)
		}
	}
	time.Sleep(time.Second)
	channel <- true
}

type NumServer struct {
	UnimplementedNumServiceServer
}

func (serverInstance *NumServer) Rnd(req *NumRequest, stream NumService_RndServer) error {
	fmt.Println(req.String())
	if req.Number <= 0 {
		return errors.New("number cannot be less or equal than 0")
	}

	if req.To <= req.From {
		return errors.New("req.to cannot be less or eqaul than req.From")
	}

	var done chan bool = make(chan bool, 1)
	go counter(req, stream, done)
	<-done
	return nil
}
