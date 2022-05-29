package config

import (
	"errors"
	"flag"
)

var (
	FilePath *string
)

func InitFlag() (err error) {
	FilePath = flag.String("configFilePath", "", "configFilePath")
	flag.Parse()

	if "" == *FilePath {
		return errors.New("please enter the configFilePath flag")
	}
	return
}
