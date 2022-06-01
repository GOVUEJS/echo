package server

import (
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/random"
	"myapp/logger"
)

func cookieForSessionMiddleware() echo.MiddlewareFunc {
	return session.Middleware(sessions.NewCookieStore([]byte("secret")))
}

func sessionMiddleware() func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, _ := session.Get("session", c)
			sess.Options = &sessions.Options{
				Path:     "/",
				MaxAge:   int(time.Hour),
				HttpOnly: true,
			}
			sess.Values["foo"] = "bar"
			sess.Save(c.Request(), c.Response())
			return c.NoContent(http.StatusOK)
		}
	}
}

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
		Output:           *logger.GetLogger(),
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
