package database

import (
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"myapp/model"
	"myapp/util"
)

var (
	rdb *gorm.DB
)

func InitRDB() {
	dsn := fmt.Sprintf(`host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Seoul`, *util.Host, *util.User, *util.Password, *util.DbName, *util.Port)

	var err error
	rdb, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(errors.New("db 연결에 실패하였습니다"))
	}

	err = autoMigrate()
	if err != nil {
		panic(err)
	}
}

func autoMigrate() error {
	if err := rdb.AutoMigrate(&model.Article{}); err != nil {
		return errors.New("rdb article auto migrate 실패")
	}
	if err := rdb.AutoMigrate(&model.User{}); err != nil {
		return errors.New("rdb user auto migrate 실패")
	}
	return nil
}

func GetRDB() *gorm.DB {
	if rdb == nil {
		InitRDB()
	}
	return rdb
}
