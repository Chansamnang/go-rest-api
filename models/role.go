package models

import "github.com/jinzhu/gorm"

type Role struct {
	gorm.Model
	Name        string       `gorm:"unique" json:"name"`
	Description string       `gorm:"type:varchar(255)" json:"description"`
	Status      bool         `gorm:"default:true" json:"status"`
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}
