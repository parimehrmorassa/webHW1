package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"

	pb "github.com/royadaneshi/webHW1/auth/DH_params"
)

type server struct {
	redisClient *redis.Client
	pb.UnimplementedDHParamsServiceServer
}

func generateRandomNumber() (int32, error) {
	num, err := rand.Int(rand.Reader, big.NewInt(1000))
	if err != nil {
		return 0, err
	}
	return int32(num.Int64()), nil
}

func redisConnect(req *pb.DHParamsRequest, client *redis.Client) (string, error) {
	ctx := context.Background()
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		return "", fmt.Errorf("Error connecting to Redis: %v", err)
	}
	fmt.Println("Connected to Redis:", pong)
	defer client.Close()

	key := fmt.Sprintf("%s:%s", req.GetNonce(), req.GetServerNonce())

	// Read data from Redis
	value, err := client.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("Error reading from Redis: %v", err)
	}
	return value, nil
}

func (s *server) ProcessRequest(ctx context.Context, req *pb.DHParamsRequest) (*pb.DHParamsResponse, error) {
	value, err := redisConnect(req, s.redisClient)
	if err != nil {
		log.Printf("Failed to get data from Redis: %v", err)
		return nil, err
	}
	fmt.Println(value, "!!!!!!!")
	fmt.Println(value[3], "***!!!!!!!")

	// Generate random number b
	b, err := generateRandomNumber()
	if err != nil {
		log.Printf("Failed to generate random number: %v", err)
		return nil, err
	}

	resp := &pb.DHParamsResponse{
		Nonce:       req.GetNonce(),
		ServerNonce: req.GetServerNonce(),
		MessageId:   req.GetMessageId(),
		B:           b,
	}

	key := fmt.Sprintf("%s:%s", req.GetNonce(), req.GetServerNonce())

	// Remove the last data of the user from Redis
	err = s.redisClient.Del(ctx, key).Err()
	if err != nil {
		log.Printf("Failed to remove data from Redis: %v", err)
	}

	return resp, nil
}

func main() {
	lis, err := net.Listen("tcp", "localhost:50054")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	s := grpc.NewServer()
	pb.RegisterDHParamsServiceServer(s, &server{redisClient: client})

	log.Println("Starting gRPC server...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
