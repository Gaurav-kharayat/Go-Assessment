🛠 Tech Stack
Golang
PostgreSQL
Gin (Web Framework)
⚙️ How to Run the Application
🔹 1. Clone the Repository
git clone https://github.com/Gaurav-kharayat/Go-Assessment.git
cd inventory-service
🔹 2. Setup Environment Variables

Create a .env file using the example:

cp .env.example .env

Update DB connection if needed:

DB_URL=postgres://postgres:password@localhost:5432/inventory_db?sslmode=disable
🔹 3. Setup Database
Create Database:
CREATE DATABASE inventory_db;
Run SQL Script:
psql -U postgres -d inventory_db -f sql/init.sql
🔹 4. Install Dependencies
go mod tidy
🔹 5. Run the Application
go run main.go

Third-Party Libraries Used
github.com/gin-gonic/gin
→ Used for building REST APIs with clean routing and middleware support.
github.com/lib/pq
→ PostgreSQL driver for Go to connect and interact with the database.
github.com/joho/godotenv
→ Loads environment variables from .env file for local development.