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
	apiToken := os.Getenv("API_TOKEN")
	token := c.GetHeader("Authorization")
	if token != apiToken {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	c.Next()
}
