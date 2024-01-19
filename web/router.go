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
	e.GET("/dummy", r.IndexHandler.Dummy)
	e.GET("/landing", r.IndexHandler.Landing)
	e.GET("/auth/:provider/callback", r.AuthHandler.GetAuthCallback)
	e.GET("/auth/:provider", r.AuthHandler.GetAuthProvider)

	authGroup := e.Group("/")
	authGroup.Use(middleware.AuthMiddleware)
	authGroup.GET("", r.IndexHandler.Get)

	authGroup.GET("logbook/log", r.LogbookHandler.GetTabLog)
	authGroup.GET("logbook/log/:year", r.LogbookHandler.GetTabLog)
	authGroup.GET("logbook/log/flight/:flight", r.LogbookHandler.GetFlight)
	authGroup.POST("logbook/log/flight", r.LogbookHandler.PostFlight)

	authGroup.GET("logbook/progression", r.LogbookHandler.GetTabProgression)

	authGroup.GET("logout", r.AuthHandler.GetLogout)
}
