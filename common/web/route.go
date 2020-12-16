package web

import (
	"github.com/labstack/echo/v4"
)

// MapableHandlerFunc 一个业务的必需组织
type MapableHandlerFunc struct {
	Path        string
	MethodConst string
	Middlewares []echo.MiddlewareFunc

	HandleFunc echo.HandlerFunc
}
