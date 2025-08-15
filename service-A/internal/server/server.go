package server

import (
	"context"
	"log"

	"service-A/internal/kafka"
	pb "service-A/internal/proto/addition"
)

const (
	KafkaHost  = "localhost:9092"
	KafkaTopic = "addition"
)

// Server implements the additionService
type Server struct {
	pb.UnimplementedAdditionServiceServer
	kafkaProducer *kafka.Producer
}

func SumServer() *Server {
	producer := kafka.NewProducer(KafkaHost, KafkaTopic)
	return &Server{
		kafkaProducer: producer,
	}
}

func (s *Server) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	log.Printf("Received Add request: %d + %d", req.GetNum1(), req.GetNum2())

	// Send both numbers to Kafka
	err1 := s.kafkaProducer.SendNumber(int(req.GetNum1()))
	err2 := s.kafkaProducer.SendNumber(int(req.GetNum2()))

	if err1 != nil || err2 != nil {
		log.Printf("Error sending numbers to Kafka: %v, %v", err1, err2)
	}

	// Still return the addition result for gRPC response
	result := int64(req.GetNum1() + req.GetNum2())

	response := &pb.AddResponse{
		Result: result,
	}

	log.Printf("Sending Add response: %d", result)
	return response, nil
}

func (s *Server) Close() error {
	return s.kafkaProducer.Close()
}
