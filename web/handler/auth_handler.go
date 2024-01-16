package handler

import (
	"context"
	"net/http"

	"github.com/AurelienS/cigare/internal/model"
	"github.com/AurelienS/cigare/internal/user"
	"github.com/AurelienS/cigare/internal/util"
	"github.com/AurelienS/cigare/web/session"
	"github.com/AurelienS/cigare/web/view/userview"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

type AuthHandler struct {
	userService user.Service
}

func NewAuthHandler(userService user.Service) AuthHandler {
	return AuthHandler{userService: userService}
}

func (h *AuthHandler) GetLogout(c echo.Context) error {
	err := session.RemoveUserFromSession(c)
	if err != nil {
		return err
	}
	err = gothic.Logout(c.Response(), c.Request())
	if err != nil {
		return err
	}
	return c.Redirect(http.StatusFound, "/")
}

func (h *AuthHandler) GetAuthCallback(c echo.Context) error {
	googleUser, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		util.Error().Err(err).Msg("Failed to complete user authentication with Google")
		return err
	}
	if googleUser.Email == "" {
		util.Error().Err(err).Msg("Failed to complete user authentication with Google (email is missing)")
		return err
	}

	user := model.User{
		GoogleID:   googleUser.UserID,
		Email:      googleUser.Email,
		Name:       googleUser.Name,
		PictureURL: googleUser.AvatarURL,
	}

	// aFIX INSERT  AND CASE WHERE ALREADY EXIST SHOULDNT BE ERROR BUT DB WILL THROW UNIQUE CONSTRAINT
	user, err = h.userService.UpsertUser(c.Request().Context(), user)
	if err != nil {
		return err
	}

	err = session.SaveUserInSession(c, user)
	if err != nil {
		return err
	}

	util.Info().Str("user", user.Email).Msg("User authenticated and session updated successfully")
	return c.Redirect(http.StatusFound, "/")
}

func (h *AuthHandler) GetAuthProvider(c echo.Context) error {
	provider := c.Param("provider")
	util.Info().Str("provider", provider).Msg("Initiating authentication with provider")

	//nolint:revive,staticcheck
	expectedReq := c.Request().WithContext(context.WithValue(c.Request().Context(), "provider", provider))
	gothic.BeginAuthHandler(c.Response(), expectedReq)

	return nil
}

func (h *AuthHandler) GetLogin(c echo.Context) error {
	util.Info().Msg("Rendering login page")
	return Render(c, userview.Login())
}
