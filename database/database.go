package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

var DB *pgx.Conn

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
