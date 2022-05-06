package model

import (
	"gorm.io/gorm"
	"time"
)

type Article struct {
	Id        int            `gorm:"primaryKey" json:"id"`
	Title     string         `json:"title" json:"title"`
	Content   string         `json:"content"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"date"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
