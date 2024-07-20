package main

import (
	"cookingapp/routes"
	"cookingapp/storage"
	"cookingapp/utils"

	"github.com/labstack/echo"
)

func main() {
	// Load .env file
	utils.LoadEnv()

	// Connect to database
	err := storage.InitDB()
	if err != nil {
		panic(err)
	}

	// Create router
	e := echo.New()

	// Use use auth middleware (api-key header)
	e.Use(routes.ApiKeyMiddleware)

	// Routes

	// Auth
	auth := e.Group("/auth")
	auth.POST("/login", routes.Login)
	auth.POST("/register", routes.Register)
	auth.GET("/token", routes.LoginWithToken)
	auth.GET("/logout", routes.Logout)

	// protected
	protected := e.Group("/protected")
	protected.Use(routes.AuthMiddleware)
	protected.Use(routes.AuthMiddleware)

	e.Logger.Fatal(e.Start(":8080"))
}
