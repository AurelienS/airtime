package middleware

import (
	"regexp"
	"strconv"
	"time"

	"github.com/AurelienS/cigare/internal/util"
	"github.com/Pallinder/go-randomdata"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// LoggerMiddleware returns a middleware that logs HTTP requests.
func LoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			requestID := generateShortID()
			util.Info().
				Str("method", c.Request().Method).
				Str("path", c.Request().URL.Path).
				Msgf("Request %s", requestID)

			err := next(c)

			var logLevel *zerolog.Event
			if err != nil {
				logLevel = util.Error()
				util.Debug().Msg(err.Error())
			} else {
				logLevel = util.Info()
			}

			logLevel.
				Str("path", c.Request().URL.Path).
				Int("status", parseErrorCode(err)).
				Dur("latency", time.Since(start)).
				Msgf("Response %s", requestID)

			return err
		}
	}
}

func parseErrorCode(err error) int {
	if err == nil {
		return 200
	}
	errStr := err.Error()
	re := regexp.MustCompile(`code=(\d+)`)
	matches := re.FindStringSubmatch(errStr)

	if len(matches) < 2 {
		return 500
	}

	code, err := strconv.Atoi(matches[1])
	if err != nil {
		return 500
	}

	return code
}

func generateShortID() string {
	return randomdata.Noun()
}
