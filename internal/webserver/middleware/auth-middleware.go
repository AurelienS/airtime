package middleware

import (
	"fmt"
	"net/http"

	"github.com/AurelienS/cigare/internal/model"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := GetUserFromContext(c)
		if err != nil {
			return c.Redirect(http.StatusFound, "/login")
		}

		return next(c)
	}
}

func GetUserFromContext(c echo.Context) (model.User, error) {
	var user model.User
	session, err := gothic.Store.Get(c.Request(), "session-name")
	if err != nil {
		return user, err // Handle error appropriately
	}

	user, ok := session.Values["user"].(model.User)
	if !ok {
		return user, fmt.Errorf("no user") // User not found in session
	}

	if user.Email == "" {
		return user, fmt.Errorf("no user")
	}

	return user, nil
}
