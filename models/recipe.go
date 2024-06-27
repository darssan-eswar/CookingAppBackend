package models

import (
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
)

type Recipe struct {
	ID          uuid.UUID    `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Ingredients []Ingredient `json:"ingredients"`
}

func CreateRecipe(db *sql.DB, name string, description string, ingredients []Ingredient) (Recipe, error) {
	recipe := Recipe{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Ingredients: ingredients,
	}

	jsonIngredients, err := json.Marshal(recipe.Ingredients)
	if err != nil {
		return Recipe{}, err
	}

	_, err = db.Exec("INSERT INTO recipes (id, name, description, ingredients) VALUES (?, ?, ?, ?)",
		recipe.ID, recipe.Name, recipe.Description, jsonIngredients)
	if err != nil {
		return Recipe{}, err
	}

	return recipe, nil
}

func GetUserRecipes(db *sql.DB, token string) ([]Recipe, error) {
	user, err := GetUser(db, uuid.MustParse(token))
	if err != nil {
		return nil, err
	}

	var recipeIds []string

	rows, err := db.Query("SELECT rid FROM user_recipe WHERE uid = ?", user.ID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var recipeId string
		if err := rows.Scan(&recipeId); err != nil {
			return nil, err
		}
		recipeIds = append(recipeIds, recipeId)
	}

	rows, err = db.Query("SELECT id, name, description, ingredients FROM recipes WHERE id IN (" + strings.Join(recipeIds, ",") + ")")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var recipes []Recipe
	for rows.Next() {
		var recipe Recipe
		var jsonIngredients string
		if err := rows.Scan(&recipe.ID, &recipe.Name, &recipe.Description, &jsonIngredients); err != nil {
			return nil, err
		}
		err := json.Unmarshal(([]byte)(jsonIngredients), &recipe.Ingredients)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, recipe)
	}

	return recipes, nil
}
