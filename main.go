package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type TimeResponse struct {
	CurrentTime string `json:"current_time"`
}

func main() {
	// Connect to MySQL database
	dsn := "root:Cyborg@09@tcp(127.0.0.1:3306)/toronto_time" // Replace 'password' with your MySQL root password
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Verify connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Create a Gin router
	router := gin.Default()

	// Define /current-time endpoint
	router.GET("/current-time", func(c *gin.Context) {
		// Get current time in Toronto timezone
		loc, err := time.LoadLocation("America/Toronto")
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to load timezone"})
			return
		}

		currentTime := time.Now().In(loc)

		// Insert current time into the database
		_, err = db.Exec("INSERT INTO time_log (timestamp) VALUES (?)", currentTime)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to log time in database"})
			return
		}

		// Return current time in JSON format
		response := TimeResponse{
			CurrentTime: currentTime.Format("2006-01-02 15:04:05 MST"),
		}
		c.JSON(200, response)
	})

	// Start the server
	router.Run(":8080")
}
