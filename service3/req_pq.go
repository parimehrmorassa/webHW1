package main

import (
	"context"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"log"
	"math/big"
	mrand "math/rand"
	"net"
	"time"

	"github.com/go-redis/redis"

	pb "github.com/royadaneshi/webHW1/service3/authservice" // Update with the correct package path
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedMyServiceServer
	redisClient *redis.Client
}

func main() {
	// Connect to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Update with your Redis server address
		Password: "",               // Set your Redis server password (if any)
		DB:       0,                // Set the Redis database index
	})
	defer redisClient.Close()

	// Initialize the random number generator
	mrand.Seed(time.Now().UnixNano())

	// Create a new gRPC server
	grpcServer := grpc.NewServer()
	s := &server{redisClient: redisClient}
	pb.RegisterMyServiceServer(grpcServer, s)

	// Start the gRPC server
	lis := listenGRPC()
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}

func (s *server) MyMethod(ctx context.Context, req *pb.MyRequest) (*pb.MyResponse, error) {
	// Generate a random nonce string
	nonce := generateRandomString(20)

	p, g := generatePrimeAndPrimitiveRoot()

	message := generateOddNumber()

	key := generateSHA1Hash(nonce + req.Nonce)

	if err := s.setRedisValue(key, nonce, time.Hour); err != nil {
		log.Printf("Failed to store data in Redis: %v", err)
	}

	resp := &pb.MyResponse{
		Nonce:       req.Nonce,
		ServerNonce: nonce,
		Message:     message,
		P:           p,
		G:           g,
	}
	return resp, nil
}

func generateOddNumber() int32 {
	for {
		num := mrand.Int31n(1000) + 1
		if num%2 != 0 {
			return num
		}
	}
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var result string
	for i := 0; i < length; i++ {
		result += string(charset[mrand.Intn(len(charset))])
	}
	return result
}

func generatePrimeAndPrimitiveRoot() (int32, int32) {
	p, err := rand.Prime(rand.Reader, 128)
	if err != nil {
		log.Fatalf("Failed to generate prime number: %v", err)
	}

	g := generatePrimitiveRoot(p)

	return int32(p.Int64()), int32(g.Int64())
}

func generateSHA1Hash(input string) string {
	hash := sha1.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}

func (s *server) setRedisValue(key string, value interface{}, expiration time.Duration) error {
	return s.redisClient.Set(key, value, expiration).Err()
}

func generatePrimitiveRoot(p *big.Int) *big.Int {
	if !p.ProbablyPrime(10) {
		log.Fatalf("Input number is not prime")
	}

	for i := big.NewInt(2); i.Cmp(p) < 0; i.Add(i, big.NewInt(1)) {
		if isPrimitiveRoot(i, p) {
			return i
		}
	}

	log.Fatalf("Failed to find a primitive root")
	return nil
}

func isPrimitiveRoot(g, p *big.Int) bool {
	one := big.NewInt(1)
	pMinusOne := new(big.Int).Sub(p, one)

	exp := new(big.Int).Div(pMinusOne, big.NewInt(2))

	modPow := new(big.Int).Exp(g, exp, p)

	return modPow.Cmp(one) != 0
}

// Listen for gRPC connections
func listenGRPC() net.Listener {
	lis, err := net.Listen("tcp", ":50051") // Update with your desired server address
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Println("gRPC server listening on port 50051")
	return lis
}
