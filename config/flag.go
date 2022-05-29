package config

import (
	"flag"
	"fmt"
	"os"
)

var (
	FilePath *string
)

func InitFlag() {
	FilePath = flag.String("configFilePath", "", "configFilePath")
	flag.Parse()

	if "" == *FilePath {
		fmt.Println("Please enter the configFilePath flag")
		os.Exit(1)
	}
}
