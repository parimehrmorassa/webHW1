package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	random1 "math/rand"
	"net"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"

	pb "github.com/royadaneshi/webHW1/auth/DH_params"
	pb1 "github.com/royadaneshi/webHW1/auth/authservice"
)

type server struct {
	redisClient *redis.Client
	pb.UnimplementedDHParamsServiceServer
}
type keys struct {
	personalKeyServer *big.Int
	publicKeyServer   *big.Int
	sharedKeyServer   *big.Int
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
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return "", fmt.Errorf("Error connecting to Redis: %v", err)
	}

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
	key1 := fmt.Sprintf("%s:%s", req.GetNonce(), req.GetServerNonce())
	//read redis to get p and g
	value, err := s.redisClient.Get(ctx, key1).Result()
	if err != nil {
		log.Printf("Failed to retrieve data from Redis for g p: %v", err)
	}
	//convert from json
	byteData := []byte(value)
	var response pb1.MyResponse
	err1 := json.Unmarshal(byteData, &response)
	if err != nil {
		return nil, err1
	}

	//calculate g^b mod p
	personal_key_b := int64(random1.Intn(10000))
	g := big.NewInt(response.G)
	b := big.NewInt(personal_key_b)
	p := big.NewInt(response.P)
	// g^b mod p:
	public_key_B := new(big.Int).Exp(g, b, p)

	resp := &pb.DHParamsResponse{
		Nonce:       req.GetNonce(),
		ServerNonce: req.GetServerNonce(),
		MessageId:   req.GetMessageId(),
		B:           public_key_B.Int64(),
	}

	//calculate Shared key
	a_client_key := big.NewInt(req.A)
	// B^a mod p:
	shared_key := new(big.Int).Exp(public_key_B, a_client_key, p)

	// Remove the last data of the user from Redis
	err = s.redisClient.Del(ctx, key1).Err()
	if err != nil {
		log.Printf("Failed to remove data from Redis: %v", err)
	}
	//save
	jsonValue, err := json.Marshal(shared_key)
	if err != nil {
		return nil, err
	}
	err = s.redisClient.Set(ctx, key1, jsonValue, 0).Err()
	if err != nil {
		log.Printf("Failed to store shared key in Redis: %v", err)
	}

	///check:
	// data, err := s.redisClient.Get(ctx, key1).Result()
	// if err != nil {
	// 	log.Printf("Failed to retrieve data from Redissssssssss: %v", err)
	// }

	// fmt.Println("Retrieved data:", data)
	// ////////

	////
	myKeys := keys{
		personalKeyServer: b,
		publicKeyServer:   public_key_B,
		sharedKeyServer:   shared_key,
	}
	fmt.Println("personal Key for server:", myKeys.personalKeyServer)
	fmt.Println("Public Key for server:", myKeys.publicKeyServer)
	fmt.Println("Shared Key:", myKeys.sharedKeyServer)
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
