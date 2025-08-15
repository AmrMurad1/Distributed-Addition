package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "service-A/internal/proto/addition"
	"service-A/internal/server"
)

const (
	defaultPort = ":50051"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	additionServer := server.SumServer()
	pb.RegisterAdditionServiceServer(s, additionServer)

	reflection.Register(s)

	log.Printf("gRPC server listening at %v", lis.Addr())

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// close Kafka producer
	if err := additionServer.Close(); err != nil {
		log.Printf("Error closing Kafka producer: %v", err)
	}

	s.GracefulStop()
	log.Println("Server stopped")
}
