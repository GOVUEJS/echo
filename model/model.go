package model

import (
	"time"
)

type Article struct {
	Id        int       `gorm:"primarykey" json:"id"`
	Title     string    `json:"title" json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"date"`
	DeletedAt time.Time `gorm:"index" json:"-"`
}
