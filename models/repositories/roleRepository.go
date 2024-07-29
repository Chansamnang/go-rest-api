package repositories

import (
	"errors"
	"github.com/jinzhu/gorm"
	"go-rest-api/internal/config"
	"go-rest-api/models"
)

var RoleRepository Role

type Role struct{}

func (r *Role) SaveRole(newRole *models.Role) error {
	return config.DB.Save(&newRole).Error
}

func (r *Role) FindByName(name string) (*models.Role, error) {
	var role models.Role
	err := config.DB.Preload("Permissions").Where("name = ?", name).First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &role, err
}
