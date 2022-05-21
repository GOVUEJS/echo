package model

import (
	"gorm.io/gorm"
	"time"
)

type Article struct {
	Id        uint           `gorm:"primaryKey;comment:글번호" json:"id"`
	Title     string         `gorm:"type:varchar(100);not null;comment:제목" json:"title"`
	Content   string         `gorm:"not null;comment:내용" json:"content"`
	CreatedAt time.Time      `gorm:"not null;comment:생성시간" json:"-"`
	UpdatedAt time.Time      `gorm:"not null;comment:수정시간" json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:삭제시간" json:"-"`
}

type User struct {
	Email string `gorm:"type:varchar(50);primaryKey" json:"email"`
	Pw    string `gorm:"type:varchar(256);not null;comment:비밀번호" json:"pw"`
}
