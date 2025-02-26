package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	// Get connection string from environment variable
	connStr := os.Getenv("POSTGRES_DSN")
	if connStr == "" {
		return nil, fmt.Errorf("POSTGRES_DSN is not set in environment variables")
	}

	// Open a connection to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %v", err)
	}

	// Verify the connection by pinging the database
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error verifying database connection: %v", err)
	}

	// Set up connection pool
	// Maximum number of connections in the pool
	db.SetMaxOpenConns(10)
	// Maximum number of idle connections in the pool
	db.SetMaxIdleConns(5)
	// Set connection max lifetime (important for long-running apps)
	db.SetConnMaxLifetime(30 * time.Minute)

	log.Println("Successfully connected to the database!")

	// Return the database connection
	return db, nil
}
