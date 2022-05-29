package postgres

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

func SignUp(user *model.User) (err error) {
	tx := rdb.Create(&user)
	return tx.Error
}

func IsEmailAvailable(email *string) (result bool) {
	user := &model.User{Email: *email}

	tx := rdb.First(user)
	if tx.Error != nil {
		return false
	}

	if tx.RowsAffected == 1 {
		return false
	}
	return true
}
