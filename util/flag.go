package util

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
)

func init() {
	Host = flag.String("host", "", "host")
	User = flag.String("user", "", "user")
	Password = flag.String("password", "", "password")
	RdbName = flag.String("dbname", "", "dbname")
	RdbPort = flag.String("port", "", "port")
	RedisPort = flag.String("redisPort", "", "redisPort")
	flag.Parse()

	if "" == *Host {
		fmt.Println("Please enter the host flag")
		os.Exit(1)
	}
	if "" == *User {
		fmt.Println("Please enter the user flag")
		os.Exit(2)
	}
	if "" == *Password {
		fmt.Println("Please enter the password flag")
		os.Exit(3)
	}
	if "" == *RdbName {
		fmt.Println("Please enter the dbname flag")
		os.Exit(4)
	}
	if "" == *RdbPort {
		fmt.Println("Please enter the port flag")
		os.Exit(5)
	}
	if "" == *RedisPort {
		fmt.Println("Please enter the redisPort flag")
		os.Exit(5)
	}
}
