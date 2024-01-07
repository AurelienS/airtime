package handler

import (
	"context"
	"io"
	"net/http"

	"github.com/AurelienS/cigare/internal/storage"
	"github.com/AurelienS/cigare/internal/webserver/middleware"
	"github.com/AurelienS/cigare/pkg/model"
	"github.com/AurelienS/cigare/web/template/flight"
	"github.com/AurelienS/cigare/web/template/page"
	goingc "github.com/ezgliding/goigc/pkg/igc"
	"github.com/labstack/echo/v4"
)

type FlightHandler struct {
	Queries storage.Queries
}

func (h *FlightHandler) GetIndexPage(c echo.Context) error {
	return c.Redirect(http.StatusFound, "/flights")
}

func (h *FlightHandler) GetGlidersPage(c echo.Context) error {
	return Render(c, page.Gliders())
}

func (h *FlightHandler) GetFlightsPage(c echo.Context) error {
	return Render(c, page.Flights())
}

func (h *FlightHandler) GetFlights(c echo.Context) error {
	user := middleware.GetUserFromContext(c)
	flights, err := h.Queries.GetFlights(context.Background(), user.ID)
	if err != nil {
		return handleError(c, err)
	}
	return Render(c, flight.FlightRecords(flights))
}

func (h *FlightHandler) GetGliders(c echo.Context) error {
	user := middleware.GetUserFromContext(c)
	gliders, err := h.Queries.GetGliders(context.Background(), user.ID)
	if err != nil {
		return handleError(c, err)
	}
	return Render(c, flight.GliderCard(gliders))
}

func (h *FlightHandler) Upload(c echo.Context) error {
	file, err := c.FormFile("igcfile")
	if err != nil {
		return handleError(c, err)
	}

	src, err := file.Open()
	if err != nil {
		return handleError(c, err)
	}
	defer src.Close()

	byteContent, err := io.ReadAll(src)
	if err != nil {
		return handleError(c, err)
	}

	// Parse the IGC file using goigc
	track, err := goingc.Parse(string(byteContent))
	if err != nil {
		// Handle parsing error
		return handleError(c, err)
	}

	user := middleware.GetUserFromContext(c)
	flight := model.ConvertToMyFlight(track)

	params := storage.InsertFlightParams{
		Date:            flight.Date,
		TakeoffLocation: flight.Points[0].Description,
		UserID:          user.ID,
		GliderID:        1,
		IgcFilePath:     "not yet",
	}
	err = h.Queries.InsertFlight(context.Background(), params)
	if err != nil {
		return handleError(c, err)
	}

	// Respond to the client
	return c.JSON(http.StatusOK, map[string]string{
		"message": "File parsed successfully",
	})
}
