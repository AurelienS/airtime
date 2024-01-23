package handler

import (
	"github.com/AurelienS/cigare/internal/service"
	"github.com/AurelienS/cigare/web/session"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return UserHandler{
		userService: userService,
	}
}

func (h UserHandler) PutTheme(c echo.Context) error {
	user := session.GetUserFromContext(c)
	theme := c.FormValue("theme")
	if theme == "" {
		theme = "dark"
	}
	user.Theme = theme
	updatedUser, err := h.userService.UpsertUser(c.Request().Context(), user)
	if err != nil {
		return err
	}
	err = session.SaveUserInSession(c, updatedUser)
	if err != nil {
		return err
	}
	return nil
}
