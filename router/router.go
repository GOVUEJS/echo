package router

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"myapp/service"
)

func InitRouter(e *echo.Echo) {
	e.GET("/", service.GetMain)

	apiV1Group := e.Group("/api/v1")
	apiV1Group.POST("/login", service.PostLogin)
	apiV1Group.GET("/logout", service.GetLogout)

	authGroup := apiV1Group.Group("")
	//authGroup.Use(getAuthWithJWT())

	articleGroup := authGroup.Group("/articles")
	articleGroup.GET("", service.GetArticleList)
	articleGroup.GET("/:id", service.GetArticle)
	articleGroup.POST("", service.PostArticle, jwtAuth())
	//articleGroup.PUT("/:id", service.PutArticle, jwtAuth())
	//articleGroup.DELETE("/articles/:id", service.DeleteArticle, jwtAuth())
}

func jwtAuth() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		TokenLookup: "query:token",
		ParseTokenFunc: func(auth string, c echo.Context) (interface{}, error) {
			keyFunc := func(t *jwt.Token) (interface{}, error) {
				if t.Method.Alg() != "HS256" {
					return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
				}
				return service.JwtKey, nil
			}

			// claims are of type `jwt.MapClaims` when token is created with `jwt.Parse`
			token, err := jwt.Parse(auth, keyFunc)
			if err != nil {
				return nil, err
			}
			if !token.Valid {
				return nil, errors.New("invalid token")
			}
			return token, nil
		},
	})
}
