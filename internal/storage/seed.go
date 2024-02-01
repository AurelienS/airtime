//nolint:gosec
package storage

import (
	"context"
	"log"
	"math/rand"
	"time"
)

func random(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func randomUniqueDays(year, month, numDays int) []time.Time {
	days := make(map[int]bool)
	var dates []time.Time

	for len(days) < numDays {
		day := random(1, 28)
		if _, exists := days[day]; !exists {
			days[day] = true
			date := time.Date(year, time.Month(month), day, random(0, 23), random(0, 59), random(0, 59), 0, time.UTC)
			dates = append(dates, date)
		}
	}
	return dates
}

func WriteSeedData() {
	client := Open()

	ctx := context.Background()
	userID := 1
	startYear := 2010
	endYear := 2019

	for year := startYear; year <= endYear; year++ {
		for month := 1; month <= 12; month++ {
			numDays := random(0, 28)
			flightDates := randomUniqueDays(year, month, numDays)

			for _, flightDate := range flightDates {
				_, err := client.Flight.
					Create().
					SetDate(flightDate).
					// SetTakeoffLocation("Lieu de décollage").
					// SetIgcFilePath("Chemin/vers/fichier.igc").
					SetPilotID(userID).
					Save(ctx)
				if err != nil {
					log.Fatalf("failed creating flight: %v", err)
				}

				// _, err = client.FlightStatistic.
				// 	Create().
				// 	SetTotalThermicTime(random(0, 20000)).
				// 	SetTotalFlightTime(random(0, 20000)).
				// 	SetMaxClimb(random(0, 4808)).
				// 	SetMaxClimbRate(float64(random(1, 9)) + rand.Float64()).
				// 	SetTotalClimb(random(0, 20000)).
				// 	SetAverageClimbRate(rand.Float64() * 5).
				// 	SetNumberOfThermals(random(0, 500)).
				// 	SetPercentageThermic(rand.Float64() * 100).
				// 	SetMaxAltitude(random(1100, 4808)).
				// 	SetFlight(f).
				// 	Save(ctx)
				// if err != nil {
				// 	log.Fatalf("failed creating flight statistic: %v", err)
				// }
			}
		}
	}
}
