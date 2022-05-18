package service

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"myapp/database"
	"myapp/model"
	"myapp/util"
	"net/http"
	"strconv"
)

var (
	rdb *gorm.DB
)

func init() {
	rdb = database.GetRDB()
}

func GetMain(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func GetArticleList(c echo.Context) error {
	response := &model.GetArticleListResponse{}

	result := rdb.Order("id desc").Find(&response.ArticleList)
	if result.RowsAffected == 0 {
		return util.Response(c, http.StatusOK, "No articles", nil)
	}

	return util.Response(c, http.StatusOK, "", response)
}

func GetArticle(c echo.Context) error {
	id := c.Param("id")

	response := &model.GetArticleResponse{}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return util.Response(c, http.StatusBadRequest, "Wrong Id", nil)
	}

	response.Article.Id = idInt

	// 읽기
	rdb.First(&response.Article, id) // primary key기준으로 Article 찾기

	return util.Response(c, http.StatusOK, "", response)
}

func PostArticle(c echo.Context) error {
	db := database.GetRDB()

	article := new(model.Article)
	if err := c.Bind(article); err != nil {
		return c.String(http.StatusBadRequest, "Wrong Parameters")
	}

	// 생성
	db.Create(&article)

	return c.String(http.StatusOK, "POST Success")
}

func PutArticle(c echo.Context) error {
	db := database.GetRDB()

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
}

func DeleteArticle(c echo.Context) error {
	db := database.GetRDB()

	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.String(http.StatusBadRequest, "Wrong Id")
	}

	// 삭제 - articleData 삭제하기
	d := db.Delete(&model.Article{}, idInt)
	_ = d

	return c.String(http.StatusOK, "DELETE Success")
}
