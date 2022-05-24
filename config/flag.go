package config

import (
	"flag"
	"fmt"
	"os"
)

var (
	FilePath *string
)

func init() {
	FilePath = flag.String("configFilePath", "", "configFilePath")
	flag.Parse()

	if "" == *FilePath {
		fmt.Println("Please enter the configFilePath flag")
		os.Exit(5)
	}
}
