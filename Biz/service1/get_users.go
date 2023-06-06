package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/big"
	"net"
	"time"

	"github.com/go-redis/redis/v8"
	pb "github.com/royadaneshi/webHW1/Biz/service1/get_userss/pb"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"

	"database/sql"

	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

type User struct {
	Id        int32 `gorm:"primarykey"`
	Name      string
	Family    string
	Age       int32
	Sex       string
	CreatedAt time.Time `gorm:"autoCreateTime:false"`
}

var (
	port = flag.Int("port", 50051, "gRPC server port")
)

type server struct {
	redisClient *redis.Client
	pb.UnimplementedGetUsersServer
}

// Function to generate sample users
func generateSampleUsers(count int) []User {
	users := make([]User, count)
	fmt.Println("innnn-----------------------")
	for i := 0; i < count; i++ {
		users[i] = User{
			Id:        int32(i + 1),
			Name:      fmt.Sprintf("User%d", i+1),
			Family:    fmt.Sprintf("Smith%d", i+1),
			Age:       25,
			Sex:       "Male",
			CreatedAt: time.Now(),
		}
	}
	fmt.Println(users[0], "    <-")

	return users
}
func DeleteAllRecords() error {
	err := DB.Exec("DELETE FROM users").Error
	if err != nil {
		return err
	}
	fmt.Println("All records deleted successfully.")
	return nil
}
func DatabaseConnection() {
	host := "webhw1_db_1"

	port := 5432
	user := "kamyar"
	password := "Karsanj1380"
	dbName := "kamyardb"

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	DB.AutoMigrate(User{})
	if err != nil {
		log.Fatal("Error connecting to the database...", err)
	}
	// Check if the table is empty
	var count int64
	DB.Model(&User{}).Count(&count)

	// delete
	// err := DeleteAllRecords()
	// if err != nil {
	// 	log.Fatal("Error deleting records...", err)
	// }
	//

	if count == 0 {
		// Insert sample records
		sampleUsers := generateSampleUsers(200)

		for _, user := range sampleUsers {
			if err := DB.Create(&user).Error; err != nil {
				log.Fatal("Error inserting sample records...", err)
			}
		}

		fmt.Println("Sample records inserted successfully.")
	}

	fmt.Println("Database connection successful...")

}

func (s *server) GetData(c context.Context, req *pb.GetDataRequest) (*pb.GetDataResponse, error) {

	var user User
	//read redis to get auth key to check validation of the recevied auth key
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis:6379", //redis port its true
		Password: "",           //dafault pass
		DB:       0,            //default db
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redissss: %v", err)
	}

	value, err := redisClient.Get(c, req.RedisKey).Result()
	if err != nil {
		log.Printf("Failed to retrieve data from Redis for get authkey: %v", err)
	}
	//convert from json
	byteData := []byte(value)
	var response *big.Int
	err1 := json.Unmarshal(byteData, &response)
	if err != nil {
		return nil, err1
	}
	authKey := new(big.Int)
	authKey.SetBytes(req.AuthKey)
	fmt.Println("client: ", authKey, "  server: ", response)
	if authKey.Cmp(response) != 0 {
		return nil, fmt.Errorf("invalid auth_key")
	} else {
		log.Printf("authentication: valid auth")

		res := DB.Find(&user, "id = ?", req.UserId)
		if res.Error != nil || string(user.Id) != string(req.UserId) {
			// return 100 first users from the table
			if res.Error == gorm.ErrRecordNotFound || string(user.Id) != string(req.UserId) {
				// Handle record not found error
				fmt.Println("get 100 first records...")
				var rows *sql.Rows
				rows, err := DB.Raw("SELECT id, name, family, age, sex, created_at FROM users LIMIT 100").Rows()
				if err != nil {
					return nil, err
				}
				defer rows.Close()
				first100Users := make([]*pb.Userr, 0)
				for rows.Next() {
					var data User

					err := rows.Scan(&data.Id, &data.Name, &data.Family, &data.Age, &data.Sex, &data.CreatedAt)

					if err != nil {
						return nil, err
					}
					first100Users = append(first100Users, &pb.Userr{
						Id:        int32(data.Id),
						Name:      data.Name,
						Family:    data.Family,
						Age:       int32(data.Age),
						Sex:       data.Sex,
						CreatedAt: data.CreatedAt.Format(time.RFC3339),
					})
				}
				if err = rows.Err(); err != nil {
					return nil, err
				}
				return &pb.GetDataResponse{
					ReturnUsers: first100Users,
					MessageId:   int32(3),
				}, nil

			} else {
				fmt.Println("another error, not get 100 first records...")
			}

		}

		messageIDResponse := int32(1)

		pbUser := &pb.Userr{
			Id:        int32(user.Id),
			Name:      user.Name,
			Family:    user.Family,
			Age:       int32(user.Age),
			Sex:       user.Sex,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		}

		return &pb.GetDataResponse{
			ReturnUsers: []*pb.Userr{pbUser},
			MessageId:   messageIDResponse,
		}, nil
	}
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
