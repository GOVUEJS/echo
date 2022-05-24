package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/random"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"myapp/router"
	"net/http"
	"os"
)

var (
	e *echo.Echo
)

func InitEcho() {
	e = echo.New()

	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return random.String(32)
		},
	}))

	writer := &lumberjack.Logger{
		Filename:   "/var/log/HWISECHO/echo.log",
		MaxSize:    500, // megabytes
		MaxBackups: 28,
		MaxAge:     28,    //days
		Compress:   false, // disabled by default
		LocalTime:  true,
	}

	multiWriter := io.MultiWriter(os.Stdout, writer)
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{` +
			`"time":"${time_rfc3339_nano}", ` +
			`"id":"${id}", ` +
			`"remote_ip":"${remote_ip}", ` +
			`"host":"${host}", ` +
			`"method":"${method}", ` +
			`"uri":"${uri}", ` +
			`"form":"${form}", ` +
			`"status":${status}, ` +
			`"error":"${error}"` +
			"}\n",
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
		Output:           multiWriter,
	}))

	e.IPExtractor = echo.ExtractIPFromXFFHeader()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Set-Cookie"},
	}))

	router.InitRouter(e)

	e.Logger.Fatal(e.Start(":1323"))
}
