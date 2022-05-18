package main

import (
	"flag"
	"fmt"
	"myapp/database"
	"myapp/router"
	"os"
)

var (
	host     *string
	user     *string
	password *string
	dbname   *string
	port     *string
)

func main() {
	err := database.InitRDB(host, user, password, dbname, port)
	if err != nil {
		panic(err)
	}

	e := router.InitRouter()
	e.Logger.Fatal(e.Start(":1323"))
}

func init() {
	host = flag.String("host", "", "host")
	user = flag.String("user", "", "user")
	password = flag.String("password", "", "password")
	dbname = flag.String("dbname", "", "dbname")
	port = flag.String("port", "", "port")
	flag.Parse()

	if "" == *host {
		fmt.Println("Please enter the host flag")
		os.Exit(1)
	}
	if "" == *user {
		fmt.Println("Please enter the user flag")
		os.Exit(2)
	}
	if "" == *password {
		fmt.Println("Please enter the password flag")
		os.Exit(3)
	}
	if "" == *dbname {
		fmt.Println("Please enter the dbname flag")
		os.Exit(4)
	}
	if "" == *port {
		fmt.Println("Please enter the port flag")
		os.Exit(5)
	}
}
