package transformer

import (
	"fmt"
	"strconv"
	"time"
)

func PrettyRate(rate float64) string {
	return fmt.Sprintf("%.1f m/s", rate)
}

func PrettyDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60

	onlyHour := hours >= 100
	if onlyHour && hours > 0 {
		return fmt.Sprintf("%d h", hours)
	}

	if hours > 0 {
		return fmt.Sprintf("%dh%02d", hours, minutes)
	}
	return fmt.Sprintf("%d min", minutes)
}

func PrettyAltitude(alt int, forceMeter bool) string {
	if forceMeter {
		return strconv.Itoa(alt) + " m"
	}
	km := alt / 1000
	m := alt % 1000

	if km > 0 {
		return fmt.Sprintf("%d km", km)
	}
	return fmt.Sprintf("%d m", m)
}
