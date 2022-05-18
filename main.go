package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"myapp/database"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"myapp/model"
)

func main() {
	host := flag.String("host", "", "host")
	user := flag.String("user", "", "user")
	password := flag.String("password", "", "password")
	dbname := flag.String("dbname", "", "dbname")
	port := flag.String("port", "", "port")
	flag.Parse()

	if "" == *host {
		fmt.Println("Please enter the host flag")
		os.Exit(1)
	}
	if "" == *user {
		fmt.Println("Please enter the user flag")
		os.Exit(2)
	}
	if "" == *password {
		fmt.Println("Please enter the password flag")
		os.Exit(3)
	}
	if "" == *dbname {
		fmt.Println("Please enter the dbname flag")
		os.Exit(4)
	}
	if "" == *port {
		fmt.Println("Please enter the port flag")
		os.Exit(5)
	}

	err := database.InitRDB(host, user, password, dbname, port)
	if err != nil {
		panic(err)
	}

	// 테이블 자동 생성
	db := database.GetRDB()
	if err := db.AutoMigrate(&model.Article{}); err != nil {
		return
	}

	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n ",
	}))

	e.Use(middleware.CORS())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	apiV1Group := e.Group("/api/v1")
	apiV1Group.GET("/articles", func(c echo.Context) error {
		var articles []model.Article
		result := db.Order("id desc").Find(&articles)
		if result.RowsAffected == 0 {
			return c.String(http.StatusOK, "No articles")
		}

		marshal, err := json.Marshal(articles)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Server Error")
		}
		return c.String(http.StatusOK, string(marshal))
	})

	apiV1Group.POST("/articles", func(c echo.Context) error {
		article := new(model.Article)
		if err = c.Bind(article); err != nil {
			return c.String(http.StatusBadRequest, "Wrong Parameters")
		}

		// 생성
		db.Create(&article)

		return c.String(http.StatusOK, "POST Success")
	})

	apiV1Group.GET("/articles/:id", func(c echo.Context) error {
		id := c.Param("id")

		idInt, err := strconv.Atoi(id)
		if err != nil {
			return c.String(http.StatusBadRequest, "Wrong Id")
		}

		article := model.Article{Id: idInt}

		// 읽기
		db.First(&article, id) // primary key기준으로 Article 찾기

		marshal, err := json.Marshal(article)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Server Error")
		}
		return c.String(http.StatusOK, string(marshal))
	})

	apiV1Group.PUT("/articles/:id", func(c echo.Context) error {
		id := c.Param("id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return c.String(http.StatusBadRequest, "Wrong Id")
		}

		articleData := new(model.Article)
		if err = c.Bind(articleData); err != nil {
			return c.String(http.StatusBadRequest, "Wrong Parameters")
		}
		articleData.Id = idInt

		// 수정 - product의 price를 200으로
		db.Model(&model.Article{Id: articleData.Id}).Updates(articleData)

		return c.String(http.StatusOK, "PUT Success")
	})

	apiV1Group.DELETE("/articles/:id", func(c echo.Context) error {
		id := c.Param("id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return c.String(http.StatusBadRequest, "Wrong Id")
		}

		// 삭제 - articleData 삭제하기
		d := db.Delete(&model.Article{}, idInt)
		_ = d

		return c.String(http.StatusOK, "DELETE Success")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
