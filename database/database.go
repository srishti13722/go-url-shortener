package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var DB *pgx.Conn
var RedisClient *redis.Client

func ConnectRedis(){
	//connect redis at localhost 6379

	RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	//Test Redis connection
	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil{
		fmt.Println("Failed to connect redis..", err)
	}else{
		fmt.Println("SuccessFully Connected to Redis..")
	}
}

func ConnectDB() {
	//load env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load env file..", err)
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DataBase url is not set..")
	}

	DB, err = pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Unable to coonect with DB", err)
	}

	fmt.Println("Successfully Coonect to DB!....")
}
