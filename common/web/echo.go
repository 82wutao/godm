package web

import (
	"github.com/labstack/echo/v4"
)

// NewEchoAppWith 构建一个新的echo应用
func NewEchoAppWith(middlewares []echo.MiddlewareFunc, handlers []*MapableHandlerFunc) *echo.Echo {
	e := echo.New()

	e.Use(middlewares...)

	for _, mapable := range handlers {
		e.Add(mapable.MethodConst, mapable.Path, mapable.HandleFunc, mapable.Middlewares...)
	}

	return e
}
