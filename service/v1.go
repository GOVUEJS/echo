package service

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
	"gorm.io/gorm"
	"myapp/database"
	"myapp/model"
	"myapp/util"
	"net/http"
	"strconv"
)

var (
	rdb          *gorm.DB
	redisClient  *redis.Client
	redisContext *context.Context
)

func init() {
	rdb = database.GetRDB()
	redisClient = database.GetRedis()
	redisContext = database.GetRedisContext()
}

func GetMain(c echo.Context) error {
	redisClient.Set(*redisContext, "test", random.String(10), 0)
	value := redisClient.Get(*redisContext, "test")
	str := value.String()
	return util.Response(c, http.StatusOK, str, nil)
}

func PostSignUp(c echo.Context) error {
	user := new(model.User)
	if err := c.Bind(user); err != nil {
		return err
	}

	rdb.Create(&user)

	return util.Response(c, http.StatusCreated, "", nil)
}

func PostLogin(c echo.Context) error {
	user := new(model.User)
	if err := c.Bind(user); err != nil {
		return err
	}

	if !database.Login(user.Email, user.Pw) {
		return echo.ErrUnauthorized
	}

	sessionId := uuid.New().String()

	accessToken, refreshToken, err := util.GetAccessRefreshToken(&user.Email, &sessionId)
	if err != nil {
		return err
	}

	ip := c.RealIP()
	redisSession := model.RedisSession{
		Email:        &user.Email,
		Ip:           &ip,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	err = database.SetRedisSession(sessionId, &redisSession)
	if err != nil {
		return err
	}

	return util.Response(c, http.StatusOK, "", map[string]interface{}{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func GetLogout(c echo.Context) error {
	return util.Response(c, http.StatusOK, "", nil)
}

func RefreshToken(c echo.Context) error {
	// TODO
	return util.Response(c, http.StatusOK, "", map[string]interface{}{})
}

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
		Select([]string{"id", "title", "TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI') date"}).
		Order("id desc").
		Limit(limit).
		Offset(offset).
		Find(&response.ArticleList)

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
		First(&response.Article, id) // primary key 기준으로 Article 찾기

	return util.Response(c, http.StatusOK, "", response)
}

func PostArticle(c echo.Context) error {
	article := new(model.Article)
	if err := c.Bind(article); err != nil {
		return util.Response(c, http.StatusBadRequest, "Wrong Parameters", nil)
	}

	// 생성
	rdb.Create(&article)

	return util.Response(c, http.StatusCreated, "POST Success", nil)
}

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
		return util.Response(c, http.StatusNotFound, "Id not found", nil)
	}

	return util.ResponseNoContent(c, http.StatusNoContent)
}
