package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID `json:"id"`
	Username        string    `json:"username"`
	Email           string    `json:"email"`
	Password        string    `json:"-"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Token           uuid.UUID `json:"token"`
	TokenExpiration time.Time `json:"token_expiration"`
}

func LoginUser(db *sql.DB, email string, password string) (User, error) {
	var user User
	var createdAt, updatedAt, tokenExpiration sql.NullTime
	var id, token sql.NullString

	err := db.QueryRow("SELECT id, email, username, password, created_at, updated_at, token, token_expiration FROM users WHERE email = ? AND password = ?", email, password).
		Scan(&id, &user.Email, &user.Username, &user.Password, &createdAt, &updatedAt, &token, &tokenExpiration)

	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, errors.New("user not found")
		}
		return User{}, err
	}

	// Parse UUID and time fields
	if id.Valid {
		user.ID, _ = uuid.Parse(id.String)
	}
	if createdAt.Valid {
		user.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		user.UpdatedAt = updatedAt.Time
	}
	if token.Valid {
		user.Token, _ = uuid.Parse(token.String)
	}
	if tokenExpiration.Valid {
		user.TokenExpiration = tokenExpiration.Time
	}

	return user, nil
}

func RegisterUser(db *sql.DB, email string, username string, password string) (User, error) {
	user := User{
		ID:              uuid.New(),
		Email:           email,
		Username:        username,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Token:           uuid.New(),
		TokenExpiration: time.Now().Add(24 * time.Hour),
	}

	_, err := db.Exec("INSERT INTO users (id, email, username, password, created_at, updated_at, token, token_expiration) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		user.ID, user.Email, user.Username, user.Password, user.CreatedAt, user.UpdatedAt, user.Token, user.TokenExpiration)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func GetUser(db *sql.DB, token uuid.UUID) (User, error) {
	var user User
	var createdAt, updatedAt, tokenExpiration sql.NullTime
	var id, tokenString sql.NullString

	err := db.QueryRow("SELECT id, email, username, created_at, updated_at, token, token_expiration FROM users WHERE token = ?", token).
		Scan(&id, &user.Email, &user.Username, &createdAt, &updatedAt, &tokenString, &tokenExpiration)

	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, errors.New("user not found")
		}
		return User{}, err
	}

	// Parse UUID and time fields
	if id.Valid {
		user.ID, _ = uuid.Parse(id.String)
	}
	if createdAt.Valid {
		user.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		user.UpdatedAt = updatedAt.Time
	}
	if tokenString.Valid {
		user.Token, _ = uuid.Parse(tokenString.String)
	}
	if tokenExpiration.Valid {
		user.TokenExpiration = tokenExpiration.Time
	}

	return user, nil
}
