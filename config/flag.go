package config

import (
	"flag"
	"fmt"
	"os"
)

var (
	Host      *string
	User      *string
	Password  *string
	RdbName   *string
	RdbPort   *string
	RedisPort *string
	FilePath  *string
)

func init() {
	Host = flag.String("host", "", "host")
	User = flag.String("user", "", "user")
	Password = flag.String("password", "", "password")
	RdbName = flag.String("dbname", "", "dbname")
	RdbPort = flag.String("port", "", "port")
	RedisPort = flag.String("redisPort", "", "redisPort")
	FilePath = flag.String("configFilePath", "", "configFilePath")
	flag.Parse()

	if "" == *FilePath {
		fmt.Println("Please enter the configFilePath flag")
		os.Exit(5)
	}
}
