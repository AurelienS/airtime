package main

import (
	"log"

	"github.com/AurelienS/cigare/internal/webserver"
)

const igcFile = "/mnt/c/Users/TheGosu/Desktop/9-9-2023--13-47.igc"
const ouputFile = "test.json"
const outputFormat = "json"

const isProd = false // Set to true when serving over https

func main() {

	// flight, _ := igcparser.Parse()
	// flight.Initialize()

	// flight.Stats.PrettyPrint()
	// flight.Draw2DMap(true)
	// flight.DrawElevation()

	// Create and configure the server
	e := webserver.NewServer(isProd)
	if e == nil {
		log.Fatal("Failed to configure the server")
		return
	}

	// Start the server
	err := e.Start(":3000")
	if err != nil {
		log.Fatal("Cannot start the server")
		return
	}

}
