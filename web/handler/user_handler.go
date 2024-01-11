package handler

import (
	"strconv"

	"github.com/AurelienS/cigare/internal/glider"
	"github.com/AurelienS/cigare/internal/user"
	"github.com/AurelienS/cigare/web/session"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService   user.Service
	gliderService glider.Service
}

func NewUserHandler(userService user.Service, gliderService glider.Service) UserHandler {
	return UserHandler{
		userService:   userService,
		gliderService: gliderService,
	}
}

func (h *UserHandler) UpdateDefaultGlider(c echo.Context) error {
	user := session.GetUserFromContext(c)
	defaultGliderId := c.QueryParam("defaultGliderId")

	gliderId, err := strconv.Atoi(defaultGliderId)
	if err != nil {
		HandleError(c, err)
	}

	err = h.userService.UpdateDefaultGlider(c.Request().Context(), gliderId, user)
	if err != nil {
		return HandleError(c, err)
	}
	user.DefaultGliderID = gliderId
	session.SaveUserInSession(c, user)

	gliderHandler := NewGliderHandler(h.gliderService)
	return gliderHandler.GetGlidersCard(c)

}
