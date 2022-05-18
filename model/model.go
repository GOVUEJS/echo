package model

import (
	"gorm.io/gorm"
	"time"
)

type ApiResponse struct {
	ErrCode *int    `json:"errCode,omitempty"`
	Message *string `json:"message,omitempty"`
}

type Article struct {
	Id        int            `gorm:"primaryKey" json:"id"`
	Title     string         `json:"title" json:"title"`
	Content   string         `json:"content"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"date"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type GetArticleListResponse struct {
	ArticleList []Article `json:"articleList"`
}

type GetArticleResponse struct {
	Article Article `json:"article"`
}
