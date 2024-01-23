package transformer

import (
	"fmt"
	"strconv"
	"time"

	"github.com/AurelienS/cigare/internal/domain"
	"github.com/AurelienS/cigare/web/viewmodel"
)

func TransformDashboardToViewModel(
	allTimeStats domain.StatsAggregated,
	currentYearStats domain.StatsAggregated,
	lastFlights []domain.Flight,
	sitesStats viewmodel.DashboardSitesStatsView,
	user domain.User,
) viewmodel.DashboardView {
	var lastFlightsView []viewmodel.DashboardFlightView
	for _, f := range lastFlights {
		lastFlightsView = append(lastFlightsView, TransformFlightToDashboardViewmodel(f))
	}

	currentYearStatsView := viewmodel.DashboardCurrentYearStatsView{
		FlightCount:       strconv.Itoa(currentYearStats.FlightCount),
		TotalFlightTime:   PrettyDuration(currentYearStats.TotalFlightTime),
		TotalDistance:     "-1 m",
		AverageFlightTime: PrettyDuration(currentYearStats.AverageFlightLength),
		MaxDuration:       PrettyDuration(currentYearStats.MaxDurationFLight.Statistic.TotalFlightTime),
		MaxDurationFlight: TransformFlightToDashboardViewmodel(currentYearStats.MaxDurationFLight),
		MaxDistance:       PrettyDistance(currentYearStats.MaxDurationFLight.Statistic.TotalDistance, false),
		MaxDistanceFlight: TransformFlightToDashboardViewmodel(currentYearStats.MaxDurationFLight),
		MaxAltitude:       PrettyDistance(currentYearStats.MaxAltitudeFlight.Statistic.MaxAltitude, true),
		MaxAltitudeFlight: TransformFlightToDashboardViewmodel(currentYearStats.MaxAltitudeFlight),
	}

	allTimeStatsView := viewmodel.DashboardStatsView{
		FlightCount:       strconv.Itoa(allTimeStats.FlightCount),
		TotalFlightTime:   PrettyDuration(allTimeStats.TotalFlightTime),
		TotalDistance:     PrettyDistance(allTimeStats.TotalDistance, false),
		AverageFlightTime: PrettyDuration(allTimeStats.AverageFlightLength),
		MaxDuration:       PrettyDuration(allTimeStats.MaxDurationFLight.Statistic.TotalFlightTime),
		MaxDurationFlight: TransformFlightToDashboardViewmodel(allTimeStats.MaxDurationFLight),
		MaxDistance:       PrettyDistance(allTimeStats.MaxDurationFLight.Statistic.TotalDistance, false),
		MaxDistanceFlight: TransformFlightToDashboardViewmodel(allTimeStats.MaxDurationFLight),
		MaxAltitude:       PrettyDistance(allTimeStats.MaxAltitudeFlight.Statistic.MaxAltitude, true),
		MaxAltitudeFlight: TransformFlightToDashboardViewmodel(allTimeStats.MaxAltitudeFlight),
	}

	return viewmodel.DashboardView{
		LastFlights:     lastFlightsView,
		SitesStats:      viewmodel.DashboardSitesStatsView{},
		User:            TransformUserToViewModel(user),
		CurrentYearStat: currentYearStatsView,
		CurrentYear:     strconv.Itoa(time.Now().Year()),
		FirstYear:       strconv.Itoa(allTimeStats.FirstFlight.Date.Year()),
		LastYear:        strconv.Itoa(allTimeStats.LastFlight.Date.Year()),
		AllTimeStats:    allTimeStatsView,
	}
}

func TransformFlightToDashboardViewmodel(flight domain.Flight) viewmodel.DashboardFlightView {
	return viewmodel.DashboardFlightView{
		Date:            flight.Date.Format("02/01/2006 15h04"),
		TakeoffLocation: flight.TakeoffLocation,
		TotalFlightTime: PrettyDuration(flight.Statistic.TotalFlightTime),
		TotalDistance:   PrettyDistance(flight.Statistic.TotalDistance, false),
		MaxAltitude:     PrettyDistance(flight.Statistic.MaxAltitude, true),
		FlightNumber:    "-1",
		Link:            fmt.Sprintf("/logbook/flight/%d", flight.ID),
	}
}
