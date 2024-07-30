package recipe

import (
	"cookingapp/api/ingredient"
	"cookingapp/storage"
	"encoding/json"

	"github.com/google/uuid"
)

type Recipe struct {
	ID          string                  `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Ingredients []ingredient.Ingredient `json:"ingredients"`
}

// this struct exists so that we can store the ingredients array in the database as bytes
type encodedRecipe struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Ingredients []byte `json:"ingredients"`
}

// creates a new recipe with random uuid
func new(name string, description string, ingredients []ingredient.Ingredient) *Recipe {
	return &Recipe{
		ID:          uuid.NewString(),
		Name:        name,
		Description: description,
		Ingredients: ingredients,
	}
}

func (r *Recipe) encode() (*encodedRecipe, error) {
	jsonIngredients, err := json.Marshal(r.Ingredients)
	if err != nil {
		return nil, err
	}

	return &encodedRecipe{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Ingredients: jsonIngredients,
	}, nil
}

func (r *encodedRecipe) decode() (*Recipe, error) {
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

// creates a new recipe, writes it to the database, and adds a record to the user_recipes
// table to associate it with the user whose token is provided
func create(token string, name string, description string, ingredients []ingredient.Ingredient) (*Recipe, error) {

	db, err := storage.GetDB()
	if err != nil {
		return nil, err
	}

	// create and encode recipe
	recipe := new(name, description, ingredients)
	encoded, err := recipe.encode()
	if err != nil {
		return nil, err
	}

	// use transaction so that if on query fails, both inserts are rolled back
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	// insert recipe
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

	// create user-recipe relation
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

	// commit transaction
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return recipe, nil
}

// reads recipes from the database associated with the provided token
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
		var encodedRecipe encodedRecipe
		err := rows.Scan(&encodedRecipe.ID, &encodedRecipe.Name, &encodedRecipe.Description, &encodedRecipe.Ingredients)
		if err != nil {
			return nil, err
		}

		recipe, err := encodedRecipe.decode()
		if err != nil {
			return nil, err
		}

		recipes = append(recipes, *recipe)
	}

	return recipes, nil
}
