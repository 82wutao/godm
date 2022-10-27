package datainput

import (
	"dm.net/datamine/common/web"
	"dm.net/datamine/datainput/restful"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func AppLaunch() {

	app := web.NewEchoAppWith([]echo.MiddlewareFunc{middleware.Logger()},
		[]*web.MapableHandlerFunc{&web.MapableHandlerFunc{
			Path:        "/about",
			MethodConst: echo.GET,
			Middlewares: nil,

			HandleFunc: restful.About,
		}})

	app.Static("/html", "visualization/view/html")
	app.Static("/img", "visualization/view/static/image")
	app.Static("/css", "visualization/view/static/css")
	app.Static("/js", "visualization/view/static/js")
	app.Start(":9090")
}
