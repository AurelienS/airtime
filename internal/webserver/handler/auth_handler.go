package handler

import (
	"context"
	"net/http"

	"github.com/AurelienS/cigare/internal/storage"
	"github.com/AurelienS/cigare/web/template/page"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

type AuthHandler struct {
	Queries storage.Queries
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
	googleUser, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return Render(c, page.Error())
	}

	session, err := getSession(c)
	if err != nil {
		return handleError(c, err)
	}

	googleUserId := googleUser.UserID

	h.Queries.UpsertUser(context.Background(), storage.UpsertUserParams{
		GoogleID:   googleUserId,
		Email:      googleUser.Email,
		Name:       googleUser.Name,
		PictureUrl: googleUser.AvatarURL,
	})

	// we need to refetch to get the actual db ID
	user, err := h.Queries.GetUserWithGoogleId(context.Background(), googleUserId)

	session.Values["user"] = user
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
