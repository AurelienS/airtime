package web

import (
	"github.com/AurelienS/cigare/web/handler"
	"github.com/AurelienS/cigare/web/middleware"
	"github.com/labstack/echo/v4"
)

type Router struct {
	authHandler        handler.AuthHandler
	logbookHandler     handler.LogbookHandler
	statisticsHandler handler.StatisticsHandler
	userHandler        handler.UserHandler
	indexHandler       handler.IndexHandler
	onboardingHandler  handler.OnboardingHandler
	dashboardHandler   handler.DashboardHandler
}

func NewRouter(
	authHandler handler.AuthHandler,
	logbookHandler handler.LogbookHandler,
	statisticsHandler handler.StatisticsHandler,
	dashboardHandler handler.DashboardHandler,
	userHandler handler.UserHandler,
	indexHandler handler.IndexHandler,
	onboardingHandler handler.OnboardingHandler,
) Router {
	return Router{
		authHandler:        authHandler,
		logbookHandler:     logbookHandler,
		statisticsHandler: statisticsHandler,
		dashboardHandler:   dashboardHandler,
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

	authGroup.PUT("user/theme", r.userHandler.PutTheme)

	authGroup.GET("logbook", r.logbookHandler.GetLogbook)
	authGroup.GET("logbook/:year", r.logbookHandler.GetLogbook)
	authGroup.GET("logbook/flight/:flight", r.logbookHandler.GetFlight)
	authGroup.POST("logbook/flight", r.logbookHandler.PostFlight)

	authGroup.GET("dashboard", r.dashboardHandler.GetIndex)

	authGroup.GET("statistics", r.statisticsHandler.GetIndex)
	authGroup.GET("statistics/count/distinct", r.statisticsHandler.GetCountDistinct)
	authGroup.GET("statistics/count/cumulative", r.statisticsHandler.GetCountCumul)
	authGroup.GET("statistics/time/distinct", r.statisticsHandler.GetTimeDistinct)
	authGroup.GET("statistics/time/cumulative", r.statisticsHandler.GetTimeCumul)

	authGroup.GET("logout", r.authHandler.GetLogout)
}
