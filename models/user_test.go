package models

import (
	"cookingapp/storage"
	"cookingapp/utils"
	"testing"
)

func TestAuthentication(t *testing.T) {

	// init

	utils.LoadEnv()
	storage.InitDB()

	// given

	email := "_testemail"
	username := "_testusername"
	password := "_testpassword"

	// test
	user, err := QueryUserByEmail(email)
	if err == nil {
		_ = DeleteUserFromDB(user.ID)
	}

	// test
	_, err = CreateUserWithEmailUsernameAndPassword(email, username, password)
	if err != nil {
		t.Error(err)
	}

	// test
	_, err = QueryUserByEmail(email)
	if err != nil {
		t.Error(err)
	}

	// test
	user, err = QueryUserByEmailAndPassword(email, password)
	if err != nil {
		t.Error(err)
	}

	// test
	if user.Email != email {
		t.Error(err)
	}

	// test
	if user.Username != username {
		t.Error(err)
	}

	// test
	if user.Token == nil {
		t.Error(err)
	}

	// test + cleanup
	err = DeleteUserFromDB(user.ID)
	if err != nil {
		t.Error(err)
	}
}
