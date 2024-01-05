package main

import (
	"time"

	"github.com/AurelienS/cigare/drawer"
	"github.com/AurelienS/cigare/enhancer"
	"github.com/AurelienS/cigare/parser"
)

const igcFile = "/mnt/c/Users/TheGosu/Desktop/9-9-2023--13-47.igc"
const ouputFile = "test.json"
const outputFormat = "json"

func main() {
	const (
		minClimbRate               = 0.2              // m/s, the threshold for considering it thermic activity
		climbRateIntegrationPeriod = 10               // Number of seconds to smooth the climbRate
		minThermalDuration         = 20 * time.Second // Minimum duration to consider a sustained climb as thermal
		allowedDownwardPoints      = 4                // Number of consecutive downward points allowed in a thermal
	)

	flight, _ := parser.Parse()
	enhancer.EnhanceWithBearing(&flight)
	flight.GenerateThermals(minClimbRate, allowedDownwardPoints, minThermalDuration, climbRateIntegrationPeriod)

	flight.Stats.PrettyPrint()

	drawer.Drawer(flight)

}
