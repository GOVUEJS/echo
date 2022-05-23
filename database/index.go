package database

import (
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"myapp/config"
	"myapp/model"
)

var (
	rdb *gorm.DB
)

func InitRDB() error {
	rdbConfig := config.Config.Rdb
	dsn := fmt.Sprintf(`host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Seoul`, rdbConfig.Ip, rdbConfig.User, rdbConfig.Password, rdbConfig.DbName, rdbConfig.Port)

	var err error
	rdb, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return errors.New("db 연결에 실패하였습니다")
	}

	if rdbConfig.AutoMigration {
		err = autoMigrate()
		if err != nil {
			return err
		}
	}

	return nil
}

func autoMigrate() error {
	if err := rdb.AutoMigrate(&model.Article{}, &model.User{}); err != nil {
		return errors.New("rdb auto migrate 실패")
	}
	return nil
}

func GetRDB() *gorm.DB {
	if rdb == nil {
		InitRDB()
	}
	return rdb
}
