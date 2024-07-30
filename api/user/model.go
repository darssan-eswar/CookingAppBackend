package user

import (
	"cookingapp/storage"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                string  `json:"id"`
	Username          string  `json:"username"`
	Email             string  `json:"email"`
	Password          string  `json:"password"`
	Token             *string `json:"token"`
	SubscriptionStart int64   `json:"subscriptionStart"`
	SubscriptionEnd   int64   `json:"subscriptionEnd"`
}

func new(email, username, password string) *User {
	token := uuid.New().String()
	u := User{
		Email:             email,
		Username:          username,
		Password:          password,
		Token:             &token,
		SubscriptionStart: 0,
		SubscriptionEnd:   0,
	}
	return &u
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func checkPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func isUsernameAvailable(email string) (bool, error) {

	db, err := storage.GetDB()
	if err != nil {
		return false, err
	}

	// select arbitrary column from users with email
	rows, err := db.Query("SELECT id FROM users WHERE email = ?", email)
	if err != nil {
		return false, err
	}

	// if no rows, return true
	return !rows.Next(), nil
}

func create(email string, username string, password string) (*User, error) {

	db, err := storage.GetDB()
	if err != nil {
		return nil, err
	}

	// check for existing email
	available, err := isUsernameAvailable(email)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, errors.New("email already in use")
	}

	// hash password for security
	passwordHash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user := new(email, username, passwordHash)

	// insert user into database
	_, err = db.Exec(`
		INSERT INTO users
		(id, email, username, password, token, subscription_start, subscription_end)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, user.ID, user.Email, user.Username, user.Password, user.Token, user.SubscriptionStart, user.SubscriptionEnd)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func queryByEmailAndPassword(email string, password string) (*User, error) {

	db, err := storage.GetDB()
	if err != nil {
		return nil, err
	}

	row := db.QueryRow(`
		SELECT
		id, email, username, password, token, subscription_start, subscription_end
		FROM users
		WHERE email = ?
	`, email)

	var user User
	err = row.Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.Token, &user.SubscriptionStart, &user.SubscriptionEnd)
	if err != nil {
		return nil, errors.New("no account with that email")
	}

	if !checkPasswordHash(password, user.Password) {
		return nil, errors.New("incorrect password")
	}

	return &user, nil
}

func queryByToken(token string) (*User, error) {

	db, err := storage.GetDB()
	if err != nil {
		return nil, err
	}

	row := db.QueryRow(`
		SELECT
		id, email, username, password, token, subscription_start, subscription_end
		FROM users
		WHERE token = ?
	`, token)

	var user User
	err = row.Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.Token, &user.SubscriptionStart, &user.SubscriptionEnd)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	return &user, nil
}

func clearToken(token string) error {

	db, err := storage.GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE users SET token = NULL WHERE token = ?", token)
	if err != nil {
		return err
	}

	return nil
}

func DeleteFromDB(id string) error {

	db, err := storage.GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}
