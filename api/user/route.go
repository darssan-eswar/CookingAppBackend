package user

import (
	"net/http"

	"github.com/labstack/echo"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// POST /auth/login
func Login(c echo.Context) error {

	var request loginRequest
	if err := c.Bind(&request); err != nil {
		return err
	}

	if request.Email == "" || request.Password == "" {
		return c.JSON(http.StatusBadRequest, "Email and password are required")
	}

	user, err := queryByEmailAndPassword(request.Email, request.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error)
	}

	if user == nil {
		return c.JSON(http.StatusBadRequest, "Incorrect password")
	}

	return c.JSON(http.StatusOK, user)
}

type registerRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// POST /auth/register
func Register(c echo.Context) error {

	var request registerRequest
	if err := c.Bind(&request); err != nil {
		return err
	}

	if request.Email == "" || request.Username == "" || request.Password == "" {
		return c.JSON(http.StatusBadRequest, "Email, username and password are required")
	}

	user, err := create(request.Email, request.Username, request.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error)
	}

	return c.JSON(http.StatusCreated, user)
}

// GET /auth/token
func LoginWithToken(c echo.Context) error {

	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusBadRequest, "Token is required")
	}

	user, err := queryByToken(token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error)
	}

	return c.JSON(http.StatusOK, user)
}

// DELETE /auth/logout
func Logout(c echo.Context) error {

	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusBadRequest, "Token is required")
	}

	err := clearToken(token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error)
	}

	return c.NoContent(http.StatusNoContent)
}
