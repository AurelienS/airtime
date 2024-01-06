package webserver

import (
	"encoding/gob"
	"log"

	"github.com/AurelienS/cigare/internal/model"
	"github.com/AurelienS/cigare/internal/storage"
	"github.com/AurelienS/cigare/internal/webserver/handler"
	"github.com/AurelienS/cigare/internal/webserver/session"
	"github.com/labstack/echo/v4"
)

func NewServer(isProd bool) *echo.Echo {
	store := session.NewStore(isProd)
	session.ConfigureGoth(store)
	gob.Register(model.User{})

	queries, err := storage.Open()
	if err != nil {
		log.Fatal("Cannot open db")
		return nil
	}

	serv := echo.New()

	authHandler := handler.AuthHandler{}
	flightHandler := handler.FlightHandler{
		Queries: queries,
	}

	router := Router{
		AuthHandler:   authHandler,
		FlightHandler: flightHandler,
	}
	router.Initialize(serv)

	return serv
}
