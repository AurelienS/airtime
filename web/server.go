package web

import (
	"github.com/AurelienS/cigare/internal/flight"
	"github.com/AurelienS/cigare/internal/squad"
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
	squadRepo := squad.NewRepository(queries)

	flightService := flight.NewService(flightRepo)
	userService := user.NewService(userRepo)
	squadService := squad.NewService(squadRepo, transactionManager)

	flightHandler := handler.NewFlightHandler(flightService)
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(userService)
	dashboardHandler := handler.NewDashboardHandler(userService, squadService, flightService)
	squadHandler := handler.NewSquadHandler(squadService)

	router := Router{
		AuthHandler:      authHandler,
		DashboardHandler: dashboardHandler,
		FlightHandler:    flightHandler,
		UserHandler:      userHandler,
		SquadHandler:     squadHandler,
	}
	router.Initialize(e)

	return &Server{Echo: e, Queries: &queries, Store: store}
}
