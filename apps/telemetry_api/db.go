package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

func initDB() (*sql.DB, error) {
	connStr := getEnv("DATABASE_URL", "user=postgres password=postgres dbname=devicedb sslmode=disable port=5432")

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	log.Println("Database connection established successfully")
	return db, nil
}
