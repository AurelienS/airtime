package user

import (
	"context"
	"fmt"

	"github.com/AurelienS/cigare/internal/auth"
	"github.com/AurelienS/cigare/internal/glider"
	"github.com/AurelienS/cigare/internal/log"
	"github.com/AurelienS/cigare/internal/util"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService   UserService
	gliderService glider.GliderService
}

func NewUserHandler(userService UserService, gliderService glider.GliderService) UserHandler {
	return UserHandler{
		userService:   userService,
		gliderService: gliderService,
	}
}

func (h *UserHandler) UpdateDefaultGlider(c echo.Context) error {
	user := auth.GetUserFromContext(c)
	defaultGliderId := c.QueryParam("defaultGliderId")

	err := h.userService.UpdateDefaultGlider(context.Background(), defaultGliderId, user)
	if err != nil {
		return util.HandleError(c, err)
	}

	user, err = h.userService.repo.queries.GetUserWithGoogleId(context.Background(), user.GoogleID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch user with Google ID")
		return util.HandleError(c, err)
	}

	auth.SaveUserInSession(c, user) // update user cause we changed it

	gliderHandler := glider.NewGliderHandler(h.gliderService)
	fmt.Println("file: user_handler.go ~ line 43 ~ WILL CALL gliderHandler.GetGlidersCard")
	return gliderHandler.GetGlidersCard(c)
}
