package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB initializes the database connection using environment variables
func InitDB() {
	connStr := os.Getenv("DB_URL")

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("DB connection error:", err)
	}
	// Verify if the database is reachable
	if err = DB.Ping(); err != nil {
		log.Fatal("DB not reachable:", err)
	}

	log.Println("✅ Connected to DB")
}
