package session

import (
	"encoding/gob"
	"errors"

	"github.com/AurelienS/cigare/internal/storage"
	"github.com/AurelienS/cigare/internal/util"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

const sessionName = "session-name"

func NewStore(isProd bool) sessions.Store {
	maxAge := 86400 * 30 // 30 days

	store := sessions.NewCookieStore([]byte(sessionName))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProd

	return store
}

func ConfigureGoth(store sessions.Store) {
	goth.UseProviders(
		google.New("267580147813-11e4e5d00rboa7udei9mbiu50hht2c7q.apps.googleusercontent.com",
			"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",
			"http://localhost:3000/auth/google/callback",
			"email",
			"profile"),
	)
	gothic.Store = store
}

func ConfigureSessionStore(isProd bool) sessions.Store {
	store := NewStore(isProd)
	ConfigureGoth(store)
	gob.Register(storage.User{})
	return store
}

func SaveUserInSession(c echo.Context, user *storage.User) error {
	session, err := getSession(c)
	if err != nil {
		return err
	}
	session.Values["user"] = user
	return saveSession(c, session)
}

func RemoveUserFromSession(c echo.Context) error {
	user := GetUserFromContext(c)
	util.Info().Msgf("Removed %s from session", user.Email)

	return SaveUserInSession(c, nil)
}

func GetUserFromContext(c echo.Context) storage.User {
	session, err := getSession(c)
	user, ok := session.Values["user"].(storage.User)

	if err != nil || !ok || user.Email == "" {
		util.Fatal().Msg("Failed to get user from session")
		panic("Failed to get user from session")
	}

	return user
}

func GetUserOrErrorFromContext(c echo.Context) (storage.User, error) {
	session, err := getSession(c)
	user, ok := session.Values["user"].(storage.User)

	if err != nil || !ok || user.Email == "" {
		return user, errors.New("Nop")
	}

	return user, nil
}

func getSession(c echo.Context) (*sessions.Session, error) {
	session, err := gothic.Store.Get(c.Request(), sessionName)
	if err != nil {
		util.Error().Msgf("Failed to get session %s", sessionName)
	}
	return session, err
}

func saveSession(c echo.Context, session *sessions.Session) error {
	err := session.Save(c.Request(), c.Response())
	if err != nil {
		util.Error().Msgf("Failed to save session %s", sessionName)
	}
	return err
}
