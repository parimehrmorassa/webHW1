package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"sync"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"

	DH_params "github.com/royadaneshi/webHW1/auth/DH_params"
	Auth_service "github.com/royadaneshi/webHW1/auth/docker1/authservice"

	grpcService_get_users "github.com/royadaneshi/webHW1/service1/get_user/pb"

	get_user_injection "github.com/royadaneshi/webHW1/service2/get_users_with_sql_inject_proto/pb"
	"google.golang.org/grpc"
)

type IPData struct {
	Count       int
	LastRequest time.Time
}

const (
	MaxRequestsPerSecond = 100            //1
	BanDuration          = 24 * time.Hour //10 * time.Second
)

var (
	AuthKey_get *big.Int
	blacklist   = make(map[string]time.Time)
	blacklistMu sync.Mutex
	ipData      = make(map[string]*IPData)
	ipDataMu    sync.Mutex
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
		log.Fatalf("Failed to connect to auth server: %v", err)
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

	// personal_key := int64(random1.Intn(10000))

	// Generate a new private key for the client side
	privateKey, err := rsa.GenerateKey(rand.Reader, 20)
	if err != nil {
		log.Fatal("Failed to generate private key b:", err)
	}

	personal_key := privateKey.D

	a := big.NewInt(personal_key.Int64())

	g := new(big.Int)
	g.SetString(response.GetG(), 10)

	p := new(big.Int)
	p.SetString(response.GetP(), 10)

	//g^a mod p:
	public_key := new(big.Int).Exp(g, a, p)
	client1 := DH_params.NewDHParamsServiceClient(conn1)
	messageidd := generateEvenNumberGreaterThan(messageID)
	request1 := &DH_params.DHParamsRequest{
		Nonce:       response.GetNonce(),
		ServerNonce: response.GetServerNonce(),
		MessageId:   messageidd,

		A: public_key.String(),
	}
	response1, err2 := client1.ProcessRequest(context.Background(), request1)
	if err2 != nil {
		log.Fatalf("Failed to call ProcessRequest auth2: %v", err2)
	}
	//calculate Shared key
	b_server_key := new(big.Int)
	b_server_key.SetString(response1.B, 10)
	// B^a mod p:
	shared_key := new(big.Int).Exp(b_server_key, a, p)

	////
	myKeys := keys{
		personalKeyClient: a,
		publicKeyClient:   public_key,
		sharedKeyClient:   shared_key,
	}

	// fmt.Println("Shared Key client:", myKeys.sharedKeyClient, " p:", p, "  g:", g, " public_key sent to server:", public_key, "  public received:", b_server_key)
	redis_key := fmt.Sprintf("%s:%s", response.GetNonce(), response.GetServerNonce())

	return myKeys.sharedKeyClient, redis_key, messageidd, nil
}

