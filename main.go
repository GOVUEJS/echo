package main

import (
	"myapp/config"
	"myapp/database"
	"myapp/server"
	"myapp/service"
)

func main() {
	err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	err = database.InitRedis()
	if err != nil {
		panic(err)
	}

	err = database.InitRDB()
	if err != nil {
		panic(err)
	}

	service.InitService()

	server.InitEcho()
}
