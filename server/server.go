package server

import (
	"github.com/labstack/echo/v4"
	"myapp/router"
)

var (
	e *echo.Echo
)

func InitEcho() {
	e = echo.New()

	e.IPExtractor = echo.ExtractIPFromXFFHeader()

	e.Use(requestIDMiddleware())
	e.Use(loggerMiddleware())
	e.Use(CORSMiddleware())

	router.InitRouter(e)
	e.Logger.Fatal(e.Start(":1323"))
}
