package session

import (
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/AurelienS/cigare/internal/domain"
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
	gob.Register(domain.User{})
	return store
}

func SaveUserInSession(c echo.Context, user domain.User) error {
	session, err := getSession(c)
	if err != nil {
		return err
	}
	session.Values["user"] = user
	fmt.Println("file: session.go ~ line 54 ~ funcSaveUserInSession ~ user : ", user)
	return saveSession(c, session)
}

func RemoveUserFromSession(c echo.Context) error {
	u := GetUserFromContext(c)
	util.Info().Msgf("Removed %s from session", u.Email)

	nilUser := domain.User{}
	return SaveUserInSession(c, nilUser)
}

func GetUserFromContext(c echo.Context) domain.User {
	session, err := getSession(c)
	user, ok := session.Values["user"].(domain.User)

	if err != nil || !ok || user.Email == "" {
		util.Fatal().Msgf("Failed to get user from session %s", err)
		panic("Failed to get user from session")
	}

	return user
}

func GetUserOrErrorFromContext(c echo.Context) (domain.User, error) {
	session, err := getSession(c)
	user, ok := session.Values["user"].(domain.User)

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
