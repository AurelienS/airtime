package main

import (
	"github.com/AurelienS/cigare/enhancer"
	"github.com/AurelienS/cigare/parser"
)

const igcFile = "/mnt/c/Users/TheGosu/Desktop/9-9-2023--13-47.igc"
const ouputFile = "test.json"
const outputFormat = "json"

func main() {
	flight, _ := parser.Parse()
	enhancer.EnhanceWithBearing(&flight)
	flight.GenerateStatistics()
	// fmt.Printf("file: main.go ~ line 16 ~ flight : %#v\n", flight.Thermals)
	flight.Thermals.Stats.PrettyPrint()

	// drawer.Drawer(flight)

}
