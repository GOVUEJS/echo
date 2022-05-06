package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type article struct {
	gorm.Model
	Id      int `gorm:"primarykey"`
	Title   string
	Content string
}

func main() {
	host := flag.String("host", "", "host")
	user := flag.String("user", "", "user")
	password := flag.String("password", "", "password")
	dbname := flag.String("dbname", "", "dbname")
	port := flag.String("port", "", "port")
	flag.Parse()

	dsn := fmt.Sprintf(`host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Seoul`, *host, *user, *password, *dbname, *port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Db 연결에 실패하였습니다. ")
	}

	// 테이블 자동 생성
	if err := db.AutoMigrate(&article{}); err != nil {
		return
	}

	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n ",
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/articles", func(c echo.Context) error {
		var articles []article
		result := db.Find(&articles)
		if result.RowsAffected == 0 {
			return c.String(http.StatusOK, "No articles")
		}

		marshal, err := json.Marshal(articles)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Server Error")
		}
		return c.String(http.StatusOK, string(marshal))
	})

	e.POST("/articles", func(c echo.Context) error {
		article := new(article)
		if err = c.Bind(article); err != nil {
			return c.String(http.StatusBadRequest, "Wrong Parameters")
		}

		// 생성
		db.Create(&article)

		return c.String(http.StatusOK, "POST Success")
	})

	e.GET("/articles/:id", func(c echo.Context) error {
		id := c.Param("id")

		idInt, err := strconv.Atoi(id)
		if err != nil {
			return c.String(http.StatusBadRequest, "Wrong Id")
		}

		article := article{Id: idInt}

		// 읽기
		db.First(&article, id) // primary key기준으로 article 찾기

		marshal, err := json.Marshal(article)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Server Error")
		}
		return c.String(http.StatusOK, string(marshal))
	})

	e.PUT("/articles/:id", func(c echo.Context) error {
		id := c.Param("id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return c.String(http.StatusBadRequest, "Wrong Id")
		}

		articleData := new(article)
		if err = c.Bind(articleData); err != nil {
			return c.String(http.StatusBadRequest, "Wrong Parameters")
		}
		articleData.Id = idInt

		// 수정 - product의 price를 200으로
		db.Model(&article{Id: articleData.Id}).Updates(articleData)

		return c.String(http.StatusOK, "PUT Success")
	})

	e.DELETE("/articles/:id", func(c echo.Context) error {
		id := c.Param("id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return c.String(http.StatusBadRequest, "Wrong Id")
		}

		// 삭제 - articleData 삭제하기
		d := db.Delete(&article{}, idInt)
		_ = d

		return c.String(http.StatusOK, "DELETE Success")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
