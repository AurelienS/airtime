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
	"github.com/rs/zerolog/log"
)

type FlightHandler struct {
	Queries storage.Queries
}

func (h *FlightHandler) GetIndexPage(c echo.Context) error {
	log.Info().Msg("Redirecting to flights page")
	return c.Redirect(http.StatusFound, "/flights")
}

func (h *FlightHandler) GetGlidersPage(c echo.Context) error {
	log.Info().Msg("Rendering gliders page")
	return Render(c, page.Gliders())
}

func (h *FlightHandler) GetFlightsPage(c echo.Context) error {
	log.Info().Msg("Rendering flights page")
	return Render(c, page.Flights())
}

func (h *FlightHandler) GetFlights(c echo.Context) error {
	user := middleware.GetUserFromContext(c)
	flights, err := h.Queries.GetFlights(context.Background(), user.ID)
	if err != nil {
		log.Error().Err(err).Str("user", user.Email).Msg("Failed to get flights")
		return handleError(c, err)
	}
	log.Info().Str("user", user.Email).Msg("Fetched flights successfully")
	return Render(c, flight.FlightRecords(flights))
}

func (h *FlightHandler) GetGliders(c echo.Context) error {
	user := middleware.GetUserFromContext(c)
	gliders, err := h.Queries.GetGliders(context.Background(), user.ID)
	if err != nil {
		log.Error().Err(err).Str("user", user.Email).Msg("Failed to get gliders")
		return handleError(c, err)
	}
	log.Info().Str("user", user.Email).Msg("Fetched gliders successfully")
	return Render(c, flight.GliderCard(gliders))
}

func (h *FlightHandler) Upload(c echo.Context) error {
	file, err := c.FormFile("igcfile")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get IGC file from form")
		return handleError(c, err)
	}

	src, err := file.Open()
	if err != nil {
		log.Error().Err(err).Str("filename", file.Filename).Msg("Failed to open IGC file")
		return handleError(c, err)
	}
	defer src.Close()

	byteContent, err := io.ReadAll(src)
	if err != nil {
		log.Error().Err(err).Str("filename", file.Filename).Msg("Failed to read IGC file")
		return handleError(c, err)
	}

	track, err := goingc.Parse(string(byteContent))
	if err != nil {
		log.Error().Err(err).Str("filename", file.Filename).Msg("Failed to parse IGC file")
		return handleError(c, err)
	}

	user := middleware.GetUserFromContext(c)
	flight := model.ConvertToMyFlight(track)

	params := storage.InsertFlightParams{
		Date:            flight.Date,
		TakeoffLocation: flight.Points[0].Description,
		UserID:          user.ID,
		GliderID:        1,         // Assuming GliderID is 1, this should be dynamically set based on your application's logic
		IgcFilePath:     "not yet", // Placeholder path, replace with actual storage path as needed
	}
	err = h.Queries.InsertFlight(context.Background(), params)
	if err != nil {
		log.Error().Err(err).Str("user", user.Email).Msg("Failed to insert flight into database")
		return handleError(c, err)
	}

	log.Info().Str("user", user.Email).Str("filename", file.Filename).Msg("File parsed and flight record created successfully")
	return c.JSON(http.StatusOK, map[string]string{
		"message": "File parsed successfully",
	})
}
