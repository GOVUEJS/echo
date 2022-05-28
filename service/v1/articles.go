package v1

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"myapp/model"
	"myapp/util"
)

// GetArticle
// @Summary Get article
// @Description 게시글 조회
// @Router /articles/{id} [GET]
// @Param id path uint false "게시글 ID"
// @Accept json
// @Produce json
// @Success 200 {object} model.GetArticleResponse
// @Failure 404 {object} model.ApiResponse
func GetArticle(c echo.Context) error {
	id := c.Param("id")

	response := &model.GetArticleResponse{}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return echo.ErrNotFound
	}

	response.Article.Id = idInt

	// 읽기
	rdb.
		Model(&model.Article{}).
		Select([]string{"id", "title", "content", "TO_CHAR(updated_at, 'YYYY-MM-DD HH24:MI') date", "writer"}).
		First(&response.Article, id) // primary key 기준으로 Article 찾기

	return util.Response(c, http.StatusOK, "", response)
}

// GetArticleList
// @Summary Get article list
// @Description 게시글 목록 조회
// @Router /articles [GET]
// @Param page query uint false "Page number"
// @Accept json
// @Produce json
// @Success 200 {object} model.GetArticleListResponse
func GetArticleList(c echo.Context) error {
	pageParam := c.QueryParam("page")
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		page = 1
	}

	response := &model.GetArticleListResponse{}

	sql := rdb.
		Model(&model.Article{})

	current, total := util.GetPagination(sql, page)
	response.Current = current
	response.TotalPage = total

	offset, limit := util.GetPageOffsetLimit(page)
	sql.
		Select([]string{"id", "title", "TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI') date", "writer"}).
		Order("id desc").
		Limit(limit).
		Offset(offset).
		Find(&response.ArticleList)

	return util.Response(c, http.StatusOK, "", response)
}

// PostArticle
// @Summary Create article
// @Description 게시글 작성
// @Router /articles [POST]
// @Param id path uint false "게시글 ID"
// @Accept json
// @Produce json
// @Success 201 {Object} model.ApiResponse
// @Failure 400 {object} model.ApiResponse
// @Failure 401 {object} model.ApiResponse
func PostArticle(c echo.Context) error {
	userEmail, ok := c.Get("email").(string)
	if !ok {
		return util.Response(c, http.StatusInternalServerError, "user", nil)
	}

	if userEmail == "" {
		return echo.ErrUnauthorized
	}

	article := &model.Article{Writer: userEmail}
	if err := c.Bind(article); err != nil {
		return echo.ErrBadRequest
	}

	if article.Title == "" {
		return util.Response(c, http.StatusBadRequest, "title is empty", nil)
	}
	if article.Content == "" {
		return util.Response(c, http.StatusBadRequest, "content is empty", nil)
	}

	// 생성
	rdb.Create(&article)

	return util.Response(c, http.StatusCreated, "POST Success", nil)
}

// PutArticle
// @Summary Update article
// @Description 게시글 수정
// @Router /articles/{id} [PUT]
// @Param id path string true "게시글 ID"
// @Param article body model.Article true "게시글 내용"
// @Accept json
// @Produce json
// @Success 200
// @Failure 400 {object} model.ApiResponse
func PutArticle(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil || idInt < 0 {
		return util.Response(c, http.StatusBadRequest, "Wrong Id", nil)
	}

	articleData := new(model.Article)
	if err = c.Bind(articleData); err != nil {
		return util.Response(c, http.StatusBadRequest, "Wrong Parameters", nil)
	}
	articleData.Id = uint(idInt)

	var writer string
	rdb.
		Model(&model.Article{}).
		Select([]string{"writer"}).
		First(&writer, id) // primary key 기준으로 Article 찾기
	if email := c.Get("email"); email != writer {
		return util.Response(c, http.StatusUnauthorized, "", nil)
	}

	rdb.Model(&model.Article{Id: articleData.Id}).Updates(articleData)

	return util.Response(c, http.StatusOK, "PUT Success", nil)
}

// DeleteArticle
// @Summary Delete article
// @Description 게시글 삭제
// @Router /articles/{id} [Delete]
// @Param id path string true "게시글 ID"
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {object} model.ApiResponse
// @Failure 401 {object} model.ApiResponse
// @Failure 404 {object} model.ApiResponse
func DeleteArticle(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return util.Response(c, http.StatusBadRequest, "Wrong Id", nil)
	}

	var writer string
	rdb.
		Model(&model.Article{}).
		Select([]string{"writer"}).
		First(&writer, id) // primary key 기준으로 Article 찾기
	if email := c.Get("email"); email != writer {
		return echo.ErrUnauthorized
	}

	// 삭제 - articleData 삭제하기
	tx := rdb.Delete(&model.Article{}, idInt)
	if tx.RowsAffected == 0 {
		return echo.ErrNotFound
	}

	return util.ResponseNoContent(c, http.StatusNoContent)
}
