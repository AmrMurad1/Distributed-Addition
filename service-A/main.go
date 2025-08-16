package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"service-A/internal/db"
	"service-A/internal/outbox"
	pb "service-A/internal/proto/addition"
	"service-A/internal/server"
)

const (
	defaultPort = ":50051"
)

func main() {
	db.InitDB()
	defer db.Close()

	outboxRepo := outbox.NewRepository(db.DB)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	additionServer := server.SumServer(outboxRepo)
	pb.RegisterAdditionServiceServer(s, additionServer)

	ctx := context.Background()
	additionServer.StartOutboxPublisher(ctx)

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

	if err := additionServer.Close(); err != nil {
		log.Printf("Error closing server resources: %v", err)
	}

	s.GracefulStop()
	log.Println("Server stopped")
}
