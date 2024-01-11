package handler

import (
	"fmt"

	"github.com/AurelienS/cigare/internal/glider"
	"github.com/AurelienS/cigare/internal/user"
	"github.com/AurelienS/cigare/web/session"
	"github.com/AurelienS/cigare/web/template/flight"
	"github.com/labstack/echo/v4"
)

type GliderHandler struct {
	GliderService glider.Service
}

func NewGliderHandler(gliderService glider.Service) GliderHandler {
	return GliderHandler{
		GliderService: gliderService,
	}
}

func (h *GliderHandler) PostGlider(c echo.Context) error {
	user := session.GetUserFromContext(c)
	gliderName := c.FormValue("gliderName")

	err := h.GliderService.AddGlider(c.Request().Context(), gliderName, user)
	if err != nil {
		return err
	}

	return h.GetGlidersCard(c)
}

func (h *GliderHandler) GetGlidersCard(c echo.Context) error {
	user := session.GetUserFromContext(c)
	gliders, err := h.GliderService.GetGliders(c.Request().Context(), user)
	if err != nil {
		return err
	}
	viewData := TransformGlidersToView(gliders, user)
	return Render(c, flight.GliderCard(viewData))
}

func TransformGlidersToView(gliders []glider.Glider, user user.User) []flight.GliderView {
	var gv []flight.GliderView
	for _, g := range gliders {
		gv = append(gv, TransformGliderToView(g, user))
	}
	return gv
}

func TransformGliderToView(glider glider.Glider, user user.User) flight.GliderView {
	isSelected := false
	if glider.ID == user.DefaultGliderID {
		isSelected = true
	}
	linkToUpdate := fmt.Sprintf("/user/%d?defaultGliderId=%d", user.ID, glider.ID)
	return flight.GliderView{
		Name:         glider.Name,
		LinkToUpdate: linkToUpdate,
		IsSelected:   isSelected,
		ID:           glider.ID,
	}
}
