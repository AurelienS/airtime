package webserver

import (
	"github.com/AurelienS/cigare/internal/auth"
	"github.com/AurelienS/cigare/internal/flight"
	"github.com/AurelienS/cigare/internal/glider"
	"github.com/AurelienS/cigare/internal/storage"
	"github.com/AurelienS/cigare/internal/user"
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
	gliderService := glider.NewGliderService(gliderRepo)
	gliderHandler := glider.NewGliderHandler(gliderService)

	flightRepo := flight.NewFlightRepository(queries, transactionManager)
	flightService := flight.NewFlightService(flightRepo)
	flightHandler := flight.NewFlightHandler(flightService, gliderService)

	userRepo := user.NewUserRepository(queries)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService, gliderService)

	authHandler := auth.NewAuthHandler(queries)

	router := Router{
		AuthHandler:   authHandler,
		FlightHandler: flightHandler,
		GliderHandler: gliderHandler,
		UserHandler:   userHandler,
	}
	router.Initialize(e)

	return &Server{Echo: e, Queries: &queries, Store: store}
}
