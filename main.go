package main

import (
	"fmt"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

type article struct {
	gorm.Model
	Seq     int
	Title   string
	Content string
}

func main() {
	host := os.Args[1]
	user := os.Args[2]
	password := os.Args[3]
	dbname := os.Args[4]
	port := os.Args[5]

	dsn := fmt.Sprintf(`host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Seoul`, host, user, password, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Db 연결에 실패하였습니다. ")
	}

	// 테이블 자동 생성
	if err := db.AutoMigrate(&article{}); err != nil {
		return
	}

	// 생성
	article := article{Title: "D42", Content: "100"}
	db.Create(&article)

	// 읽기
	db.First(&article, 1)                 // primary key기준으로 article 찾기
	db.First(&article, "code = ?", "D42") // code가 D42인 article 찾기

	// 수정 - product의 price를 200으로
	//db.Model(&article).Update("Price", 200)
	//// 수정 - 여러개의 필드를 수정하기
	//db.Model(&article).Updates(article{Price: 200, Code: "F42"})
	//db.Model(&article).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})
	//
	//// 삭제 - article 삭제하기
	//db.Delete(&article, 1)

	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n ",
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
