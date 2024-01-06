package middleware

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AurelienS/cigare/internal/model"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

const UserContextKey = "user"

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := getUserFromSession(c)
		if err != nil {
			return c.Redirect(http.StatusFound, "/login")
		}

		c.Set(UserContextKey, user)

		return next(c)
	}
}

func getUserFromSession(c echo.Context) (model.User, error) {
	var user model.User
	session, err := gothic.Store.Get(c.Request(), "session-name")
	if err != nil {
		return user, err
	}

	if tempUser, ok := session.Values["user"].(model.User); ok {
		user = tempUser
		if user.Email == "" {
			return user, fmt.Errorf("no user email")
		}
	} else {
		return user, fmt.Errorf("no user in session")
	}

	return user, nil
}

func GetUserFromContext(c echo.Context) model.User {
	user, ok := c.Get(UserContextKey).(model.User)
	if !ok {
		log.Fatal("no user")
	}
	return user
}
