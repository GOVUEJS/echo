package main

import (
	"myapp/database"
	"myapp/server"
)

func main() {
	database.InitRDB()

	server.InitEcho()
}
