package transformer

import (
	"fmt"
	"strconv"

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
		lastFlightsView = append(lastFlightsView, viewmodel.DashboardFlightView{
			Date:            f.Date.Format("02/01/2006 15h04"),
			TakeoffLocation: f.TakeoffLocation,
			TotalFlightTime: PrettyDuration(f.Statistic.TotalFlightTime),
			TotalDistance:   PrettyDistance(f.Statistic.TotalDistance, false),
			FlightNumber:    "-1",
			Link:            fmt.Sprintf("/logbook/flight/%d", f.ID),
		})
	}

	currentYearStatsView := viewmodel.DashboardCurrentYearStatsView{
		FlightCount:       strconv.Itoa(currentYearStats.FlightCount),
		TotalFlightTime:   PrettyDuration(currentYearStats.TotalFlightTime),
		TotalDistance:     "-1 m",
		AverageFlightTime: PrettyDuration(currentYearStats.AverageFlightLength),
		MaxFlightTime:     PrettyDuration(currentYearStats.MaxFlightLength),
		MaxDistance:       "-1 m",
	}

	allTimeStatsView := viewmodel.DashboardStatsView{
		FlightCount:       strconv.Itoa(allTimeStats.FlightCount),
		TotalFlightTime:   PrettyDuration(allTimeStats.TotalFlightTime),
		TotalDistance:     PrettyDistance(allTimeStats.TotalDistance, false),
		AverageFlightTime: PrettyDuration(allTimeStats.AverageFlightLength),
		MaxFlightTime:     PrettyDuration(allTimeStats.MaxFlightLength),
		MaxDistance:       PrettyDistance(allTimeStats.MaxDistance, false),
	}

	return viewmodel.DashboardView{
		LastFlights:     lastFlightsView,
		SitesStats:      viewmodel.DashboardSitesStatsView{},
		User:            TransformUserToViewModel(user),
		CurrentYearStat: currentYearStatsView,
		AllTimeStats:    allTimeStatsView,
	}
}
