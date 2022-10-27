package visualization

import (
	"dm.net/datamine/common/web"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func AppLaunch() {

	app := web.NewEchoAppWith([]echo.MiddlewareFunc{middleware.Logger()},
		[]*web.MapableHandlerFunc{&web.MapableHandlerFunc{
			Path:        "/about",
			MethodConst: echo.GET,
			Middlewares: nil,

			HandleFunc: about,
		}})

	app.Static("/html", "visualization/view/html")
	app.Static("/img", "visualization/view/static/image")
	app.Static("/css", "visualization/view/static/css")
	app.Static("/js", "visualization/view/static/js")
	app.Start(":8080")
}
