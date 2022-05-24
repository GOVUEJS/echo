package config

import (
	"github.com/BurntSushi/toml"
)

type config struct {
	Phase string      `toml:"phase"`
	Jwt   jwt         `toml:"jwt"`
	Echo  echoConfig  `toml:"echo"`
	Rdb   rdbConfig   `toml:"rdb"`
	Redis redisConfig `toml:"redis"`
}

type jwt struct {
	Key       []byte `toml:"-"`
	keyString string `toml:"key"`
}

type echoConfig struct {
	Port uint `toml:"port"`
}

type rdbConfig struct {
	Ip            string `toml:"ip"`
	Port          uint   `toml:"port"`
	User          string `toml:"user"`
	Password      string `toml:"password"`
	DbName        string `toml:"dbName"`
	AutoMigration bool   `toml:"auto_migration"`
}

type redisConfig struct {
	Ip       string `toml:"ip"`
	Port     uint   `toml:"port"`
	Password string `toml:"password"`
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
