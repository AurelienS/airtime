package handler

import (
	"context"
	"net/http"

	"github.com/AurelienS/cigare/internal/model"
	"github.com/AurelienS/cigare/web/template/page"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

type AuthHandler struct {
}

func (h *AuthHandler) GetLogout(c echo.Context) error {
	session, err := gothic.Store.Get(c.Request(), "session-name")
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	session.Values["user"] = nil

	err = session.Save(c.Request(), c.Response())
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	gothic.Logout(c.Response(), c.Request())

	return c.Redirect(http.StatusFound, "/flights")
}

func (h *AuthHandler) GetAuthCallback(c echo.Context) error {

	user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return Render(c, page.Error())
	}

	session, err := gothic.Store.Get(c.Request(), "session-name")
	if err != nil {
		// Handle error - could not get or create the session
		return Render(c, page.Error())
	}

	// Save the user in the session
	session.Values["user"] = model.User{Email: user.Email} // Storing the user object in the session
	err = session.Save(c.Request(), c.Response())          // Important: Save the session!
	if err != nil {
		return Render(c, page.Error())
	}

	return c.Redirect(http.StatusFound, "/flights")
}

func (h *AuthHandler) GetAuthProvider(c echo.Context) error {
	// echo does not setup request with what gothic expect
	expectedReq := c.Request().WithContext(context.WithValue(context.Background(), "provider", c.Param("provider")))

	gothic.BeginAuthHandler(c.Response(), expectedReq)
	return nil
}

func (h *AuthHandler) GetLogin(c echo.Context) error {
	return Render(c, page.Login())
}
