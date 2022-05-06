package main

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	Id      int `gorm:"primarykey"`
	Title   string
	Content string
}
