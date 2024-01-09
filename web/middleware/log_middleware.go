package middleware

import (
	"time"

	"github.com/AurelienS/cigare/internal/util"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// LoggerMiddleware returns a middleware that logs HTTP requests.
func LoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c) // Call the next handler.

			var logLevel *zerolog.Event
			if err != nil {
				logLevel = util.Warn()
			} else {
				logLevel = util.Info()
			}

			logLevel.Str("method", c.Request().Method).
				Str("path", c.Request().URL.Path).
				Int("status", c.Response().Status).
				Dur("latency", time.Since(start)).
				Msg("request handled")

			return err // Return any errors from the handler.
		}
	}
}
