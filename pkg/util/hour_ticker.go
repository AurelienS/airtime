package util

import (
	"math"
	"time"

	"gonum.org/v1/plot"
)

type HourTicker struct {
	StartTime time.Time
}

func (t HourTicker) Ticks(min, max float64) []plot.Tick {
	var ticks []plot.Tick
	step := math.Round((max - min) / 10)

	for x := min; x <= max; x += step {
		tickTime := t.StartTime.Add(time.Duration(x) * time.Minute)

		label := tickTime.Format("15:04")

		ticks = append(ticks, plot.Tick{Value: x, Label: label})
	}
	return ticks
}
