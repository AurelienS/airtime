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
	session, err := getSession(c)
	if err != nil {
		return handleError(c, err)
	}

	session.Values["user"] = nil
	if err := saveSession(c, session); err != nil {
		return handleError(c, err)
	}

	gothic.Logout(c.Response(), c.Request())
	return c.Redirect(http.StatusFound, "/")
}

func (h *AuthHandler) GetAuthCallback(c echo.Context) error {
	user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return Render(c, page.Error())
	}

	session, err := getSession(c)
	if err != nil {
		return handleError(c, err)
	}

	session.Values["user"] = model.User{Email: user.Email}
	if err := saveSession(c, session); err != nil {
		return handleError(c, err)
	}

	return c.Redirect(http.StatusFound, "/")
}

func (h *AuthHandler) GetAuthProvider(c echo.Context) error {
	provider := c.Param("provider")
	expectedReq := c.Request().WithContext(context.WithValue(context.Background(), "provider", provider))

	gothic.BeginAuthHandler(c.Response(), expectedReq)
	return nil
}

func (h *AuthHandler) GetLogin(c echo.Context) error {
	return Render(c, page.Login())
}
