package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "service-A/internal/proto/addition"
	"service-A/internal/server"
)

const (
	port = ":50051"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = ":50051"
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

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
