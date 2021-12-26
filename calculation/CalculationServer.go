package calculation

import (
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type CalculationServer struct {
	UnimplementedCalculationServiceServer
}

func (a CalculationServer) SumStream(stream CalculationService_SumStreamServer) error {
	for {
		sumStreamRequest, err := stream.Recv()
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
