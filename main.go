package main

import (
	"github.com/AurelienS/cigare/drawer"
	"github.com/AurelienS/cigare/enhancer"
	"github.com/AurelienS/cigare/parser"
	"github.com/AurelienS/cigare/phaser"
)

const igcFile = "/mnt/c/Users/TheGosu/Desktop/9-9-2023--13-47.igc"
const ouputFile = "test.json"
const outputFormat = "json"

func main() {
	flight, _ := parser.Parse()
	enhancer.EnhanceWithBearing(&flight)
	phaser.Phase(flight)

	drawer.Drawer(flight)

}
