package web

import (
	"github.com/AurelienS/cigare/web/handler"
	"github.com/AurelienS/cigare/web/middleware"
	"github.com/labstack/echo/v4"
)

type Router struct {
	authHandler        handler.AuthHandler
	logbookHandler     handler.LogbookHandler
	progressionHandler handler.ProgressionHandler
	userHandler        handler.UserHandler
	indexHandler       handler.IndexHandler
	onboardingHandler  handler.OnboardingHandler
}

func NewRouter(
	authHandler handler.AuthHandler,
	logbookHandler handler.LogbookHandler,
	progressionHandler handler.ProgressionHandler,
	userHandler handler.UserHandler,
	indexHandler handler.IndexHandler,
	onboardingHandler handler.OnboardingHandler,
) Router {
	return Router{
		authHandler:        authHandler,
		logbookHandler:     logbookHandler,
		progressionHandler: progressionHandler,
		userHandler:        userHandler,
		indexHandler:       indexHandler,
		onboardingHandler:  onboardingHandler,
	}
}

func (r *Router) Initialize(e *echo.Echo) {
	e.Use(middleware.LoggerMiddleware())

	e.Static("/static", "web/static/")

	e.GET("/login", r.authHandler.GetLogin)
	e.GET("/dummy", r.indexHandler.Dummy)
	e.GET("/landing", r.indexHandler.Landing)
	e.GET("/auth/:provider/callback", r.authHandler.GetAuthCallback)
	e.GET("/auth/:provider", r.authHandler.GetAuthProvider)

	authGroup := e.Group("/")
	authGroup.Use(middleware.AuthMiddleware)
	authGroup.GET("", r.indexHandler.Get)

	authGroup.GET("onboarding", r.onboardingHandler.Get)
	authGroup.GET("logbook", r.logbookHandler.GetLogbook)
	authGroup.GET("logbook/:year", r.logbookHandler.GetLogbook)
	authGroup.GET("logbook/flight/:flight", r.logbookHandler.GetFlight)
	authGroup.POST("logbook/flight", r.logbookHandler.PostFlight)

	authGroup.GET("progression", r.progressionHandler.GetProgression)

	authGroup.GET("logout", r.authHandler.GetLogout)
}
