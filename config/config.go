package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetDBConnectionString() string {
	if err := godotenv.Load(); err != nil {
		log.Println("Failed load .env")
	}

	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASS")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	dbName := os.Getenv("POSTGRES_DB")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, dbName)
}
