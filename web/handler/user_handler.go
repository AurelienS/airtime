package handler

import (
	"errors"

	"github.com/AurelienS/cigare/internal/glider"
	"github.com/AurelienS/cigare/internal/user"
	"github.com/AurelienS/cigare/web/session"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService   user.UserService
	gliderService glider.GliderService
}

func NewUserHandler(userService user.UserService, gliderService glider.GliderService) UserHandler {
	return UserHandler{
		userService:   userService,
		gliderService: gliderService,
	}
}

func (h *UserHandler) UpdateDefaultGlider(c echo.Context) error {
	user := session.GetUserFromContext(c)
	defaultGliderId := c.QueryParam("defaultGliderId")

	err := h.userService.UpdateDefaultGlider(c.Request().Context(), defaultGliderId, user)
	if err != nil {
		return HandleError(c, err)
	}

	// //TODO
	// user, err = h.userService.repo.queries.GetUserWithGoogleId(context.Background(), user.GoogleID)
	// if err != nil {
	// 	util.Error().Err(err).Msg("Failed to fetch user with Google ID")
	// 	return HandleError(c, err)
	// }

	// session.SaveUserInSession(c, &user)

	// gliderHandler := NewGliderHandler(h.gliderService)
	// return gliderHandler.GetGlidersCard(c)

	return errors.New("TODO")
}
