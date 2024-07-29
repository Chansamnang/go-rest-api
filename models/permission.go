package models

import "github.com/jinzhu/gorm"

type Permission struct {
	gorm.Model
	Name   string `gorm:"type:varchar(255);not null" json:"name"`
	Method string `gorm:"type:varchar(255);not null" json:"method"`
	Route  string `gorm:"type:varchar(255);not null" json:"route"`
}
