// GOLANG CONCEPT - WHAT IS "package main"?
// In Go, code is organized into packages. The package name "main" is special:
// It tells the Go compiler that this package should compile into an independent executable binary
// rather than a shared library. Every runnable Go program must have one "package main".
package main

import (
	"log"
	"net/http"
	"os"

	"crud-employee/config"
	"crud-employee/routes"

	"github.com/gin-gonic/gin"
)

// GOLANG CONCEPT - WHAT IS "func main()"?
// The main() function is the main entry point of the application.
// When you run `go run main.go` or execute the compiled binary, the Go runtime
// starts execution here. When func main() finishes or exits, the entire program terminates.
func main() {
	// 1. Initialize MySQL database connection
	db := config.ConnectDB()
	// defer ensures db.Close() runs when the main function terminates
	defer db.Close()

	// 2. Initialize Gin router
	// gin.Default() creates a Gin engine equipped with standard Logger and Recovery middlewares
	router := gin.Default()

	// 3. Configure CORS (Cross-Origin Resource Sharing) middleware
	// This enables the frontend to communicate with the API even if hosted on a different port or origin
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		// Handle preflight OPTIONS requests from browsers
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// 4. Serve Frontend Static Files
	// This allows users to access the frontend directly at http://localhost:8080/
	router.StaticFile("/", "./frontend/index.html")
	router.StaticFile("/style.css", "./frontend/style.css")
	router.StaticFile("/script.js", "./frontend/script.js")
	router.Static("/assets", "./frontend")

	// 5. Register REST API Routes
	routes.SetupEmployeeRoutes(router)

	// 6. Read SERVER_PORT from environment variable or default to 8080
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	// 7. Start HTTP Server
	log.Printf("==================================================")
	log.Printf("🚀 Server running on http://localhost:%s", serverPort)
	log.Printf("💻 Open frontend at: http://localhost:%s/", serverPort)
	log.Printf("📡 API Base URL at: http://localhost:%s/api/employees", serverPort)
	log.Printf("==================================================")

	err := router.Run(":" + serverPort)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
