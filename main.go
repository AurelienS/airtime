package main

import (
	"fmt"
	"time"

	"github.com/AurelienS/cigare/drawer"
	"github.com/AurelienS/cigare/parser"
)

const igcFile = "/mnt/c/Users/TheGosu/Desktop/9-9-2023--13-47.igc"
const ouputFile = "test.json"
const outputFormat = "json"

func main() {

	start := time.Now()
	flight, _ := parser.Parse()
	flight.Initialize()

	flight.Stats.PrettyPrint()
	drawer.Draw2DMap(flight, true)
	drawer.DrawElevation(flight)
	fmt.Println("file: main.go ~ line 25 ~ elapsed : ", time.Now().Sub(start))

}
