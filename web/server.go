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

	flightService := service.NewLogbookService(flightRepo)
	userService := service.NewUserService(userRepo)

	logbookHandler := handler.NewLogbookHandler(flightService)
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(userService)

	router := Router{
		AuthHandler:    authHandler,
		LogbookHandler: logbookHandler,
		UserHandler:    userHandler,
	}
	router.Initialize(e)

	return &Server{Echo: e, Store: store}
}
