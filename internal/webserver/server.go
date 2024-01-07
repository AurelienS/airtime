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

	flightRepo := flight.SQLFlightRepository{
		Queries: *queries,
	}
	flightService := flight.FlightService{
		Repo: flightRepo,
	}
	flightHandler := flight.FlightHandler{
		FlightService: flightService,
	}

	gliderRepo := glider.SQLGliderRepository{
		Queries: *queries,
	}
	gliderService := glider.GliderService{
		Repo: gliderRepo,
	}
	gliderHandler := glider.GliderHandler{
		GliderService: gliderService,
	}
	authHandler := auth.AuthHandler{Queries: *queries}

	router := Router{
		AuthHandler:   authHandler,
		FlightHandler: flightHandler,
		GliderHandler: gliderHandler,
	}
	router.Initialize(e)

	return &Server{Echo: e, Queries: queries, Store: store}
}
