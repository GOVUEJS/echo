package service

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
	"gorm.io/gorm"
	"myapp/database"
	"myapp/model"
	"myapp/util"
	"net/http"
	"strconv"
	"time"
)

var (
	rdb    *gorm.DB
	redis  *database.Redis
	jwtKey []byte
)

func init() {
	rdb = database.GetRDB()
	redis = database.GetRedis()
	jwtKey = []byte(random.String(32))
}

func GetMain(c echo.Context) error {
	redis.Set("test", "test", 0)
	value, err := redis.Get("test")
	if err != nil {
		return util.Response(c, http.StatusInternalServerError, "Redis Error", nil)
	}

	return util.Response(c, http.StatusOK, value, nil)
}

func PostLogin(c echo.Context) error {
	user := new(model.User)
	if err := c.Bind(user); err != nil {
		return err
	}

	if !database.Login(user.Email, user.Pw) {
		return echo.ErrUnauthorized
	}

	// Set custom claims
	claims := &model.JwtCustomClaims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
	}

	cookie := new(http.Cookie)
	cookie.Name = "accessToken"
	cookie.Value = user.Email
	cookie.Expires = time.Now().Add(time.Hour)
	c.SetCookie(cookie)

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString(jwtKey)
	if err != nil {
		return err
	}

	return util.Response(c, http.StatusOK, "", echo.Map{"accessToken": t})
}

func GetLogout(c echo.Context) error {

	cookie := new(http.Cookie)
	cookie.Name = "logout"
	cookie.Value = "test"
	cookie.Expires = time.Now().Add(1 * time.Hour)
	c.SetCookie(cookie)

	return util.ResponseNoContent(c, http.StatusOK)
}

func GetArticleList(c echo.Context) error {
	pageParam := c.QueryParam("page")
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		page = 1
	}

	response := &model.GetArticleListResponse{}

	sql := rdb.
		Model(&model.Article{}).
		Select([]string{"id", "title", "TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI') date"}).
		Order("id desc")

	current, total := util.GetPagination(sql, page)
	response.Current = current
	response.TotalPage = total

	offset, limit := util.GetPageOffsetLimit(page)
	sql.Limit(limit).
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
