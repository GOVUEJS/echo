package util

import (
	"flag"
	"fmt"
	"os"
)

var (
	Host     *string
	User     *string
	Password *string
	DbName   *string
	Port     *string
)

func init() {
	Host = flag.String("host", "", "host")
	User = flag.String("user", "", "user")
	Password = flag.String("password", "", "password")
	DbName = flag.String("dbname", "", "dbname")
	Port = flag.String("port", "", "port")
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
	if "" == *DbName {
		fmt.Println("Please enter the dbname flag")
		os.Exit(4)
	}
	if "" == *Port {
		fmt.Println("Please enter the port flag")
		os.Exit(5)
	}
}
