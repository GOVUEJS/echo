package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"myapp/router"
)

var (
	e *echo.Echo
)

func init() {
	e = echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n ",
	}))
	e.Use(middleware.CORS())
}

func InitEcho() {
	router.InitRouter(e)

	e.Logger.Fatal(e.Start(":1323"))
}
