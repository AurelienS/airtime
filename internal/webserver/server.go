package webserver

import (
	"github.com/AurelienS/cigare/internal/storage"
	"github.com/AurelienS/cigare/internal/webserver/handler"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

type Server struct {
	*echo.Echo
	Queries *storage.Queries
	Store   sessions.Store
}

func NewServer(queries *storage.Queries, store sessions.Store) *Server {
	e := echo.New()

	authHandler := handler.AuthHandler{Queries: *queries}
	flightHandler := handler.FlightHandler{Queries: *queries}

	router := Router{
		AuthHandler:   authHandler,
		FlightHandler: flightHandler,
	}
	router.Initialize(e)

	return &Server{Echo: e, Queries: queries, Store: store}
}
