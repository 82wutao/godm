package visualization

import "github.com/labstack/echo/v4"

func about(app echo.Context) error {
	var about = struct {
		Name    string
		Version string
	}{
		Name:    "visualization",
		Version: "0.0.1",
	}
	return app.JSON(200, about)
}
