package main

import (
	"github.com/AurelienS/cigare/internal/app"
	ourLogger "github.com/AurelienS/cigare/internal/log"
	"github.com/rs/zerolog/log"
)

const isProd = false // Set to true when serving over https

func main() {
	log.Logger = ourLogger.SetupLogger()

	server := app.Initialize(isProd)
	if server == nil {
		return
	}

	server.Logger.Fatal(server.Start(":3000"))
}
