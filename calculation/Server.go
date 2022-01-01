package calculation

import (
	"fmt"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type Server struct {
	UnimplementedCalculationServiceServer
}

func (a Server) SumStream(stream CalculationService_SumStreamServer) error {
	for {
		sumStreamRequest, err := stream.Recv()
		fmt.Println("server got ", sumStreamRequest.String())
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
		err = stream.Send(&SumStreamResponse{
			Sum: sumStreamRequest.A + sumStreamRequest.B,
		})
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}
}
