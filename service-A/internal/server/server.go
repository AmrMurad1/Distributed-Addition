package server

import (
	"context"
	"log"

	pb "service-A/internal/proto/addition"
)

// Server implements the additionService
type Server struct {
	pb.UnimplementedAdditionServiceServer
}

func SumServer() *Server {
	return &Server{}
}

func (s *Server) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	log.Printf("Received Add request: %d + %d", req.GetNum1(), req.GetNum2())

	result := int64(req.GetNum1() + req.GetNum2())

	response := &pb.AddResponse{
		Result: result,
	}

	log.Printf("Sending Add response: %d", result)
	return response, nil
}
