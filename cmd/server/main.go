package main

import (
	"log"

	"github.com/AurelienS/cigare/internal/storage"
	"github.com/AurelienS/cigare/internal/webserver"
	"github.com/labstack/echo/v4"
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

	queries, err := storage.Open()
	if err != nil {
		log.Fatal("Cannot open db")
		return
	}

	e := echo.New()
	router := webserver.Router{
		Handler: webserver.Handler{
			Queries: queries,
		},
	}
	router.Initialize(e)

	e.Logger.Fatal(e.Start(":3000"))

}
