package main

import (
	"github.com/labstack/gommon/log"
	"myapp/config"
	"myapp/database"
	"myapp/logger"
	"myapp/server"
	"myapp/service"
)

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
