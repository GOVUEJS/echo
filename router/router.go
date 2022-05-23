package router

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"myapp/config"
	"myapp/service"
)

func InitRouter(e *echo.Echo) {
	e.GET("/", service.GetMain)

	apiV1Group := e.Group("/api/v1")
	apiV1Group.POST("/signup", service.PostSignUp)
	apiV1Group.POST("/login", service.PostLogin)
	apiV1Group.GET("/logout", service.GetLogout)
	apiV1Group.POST("/token/refresh", service.RefreshToken)

	articleGroup := apiV1Group.Group("/articles")
	articleGroup.GET("", service.GetArticleList)
	articleGroup.GET("/:id", service.GetArticle)
	articleGroup.POST("", service.PostArticle, jwtAuth())
	//articleGroup.PUT("/:id", service.PutArticle, jwtAuth())
	//articleGroup.DELETE("/articles/:id", service.DeleteArticle, jwtAuth())
}

func jwtAuth() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		ParseTokenFunc: func(auth string, c echo.Context) (interface{}, error) {
			keyFunc := func(t *jwt.Token) (interface{}, error) {
				if t.Method.Alg() != "HS256" {
					return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
				}
				return config.Config.Jwt.Key, nil
			}

			// claims are of type `jwt.MapClaims` when token is created with `jwt.Parse`
			token, err := jwt.Parse(auth, keyFunc)
			if err != nil {
				return nil, err
			}
			if !token.Valid {
				return nil, errors.New("invalid token")
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				c.Set("sessionId", claims["sessionId"])
				c.Set("email", claims["email"])
			}
			return token, nil
		},
	})
}
