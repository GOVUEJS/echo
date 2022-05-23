package config

import (
	"github.com/BurntSushi/toml"
)

type config struct {
	Phase string
	Jwt   jwt
	Echo  echoConfig
	Rdb   rdbConfig
	Redis redisConfig
}

type jwt struct {
	Key       []byte `toml:"-"`
	keyString string `toml:"key"`
}

type echoConfig struct {
	Port uint
}

type rdbConfig struct {
	Ip            string
	Port          uint
	User          string
	Password      string
	DbName        string
	AutoMigration bool
}

type redisConfig struct {
	Ip       string
	Port     uint
	Password string
}

var (
	Config config
)

func InitConfig() error {
	_, err := toml.DecodeFile(*FilePath, &Config)
	Config.Jwt.Key = []byte(Config.Jwt.keyString)
	if err != nil {
		return err
	}

	return nil
}
