package web

import (
	"github.com/AurelienS/cigare/internal/repository"
	"github.com/AurelienS/cigare/internal/service"
	"github.com/AurelienS/cigare/internal/storage/ent"
	"github.com/AurelienS/cigare/web/handler"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

type Server struct {
	*echo.Echo
	Store sessions.Store
}

func NewServer(client *ent.Client, store sessions.Store) *Server {
	e := echo.New()

	flightRepo := repository.NewFlightRepository(client)
	userRepo := repository.NewUserRepository(client)

	logbookService := service.NewLogbookService(flightRepo)
	flightService := service.NewFlightService(flightRepo)
	statisticService := service.NewStatisticService(flightRepo, flightService)
	userService := service.NewUserService(userRepo)

	indexHandler := handler.NewIndexHandler(flightService)
	logbookHandler := handler.NewLogbookHandler(logbookService, statisticService, flightService)
	progressionHandler := handler.NewProgressionHandler(statisticService)
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(userService)
	onboardingHandler := handler.NewOnboardingHandler(flightService)

	router := NewRouter(
		authHandler,
		logbookHandler,
		progressionHandler,
		userHandler,
		indexHandler,
		onboardingHandler,
	)
	router.Initialize(e)

	return &Server{Echo: e, Store: store}
}
