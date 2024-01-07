package flight

import (
	"io"
	"net/http"

	"github.com/AurelienS/cigare/internal/webserver/handler"
	"github.com/AurelienS/cigare/internal/webserver/middleware"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func (h *FlightHandler) Upload(c echo.Context) error {
	file, err := c.FormFile("igcfile")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get IGC file from form")
		return handler.HandleError(c, err)
	}

	src, err := file.Open()
	if err != nil {
		log.Error().Err(err).Str("filename", file.Filename).Msg("Failed to open IGC file")
		return handler.HandleError(c, err)
	}
	defer src.Close()

	byteContent, err := io.ReadAll(src)
	if err != nil {
		log.Error().Err(err).Str("filename", file.Filename).Msg("Failed to read IGC file")
		return handler.HandleError(c, err)
	}

	user := middleware.GetUserFromContext(c)

	err = h.FlightService.UploadFlight(c.Request().Context(), byteContent, user)
	if err != nil {
		log.Error().Err(err).Str("user", user.Email).Msg("Failed to insert flight into database")
		return handler.HandleError(c, err)
	}

	log.Info().Str("user", user.Email).Str("filename", file.Filename).Msg("File parsed and flight record created successfully")
	return c.JSON(http.StatusOK, map[string]string{
		"message": "File parsed successfully",
	})
}
