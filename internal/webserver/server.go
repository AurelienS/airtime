package webserver

import (
	"github.com/AurelienS/cigare/internal/auth"
	"github.com/AurelienS/cigare/internal/flight"
	"github.com/AurelienS/cigare/internal/glider"
	"github.com/AurelienS/cigare/internal/storage"

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

	gliderRepo := glider.NewSQLGliderRepository(queries)
	gliderService := glider.NewGliderService(gliderRepo)
	gliderHandler := glider.NewGliderHandler(gliderService)

	flightRepo := flight.NewSQLFlightRepository(*queries)
	flightService := flight.NewFlightService(flightRepo)
	flightHandler := flight.NewFlightHandler(flightService, gliderService)

	authHandler := auth.NewAuthHandler(queries)

	router := Router{
		AuthHandler:   *authHandler,
		FlightHandler: *flightHandler,
		GliderHandler: *gliderHandler,
	}
	router.Initialize(e)

	return &Server{Echo: e, Queries: queries, Store: store}
}
