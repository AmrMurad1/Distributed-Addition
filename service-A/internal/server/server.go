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

	// add the numbers before send
	result := int64(req.GetNum1() + req.GetNum2())

	err := s.kafkaProducer.SendNumber(int(result))
	if err != nil {
		log.Printf("Error sending sum to Kafka: %v", err)
	}

	response := &pb.AddResponse{
		Result: result,
	}

	log.Printf("Sending Add response: %d", result)
	return response, nil
}

func (s *Server) Close() error {
	return s.kafkaProducer.Close()
}
