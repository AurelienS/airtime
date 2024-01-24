package transformer

import (
	"fmt"
	"strconv"
	"time"

	"github.com/AurelienS/cigare/internal/domain"
	"github.com/AurelienS/cigare/web/viewmodel"
)

func TransformMultipleStatsToViewModel(stats domain.MultipleFlightStats) viewmodel.DashboardStatsView {
	fmt.Println(
		"file: dashboard_transformer.go ~ line 15 ~ funcTransformMultipleStatsToViewModel ~ strconv.Itoa(len(stats.Flights)) : ",
		strconv.Itoa(len(stats.Flights)),
	)
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
	allTimeStats domain.MultipleFlightStats,
	currentYearStats domain.MultipleFlightStats,
	lastFlights []domain.Flight,
	sitesStats viewmodel.DashboardSitesStatsView,
	user domain.User,
) viewmodel.DashboardView {
	var lastFlightsView []viewmodel.FlightView
	for _, f := range lastFlights {
		lastFlightsView = append(lastFlightsView, TransformFlightToViewmodel(f))
	}

	currentYearStatsView := TransformMultipleStatsToViewModel(currentYearStats)
	allTimeStatsView := TransformMultipleStatsToViewModel(allTimeStats)

	totalFlightCount := len(allTimeStats.Flights)
	firstYear := time.Now().Year()
	lastYear := time.Now().Year()
	if totalFlightCount > 0 {
		lastYear = allTimeStats.Flights[0].Date.Year()
	}
	if totalFlightCount > 1 {
		firstYear = allTimeStats.Flights[totalFlightCount-1].Date.Year()
	}

	return viewmodel.DashboardView{
		LastFlights:     lastFlightsView,
		SitesStats:      viewmodel.DashboardSitesStatsView{},
		User:            TransformUserToViewModel(user),
		CurrentYearStat: currentYearStatsView,
		CurrentYear:     strconv.Itoa(time.Now().Year()),
		FirstYear:       strconv.Itoa(firstYear),
		LastYear:        strconv.Itoa(lastYear),
		AllTimeStats:    allTimeStatsView,
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
