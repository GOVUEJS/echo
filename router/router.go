package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"myapp/service"
)

func InitRouter() (e *echo.Echo) {
	e = echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n ",
	}))
	e.Use(middleware.CORS())

	e.GET("/", service.GetMain)

	apiV1Group := e.Group("/api/v1")
	apiV1Group.GET("/articles", service.GetArticleList)
	apiV1Group.POST("/articles", service.PostArticle)
	apiV1Group.GET("/articles/:id", service.GetArticle)
	apiV1Group.PUT("/articles/:id", service.PutArticle)
	apiV1Group.DELETE("/articles/:id", service.DeleteArticle)

	return
}
