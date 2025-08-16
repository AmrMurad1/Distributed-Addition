package server

import (
	"context"
	"log"

	"service-A/internal/kafka"
	"service-A/internal/outbox"
	pb "service-A/internal/proto/addition"
)

const (
	KafkaHost  = "localhost:9092"
	KafkaTopic = "addition"
)

type Server struct {
	pb.UnimplementedAdditionServiceServer
	kafkaProducer   *kafka.Producer
	outboxRepo      *outbox.Repository
	outboxPublisher *outbox.Publisher
}

func SumServer(outboxRepo *outbox.Repository) *Server {
	producer := kafka.NewProducer(KafkaHost, KafkaTopic)

	outboxPublisher := outbox.NewPublisher(outboxRepo, producer, 5)

	return &Server{
		kafkaProducer:   producer,
		outboxRepo:      outboxRepo,
		outboxPublisher: outboxPublisher,
	}
}

func (s *Server) StartOutboxPublisher(ctx context.Context) {
	s.outboxPublisher.Start(ctx)
}

func (s *Server) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	log.Printf("Received Add request: %d + %d", req.GetNum1(), req.GetNum2())

	result := int64(req.GetNum1() + req.GetNum2())

	additionEvent := outbox.AdditionEvent{
		Number: int(result),
	}

	err := s.outboxRepo.SaveEvent(ctx, "addition_result", additionEvent)
	if err != nil {
		log.Printf("Error saving event to outbox: %v", err)
	}

	response := &pb.AddResponse{
		Result: result,
	}

	log.Printf("Sending Add response: %d", result)
	return response, nil
}

func (s *Server) Close() error {
	s.outboxPublisher.Stop()
	return s.kafkaProducer.Close()
}
