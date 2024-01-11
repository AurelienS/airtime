package web

import (
	"github.com/AurelienS/cigare/internal/flight"
	"github.com/AurelienS/cigare/internal/glider"
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

	gliderRepo := glider.NewSQLGliderRepository(queries)
	flightRepo := flight.NewFlightRepository(queries, transactionManager)
	userRepo := user.NewUserRepository(queries)

	gliderService := glider.NewGliderService(gliderRepo)
	flightService := flight.NewFlightService(flightRepo, gliderService)
	userService := user.NewUserService(userRepo)

	gliderHandler := handler.NewGliderHandler(gliderService)
	flightHandler := handler.NewFlightHandler(flightService, gliderService)
	userHandler := handler.NewUserHandler(userService, gliderService)
	authHandler := handler.NewAuthHandler(userService)

	router := Router{
		AuthHandler:   authHandler,
		FlightHandler: flightHandler,
		GliderHandler: gliderHandler,
		UserHandler:   userHandler,
	}
	router.Initialize(e)

	return &Server{Echo: e, Queries: &queries, Store: store}
}
