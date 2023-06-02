package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net"
	"time"
	mrand "math/rand"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"

	pb "github.com/royadaneshi/webHW1/auth/authservice"
)

type server struct {
	redisClient *redis.Client
	pb.UnimplementedMyServiceServer
}

func generateRandomGenerator(p *big.Int) (*big.Int, error) {
	two := big.NewInt(2)
	// Choose a random value for g between 2 and p-2
	g, err := rand.Int(rand.Reader, new(big.Int).Sub(p, two))
	if err != nil {
		return nil, err
	}
	// make sure g is at least 2
	g.Add(g, two)
	return g, nil
}
func generateRandomPrime(bits int) (*big.Int, error) {
	p, err := rand.Prime(rand.Reader, bits)
	if err != nil {
		return nil, err
	}
	return p, nil
}
func generateOddNumber() int32 {
	for {
		num := mrand.Int31n(1000) + 1
		if num%2 != 0 {
			return num
		}
	}
}
func (s *server) ProcessRequest(ctx context.Context, req *pb.MyRequest) (*pb.MyResponse, error) {
	if req.MessageId%2 != 0 || req.MessageId <= 0 {
		return nil, fmt.Errorf("Invalid message ID")
	}
	if len(req.Nonce) != 20 {
		return nil, fmt.Errorf("Invalid nonce length")
	}

	p, err := generateRandomPrime(2048)
	if err != nil {
		return nil, fmt.Errorf("Error:", err)

	}

	g, err := generateRandomGenerator(p)
	if err != nil {
		return nil, fmt.Errorf("Error:", err)
	}

	resp := &pb.MyResponse{
		Nonce:       req.GetNonce(),
		ServerNonce: generateNonce(),
		MessageId:   generateOddNumber(),
		P:           p.Int64(),
		G:           int32(g.Int64()),
	}

	// save to redis
	key := fmt.Sprintf("%s:%s", resp.Nonce, resp.ServerNonce)
	err = s.redisClient.Set(ctx, key, "", 20*time.Minute).Err()
	if err != nil {
		log.Printf("Failed to store data in Redis: %v", err)
	}

	return resp, nil
}

func generateNonce() string {
	const nonceLength = 20

	nonce := make([]byte, nonceLength)
	_, err := rand.Read(nonce)
	if err != nil {
		log.Fatal("Failed to generate nonce:", err)
	}

	nonceString := ""
	for _, b := range nonce {
		nonceString += fmt.Sprintf("%02x", b)
	}

	return nonceString
}

func main() {
	// Create a Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", //redis port its true
		Password: "",               //dafault pass
		DB:       0,                //default db
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	lis, err := net.Listen("tcp", "localhost:50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMyServiceServer(s, &server{redisClient: redisClient})

	log.Println("Starting gRPC server...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
