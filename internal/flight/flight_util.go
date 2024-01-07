package flight

import (
	"math"
	"time"

	"gonum.org/v1/plot"
)

func Average(numbers []float64) float64 {
	total := 0.0
	for _, number := range numbers {
		total += number
	}
	if len(numbers) == 0 {
		return 0
	}
	return total / float64(len(numbers))
}

func UpdateRateOfClimbHistory(history []float64, rateOfClimb float64, period int) []float64 {
	history = append(history, rateOfClimb)
	if len(history) > period {
		history = history[1:]
	}
	return history
}

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
