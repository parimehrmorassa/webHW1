package main

import (
	"context"
	"fmt"
	"log"

	authpb "/Users/parimehrsmacbook/webHW1/service3" // Import the generated protobuf code for the Auth service

	"google.golang.org/grpc"
)

type authServer struct{}

func (*authServer) GetDiffieHellmanParams(ctx context.Context, req *authpb.RequestParams) (*authpb.DiffieHellmanParams, error) {
	// Generate server_nonce and message ID
	serverNonce, err := generateNonce()
	if err != nil {
		return nil, err
	}
	messageID := 3 // Odd number

	// Retrieve the p and g values from Redis using the provided nonce
	p, err := rdb.Get(ctx, fmt.Sprintf("p:%s", req.Nonce)).Result()
	if err != nil {
		return nil, err
	}
	g, err := rdb.Get(ctx, fmt.Sprintf("g:%s", req.Nonce)).Result()
	if err != nil {
		return nil, err
	}

	// Create response params
	res := &authpb.DiffieHellmanParams{
		Nonce:       req.Nonce,
		ServerNonce: serverNonce,
		P:           p,
		G:           g,
	}

	// Cache the response parameters in Redis
	err = cacheParamsInRedis(req.Nonce, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (*authServer) ProcessDiffieHellmanParams(ctx context.Context, req *authpb.ResponseParams) (*authpb.ResponseParams, error) {
	// Process the Diffie-Hellman parameters and return the response
	// You can implement the Diffie-Hellman key exchange algorithm here

	// For demonstration purposes, simply return the received params with an updated message ID
	req.MessageId = req.MessageId + 1

	return req, nil
}

func generateNonce() (string, error) {
	// Same nonce generation code as mentioned above
	// ...

	return nonce, nil
}

func main() {
	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register Auth service server
	authpb.RegisterAuthServiceServer(grpcServer, &authServer{})

	// Start gRPC server
	log.Println("Starting gRPC server on port 50051...")
	if err := grpcServer.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
