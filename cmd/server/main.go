package main

import (
	"github.com/AurelienS/cigare/internal/app"
)

const isProd = false // Set to true when serving over https

func main() {
	server := app.Initialize(isProd)
	if server == nil {
		return
	}

	server.Logger.Fatal(server.Start(":3000"))
}
