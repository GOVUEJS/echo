package v1

import (
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"myapp/database/rdb/postgres"
	"myapp/database/session"
	"myapp/model"
	"myapp/util"
)

var (
	rdb         *gorm.DB
	redisClient *redis.Client
)

func InitService() {
	rdb = postgres.GetRDB()
	redisClient = session.GetRedis()
}

// PostSignUp
// @Summary Sign up
// @Description 회원가입
// @Router /signup [POST]
// @Param user body model.User true "회원가입 유저 정보"
// @Accept json
// @Produce json
// @Success 200 {object} model.ApiResponse
func PostSignUp(c echo.Context) error {
	user := new(model.User)
	if err := c.Bind(user); err != nil {
		return err
	}

	// TODO 해당 회원이 존재하는지 확인

	rdb.Create(&user)

	return util.Response(c, http.StatusCreated, "", nil)
}

// PostLogin
// @Summary Login
// @Description 로그인
// @Router /login [POST]
// @Param loginInfo body model.PostLoginRequest true "로그인 정보"
// @Accept json
// @Produce json
// @Success 200 {object} model.PostLoginResponse
// @Failure 400 {object} model.ApiResponse
// @Failure 401 {object} model.ApiResponse
func PostLogin(c echo.Context) error {
	user := new(model.User)
	if err := c.Bind(user); err != nil {
		return err
	}

	if user.Email == "" {
		return util.Response(c, http.StatusBadRequest, "email is empty", nil)
	}
	if user.Pw == "" {
		return util.Response(c, http.StatusBadRequest, "pw is empty", nil)
	}

	if !postgres.Login(user.Email, user.Pw) {
		return echo.ErrUnauthorized
	}

	sessionId := uuid.New().String()

	accessToken, refreshToken, err := util.GetAccessRefreshToken(&user.Email, &sessionId)
	if err != nil {
		return echo.ErrInternalServerError
	}

	ip := c.RealIP()
	redisSession := model.RedisSession{
		Email:        &user.Email,
		Ip:           &ip,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	err = session.SetRedisSession(sessionId, &redisSession)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return util.Response(c, http.StatusOK, "", model.PostLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// GetLogout
// @Summary Logout
// @Description 로그아웃
// @Router /logout [GET]
// @Accept json
// @Produce json
// @Success 200 {object} model.ApiResponse
func GetLogout(c echo.Context) error {
	// TODO 로그아웃 시 세션 삭제해야함
	return util.Response(c, http.StatusOK, "", nil)
}

// PostRefreshToken
// @Summary Refresh token
// @Description 토큰 재발급
// @Router /token/refresh [POST]
// @Param accessToken body string true "AccessToken"
// @Param refreshToken body string true "RefreshToken"
// @Accept json
// @Produce json
// @Success 200 {object} model.PostLoginResponse
// @Failure 401 {object} model.ApiResponse
func PostRefreshToken(c echo.Context) error {
	tokens := new(model.Tokens)
	if err := c.Bind(tokens); err != nil {
		return err
	}

	requestAccessTokenClaims, _, err := util.CheckRefreshToken(tokens)
	if err != nil {
		return err
	}

	email, ok := requestAccessTokenClaims["email"].(string)
	if !ok {
		return util.Response(c, http.StatusBadRequest, "jwt email error", nil)
	}
	sessionId, ok := requestAccessTokenClaims["sessionId"].(string)
	if !ok {
		return util.Response(c, http.StatusBadRequest, "jwt sessionId error", nil)
	}

	if tokens.AccessToken == nil || tokens.RefreshToken == nil {
		return util.Response(c, http.StatusBadRequest, "wrong parameter", nil)
	}

	accessToken, refreshToken, err := util.GetAccessRefreshToken(&email, &sessionId)
	if err != nil {
		return err
	}

	ip := c.RealIP()
	redisSession := model.RedisSession{
		Email:        &email,
		Ip:           &ip,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	err = session.SetRedisSession(sessionId, &redisSession)
	if err != nil {
		return err
	}

	return util.Response(c, http.StatusOK, "", map[string]interface{}{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}
