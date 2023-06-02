package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"

	"google.golang.org/grpc"

	pb "github.com/royadaneshi/webHW1/auth/authservice"
)

func generateNonce(length int) (string, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	randomString := base64.URLEncoding.EncodeToString(randomBytes)
	randomString = randomString[:length]

	return randomString, nil
}

func generateMessageID() (int32, error) {

	max := big.NewInt(50)
	randomInt, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}
	messageID := randomInt.Int64() * 2

	return int32(messageID), nil
}
func main() {
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := pb.NewMyServiceClient(conn)

	nonce_gen, err := generateNonce(20)
	if err != nil {
		log.Fatalf("Failed to generate nonce: %v", err)
	}

	messageID, err := generateMessageID()
	if err != nil {
		log.Fatalf("Failed to generate message ID: %v", err)
	}

	request := &pb.MyRequest{
		MessageId: messageID,
		Nonce:     nonce_gen,
	}

	response, err := client.ProcessRequest(context.Background(), request)
	if err != nil {
		log.Fatalf("Failed to call ProcessRequest: %v", err)
	}

	fmt.Printf("Response:\nNonce: %s\nServer Nonce: %s\nMessage ID: %d\nP: %d\nG: %d\n",
		response.GetNonce(), response.GetServerNonce(), response.GetMessageId(), response.GetP(), response.GetG())

}
