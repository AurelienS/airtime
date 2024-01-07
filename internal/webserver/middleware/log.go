package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// LoggerMiddleware returns a middleware that logs HTTP requests.
func LoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c) // Call the next handler.

			// After the handler has finished, log the request details.
			log.Info().
				Str("method", c.Request().Method).
				Str("path", c.Request().URL.Path).
				Int("status", c.Response().Status).
				Dur("latency", time.Since(start)).
				Msg("request handled")

			return err // Return any errors from the handler.
		}
	}
}
