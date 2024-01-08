package auth

import (
	"errors"
	"net/http"

	"github.com/AurelienS/cigare/internal/log"
	"github.com/AurelienS/cigare/internal/storage"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

const UserContextKey = "user"

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := getUserFromSession(c)
		if err != nil {
			log.Error().Err(err).Str("endpoint", c.Request().RequestURI).Msg("Failed to get user from session")
			return c.Redirect(http.StatusFound, "/login")
		}

		c.Set(UserContextKey, user)
		return next(c)
	}
}

func getUserFromSession(c echo.Context) (storage.User, error) {
	var user storage.User
	session, err := gothic.Store.Get(c.Request(), "session-name")
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving session")
		return user, err
	}

	if tempUser, ok := session.Values["user"].(storage.User); ok {
		user = tempUser
		if user.Email == "" {
			errMsg := "user email not found in session"
			log.Error().Msg(errMsg)
			return user, errors.New(errMsg)
		}
	} else {
		errMsg := "no user found in session"
		log.Error().Msg(errMsg)
		return user, errors.New(errMsg)
	}

	return user, nil
}

func GetUserFromContext(c echo.Context) storage.User {
	user, err := getUserFromSession(c)
	if err != nil {
		log.Fatal().Str("endpoint", c.Request().RequestURI).Msg("No user in context")
		// log.Fatal terminates the program; consider if this is the desired behavior.
		// In a real-world application, you might want to handle this more gracefully.
	}
	return user
}
