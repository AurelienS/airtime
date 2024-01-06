package main

import (
	"encoding/gob"
	"log"

	"github.com/AurelienS/cigare/internal/model"
	"github.com/AurelienS/cigare/internal/storage"
	"github.com/AurelienS/cigare/internal/webserver"
	"github.com/AurelienS/cigare/internal/webserver/handler"
	"github.com/labstack/echo/v4"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

const igcFile = "/mnt/c/Users/TheGosu/Desktop/9-9-2023--13-47.igc"
const ouputFile = "test.json"
const outputFormat = "json"

func main() {

	// flight, _ := igcparser.Parse()
	// flight.Initialize()

	// flight.Stats.PrettyPrint()
	// flight.Draw2DMap(true)
	// flight.DrawElevation()

	key := "session-name" // Replace with your SESSION_SECRET or similar
	maxAge := 86400 * 30  // 30 days
	isProd := false       // Set to true when serving over https

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = isProd

	gothic.Store = store
	gob.Register(model.User{})

	goth.UseProviders(
		google.New("267580147813-11e4e5d00rboa7udei9mbiu50hht2c7q.apps.googleusercontent.com",
			"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",
			"http://localhost:3000/auth/google/callback",
			"email",
			"profile"),
	)

	queries, err := storage.Open()
	if err != nil {
		log.Fatal("Cannot open db")
		return
	}

	e := echo.New()

	authHandler := handler.AuthHandler{}
	flightHandler := handler.FlightHandler{
		Queries: queries,
	}

	router := webserver.Router{
		AuthHandler:   authHandler,
		FlightHandler: flightHandler,
	}
	router.Initialize(e)

	e.Logger.Fatal(e.Start(":3000"))

}
