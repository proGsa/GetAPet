package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	httpserver "getapet-backend/internal/delivery/http"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func getDB() (*sql.DB, error) {
	_ = godotenv.Load()

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	return sql.Open("postgres", connStr)
}

// @title GetAPet Backend API
// @version 1.0
// @description An aggregator of pet lists from shelters and private individuals.
func main() {
	db, err := getDB()
	if err != nil {
		log.Fatal("Connection error:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("DB not reachable:", err)
	}

	fmt.Println("DB connected! Starting HTTP server on :8080")

	server := httpserver.NewHTTPServer(":8080")
	if err := server.Start(db); err != nil {
		log.Fatal("HTTP server failed:", err)
	}
}
