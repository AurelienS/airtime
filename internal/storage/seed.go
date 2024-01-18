package storage

import (
	"context"
	"log"
	"math/rand"
	"os"
	"time"
)

// random génère un nombre aléatoire entre min et max.
func random(min, max int) int {
	return rand.Intn(max-min+1) + min
}

// randomUniqueDays génère un ensemble unique de jours dans un mois.
func randomUniqueDays(year, month, numDays int) []time.Time {
	days := make(map[int]bool)
	var dates []time.Time

	for len(days) < numDays {
		day := random(1, 28) // Assurez-vous de respecter les jours du mois
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

	fileF, err := os.Create("seed_flights.txt")
	if err != nil {
		panic(err)
	}
	defer fileF.Close()
	fileFs, err := os.Create("seed_flights_stat.txt")
	if err != nil {
		panic(err)
	}
	defer fileFs.Close()

	for year := startYear; year <= endYear; year++ {
		for month := 1; month <= 12; month++ {
			numDays := random(0, 28)
			flightDates := randomUniqueDays(year, month, numDays)

			for _, flightDate := range flightDates {
				f, err := client.Flight.
					Create().
					SetDate(flightDate).
					SetTakeoffLocation("Lieu de décollage").
					SetIgcFilePath("Chemin/vers/fichier.igc").
					SetPilotID(userID).
					Save(ctx)
				if err != nil {
					log.Fatalf("failed creating flight: %v", err)
				}

				_, err = client.FlightStatistic.
					Create().
					SetTotalThermicTime(random(100, 1000)).
					SetTotalFlightTime(random(120, 1000)).
					SetMaxClimb(random(1100, 4808)).
					SetMaxClimbRate(float64(random(1, 15)) + rand.Float64()).
					SetTotalClimb(random(3000, 10000)).
					SetAverageClimbRate(rand.Float64() * 5).
					SetNumberOfThermals(random(1, 50)).
					SetPercentageThermic(rand.Float64() * 100).
					SetMaxAltitude(random(1500, 4500)).
					SetFlight(f).
					Save(ctx)
				if err != nil {
					log.Fatalf("failed creating flight statistic: %v", err)
				}
			}
		}
	}
}
