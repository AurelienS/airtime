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
	authGroup.GET("", r.FlightHandler.GetIndex)
	authGroup.GET("gliders", r.FlightHandler.GetGliders)
	authGroup.GET("flights", r.FlightHandler.GetFlights)
	authGroup.GET("logout", r.AuthHandler.GetLogout)

	e.GET("/auth/:provider/callback", r.AuthHandler.GetAuthCallback)
	e.GET("/auth/:provider", r.AuthHandler.GetAuthProvider)

}
