package routes

import (
	"CookingApp/models"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createRecipeBody struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Ingredients []models.Ingredient `json:"ingredients"`
}

func CreateRecipe(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	var body createRecipeBody
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	recipe, err := models.CreateRecipe(db, body.Name, body.Description, body.Ingredients)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, recipe)
}

func GetUserRecipes(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	token := c.MustGet("token").(string)
	recipes, err := models.GetUserRecipes(db, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, recipes)
}
