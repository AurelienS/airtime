package webserver

import (
	"github.com/labstack/echo/v4"
)

type Router struct {
	Handler Handler
}

func (r *Router) Initialize(e *echo.Echo) {
	e.GET("/", r.Handler.GetIndex)
}
