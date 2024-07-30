package recipe

import (
	"cookingapp/api/ingredient"
	"net/http"

	"github.com/labstack/echo"
)

type getRecipesRequest struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}

// GET /protected/recipes?limit=<int>&offset=<int>
func GetRecipes(e echo.Context) error {

	var request getRecipesRequest
	if err := e.Bind(&request); err != nil {
		e.String(http.StatusBadRequest, err.Error())
	}

	if request.Limit == 0 || request.Limit > 30 {
		request.Limit = 10
	}

	token := e.Request().Header.Get("Authorization")

	recipes, err := readWithToken(token, request.Limit, request.Offset)
	if err != nil {
		e.String(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, recipes)
}

type createRecipeRequest struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Ingredients []ingredient.Ingredient `json:"ingredients"`
}

// POST /protected/recipes
func CreateRecipe(e echo.Context) error {

	var request createRecipeRequest
	if err := e.Bind(&request); err != nil {
		e.String(http.StatusBadRequest, err.Error())
	}

	token := e.Request().Header.Get("Authorization")

	recipe, err := create(token, request.Name, request.Description, request.Ingredients)
	if err != nil {
		e.String(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusCreated, recipe)
}
