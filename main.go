package main

import (
	"cookingapp/api/middleware"
	"cookingapp/api/recipe"
	"cookingapp/api/user"
	"cookingapp/storage"
	"cookingapp/util"

	"github.com/labstack/echo"
)

func main() {
	// Load .env file
	util.LoadEnvFromPath("./.env")

	// Connect to database
	err := storage.InitDB()
	if err != nil {
		panic(err)
	}

	// Create router
	e := echo.New()

	// e.Use(middleware.ApiKeyMiddleware) // check for api key header

	// Routes

	// auth
	auth := e.Group("/auth")
	auth.POST("/login", user.Login)
	auth.POST("/register", user.Register)
	auth.GET("/token", user.LoginWithToken)
	auth.DELETE("/logout", user.Logout)

	// protected
	protected := e.Group("/protected")
	protected.Use(middleware.AuthMiddleware) // check for token header
	// recipes
	recipes := protected.Group("/recipes")
	recipes.GET("/", recipe.GetRecipes)
	recipes.POST("/", recipe.CreateRecipe)

	e.Logger.Fatal(e.Start(":8080"))
}
