package transformer

import (
	"fmt"
	"strconv"
	"time"

	"github.com/AurelienS/cigare/internal/domain"
	"github.com/AurelienS/cigare/web/viewmodel"
)

func TransformMultipleStatsToViewModel(stats domain.MultipleFlightStats) viewmodel.DashboardStatsView {
	return viewmodel.DashboardStatsView{
		FlightCount:       strconv.Itoa(len(stats.Flights)),
		TotalDuration:     PrettyDuration(stats.DurationTotal),
		TotalDistance:     PrettyDistance(stats.DistanceTotal, false),
		AverageDuration:   PrettyDuration(stats.AverageDuration),
		DurationMax:       PrettyDuration(stats.DurationMaxFlight.Duration),
		DurationMaxFlight: TransformFlightToViewmodel(stats.DurationMaxFlight),
		DistanceMax:       PrettyDistance(stats.DistanceMaxFlight.Distance, false),
		DistanceMaxFlight: TransformFlightToViewmodel(stats.DistanceMaxFlight),
		AltitudeMax:       PrettyDistance(stats.AltitudeMaxFlight.AltitudeMax, true),
		AltitudeMaxFlight: TransformFlightToViewmodel(stats.AltitudeMaxFlight),
	}
}

func TransformDashboardToViewModel(
	allTimeStats, currentYearStats domain.MultipleFlightStats,
	lastFlights []domain.Flight,
	sitesStats viewmodel.DashboardSitesStatsView,
	user domain.User,
) viewmodel.DashboardView {
	lastFlightsView := make([]viewmodel.FlightView, len(lastFlights))
	for i, f := range lastFlights {
		lastFlightsView[i] = TransformFlightToViewmodel(f)
	}

	currentYear, firstYear, lastYear := deriveYears(allTimeStats)

	showAlltime := currentYear != firstYear
	allTimeTitle := fmt.Sprintf("%d - %d", firstYear, lastYear)
	if firstYear == lastYear {
		allTimeTitle = fmt.Sprintf("%d", firstYear)
	}

	return viewmodel.DashboardView{
		LastFlights:     lastFlightsView,
		SitesStats:      sitesStats,
		User:            TransformUserToViewModel(user),
		CurrentYearStat: TransformMultipleStatsToViewModel(currentYearStats),
		CurrentYear:     strconv.Itoa(currentYear),
		FirstYear:       strconv.Itoa(firstYear),
		LastYear:        strconv.Itoa(lastYear),
		AllTimeStats:    TransformMultipleStatsToViewModel(allTimeStats),
		ShowAllTime:     showAlltime,
		AllTimeTitle:    allTimeTitle,
	}
}

func TransformFlightToViewmodel(flight domain.Flight) viewmodel.FlightView {
	return viewmodel.FlightView{
		Fulldate:    flight.Date.Format("02/01/2006 15h04"),
		Date:        flight.Date.Format("02/01/2006"),
		Location:    flight.Location,
		Duration:    PrettyDuration(flight.Duration),
		Distance:    PrettyDistance(flight.Distance, false),
		AltitudeMax: PrettyDistance(flight.AltitudeMax, true),
		Link:        fmt.Sprintf("/logbook/flight/%d", flight.ID),
		ID:          strconv.Itoa(flight.ID),
	}
}

func deriveYears(stats domain.MultipleFlightStats) (currentYear, firstYear, lastYear int) {
	totalFlightCount := len(stats.Flights)
	currentYear = time.Now().Year()

	if totalFlightCount > 0 {
		lastYear = stats.Flights[0].Date.Year()
		firstYear = lastYear
	}

	if totalFlightCount > 1 {
		firstYear = stats.Flights[totalFlightCount-1].Date.Year()
	}

	return currentYear, firstYear, lastYear
}
