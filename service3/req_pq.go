package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/kamyartaeb/webHW1/service3/authservice" // Update with the correct package path
)

type myServiceServer struct{}

func (s *myServiceServer) ProcessRequest(ctx context.Context, req *pb.MyRequest) (*pb.MyResponse, error) {
	// Validate the input
	if req.MessageId%2 != 0 || req.MessageId <= 0 {
		return nil, fmt.Errorf("Invalid message ID")
	}
	if len(req.Nonce) != 20 {
		return nil, fmt.Errorf("Invalid nonce length")
	}

	// Process the request
	// Replace the following code with your own implementation
	response := &pb.MyResponse{
		// Set your response message here
	}

	return response, nil
}

func main() {
	// Create a TCP listener on port 50052
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register the service implementation
	pb.RegisterMyServiceServer(grpcServer, &myServiceServer{})

	// Start the server
	log.Println("Starting gRPC server on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
