package restful

import (
	"github.com/labstack/echo/v4"
)

func NewEchoApp() *echo.Echo {
	return echo.New()
}
