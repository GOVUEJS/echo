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
	"time"
)

var (
	e *echo.Echo
)

func init() {
	e = echo.New()

	e.IPExtractor = echo.ExtractIPDirect()

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

	go func() {
		for {
			time.Sleep(time.Hour * 24)
			err := writer.Rotate()
			if err != nil {
				panic(err)
			}
		}
	}()

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

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		ExposeHeaders:    []string{"Set-Cookie"},
		AllowCredentials: true,
	}))

	router.InitRouter(e)
}

func InitEcho() {
	e.Logger.Fatal(e.Start(":1323"))
}
