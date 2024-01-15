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
	SquadHandler   handler.SquadHandler
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

	authGroup.GET("logbook", r.LogbookHandler.GetPage)
	authGroup.POST("logbook/flight", r.LogbookHandler.PostFlight)

	authGroup.GET("create-squad", r.SquadHandler.GetCreateSquad)
	authGroup.POST("squad", r.SquadHandler.PostSquad)

	authGroup.GET("logout", r.AuthHandler.GetLogout)
}
