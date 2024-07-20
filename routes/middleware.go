package routes

import (
	"net/http"
	"os"

	"github.com/labstack/echo"
)

// Used for entire backend to protect from unauthorized access
func ApiKeyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		envKey := os.Getenv("API_KEY")
		key := c.Request().Header.Get("api-key")
		if key != envKey {
			return c.JSON(http.StatusBadRequest, "Unauthorized")
		}

		return next(c)
	}
}

// Used for protected routes where user's token is requried
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return c.JSON(http.StatusBadRequest, "Token is required")
		}

		return next(c)
	}
}
