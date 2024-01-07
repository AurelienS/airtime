package webserver

import (
	"github.com/AurelienS/cigare/internal/service"
	"github.com/AurelienS/cigare/internal/storage"
	repo "github.com/AurelienS/cigare/internal/storage/repository"

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

	flightRepo := repo.SQLFlightRepository{
		Queries: *queries,
	}
	gliderRepo := repo.SQLGliderRepository{
		Queries: *queries,
	}
	flightService := service.FlightService{
		Repo: flightRepo,
	}
	gliderService := service.GliderService{
		Repo: gliderRepo,
	}
	authHandler := handler.AuthHandler{Queries: *queries}

	flightHandler := handler.FlightHandler{
		FlightService: flightService,
		GliderService: gliderService,
	}

	router := Router{
		AuthHandler:   authHandler,
		FlightHandler: flightHandler,
	}
	router.Initialize(e)

	return &Server{Echo: e, Queries: queries, Store: store}
}
