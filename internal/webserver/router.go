package webserver

import (
	"github.com/AurelienS/cigare/internal/auth"
	"github.com/AurelienS/cigare/internal/flight"
	"github.com/AurelienS/cigare/internal/glider"
	"github.com/AurelienS/cigare/internal/log"
	"github.com/labstack/echo/v4"
)

type Router struct {
	AuthHandler   auth.AuthHandler
	FlightHandler flight.FlightHandler
	GliderHandler glider.GliderHandler
}

func (r *Router) Initialize(e *echo.Echo) {
	e.Use(log.LoggerMiddleware())
	e.GET("/login", r.AuthHandler.GetLogin)

	authGroup := e.Group("/")
	authGroup.Use(auth.AuthMiddleware)
	authGroup.GET("", r.FlightHandler.GetIndexPage)
	authGroup.GET("gliders", r.FlightHandler.GetGlidersPage)
	authGroup.GET("flights", r.FlightHandler.GetFlightsPage)
	authGroup.GET("flights/all", r.FlightHandler.GetFlights)
	authGroup.POST("flight", r.FlightHandler.Upload)
	authGroup.GET("glidersCard", r.GliderHandler.GetGlidersCard)

	authGroup.GET("logout", r.AuthHandler.GetLogout)

	e.GET("/auth/:provider/callback", r.AuthHandler.GetAuthCallback)
	e.GET("/auth/:provider", r.AuthHandler.GetAuthProvider)

}
