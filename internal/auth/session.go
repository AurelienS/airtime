package auth

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

const Key = "session-name"

func NewStore(isProd bool) sessions.Store {
	maxAge := 86400 * 30 // 30 days

	store := sessions.NewCookieStore([]byte(Key))
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

const sessionName = "session-name"

func getSession(c echo.Context) (*sessions.Session, error) {
	return gothic.Store.Get(c.Request(), sessionName)
}

func saveSession(c echo.Context, session *sessions.Session) error {
	return session.Save(c.Request(), c.Response())
}
