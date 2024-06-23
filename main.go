package main

import (
	"CookingApp/routes"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	// Connect to Turso database
	envErr := godotenv.Load(".env")
	if envErr != nil {
		log.Fatalf("Error loading .env file: %s", envErr)
	}

	dbUrl := os.Getenv("DB_URL")
	dbToken := os.Getenv("DB_TOKEN")

	db, err := sql.Open("libsql", fmt.Sprintf("%s?authToken=%s", dbUrl, dbToken))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create Gin router
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	// Define routes
	router.POST("/auth/login", routes.Login)
	router.POST("/auth/register", routes.Register)

	// Start server
	router.Run(":8080")
}
