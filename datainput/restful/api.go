package restful

import "github.com/labstack/echo/v4"

// About doc about this these api
func About(app echo.Context) error {
	var about = struct {
		Name    string
		Version string
	}{
		Name:    "Datainput Restful Api",
		Version: "0.0.1",
	}
	return app.JSON(200, about)
}
