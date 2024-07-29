package models

import (
	"cookingapp/storage"
	"database/sql"
	"encoding/json"
)

type Recipe struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Ingredients []Ingredient `json:"ingredients"`
}

func newRecipe(name string, description string, ingredients []Ingredient) *Recipe {
	return &Recipe{
		Name:        name,
		Description: description,
		Ingredients: ingredients,
	}
}

func getRecipesFromRows(rows *sql.Rows) ([]Recipe, error) {

	var recipes []Recipe
	for rows.Next() {
		var encodedRecipe EncodedRecipe

		err := rows.Scan(&encodedRecipe.ID, &encodedRecipe.Name, &encodedRecipe.Description, &encodedRecipe.Ingredients)
		if err != nil {
			return nil, err
		}

		recipe, err := encodedRecipe.Decode()
		if err != nil {
			return nil, err
		}

		recipes = append(recipes, *recipe)
	}

	return recipes, nil
}

// Create

func CreateRecipeInDB(token string, name string, description string, ingredients []Ingredient) error {

	db, err := storage.GetDB()
	if err != nil {
		return err
	}

	recipe := newRecipe(name, description, ingredients)
	encodedRecipe, err := recipe.Encode()
	if err != nil {
		return err
	}

	row := db.QueryRow(`
		INSERT INTO recipes
		(id, name, description, ingredients)
		VALUES (?, ?, ?)
		RETURNING id
	`, encodedRecipe.Name, encodedRecipe.Description, encodedRecipe.Ingredients)

	var id int
	row.Scan(&id)

	err = CreateUserRecipeRelation(db, token, id)
	if err != nil {
		return err
	}

	return nil
}

// Read

func ReadRecipesFromDBWithToken(token string, limit int, offset int) ([]Recipe, error) {

	db, err := storage.GetDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(`
		SELECT r.id, r.name, r.description, r.ingredients
		FROM users u
		JOIN user_recipes ur ON u.id = ur.user_id
		JOIN recipes r ON ur.recipe_id = r.id
		WHERE u.token = ?
		LIMIT ? OFFSET ?
	`, token, limit, offset)
	if err != nil {
		return nil, err
	}

	return getRecipesFromRows(rows)
}

type EncodedRecipe struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Ingredients []byte `json:"ingredients"`
}

func (r *Recipe) Encode() (*EncodedRecipe, error) {

	jsonIngredients, err := json.Marshal(r.Ingredients)
	if err != nil {
		return nil, err
	}

	return &EncodedRecipe{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Ingredients: jsonIngredients,
	}, nil
}

func (r *EncodedRecipe) Decode() (*Recipe, error) {

	var ingredients []Ingredient
	err := json.Unmarshal(r.Ingredients, &ingredients)
	if err != nil {
		return nil, err
	}

	return &Recipe{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Ingredients: ingredients,
	}, nil
}
