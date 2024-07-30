package recipe

import (
	"cookingapp/api/ingredient"
	"cookingapp/storage"
	"encoding/json"
)

type Recipe struct {
	ID          string                  `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Ingredients []ingredient.Ingredient `json:"ingredients"`
}

type EncodedRecipe struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Ingredients []byte `json:"ingredients"`
}

func new(name string, description string, ingredients []ingredient.Ingredient) *Recipe {
	return &Recipe{
		Name:        name,
		Description: description,
		Ingredients: ingredients,
	}
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
	var ingredients []ingredient.Ingredient
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

func create(token string, name string, description string, ingredients []ingredient.Ingredient) (*Recipe, error) {

	db, err := storage.GetDB()
	if err != nil {
		return nil, err
	}

	recipe := new(name, description, ingredients)
	encoded, err := recipe.Encode()
	if err != nil {
		return nil, err
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(`
		INSERT INTO recipes
		(id, name, description, ingredients)
		VALUES
		(?, ?, ?, ?)
	`, encoded.ID, encoded.Name, encoded.Description, encoded.Ingredients)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = tx.Exec(`
		INSERT INTO user_recipes
		(user_id, recipe_id)
		VALUES
		((SELECT id FROM users WHERE token = ?), (SELECT id FROM recipes WHERE id = ?))
	`, token, encoded.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return recipe, nil
}

func readWithToken(token string, limit int, offset int) ([]Recipe, error) {

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
