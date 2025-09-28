package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func initDB() (*sql.DB, error) {
	connStr := getEnv("DATABASE_URL", "user=postgres password=postgres dbname=devicesdb sslmode=disable port=5432")

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	log.Println("Database connection established successfully")
	return db, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
