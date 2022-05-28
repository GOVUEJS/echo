package main

import (
	"github.com/labstack/gommon/log"
	"myapp/config"
	"myapp/database"
	_ "myapp/docs"
	"myapp/logger"
	"myapp/server"
	"myapp/service"
)

// @title Hwisaek's server
// @version 1.0
// @description This is a Hwisaek's server.
// @contact.name API Support
// @contact.email dia_changmin@naver.com
// @host 211.34.36.139:1323
// @BasePath /api/v1
func main() {
	logger.InitLogger()

	err := config.InitConfig()
	if err != nil {
		log.Fatalf("Decode toml error: %s", err)
		panic(err)
	}

	err = database.InitRedis()
	if err != nil {
		log.Fatalf("InitRedis error: %s", err)
		panic(err)
	}

	err = database.InitRDB()
	if err != nil {
		log.Fatalf("InitRDB error: %s", err)
		panic(err)
	}

	service.InitService()

	server.InitEcho()
}
