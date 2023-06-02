package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"sync"
	"time"

	"context"
	"fmt"
	"log"
	"math/big"
	random1 "math/rand"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	DH_params "github.com/royadaneshi/webHW1/auth/DH_params"
	Auth_service "github.com/royadaneshi/webHW1/auth/authservice"
	grpcService_get_users "github.com/royadaneshi/webHW1/service1/get_user/pb"
)

const (
	MaxRequestsPerSecond = 100
	BanDuration          = 24 * time.Hour
)

var (
	blacklist   = make(map[string]time.Time)
	blacklistMu sync.Mutex

	client      grpcService_get_users.GetUsersClient
	clientAuth  Auth_service.AuthClient
	AuthKey_get *big.Int
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

type keys struct {
	personalKeyClient *big.Int
	publicKeyClient   *big.Int
	sharedKeyClient   *big.Int
}

func getAuthKey() (*big.Int, string, int32, error) {
	// call  Auth service
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := Auth_service.NewMyServiceClient(conn)

	nonce_gen, err := generateNonce(20)
	if err != nil {
		log.Fatalf("Failed to generate nonce: %v", err)
	}

	messageID, err := generateMessageID()
	if err != nil {
		log.Fatalf("Failed to generate message ID: %v", err)
	}

	request := &Auth_service.MyRequest{
		MessageId: messageID,
		Nonce:     nonce_gen,
	}

	response, err := client.ProcessRequest(context.Background(), request)
	if err != nil {
		log.Fatalf("Failed to call ProcessRequest: %v", err)
	}
	conn.Close()

	// call the next service to get key
	conn1, err1 := grpc.Dial("localhost:50054", grpc.WithInsecure())
	if err1 != nil {
		log.Fatalf("Failed to connect to server auth2: %v", err1)
	}
	defer conn1.Close()

	personal_key := int64(random1.Intn(10000))
	g := big.NewInt(response.GetG())
	a := big.NewInt(personal_key)
	p := big.NewInt(response.GetP())
	//g^a mod p:
	public_key := new(big.Int).Exp(g, a, p)
	client1 := DH_params.NewDHParamsServiceClient(conn1)
	messageidd := generateEvenNumberGreaterThan(messageID)
	request1 := &DH_params.DHParamsRequest{
		Nonce:       response.GetNonce(),
		ServerNonce: response.GetServerNonce(),
		MessageId:   messageidd,

		A: public_key.Int64(),
	}
	response1, err2 := client1.ProcessRequest(context.Background(), request1)
	if err2 != nil {
		log.Fatalf("Failed to call ProcessRequest auth2: %v", err2)
	}
	//calculate Shared key
	b_server_key := big.NewInt(response1.B)
	// B^a mod p:
	shared_key := new(big.Int).Exp(public_key, b_server_key, p)

	////
	myKeys := keys{
		personalKeyClient: a,
		publicKeyClient:   public_key,
		sharedKeyClient:   shared_key,
	}

	fmt.Println("Shared Key:", myKeys.sharedKeyClient)
	redis_key := fmt.Sprintf("%s:%s", response.GetNonce(), response.GetServerNonce())

	return myKeys.sharedKeyClient, redis_key, messageidd, nil
}

func authenticateIP(c *gin.Context) {
	ip := c.ClientIP()
	blacklistMu.Lock()
	if banTime, ok := blacklist[ip]; ok && banTime.After(time.Now()) {
		c.AbortWithStatus(http.StatusForbidden)
		blacklistMu.Unlock()
		return
	}
	blacklistMu.Unlock()

	// rate limiting by tick
	limiter := time.Tick(time.Second / MaxRequestsPerSecond)
	select {
	case <-limiter:
		// its okk
		c.Next()
	default:
		// add IP to the blacklist there was many requests
		blacklistMu.Lock()
		blacklist[ip] = time.Now().Add(BanDuration)
		blacklistMu.Unlock()
		c.AbortWithStatus(http.StatusTooManyRequests)
	}
}

func gatewayHandler(c *gin.Context) {
	request := &grpcService_get_users.GetDataRequest{
		UserId:  10,
		AuthKey: AuthKey_get,
	}

	ctx := c.Request.Context()
	response, err := client.GetData(ctx, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response.ReturnUsers})
}
func cleanupBlacklist() {
	for {
		time.Sleep(time.Minute)
		currentTime := time.Now()

		blacklistMu.Lock()
		for ip, banTime := range blacklist {
			if banTime.Before(currentTime) {
				delete(blacklist, ip)
			}
		}
		blacklistMu.Unlock()
	}
}
func BizService(redis_key string, message int32) {
	grpcAddress := "localhost:50051"
	conn, err := grpc.Dial(grpcAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()
	client := grpcService_get_users.NewGetUsersClient(conn)
	request := &grpcService_get_users.GetDataRequest{
		UserId:    10000000,
		AuthKey:   AuthKey_get,
		MessageId: message,
		RedisKey:  redis_key,
	}
	response, err := client.GetData(context.Background(), request)
	if err != nil {
		log.Fatalf("Failed to get data from Biz service : %v", err)
	}
	if response.MessageId == 1 {
		user := response.ReturnUsers[0]
		fmt.Println("User ID: %d\n", user.Id)
		fmt.Printf("Name: %s\n", user.Name)
		fmt.Printf("Family: %s\n", user.Family)
		fmt.Printf("Age: %d\n", user.Age)
		fmt.Printf("Sex: %s\n", user.Sex)
		fmt.Printf("Created At: %s\n", user.CreatedAt)
	} else if response.MessageId == 3 {
		for _, user := range response.ReturnUsers {
			fmt.Printf("User ID: %d\n", user.Id)
			fmt.Printf("Name: %s\n", user.Name)
			fmt.Printf("Family: %s\n", user.Family)
			fmt.Printf("Age: %d\n", user.Age)
			fmt.Printf("Sex: %s\n", user.Sex)
			fmt.Printf("Created At: %s\n", user.CreatedAt)
			fmt.Println("------")
		}
	} else {
		fmt.Println("Unknown response from server")
	}

}
func BizServiceWithSqlInject(redis_key string, message int32) {

}
func main() {
	router := gin.Default()
	go cleanupBlacklist()

	// Connect to the Auth service
	grpcAddressAuth := "localhost:50052" //auth port
	connAuth, errAuth := grpc.Dial(grpcAddressAuth, grpc.WithInsecure())
	if errAuth != nil {
		log.Fatalf("Failed to connect to Auth server: %v", errAuth)
	}
	defer connAuth.Close()

	clientAuth = Auth_service.NewAuthClient(connAuth)

	router.Use(authenticateIP)

	x, redis_key, message, y := getAuthKey()
	AuthKey_get = x
	err := y
	if err != nil {
		log.Fatalf("Failed to get the auth key: %v", err)
	}

	// Connect to the get_users service
	BizService(redis_key, message)
	//////////////////////////
	BizServiceWithSqlInject(redis_key, message)
	//////////////////////////

	router.GET("/gateway", gatewayHandler)
	router.Run(":8080")
}
