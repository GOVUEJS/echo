package service

import (
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
	"gorm.io/gorm"
	"myapp/database/rdb/postgres"
	"myapp/database/session"
	v1 "myapp/service/v1"
	"myapp/util"
)

var (
	rdb         *gorm.DB
	redisClient *redis.Client
)

func InitService() {
	rdb = postgres.GetRDB()
	redisClient = session.GetRedis()

	v1.InitService()
}

// GetMain
// @Summary Redis Test API
// @Description 레디스 테스트 API
// @Router / [GET]
// @Accept json
// @Produce json
// @Success 200 {object} model.ApiResponse
func GetMain(c echo.Context) error {
	redisClient.Set(*session.GetRedisContext(), "test", random.String(10), 0)
	value := redisClient.Get(*session.GetRedisContext(), "test")
	str := value.String()
	return util.Response(c, http.StatusOK, str, nil)
}
