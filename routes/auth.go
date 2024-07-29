package routes

import (
	"cookingapp/models"
	"net/http"

	"github.com/labstack/echo"
)

type responseBody struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

type loginRequest struct {
	Email    string `body:"email"`
	Password string `body:"email"`
}

func Login(e echo.Context) error {

	var request loginRequest
	if err := e.Bind(&request); err != nil {
		return err
	}

	if request.Email == "" || request.Password == "" {
		return e.JSON(http.StatusBadRequest, "Email and password are required")
	}

	user, err := models.QueryUserByEmailAndPassword(request.Password, request.Password)
	if err != nil {
		return e.JSON(http.StatusBadRequest, "Invalid email or password")
	}

	return e.JSON(http.StatusOK, responseBody{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Token:    *user.Token,
	})
}

type registerRequest struct {
	Email    string `body:"email"`
	Username string `body:"username"`
	Password string `body:"password"`
}

func Register(e echo.Context) error {

	var request registerRequest

	if err := e.Bind(&request); err != nil {
		return err
	}

	if request.Email == "" || request.Username == "" || request.Password == "" {
		return e.JSON(http.StatusBadRequest, "Email, username and password are required")
	}

	user, err := models.CreateUserWithEmailUsernameAndPassword(request.Email, request.Username, request.Password)
	if err != nil {
		return e.JSON(http.StatusBadRequest, "Error creating user")
	}

	return e.JSON(http.StatusOK, responseBody{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Token:    *user.Token,
	})
}

func LoginWithToken(e echo.Context) error {

	token := e.Request().Header.Get("Authorization")

	if token == "" {
		return e.JSON(http.StatusBadRequest, "Token is required")
	}

	user, err := models.QueryUserByToken(token)
	if err != nil {
		return e.JSON(http.StatusBadRequest, "Invalid token")
	}

	return e.JSON(http.StatusOK, responseBody{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Token:    *user.Token,
	})
}

func Logout(e echo.Context) error {

	token := e.Request().Header.Get("Authorization")
	if token == "" {
		return e.JSON(http.StatusBadRequest, "Token is required")
	}

	err := models.ClearToken(token)
	if err != nil {
		return e.JSON(http.StatusBadRequest, "Invalid token")
	}

	return e.JSON(http.StatusOK, "Logged out")
}
