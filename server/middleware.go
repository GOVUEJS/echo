package server

import (
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/random"
	"myapp/logger"
)

func requestIDMiddleware() echo.MiddlewareFunc {
	return middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return random.String(32)
		},
	})
}

func loggerMiddleware() echo.MiddlewareFunc {
	return middleware.LoggerWithConfig(middleware.LoggerConfig{
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
		Output:           io.MultiWriter(os.Stdout, *logger.GetLogger()),
	})
}

func CORSMiddleware() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Set-Cookie"},
	})
}
