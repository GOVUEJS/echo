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
		Format: `{` +
			`"time":"${time_rfc3339_nano}"` +
			`"id":"${id}","remote_ip":"${remote_ip}",` +
			`"host":"${host}",` +
			`"method":"${method}",` +
			`"uri":"${uri}",` +
			`"form":"${form}",` +
			`"user_agent":"${user_agent}",` +
			`"status":${status},` +
			`"error":"${error}",` +
			`"latency":${latency},` +
			`"latency_human":"${latency_human}",` +
			`"bytes_in":${bytes_in},` +
			`"bytes_out":${bytes_out}` +
			"}\n",
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
	}))

	e.Use(middleware.CORS())
}

func InitEcho() {
	router.InitRouter(e)

	e.Logger.Fatal(e.Start(":1323"))
}
