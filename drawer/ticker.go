package drawer

import (
	"math"
	"time"

	"gonum.org/v1/plot"
)

// CustomTicker is a Ticker that formats the labels as hours and minutes.
type CustomTicker struct {
	StartTime time.Time // The start time of the flight
}

// Ticks computes the Ticks in a specified range
func (t CustomTicker) Ticks(min, max float64) []plot.Tick {
	var ticks []plot.Tick
	// Choose a reasonable step size depending on the duration of your flight
	step := math.Round((max - min) / 10) // Adjust this depending on your time range

	for x := min; x <= max; x += step {
		// Convert the float value back into a time.Time
		tickTime := t.StartTime.Add(time.Duration(x) * time.Minute) // Assuming x is in minutes

		// Generate the label for the tick
		label := tickTime.Format("15:04") // Hours:Minutes format

		ticks = append(ticks, plot.Tick{Value: x, Label: label})
	}
	return ticks
}
