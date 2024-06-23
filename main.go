package main

import (
	"CookingApp/routes"
	"CookingApp/utils"

	"github.com/gin-gonic/gin"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	// Load .env file
	utils.LoadEnv()

	// Create Gin router
	router := gin.Default()

	// Middleware
	router.Use(routes.AuthMiddleware)
	router.Use(routes.DBMiddleware)

	// Define routes
	router.POST("/auth/login", routes.Login)
	router.POST("/auth/register", routes.Register)

	// Start server
	router.Run(":8080")
}
