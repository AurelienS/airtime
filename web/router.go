package web

import (
	"github.com/AurelienS/cigare/web/handler"
	"github.com/AurelienS/cigare/web/middleware"
	"github.com/labstack/echo/v4"
)

type Router struct {
	AuthHandler   handler.AuthHandler
	FlightHandler handler.FlightHandler
	GliderHandler handler.GliderHandler
	UserHandler   handler.UserHandler
}

func (r *Router) Initialize(e *echo.Echo) {
	e.Use(middleware.LoggerMiddleware())

	e.GET("/login", r.AuthHandler.GetLogin)
	e.GET("/auth/:provider/callback", r.AuthHandler.GetAuthCallback)
	e.GET("/auth/:provider", r.AuthHandler.GetAuthProvider)

	authGroup := e.Group("/")
	authGroup.Use(middleware.AuthMiddleware)

	authGroup.GET("", r.FlightHandler.GetIndexPage)
	authGroup.GET("logout", r.AuthHandler.GetLogout)

	authGroup.POST("flight", r.FlightHandler.PostFlight)
	authGroup.POST("glider", r.GliderHandler.PostGlider)

	authGroup.PUT("user/:userId", r.UserHandler.UpdateDefaultGlider)

}
