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

	createTable()
	seedData()
}

// createTable creates the inventory table if it does not already exist
func createTable() {
	query := `
	CREATE TABLE IF NOT EXISTS inventory (
		id SERIAL PRIMARY KEY,
		item_name VARCHAR(255),
		sku VARCHAR(100) UNIQUE NOT NULL,
		stock_count INT NOT NULL,
		price DECIMAL(10,2),
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Error creating table:", err)
	}

	log.Println("✅ Inventory table ensured")
}

// seedData inserts initial sample data if the table is empty
func seedData() {
	var count int

	err := DB.QueryRow("SELECT COUNT(*) FROM inventory").Scan(&count)
	if err != nil {
		log.Fatal("Error checking inventory count:", err)
	}

	if count > 0 {
		log.Println("ℹ️ Seed data already exists, skipping insert")
		return
	}

	query := `
	INSERT INTO inventory (item_name, sku, stock_count, price) VALUES
	('Wireless Mouse', 'MOUSE-001', 25, 499.99),
	('Mechanical Keyboard', 'KEY-001', 15, 1999.99),
	('USB-C Cable', 'CABLE-001', 50, 199.99),
	('Laptop Stand', 'STAND-001', 8, 999.99),
	('Webcam', 'CAM-001', 5, 1499.99)
	ON CONFLICT (sku) DO NOTHING;
	`

	_, err = DB.Exec(query)
	if err != nil {
		log.Fatal("Error inserting seed data:", err)
	}

	log.Println("✅ Seed data inserted")
}
