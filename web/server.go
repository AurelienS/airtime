package web

import (
	"github.com/AurelienS/cigare/internal/logbook"
	"github.com/AurelienS/cigare/internal/storage/ent"
	"github.com/AurelienS/cigare/internal/user"
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

	flightRepo := logbook.NewRepository(client)
	userRepo := user.NewRepository(client)

	flightService := logbook.NewService(flightRepo)
	userService := user.NewService(userRepo)

	flightHandler := handler.NewLogbookHandler(flightService)
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(userService)

	router := Router{
		AuthHandler:    authHandler,
		LogbookHandler: flightHandler,
		UserHandler:    userHandler,
	}
	router.Initialize(e)

	return &Server{Echo: e, Store: store}
}
