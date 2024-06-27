package routes

import (
	"CookingApp/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func DBMiddleware(c *gin.Context) {
	db, err := utils.DBConnect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.Set("db", db)
	c.Next()
}

func AuthMiddleware(c *gin.Context) {
	apiToken := os.Getenv("API_KEY")
	authToken := c.GetHeader("Authorization")
	if authToken != apiToken {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	token := c.GetHeader("Token")
	c.Set("token", token)
	c.Next()
}
