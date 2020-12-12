package web

import "github.com/labstack/echo/v4"

type LoadableHandlerFunc interface {
	LoadedByEcho(echoApp *echo.Echo)
}
