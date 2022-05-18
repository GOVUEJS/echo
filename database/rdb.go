package database

import (
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"myapp/model"
)

var rdb *gorm.DB

func InitRDB(host *string, user *string, password *string, dbname *string, port *string) error {
	dsn := fmt.Sprintf(`host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Seoul`, *host, *user, *password, *dbname, *port)

	var err error
	rdb, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return errors.New("db 연결에 실패하였습니다")
	}

	err = autoMigrate()
	if err != nil {
		return err
	}

	return err
}

func autoMigrate() error {
	if err := rdb.AutoMigrate(&model.Article{}); err != nil {
		return errors.New("rdb auto migrate 실패")
	}
	return nil
}

func GetRDB() *gorm.DB {
	return rdb
}
