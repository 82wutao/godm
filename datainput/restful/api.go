package restful

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// LoadRouters 加载路由
func LoadRouters(echoApp *echo.Echo) {
	echoApp.GET("/", hello)
}

// LoadMiddleWares 加载中间件
func LoadMiddleWares(echoApp *echo.Echo) {
	echoApp.Use(middleware.Logger())
	echoApp.Use(middleware.Recover())
}

func NewEchoApp() *echo.Echo {
	return echo.New()
}
