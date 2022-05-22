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
	JwtKey []byte
)

func init() {
	rdb = database.GetRDB()
	redis = database.GetRedis()
	JwtKey = []byte(random.String(32))
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

	// Set custom accessTokenClaims
	accessTokenClaims := &model.JwtCustomClaims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
	}

	// Create token with accessTokenClaims
	accessTokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)

	// Generate encoded token and send it as response.
	accessToken, err := accessTokenJWT.SignedString(JwtKey)
	if err != nil {
		return err
	}

	c.SetCookie(&http.Cookie{
		Name:    "accessToken",
		Value:   accessToken,
		Path:    "/",
		Expires: time.Now().Add(1 * time.Hour),
	})

	// Set custom refreshTokenClaims
	refreshTokenClaims := &model.JwtCustomClaims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	// Create token with refreshTokenClaims
	refreshTokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	// Generate encoded token and send it as response.
	refreshToken, err := refreshTokenJWT.SignedString(JwtKey)
	if err != nil {
		return err
	}

	c.SetCookie(&http.Cookie{
		Name:   "refreshToken",
		Value:  refreshToken,
		Path:   "/",
		MaxAge: int(24 * time.Hour / time.Second),
	})

	return util.Response(c, http.StatusOK, "", map[string]interface{}{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func GetLogout(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:   "accessToken",
		Path:   "/",
		MaxAge: -1,
	})
	c.SetCookie(&http.Cookie{
		Name:   "refreshToken",
		Path:   "/",
		MaxAge: -1,
	})
	return util.Response(c, http.StatusOK, "", nil)
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
