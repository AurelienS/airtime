package main

import (
	logconfig "github.com/AurelienS/cigare/internal"
	"github.com/AurelienS/cigare/internal/app"
	"github.com/rs/zerolog/log" // Global logger
)

const isProd = false // Set to true when serving over https

func main() {
	log.Logger = logconfig.SetupLogger()

	server := app.Initialize(isProd)
	if server == nil {
		return
	}

	server.Logger.Fatal(server.Start(":3000"))
}
