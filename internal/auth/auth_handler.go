package auth

import (
	"context"
	"net/http"

	"github.com/AurelienS/cigare/internal/log"
	"github.com/AurelienS/cigare/internal/storage"
	"github.com/AurelienS/cigare/internal/util"
	"github.com/AurelienS/cigare/web/template/page"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

type AuthHandler struct {
	Queries storage.Queries
}

func NewAuthHandler(queries *storage.Queries) *AuthHandler {
	return &AuthHandler{Queries: *queries}
}

func (h *AuthHandler) GetLogout(c echo.Context) error {
	session, err := getSession(c)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get session during logout")
		return util.HandleError(c, err)
	}

	session.Values["user"] = nil
	if err := saveSession(c, session); err != nil {
		log.Error().Err(err).Msg("Failed to save session during logout")
		return util.HandleError(c, err)
	}

	gothic.Logout(c.Response(), c.Request())
	log.Info().Msg("User logged out successfully")
	return c.Redirect(http.StatusFound, "/")
}

func (h *AuthHandler) GetAuthCallback(c echo.Context) error {
	googleUser, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		log.Error().Err(err).Msg("Failed to complete user authentication with Google")
		return util.Render(c, page.Error())
	}

	session, err := getSession(c)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get session during auth callback")
		return util.HandleError(c, err)
	}

	googleUserId := googleUser.UserID

	err = h.Queries.UpsertUser(context.Background(), storage.UpsertUserParams{
		GoogleID:   googleUserId,
		Email:      googleUser.Email,
		Name:       googleUser.Name,
		PictureUrl: googleUser.AvatarURL,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to upsert user")
		return util.HandleError(c, err)
	}

	// we need to refetch to get the actual db ID
	user, err := h.Queries.GetUserWithGoogleId(context.Background(), googleUserId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch user with Google ID")
		return util.HandleError(c, err)
	}

	session.Values["user"] = user
	if err := saveSession(c, session); err != nil {
		log.Error().Err(err).Msg("Failed to save session after getting auth callback")
		return util.HandleError(c, err)
	}

	log.Info().Str("user", user.Email).Msg("User authenticated and session updated successfully")
	return c.Redirect(http.StatusFound, "/")
}

type contextKey string

const providerKey contextKey = "provider"

func (h *AuthHandler) GetAuthProvider(c echo.Context) error {
	provider := c.Param("provider")
	log.Info().Str("provider", provider).Msg("Initiating authentication with provider")
	expectedReq := c.Request().WithContext(context.WithValue(context.Background(), providerKey, provider))

	gothic.BeginAuthHandler(c.Response(), expectedReq)

	return nil
}

func (h *AuthHandler) GetLogin(c echo.Context) error {
	log.Info().Msg("Rendering login page")
	return util.Render(c, page.Login())
}
