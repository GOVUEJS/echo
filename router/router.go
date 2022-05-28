package router

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"myapp/service"
	"myapp/service/v1"
)

func InitRouter(e *echo.Echo) {
	e.GET("/", service.GetMain)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	apiV1Group := e.Group("/api/v1")
	apiV1Group.POST("/signup", v1.PostSignUp)
	apiV1Group.POST("/login", v1.PostLogin)
	apiV1Group.GET("/logout", v1.GetLogout)
	apiV1Group.POST("/token/refresh", v1.PostRefreshToken)

	articleGroup := apiV1Group.Group("/articles")
	articleGroup.GET("", v1.GetArticleList)
	articleGroup.GET("/:id", v1.GetArticle)
	articleGroup.POST("", v1.PostArticle, jwtAuth())
	articleGroup.PUT("/:id", v1.PutArticle, jwtAuth())
	articleGroup.DELETE("/:id", v1.DeleteArticle, jwtAuth())
}
