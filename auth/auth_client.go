package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"

	pb1 "github.com/royadaneshi/webHW1/auth/DH_params"
	pb "github.com/royadaneshi/webHW1/auth/authservice"
	"google.golang.org/grpc"
)

type keys struct {
	personalKeyClient *big.Int
	publicKeyClient   *big.Int
	sharedKeyClient   *big.Int
}

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
	if messageID == 0 {
		messageID = messageID + 2
	}

	return int32(messageID), nil
}

func generateEvenNumberGreaterThan(x int32) int32 {
	evenNumber := x + 1
	if evenNumber%2 != 0 {
		evenNumber++
	}
	return evenNumber
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

	// fmt.Printf("Response:\nNonce: %s\nServer Nonce: %s\nMessage ID: %d\nP: %d\nG: %d\n",
	// 	response.GetNonce(), response.GetServerNonce(), response.GetMessageId(), response.GetP(), response.GetG())

	conn.Close()

	// call the next service to get key
	conn1, err1 := grpc.Dial("localhost:50054", grpc.WithInsecure())
	if err1 != nil {
		log.Fatalf("Failed to connect to server auth2: %v", err1)
	}
	defer conn1.Close()

	// personal_key := int64(random1.Intn(10000))

	// Generate a new private key for the client side
	privateKey, err := rsa.GenerateKey(rand.Reader, 20)
	if err != nil {
		log.Fatal("Failed to generate private key b:", err)
	}

	personal_key := privateKey.D

	g := new(big.Int)
	g.SetString(response.GetG(), 10)

	p := new(big.Int)
	p.SetString(response.GetP(), 10)

	a := big.NewInt(personal_key.Int64())
	//g^a mod p:
	public_key := new(big.Int).Exp(g, a, p)
	client1 := pb1.NewDHParamsServiceClient(conn1)
	request1 := &pb1.DHParamsRequest{
		Nonce:       response.GetNonce(),
		ServerNonce: response.GetServerNonce(),
		MessageId:   generateEvenNumberGreaterThan(messageID),

		A: public_key.String(),
	}
	response1, err2 := client1.ProcessRequest(context.Background(), request1)
	if err2 != nil {
		log.Fatalf("Failed to call ProcessRequest auth2: %v", err2)
	}
	//calculate Shared key
	// b_server_key := big.NewInt(response1.B)
	b_server_key := new(big.Int)
	b_server_key.SetString(response1.B, 10)
	// B^a mod p:
	shared_key := new(big.Int).Exp(public_key, b_server_key, p)

	////
	myKeys := keys{
		personalKeyClient: a,
		publicKeyClient:   public_key,
		sharedKeyClient:   shared_key,
	}
	fmt.Println("Personal Key for client:", myKeys.personalKeyClient)
	fmt.Println("Public Key for client:", myKeys.publicKeyClient)
	fmt.Println("Shared Key:", myKeys.sharedKeyClient)
}
