package main

import (
	"github.com/AurelienS/airtime/internal/app"
	"github.com/AurelienS/airtime/internal/util"
)

const isProd = false // Set to true when serving over https

func main() {
	util.SetupLogger()
	server, err := app.Initialize(isProd)
	if err != nil {
		util.Fatal().Msg("Cannot initialize server")
		return
	}

	err = server.Start(":3000")
	if err != nil {
		util.Fatal().Msg("Cannot start server")
		return
	}
}
