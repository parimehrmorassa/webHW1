package main

import (
	"context"
	"log"
	"math/rand"

	"google.golang.org/grpc"

	pb "github.com/royadaneshi/webHW1/service3/authservice" // Update with the correct package path
)

func main() {
	// Set up a connection to the gRPC server
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a new gRPC client
	client := pb.NewMyServiceClient(conn)

	// Prepare the request
	request := &pb.MyRequest{
		MessageId: 1234, // Replace with your own message ID
		Nonce: string(func(l int) []byte {
			b := make([]byte, l)
			for i := range b {
				b[i] = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"[rand.Intn(len("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"))]
			}
			return b
		}(rand.Intn(20) + 1)), // Replace with your own generated nonce
	}

	// Send the request to the server
	response, err := client.ProcessRequest(context.Background(), request)
	if err != nil {
		log.Fatalf("Failed to process request: %v", err)
	}

	// Process the response
	// Replace the following code with your own logic
	log.Printf("Response: %v", response)
}
