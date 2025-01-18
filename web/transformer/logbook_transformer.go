package transformer

import (
	"strconv"

	"github.com/AurelienS/airtime/internal/domain"
	"github.com/AurelienS/airtime/web/viewmodel"
)

func TransformLogbookToViewModel(
	flights []domain.Flight,
	flyingYears []int,
	year int,
) viewmodel.LogbookView {
	var flightViews []viewmodel.FlightView
	for _, f := range flights {
		flightViews = append(flightViews, TransformFlightToViewmodel(f))
	}

	flyingYearsString := make([]string, 0, len(flyingYears))
	for _, y := range flyingYears {
		flyingYearsString = append(flyingYearsString, strconv.Itoa(y))
	}

	return viewmodel.LogbookView{
		CurrentYear:    strconv.Itoa(year),
		AvailableYears: flyingYearsString,
		Flights:        flightViews,
	}
}
