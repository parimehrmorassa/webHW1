package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	Auth_service "github.com/royadaneshi/webHW1/auth_service/auth/pb" //should changeee
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
	AuthKey_get string
)

func getAuthKey() (string, error) {
	// call  Auth service
	authResponse, err := clientAuth.GetAuthKey(context.Background(), &Auth_service.AuthRequest{})
	if err != nil {
		return "", err
	}
	return authResponse.AuthKey, nil
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

func main() {
	router := gin.Default()
	go cleanupBlacklist()

	// Connect to the Auth service
	grpcAddressAuth := "localhost:50054" //should changeee
	connAuth, errAuth := grpc.Dial(grpcAddressAuth, grpc.WithInsecure())
	if errAuth != nil {
		log.Fatalf("Failed to connect to Auth server: %v", errAuth)
	}
	defer connAuth.Close()

	clientAuth = Auth_service.NewAuthClient(connAuth)

	// rate limiting and IP banning
	router.Use(authenticateIP)

	// Call the Auth service to get the auth key
	x, y := getAuthKey()
	AuthKey_get = x
	err := y
	if err != nil {
		log.Fatalf("Failed to get the auth key: %v", err)
	}

	// Connect to the get_users service
	grpcAddress := "localhost:50051"
	conn, err := grpc.Dial(grpcAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client = grpcService_get_users.NewGetUsersClient(conn)

	router.GET("/gateway", gatewayHandler)
	router.Run(":8080")
}
