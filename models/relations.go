package models

import "database/sql"

func CreateUserRecipeRelation(db *sql.DB, token string, recipeId int) error {

	_, err := db.Exec(`
		INSERT INTO user_recipe
		(uid, rid)
		VALUES
		((SELECT id FROM users WHERE token = ?), ?)
	`, token, recipeId)

	return err
}
