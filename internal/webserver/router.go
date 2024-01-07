package webserver

import (
	"github.com/AurelienS/cigare/internal/webserver/handler"
	"github.com/AurelienS/cigare/internal/webserver/middleware"
	"github.com/labstack/echo/v4"
)

type Router struct {
	AuthHandler   handler.AuthHandler
	FlightHandler handler.FlightHandler
}

func (r *Router) Initialize(e *echo.Echo) {
	e.GET("/login", r.AuthHandler.GetLogin)

	authGroup := e.Group("/")
	authGroup.Use(middleware.AuthMiddleware)
	authGroup.GET("", r.FlightHandler.GetIndexPage)
	authGroup.GET("gliders", r.FlightHandler.GetGlidersPage)
	authGroup.GET("flights", r.FlightHandler.GetFlightsPage)
	authGroup.GET("flights/all", r.FlightHandler.GetFlights)
	authGroup.POST("flight", r.FlightHandler.Upload)
	authGroup.GET("gliders/all", r.FlightHandler.GetGliders)

	authGroup.GET("logout", r.AuthHandler.GetLogout)

	e.GET("/auth/:provider/callback", r.AuthHandler.GetAuthCallback)
	e.GET("/auth/:provider", r.AuthHandler.GetAuthProvider)

}
