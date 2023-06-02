package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/royadaneshi/webHW1/service2/get_users_with_sql_inject_proto/pb"

	"google.golang.org/grpc"
)

func main() {
	serverAddress := "localhost:50053"
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewGetUsersInjectClient(conn)
	request := &pb.GetDataRequestInject{
		UserId: "1000000",
	}

	response, err := client.GetData(context.Background(), request)
	if err != nil {
		log.Fatalf("Failed to get data: %v", err)
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
