package test

import (
	"flag"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"myapp/config"
	"myapp/database/rdb/postgres"
	"myapp/database/session"
	"myapp/logger"
	"myapp/service"
)

var (
	filePath = flag.String("configFilePath", "", "configFilePath")
)

func NewEchoForTest() *echo.Echo {
	logger.InitLogger()

	err := config.InitConfig(*filePath)
	if err != nil {
		log.Fatalf("Decode toml error: %s", err)
		panic(err)
	}

	err = session.InitRedis()
	if err != nil {
		log.Fatalf("InitRedis error: %s", err)
		panic(err)
	}

	err = postgres.InitRDB()
	if err != nil {
		log.Fatalf("InitPostgres error: %s", err)
		panic(err)
	}

	service.InitService()

	return echo.New()
}
