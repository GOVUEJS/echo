package router

import (
	"github.com/labstack/echo/v4"
	"myapp/service"
)

func InitRouter(e *echo.Echo) {
	e.GET("/", service.GetMain)

	apiV1Group := e.Group("/api/v1")
	apiV1Group.GET("/articles", service.GetArticleList)
	apiV1Group.POST("/articles", service.PostArticle)
	apiV1Group.GET("/articles/:id", service.GetArticle)
	apiV1Group.PUT("/articles/:id", service.PutArticle)
	//apiV1Group.DELETE("/articles/:id", service.DeleteArticle)
}
