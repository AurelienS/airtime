package main

import (
	"github.com/AurelienS/cigare/internal/app"
	"github.com/AurelienS/cigare/internal/log"
)

const isProd = false // Set to true when serving over https

func main() {

	log.SetupLogger()
	server := app.Initialize(isProd)
	if server == nil {
		return
	}

	server.Logger.Fatal(server.Start(":3000"))
}
