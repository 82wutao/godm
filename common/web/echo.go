package web

import "github.com/labstack/echo/v4"

func NewEchoAppWith(middlewares []echo.MiddlewareFunc, handlers []LoadableHandlerFunc) *echo.Echo {
	e := echo.New()

	e.Use(middlewares...)

	for _, loadable := range handlers {
		loadable.LoadedByEcho(e)
	}

	return e
}
