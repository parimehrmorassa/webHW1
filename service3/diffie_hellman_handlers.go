// package main

// import (
// 	"context"
// 	"crypto/rand"
// 	"log"
// 	"math/big"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"google.golang.org/grpc"

// 	authpb "/Users/parimehrsmacbook/webHW1/service3/req_pq.proto" // Import the generated protobuf code for authentication
// 	pb "your_package_path/pb"          // Import the generated protobuf code for other services

// 	"encoding/base64"

// 	"github.com/go-redis/redis/v8"
// )

// const (
// 	redisKeyDiffieHellmanParams = "diffie_hellman_params"
// )

// // AuthServer implements the gRPC server for authentication service
// type AuthServer struct {
// 	authpb.UnimplementedAuthServiceServer
// 	redisClient *redis.Client
// }

// // PerformKeyExchange is the implementation of the PerformKeyExchange RPC method
// func (s *AuthServer) PerformKeyExchange(ctx context.Context, req *authpb.DiffieHellmanKeyExchange) (*authpb.DiffieHellmanKeyExchange, error) {
// 	// Generate server's Diffie-Hellman key pair
// 	p, _ := rand.Prime(rand.Reader, 256) // Use the appropriate size for your use case
// 	g := big.NewInt(2)
// 	privateKeyServer := new(big.Int).Rand(rand.Reader, p)
// 	publicKeyServer := new(big.Int).Exp(g, privateKeyServer, p)

// 	// Compute shared secret
// 	clientPublicKey := new(big.Int).SetBytes(req.PublicKey)
// 	sharedSecret := new(big.Int).Exp(clientPublicKey, privateKeyServer, p)

// 	// Generate server's public key
// 	publicKeyServerBytes := publicKeyServer.Bytes()

// 	// Create response message with server's public key
// 	res := &authpb.DiffieHellmanKeyExchange{
// 		PublicKey: publicKeyServerBytes,
// 	}

// 	return res, nil
// }

// // Authenticate is the implementation of the Authenticate RPC method
// func (s *AuthServer) Authenticate(ctx context.Context, req *authpb.AuthenticateRequest) (*authpb.AuthenticateResponse, error) {
// 	// Perform authentication logic here
// 	// Example: check username and password against database or other authentication mechanism
// 	// Return the appropriate authentication response

// 	// Placeholder implementation
// 	res := &authpb.AuthenticateResponse{
// 		Nonce:       req.Nonce,
// 		ServerNonce: generateNonce(),
// 		MessageId:   req.MessageId,
// 	}

// 	return res, nil
// }

// // GetDiffieHellmanParameters is the implementation of the GetDiffieHellmanParameters RPC method
// func (s *AuthServer) GetDiffieHellmanParameters(ctx context.Context, req *pb.Empty) (*authpb.DiffieHellmanParameters, error) {
// 	// Check if the Diffie-Hellman parameters are cached in Redis
// 	if s.redisClient.Exists(ctx, redisKeyDiffieHellmanParams).Val() > 0 {
// 		// Retrieve the parameters from Redis
// 		p, _ := s.redisClient.Get(ctx, "p").Result()
// 		g, _ := s.redisClient.Get(ctx, "g").Result()

// 		return &authpb.DiffieHellmanParameters{
// 			P: p,
// 			G: g,
// 		}, nil
// 	}

// 	// Generate new Diffie-Hellman parameters
// 	p, _ := rand.Prime(rand.Reader, 256) // Use the appropriate size for your use case
// 	g := big.NewInt(2)

// 	// Store the parameters in Redis
// 	s.redisClient.Set(ctx, "p", p.String(), 0)
// 	s.redisClient.Set(ctx, "g", g.String(), 0)

// 	return &authpb.DiffieHellmanParameters{
// 		P: p.String(),
// 		G: g.String(),
// 	}, nil
// }

// func generateNonce() string {
// 	const nonceLength = 20

// 	nonceBytes := make([]byte, nonceLength)
// 	rand.Read(nonceBytes)

// 	return base64.StdEncoding.EncodeToString(nonceBytes)
// }

// func main() {
// 	// Create a new Gin router
// 	router := gin.Default()

// 	// Create a new gRPC server
// 	grpcServer := grpc.NewServer()

// 	// Create a Redis client
// 	redisClient := redis.NewClient(&redis.Options{
// 		Addr: "localhost:6379",
// 	})

// 	// Initialize the AuthServer with the Redis client
// 	authServer := &AuthServer{
// 		redisClient: redisClient,
// 	}

// 	// Register your gRPC server implementation with the gRPC server
// 	authpb.RegisterAuthServiceServer(grpcServer, authServer)

// 	// Define an endpoint for handling gRPC requests via HTTP
// 	router.POST("/auth/perform_key_exchange", func(c *gin.Context) {
// 		// Forward the request to the gRPC server
// 		grpcServer.ServeHTTP(c.Writer, c.Request)
// 	})

// 	router.POST("/auth/authenticate", func(c *gin.Context) {
// 		// Forward the request to the gRPC server
// 		grpcServer.ServeHTTP(c.Writer, c.Request)
// 	})

// 	router.POST("/auth/get_diffie_hellman_parameters", func(c *gin.Context) {
// 		// Forward the request to the gRPC server
// 		grpcServer.ServeHTTP(c.Writer, c.Request)
// 	})

// 	// Start the HTTP server
// 	err := http.ListenAndServe(":8080", router)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

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

	authpb "/Users/parimehrsmacbook/webHW1/service3" // Import the generated protobuf code for the Auth service
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
