package webserver

import (
	"github.com/AurelienS/cigare/internal/auth"
	"github.com/AurelienS/cigare/internal/flight"
	"github.com/AurelienS/cigare/internal/glider"
	"github.com/AurelienS/cigare/internal/log"
	"github.com/AurelienS/cigare/internal/user"
	"github.com/labstack/echo/v4"
)

type Router struct {
	AuthHandler   auth.AuthHandler
	FlightHandler flight.FlightHandler
	GliderHandler glider.GliderHandler
	UserHandler   user.UserHandler
}

func (r *Router) Initialize(e *echo.Echo) {
	e.Use(log.LoggerMiddleware())

	e.GET("/login", r.AuthHandler.GetLogin)
	e.GET("/auth/:provider/callback", r.AuthHandler.GetAuthCallback)
	e.GET("/auth/:provider", r.AuthHandler.GetAuthProvider)

	authGroup := e.Group("/")
	authGroup.Use(auth.AuthMiddleware)

	authGroup.GET("", r.FlightHandler.GetIndexPage)
	authGroup.GET("logout", r.AuthHandler.GetLogout)

	authGroup.POST("flight", r.FlightHandler.PostFlight)
	authGroup.POST("glider", r.GliderHandler.PostGlider)

	authGroup.PUT("user/:userId", r.UserHandler.UpdateDefaultGlider)

}
