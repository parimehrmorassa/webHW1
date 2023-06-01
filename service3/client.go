package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	pb "github.com/royadaneshi/webHW1/service3/authservice" // Update with the correct package path
	"google.golang.org/grpc"
)

func main() {
	// Set up a connection to the server
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a new client
	client := pb.NewMyServiceClient(conn)

	// Generate a random even number
	messageID := generateEvenNumber()

	// Generate a random nonce string
	nonce := generateRandomString(20)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Send the gRPC request to the server
	response, err := client.MyMethod(ctx, &pb.MyRequest{MessageId: messageID, Nonce: nonce})
	if err != nil {
		log.Fatalf("Failed to call MyMethod: %v", err)
	}

	// Print the response
	log.Printf("Nonce sent by client: %s", response.Nonce)
	log.Printf("Server Nonce: %s", response.ServerNonce)
	log.Printf("Message: %d", response.Message)
	log.Printf("P: %d", response.P)
	log.Printf("G: %d", response.G)
}

// Generate a random even number greater than 0
func generateEvenNumber() int32 {
	for {
		num := rand.Int31n(1000) + 1 // Generate a random number between 1 and 1000
		if num%2 == 0 {
			return num
		}
	}
}

// Generate a random string of the specified length
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var result string
	for i := 0; i < length; i++ {
		result += string(charset[rand.Intn(len(charset))])
	}
	return result
}

// package main

// import (
// 	"context"
// 	"log"
// 	"math/rand"

// 	"google.golang.org/grpc"

// 	pb "github.com/royadaneshi/webHW1/royadaneshi/webHW1/service3/authservice" // Update with the correct package path
// )

// func main() {
// 	// Set up a connection to the gRPC server
// 	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
// 	if err != nil {
// 		log.Fatalf("Failed to connect: %v", err)
// 	}
// 	defer conn.Close()

// 	// Create a new gRPC client
// 	client := pb.NewMyServiceClient(conn)

// 	// Prepare the request
// 	request := &pb.MyRequest{
// 		MessageId: 1234,
// 		Nonce: string(func(l int) []byte {
// 			b := make([]byte, l)
// 			for i := range b {
// 				b[i] = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"[rand.Intn(len("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"))]
// 			}
// 			return b
// 		}(rand.Intn(20) + 1)),
// 	}

// 	// Send the request to the server
// 	response, err := client.ProcessRequest(context.Background(), request)
// 	if err != nil {
// 		log.Fatalf("Failed to process request: %v", err)
// 	}

// 	log.Printf("Response: %v", response)
// }
