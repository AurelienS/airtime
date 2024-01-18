package web

import (
	"github.com/AurelienS/cigare/web/handler"
	"github.com/AurelienS/cigare/web/middleware"
	"github.com/labstack/echo/v4"
)

type Router struct {
	AuthHandler    handler.AuthHandler
	LogbookHandler handler.LogbookHandler
	UserHandler    handler.UserHandler
	IndexHandler   handler.IndexHandler
}

func (r *Router) Initialize(e *echo.Echo) {
	e.Use(middleware.LoggerMiddleware())

	e.Static("/static", "web/static/")

	e.GET("/login", r.AuthHandler.GetLogin)
	e.GET("/auth/:provider/callback", r.AuthHandler.GetAuthCallback)
	e.GET("/auth/:provider", r.AuthHandler.GetAuthProvider)

	authGroup := e.Group("/")
	authGroup.Use(middleware.AuthMiddleware)
	authGroup.GET("", r.IndexHandler.Get)

	authGroup.GET("logbook/log", r.LogbookHandler.Get)
	authGroup.GET("logbook/progression", r.LogbookHandler.GetProgression)
	authGroup.GET("logbook/log/:year", r.LogbookHandler.Get)
	authGroup.POST("logbook/log/flight", r.LogbookHandler.PostFlight)

	authGroup.GET("logout", r.AuthHandler.GetLogout)
}
