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

	app.Static("/static", "view/static")
	app.Start(":8080")
}
