package database

import (
	"myapp/model"
)

func Login(email, pw string) (result bool) {
	var count int64

	rdb.Model(&model.User{}).
		Where(&model.User{Email: email, Pw: pw}).
		Count(&count)

	if count == 1 {
		return true
	}
	return false
}
