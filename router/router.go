package router

import (
	"github.com/labstack/echo/v4"
	"myapp/service"
)

func InitRouter(e *echo.Echo) {
	e.GET("/", service.GetMain)

	apiV1Group := e.Group("/api/v1")
	apiV1Group.POST("/login", service.PostLogin)
	apiV1Group.GET("/logout", service.GetLogout)

	articleGroup := apiV1Group.Group("/articles")
	articleGroup.GET("", service.GetArticleList)
	articleGroup.POST("", service.PostArticle)
	articleGroup.GET("/:id", service.GetArticle)
	articleGroup.PUT("/:id", service.PutArticle)
	//articleGroup.DELETE("/articles/:id", service.DeleteArticle)
}
