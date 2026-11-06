package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	// GOLANG CONCEPT - BLANK IMPORT (_):
	// The underscore (_) before the import is a blank identifier import.
	// It tells Go to execute the package's init() function without directly referencing
	// its exported names in our code.
	// Here, "github.com/go-sql-driver/mysql" registers itself as a driver named "mysql"
	// with the standard Go "database/sql" package.
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// GOLANG CONCEPT - GLOBAL DB INSTANCE:
// To keep this beginner-friendly and avoid overly complex dependency injection patterns,
// we export the DB pointer so our controllers can directly execute SQL queries.
var DB *sql.DB

// ConnectDB initializes and tests the connection to MySQL.
//
// GOLANG CONCEPT - HOW GOLANG CONNECTS TO MYSQL:
// 1. Golang uses the built-in "database/sql" package as a generic database interface.
// 2. A database-specific driver (like go-sql-driver/mysql) translates generic SQL operations
//    into MySQL-specific network protocols.
// 3. We construct a Data Source Name (DSN) string with credentials, host, port, and database name.
// 4. sql.Open() creates the database handle / connection pool.
// 5. db.Ping() checks if the connection actually works.
func ConnectDB() *sql.DB {
	// Attempt to load .env file; if it fails (e.g., in production or container),
	// we fall back to system environment variables.
	err := godotenv.Load()
	if err != nil {
		log.Println("Notice: No .env file found or error loading it. Using environment variables.")
	}

	// Read environment variables with sensible defaults
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "crud-employee")

	// DSN (Data Source Name) Format for MySQL:
	// username:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	var dsn string
	if dbPassword != "" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)
	} else {
		// When password is empty
		dsn = fmt.Sprintf("%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbHost, dbPort, dbName)
	}

	// GOLANG CONCEPT - sql.Open() vs db.Ping():
	// sql.Open() does NOT open an active network connection to MySQL immediately.
	// It only validates the DSN string format and initializes the connection pool manager.
	var openErr error
	DB, openErr = sql.Open("mysql", dsn)
	if openErr != nil {
		log.Fatalf("Failed to initialize database driver handle: %v", openErr)
	}

	// GOLANG CONCEPT - db.Ping():
	// Ping() actually sends a packet over the network to MySQL to verify that:
	// 1. The database server is running and reachable.
	// 2. The username and password are correct.
	// 3. The specified database exists.
	pingErr := DB.Ping()
	if pingErr != nil {
		log.Fatalf("Failed to connect to MySQL database: %v", pingErr)
	}

	// Recommended connection pool settings
	DB.SetMaxOpenConns(25)                 // Maximum number of open connections to the database
	DB.SetMaxIdleConns(10)                 // Maximum number of connections in the idle connection pool
	DB.SetConnMaxLifetime(5 * time.Minute) // Maximum amount of time a connection may be reused

	log.Printf("Successfully connected to MySQL database: %s\n", dbName)
	return DB
}

// Helper function to read an environment variable or return a fallback default value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}