func BizService(redis_key string, message int32, c *gin.Context, userID int32) {
	// fmt.Println(userID, "iddddddddddd")
	grpcAddress := "localhost:50051"
	conn, err := grpc.Dial(grpcAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()
	client := grpcService_get_users.NewGetUsersClient(conn)

	request := &grpcService_get_users.GetDataRequest{
		UserId:    userID,
		AuthKey:   AuthKey_get.Bytes(),
		MessageId: message,
		RedisKey:  redis_key,
	}
	response, err := client.GetData(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Fatalf("Failed to get data from Biz service : %v", err)
	}
	c.JSON(http.StatusOK, gin.H{"data": response.ReturnUsers})

	if response.MessageId == 1 {
		user := response.ReturnUsers[0]
		fmt.Println("User ID: ", user.Id)
		fmt.Printf("Name: %s\n", user.Name)
		fmt.Printf("Family: %s\n", user.Family)
		fmt.Printf("Age: %d\n", user.Age)
		fmt.Printf("Sex: %s\n", user.Sex)
		fmt.Printf("Created At: %s\n", user.CreatedAt)
	} else if response.MessageId == 3 {
		for _, user := range response.ReturnUsers {
			fmt.Printf("User ID: ", user.Id)
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
func BizServiceWithSqlInject(redis_key string, message int32, c *gin.Context, userID string) {
	grpcAddress := "localhost:50053"
	conn, err := grpc.Dial(grpcAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := get_user_injection.NewGetUsersInjectClient(conn)
	request := &get_user_injection.GetDataRequestInject{
		UserId:    userID,
		AuthKey:   AuthKey_get.Bytes(),
		MessageId: message,
		RedisKey:  redis_key,
	}

	response, err := client.GetData(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Fatalf("Failed to get data from Biz service : %v", err)
	}
	c.JSON(http.StatusOK, gin.H{"data": response.ReturnUsers})

	if response.MessageId == 1 {
		user := response.ReturnUsers[0]
		fmt.Println("User ID: ", user.Id)
		fmt.Printf("Name: %s\n", user.Name)
		fmt.Printf("Family: %s\n", user.Family)
		fmt.Printf("Age: %d\n", user.Age)
		fmt.Printf("Sex: %s\n", user.Sex)
		fmt.Printf("Created At: %s\n", user.CreatedAt)
	} else if response.MessageId == 3 {
		for _, user := range response.ReturnUsers {
			fmt.Printf("User ID: ", user.Id)
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

func gatewayHandler(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err1 := strconv.ParseInt(userIDStr, 10, 32)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}
	x, redis_key, message, y := getAuthKey()
	AuthKey_get = x
	err := y
	if err != nil {
		log.Fatalf("Failed to get the auth key: %v", err)
	}
	if c.Writer.Status() == http.StatusTooManyRequests {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too Many Requests"})
		return
	} else if c.Writer.Status() == http.StatusForbidden {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	// Connect to the get_users service
	BizService(redis_key, message, c, int32(userID))
}

func gatewayHandlerSqlInject(c *gin.Context) {
	userID := c.Param("user_id")
	x, redis_key, message, y := getAuthKey()
	AuthKey_get = x
	err := y
	if err != nil {
		log.Fatalf("Failed to get the auth key: %v", err)
	}

	if c.Writer.Status() == http.StatusTooManyRequests {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too Many Requests"})
		return
	} else if c.Writer.Status() == http.StatusForbidden {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	// Connect to the get_users_with_sql_inject service
	BizServiceWithSqlInject(redis_key, message, c, userID)
}

func authenticateIP(c *gin.Context) {
	ip := c.ClientIP()
	// Check if the IP is blacklisted
	blacklistMu.Lock()
	if banTime, ok := blacklist[ip]; ok {
		if banTime.After(time.Now()) {
			c.AbortWithStatus(http.StatusForbidden)
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			blacklistMu.Unlock()
			return
		}
		delete(blacklist, ip)
	}
	blacklistMu.Unlock()

	ipDataMu.Lock()
	if _, ok := ipData[ip]; !ok {
		ipData[ip] = &IPData{
			Count:       0,
			LastRequest: time.Now(),
		}
	}
	ipDataMu.Unlock()

	ipDataMu.Lock()
	data := ipData[ip]
	data.Count++
	if data.Count > MaxRequestsPerSecond {
		blacklistMu.Lock()
		blacklist[ip] = time.Now().Add(BanDuration)
		blacklistMu.Unlock()
		c.AbortWithStatus(http.StatusTooManyRequests)
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too Many Requests"})
		ipDataMu.Unlock()
		return
	}
	data.LastRequest = time.Now()
	ipDataMu.Unlock()

	c.Next()
}

func resetRequestCount() {
	for {
		time.Sleep(1 * time.Second) //to count just for each second  //time.Millisecond
		ipDataMu.Lock()
		for _, data := range ipData {
			if time.Since(data.LastRequest) >= time.Second {
				data.Count = 0
			}
		}
		ipDataMu.Unlock()
	}
}

func cleanupBlacklist() {
	for {
		time.Sleep(3 * time.Second)
		currentTime := time.Now()
		blacklistMu.Lock()
		for ip, banTime := range blacklist {
			if banTime.Before(currentTime) {
				fmt.Println("remove from blacklist  --ip:  ", ip)
				delete(blacklist, ip)
			}
		}
		blacklistMu.Unlock()

		ipDataMu.Lock()
		for ip, data := range ipData {
			if data.LastRequest.Add(BanDuration).Before(currentTime) {
				delete(ipData, ip)
			}
		}
		ipDataMu.Unlock()
	}
}

func main() {
	router := gin.Default()
	go resetRequestCount()

	go cleanupBlacklist()
	router.Use(authenticateIP)

	router.GET("/gateway/get_users/:user_id", gatewayHandler)

	router.GET("/gateway/get_users_with_sql_inject/:user_id", gatewayHandlerSqlInject)
	router.Run(":8080")
}
