package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	random1 "math/rand"

	pb1 "github.com/royadaneshi/webHW1/auth/DH_params"
	pb "github.com/royadaneshi/webHW1/auth/authservice"
	"google.golang.org/grpc"
)

func generateNonce(length int) (string, error) {
	   randomBytes := make([]byte, length)
	_ , err        := rand.Read(randomBytes)
	if err         != nil {
		return "", err
	}

	randomString := base64.URLEncoding.EncodeToString(randomBytes)
	randomString  = randomString[:length]

	return randomString, nil
}

func generateMessageID() (int32, error) {

	          max  := big.NewInt(50)
	randomInt, err := rand.Int(rand.Reader, max)
	if        err  != nil {
		return 0, err
	}
	messageID := randomInt.Int64() * 2

	return int32(messageID), nil
}

func generateEvenNumberGreaterThan(x int32) int32 {
	   evenNumber   := x + 1
	if evenNumber%2 != 0 {
		evenNumber++
	}
	return evenNumber
}

func main() {
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if   err  != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := pb.NewMyServiceClient(conn)

	nonce_gen, err := generateNonce(20)
	if        err  != nil {
		log.Fatalf("Failed to generate nonce: %v", err)
	}

	messageID, err := generateMessageID()
	if        err  != nil {
		log.Fatalf("Failed to generate message ID: %v", err)
	}

	request := &pb.MyRequest{
		MessageId: messageID,
		Nonce    : nonce_gen,
	}

	response, err := client.ProcessRequest(context.Background(), request)
	if       err  != nil {
		log.Fatalf("Failed to call ProcessRequest: %v", err)
	}

	fmt.Printf("Response:\nNonce: %s\nServer Nonce: %s\nMessage ID: %d\nP: %d\nG: %d\n",
		response.GetNonce(), response.GetServerNonce(), response.GetMessageId(), response.GetP(), response.GetG())

	conn.Close()

	  // call the next service to get key
	conn1, err1 := grpc.Dial("localhost:50054", grpc.WithInsecure())
	if    err1  != nil {
		log.Fatalf("Failed to connect to server auth2: %v", err1)
	}
	defer conn1.Close()

	personal_key := int64(random1.Intn(10000))
	g            := big.NewInt(response.GetG())
	a            := big.NewInt(personal_key)
	p            := big.NewInt(response.GetP())
	  //g^a mod p:
	public_key := new(big.Int).Exp(g, a, p)
	client1    := pb1.NewMyServiceClient(conn1)
	request1   := &pb1.DHParamsRequest{
		Nonce      : response.GetNonce(),
		ServerNonce: response.GetServerNonce(),
		MessageId  : generateEvenNumberGreaterThan(messageID),

		A: int32(public_key.Int64()),
	}
	response1, err2 := client1.ProcessRequest(context.Background(), request1)
	if        err2  != nil {
		log.Fatalf("Failed to call ProcessRequest: %v", err2)
	}
	fmt.Println(response1)
}
