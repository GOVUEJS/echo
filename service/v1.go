package service

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"myapp/database"
	"myapp/model"
	"myapp/util"
	"net/http"
	"strconv"
)

func GetMain(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func GetArticleList(c echo.Context) error {
	response := &model.GetArticleListResponse{}
	db := database.GetRDB()

	result := db.Order("id desc").Find(&response.ArticleList)
	if result.RowsAffected == 0 {
		return util.Response(c, http.StatusOK, "No articles", nil)
	}

	return util.Response(c, http.StatusOK, "", response)
}

func GetArticle(c echo.Context) error {
	db := database.GetRDB()

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
