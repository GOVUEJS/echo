package database

import (
	"myapp/model"
)

func Login(email, pw string) (result bool) {
	tx := rdb.Where(&model.User{Email: email, Pw: pw}).
		Find(&model.User{})

	if tx.RowsAffected == 1 {
		return true
	}
	return false
}
