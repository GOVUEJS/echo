package main

import (
	"myapp/database"
	"myapp/server"
)

func main() {
	database.InitRedis()

	database.InitRDB()

	server.InitEcho()
}
