package web

import (
	"github.com/AurelienS/cigare/internal/flight"
	"github.com/AurelienS/cigare/internal/storage"
	"github.com/AurelienS/cigare/internal/user"
	"github.com/AurelienS/cigare/web/handler"
	"github.com/jackc/pgx/v5"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

type Server struct {
	*echo.Echo
	Queries *storage.Queries
	Store   sessions.Store
}

func NewServer(queries storage.Queries, db *pgx.Conn, store sessions.Store) *Server {
	e := echo.New()

	transactionManager := storage.NewTransactionManager(db)

	flightRepo := flight.NewRepository(queries, transactionManager)
	userRepo := user.NewRepository(queries)

	flightService := flight.NewService(flightRepo)
	userService := user.NewService(userRepo)

	flightHandler := handler.NewFlightHandler(flightService)
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(userService)

	router := Router{
		AuthHandler:   authHandler,
		FlightHandler: flightHandler,
		UserHandler:   userHandler,
	}
	router.Initialize(e)

	return &Server{Echo: e, Queries: &queries, Store: store}
}
