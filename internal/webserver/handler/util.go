package handler

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

func Render(c echo.Context, component templ.Component) error {
	return component.Render(c.Request().Context(), c.Response())
}

func HandleError(c echo.Context, err error) error {
	fmt.Printf("Error encountered: %s\n", err)
	return c.String(http.StatusInternalServerError, "Internal Server Error")
}

const sessionName = "session-name"

func getSession(c echo.Context) (*sessions.Session, error) {
	return gothic.Store.Get(c.Request(), sessionName)
}

func saveSession(c echo.Context, session *sessions.Session) error {
	return session.Save(c.Request(), c.Response())
}
