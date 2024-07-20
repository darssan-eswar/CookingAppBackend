package models

import (
	"cookingapp/storage"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                int     `json:"id"`
	Username          string  `json:"username"`
	Email             string  `json:"email"`
	Password          string  `json:"password"`
	Token             *string `json:"token"`
	SubscriptionStart int64   `json:"subscriptionStart"`
	SubscriptionEnd   int64   `json:"subscriptionEnd"`
}

// create new user (should only be used in register)
func newUser(username, email, password string) *User {
	token := uuid.New().String()

	u := User{
		Username:          username,
		Email:             email,
		Password:          password,
		Token:             &token,
		SubscriptionStart: 0,
		SubscriptionEnd:   0,
	}
	return &u
}

// parse user from row
func userFromRow(row *sql.Row) (*User, error) {
	var u User

	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Token, &u.SubscriptionStart, &u.SubscriptionEnd)
	if err != nil {
		return nil, err
	}

	return &u, err
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

func (u *User) GetSubscriptionStart() *time.Time {
	if u.SubscriptionStart == 0 {
		return nil
	}

	t := time.Unix(u.SubscriptionStart, 0)
	return &t
}

func (u *User) GetSubscriptionEnd() *time.Time {
	if u.SubscriptionEnd == 0 {
		return nil
	}

	t := time.Unix(u.SubscriptionEnd, 0)
	return &t
}

// CREATE

func CreateUserWithEmailUsernameAndPassword(email string, username string, password string) (*User, error) {

	db, err := storage.GetDB()
	if err != nil {
		return nil, err
	}

	passwordHash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user := newUser(username, email, passwordHash)

	_, err = db.Exec(`
					INSERT INTO users
					(username, email, password, token, subscription_start, subscription_end)
					VALUES (?, ?, ?, ?, ?, ?)
	`, user.Username, user.Email, user.Password, user.Token, user.SubscriptionStart, user.SubscriptionEnd)
	if err != nil {
		return nil, err
	}

	user, err = QueryUserByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// READ

func QueryUserByEmail(email string) (*User, error) {

	db, err := storage.GetDB()
	if err != nil {
		return nil, err
	}

	row := db.QueryRow(`
					SELECT
					id, username, email, password, token, subscription_start, subscription_end
					FROM users
					WHERE email = ?
	`, email)

	return userFromRow(row)
}

func QueryUserByToken(token string) (*User, error) {

	db, err := storage.GetDB()
	if err != nil {
		return nil, err
	}

	row := db.QueryRow(`
					SELECT
					id, username, email, password, token, subscription_start, subscription_end
					FROM users
					WHERE token = ?
	`, token)

	return userFromRow(row)
}

func QueryUserByEmailAndPassword(email string, password string) (*User, error) {

	db, err := storage.GetDB()
	if err != nil {
		return nil, err
	}

	row := db.QueryRow(`
					SELECT
					id, username, email, password, token, subscription_start, subscription_end
					FROM users
					WHERE email = ?
	`, email)

	user, err := userFromRow(row)
	if err != nil {
		return nil, err
	}

	if !checkPasswordHash(password, user.Password) {
		return nil, err
	}

	return user, nil
}

// UPDATE

func UpdateUser(user *User) error {

	db, err := storage.GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec(`
					UPDATE users
					SET
					username = ?,
					email = ?,
					password = ?,
					token = ?,
					subscription_start = ?,
					subscription_end = ?
					WHERE id = ?
	`, user.Username, user.Email, user.Password, user.Token, user.SubscriptionStart, user.SubscriptionEnd, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func ClearToken(token string) error {

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

// DELETE

func DeleteUserFromDB(id int) error {

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
