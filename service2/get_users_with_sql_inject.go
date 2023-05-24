package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	pb "service2/get_users_with_sql_inject_proto/pb"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"

	"database/sql"

	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

type User struct {
	ID        string `gorm:"primarykey"`
	Name      string
	Family    string
	Age       int32
	Sex       string
	CreatedAt time.Time `gorm:"autoCreateTime:false"`
}

var (
	port = flag.Int("port", 50052, "gRPC server port")
)

type server struct {
	pb.UnimplementedGetUsersServer
}

func DatabaseConnection() {
	host := "localhost"
	port := 5432
	user := "postgres"
	password := "web14022"
	dbName := "hw1"

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	DB.AutoMigrate(User{})
	if err != nil {
		log.Fatal("Error connecting to the database...", err)
	}
	fmt.Println("Database connection successful...")

}
func (*server) GetData(c context.Context, req *pb.GetDataRequest) (*pb.GetDataResponse, error) {
	fmt.Print("get request", req.UserId)
	var user User

	res := DB.Find(&user, "id = "+req.UserId)

	if res.Error != nil {
		// return 100 first users from the table
		if res.Error == gorm.ErrRecordNotFound {
			// Handle record not found error
			fmt.Println("get 100 first records...")
		} else {
			// Handle other errors
			fmt.Println("another error, not get 100 first records...")
		}

		var rows *sql.Rows
		rows, err := DB.Raw("SELECT * FROM users LIMIT 100").Rows()
		if err != nil {
			fmt.Println("111111111111")
			return nil, err
		}
		defer rows.Close()

		first100Users := make([]*pb.User, 0)
		for rows.Next() {
			var data User
			err := rows.Scan(&data.ID, &data.Name, &data.Family, &data.Age, &data.Sex, &data.CreatedAt)
			if err != nil {
				fmt.Println("222222222222222")

				return nil, err
			}
			first100Users = append(first100Users, &pb.User{
				Id:        string(data.ID),
				Name:      data.Name,
				Family:    data.Family,
				Age:       int32(data.Age),
				Sex:       data.Sex,
				CreatedAt: data.CreatedAt.Format(time.RFC3339),
			})
		}
		if err = rows.Err(); err != nil {
			fmt.Println("333333333333")

			return nil, err
		}
		return &pb.GetDataResponse{
			ReturnUsers: first100Users,
			MessageId:   3,
		}, nil
	}

	messageIDResponse := int32(1)

	pbUser := &pb.User{
		Id:        string(user.ID),
		Name:      user.Name,
		Family:    user.Family,
		Age:       int32(user.Age),
		Sex:       user.Sex,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}

	return &pb.GetDataResponse{
		ReturnUsers: []*pb.User{pbUser},
		MessageId:   messageIDResponse,
	}, nil
}

func main() {
	fmt.Println("gRPC server running ...")
	DatabaseConnection()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterGetUsersServer(s, &server{})

	log.Printf("Server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve : %v", err)
	}

}
