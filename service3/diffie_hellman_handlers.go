package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	"github.com/go-redis/redis/v8"

	authpb "github.com/parimehrmorassa/webHW1/service3/auth/authpb" // Import the generated protobuf code for the Auth service
)

var rdb *redis.Client

func main() {
	// Initialize Redis client
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Initialize Gin router
	router := gin.Default()

	// Register gRPC handler
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	authClient := authpb.NewAuthServiceClient(conn)

	// Define the route handler
	router.GET("/get_diffie_hellman_params", func(c *gin.Context) {
		// Generate nonce and message ID
		nonce, err := generateNonce()
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to generate nonce"})
			return
		}
		messageID := 2 // Even number

		// Create request params
		req := &authpb.RequestParams{
			Nonce:     nonce,
			MessageId: int32(messageID),
		}

		// Call gRPC method to get Diffie-Hellman parameters
		params, err := authClient.GetDiffieHellmanParams(context.Background(), req)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to get Diffie-Hellman parameters"})
			return
		}

		// Store parameters in Redis
		err = cacheParamsInRedis(nonce, params)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to cache parameters in Redis"})
			return
		}

		// Send the nonce and server_nonce in the response
		c.JSON(200, gin.H{
			"nonce":        nonce,
			"server_nonce": params.ServerNonce,
		})
	})

	// Run the Gin server
	router.Run(":8080")
}

func generateNonce() (string, error) {
	nonceLength := 20
	const nonceAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var nonceBuilder strings.Builder
	for i := 0; i < nonceLength; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(nonceAlphabet))))
		if err != nil {
			return "", err
		}
		nonceChar := nonceAlphabet[randomIndex.Int64()]
		nonceBuilder.WriteByte(nonceChar)
	}
	return nonceBuilder.String(), nil
}

func cacheParamsInRedis(nonce string, params *authpb.DiffieHellmanParams) error {
	pipeline := rdb.Pipeline()
	ctx := context.Background()

	// Store the parameters in Redis
	pipeline.Set(ctx, fmt.Sprintf("nonce:%s", nonce), nonce, 0)
	pipeline.Set(ctx, fmt.Sprintf("server_nonce:%s", nonce), params.ServerNonce, 0)
	pipeline.Set(ctx, fmt.Sprintf("p:%s", nonce), params.P, 0)
	pipeline.Set(ctx, fmt.Sprintf("g:%s", nonce), params.G, 0)

	_, err := pipeline.Exec(ctx)
	return err
}
