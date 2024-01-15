package handler

import (
	"fmt"
	"strconv"

	"github.com/AurelienS/cigare/internal/logbook"
	"github.com/AurelienS/cigare/internal/util"
	"github.com/AurelienS/cigare/web/session"
	flightView "github.com/AurelienS/cigare/web/template/flight"
	"github.com/AurelienS/cigare/web/template/page"
	"github.com/labstack/echo/v4"
)

type LogbookHandler struct {
	LogbookService logbook.Service
}

func NewLogbookHandler(logbookService logbook.Service) LogbookHandler {
	return LogbookHandler{
		LogbookService: logbookService,
	}
}

func (h *LogbookHandler) GetPage(c echo.Context) error {
	user := session.GetUserFromContext(c)

	var viewData flightView.DashboardView // := TransformDashboardToView(data)
	viewData.Img = user.PictureURL
	return Render(c, page.Flights())
}

func TransformDashboardToView(data logbook.DashboardData) flightView.DashboardView {
	var fv []flightView.FlightView
	for _, f := range data.Flights {
		fv = append(fv, flightView.FlightView{
			TakeoffLocation: f.TakeoffLocation,
			Date:            f.Date.Format("02/01 15h04"),
		})
	}

	return flightView.DashboardView{
		Flights:         fv,
		NumberOfFlight:  strconv.Itoa(len(data.Flights)),
		TotalFlightTime: fmt.Sprintf("%d", int(data.TotalFlightTime.Hours())),
	}
}

func (h *LogbookHandler) PostFlight(c echo.Context) error {
	fmt.Println("file: logbook_handler.go ~ line 55 ~ func ~ PostFlight : ")
	file, err := c.FormFile("igcfile")
	if err != nil {
		util.Error().Err(err).Msg("Failed to get IGC file from form")
		return err
	}

	user := session.GetUserFromContext(c)

	err = h.LogbookService.ProcessAndAddFlight(c.Request().Context(), file, user)
	if err != nil {
		util.Error().Err(err).Str("user", user.Email).Msg("Failed to process and insert flight")
		return err
	}

	return h.GetPage(c)
}
