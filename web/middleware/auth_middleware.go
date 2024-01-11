package middleware

import (
	"net/http"

	"github.com/AurelienS/cigare/internal/util"
	"github.com/AurelienS/cigare/web/session"
	"github.com/labstack/echo/v4"
)

const UserContextKey = "user"

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := session.GetUserOrErrorFromContext(c)
		if err != nil {
			util.Warn().Msg("user not logged in. Will be redirected soon")
			return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")
		}

		c.Set(UserContextKey, user)
		return next(c)
	}
}
