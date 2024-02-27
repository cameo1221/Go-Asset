package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// Database represents the PostgreSQL database connection
type Database struct {
	Conn *sql.DB
}

// Connect connects to the PostgreSQL database
func Connect() (*Database, error) {
	// Construct the connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	// Open a connection to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Ping the database to ensure connection is established
	if err = db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Connected to the database")

	return &Database{Conn: db}, nil
}

// Close closes the database connection
func (db *Database) Close() error {
	return db.Conn.Close()
}
