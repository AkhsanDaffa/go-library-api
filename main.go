package main

import (
	"database/sql"
	"go-library-api/config"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	log.Println("Memulai Aplikasi Perpustakaan")

	connStr := config.GetDBConnectionString()

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Failed Open DB: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed Connect DB: %v", err)
	}

	log.Println("Connect DB Successful")
}
