package handler

import (
	"github.com/AurelienS/airtime/internal/service"
	"github.com/AurelienS/airtime/web/session"
	"github.com/AurelienS/airtime/web/transformer"
	"github.com/AurelienS/airtime/web/view/onboarding"
	"github.com/labstack/echo/v4"
)

type OnboardingHandler struct {
	flightService service.FlightService
}

func NewOnboardingHandler(flightService service.FlightService) OnboardingHandler {
	return OnboardingHandler{flightService: flightService}
}

func (h OnboardingHandler) Get(c echo.Context) error {
	user := session.GetUserFromContext(c)
	userview := transformer.TransformUserToViewModel(user)
	return Render(c, onboarding.Index(userview))
}
