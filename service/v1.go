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

	result := rdb.
		Model(&model.Article{}).
		Select([]string{"id", "title", "TO_CHAR(updated_at, 'YYYY-MM-DD HH24:MI') date"}).
		Order("id desc").
		Find(&response.ArticleList)

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
	rdb.
		Model(&model.Article{}).
		Select([]string{"id", "title", "content", "TO_CHAR(updated_at, 'YYYY-MM-DD HH24:MI') date"}).
		First(&response.Article, id) // primary key기준으로 Article 찾기

	return util.Response(c, http.StatusOK, "", response)
}

func PostArticle(c echo.Context) error {
	article := new(model.Article)
	if err := c.Bind(article); err != nil {
		return util.Response(c, http.StatusBadRequest, "Wrong Parameters", nil)
	}

	// 생성
	rdb.Create(&article)

	return util.Response(c, http.StatusOK, "POST Success", nil)
}

func PutArticle(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return util.Response(c, http.StatusBadRequest, "Wrong Id", nil)
	}

	articleData := new(model.Article)
	if err = c.Bind(articleData); err != nil {
		return util.Response(c, http.StatusBadRequest, "Wrong Parameters", nil)
	}
	articleData.Id = idInt

	// 수정 - product의 price를 200으로
	rdb.Model(&model.Article{Id: articleData.Id}).Updates(articleData)

	return util.Response(c, http.StatusOK, "PUT Success", nil)
}

func DeleteArticle(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return util.Response(c, http.StatusBadRequest, "Wrong Id", nil)
	}

	// 삭제 - articleData 삭제하기
	tx := rdb.Delete(&model.Article{}, idInt)
	if tx.RowsAffected == 0 {
		return util.Response(c, http.StatusBadRequest, "Id not found", nil)
	}

	return util.Response(c, http.StatusOK, "DELETE Success", nil)
}
