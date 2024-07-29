package repositories

import (
	"errors"
	"github.com/jinzhu/gorm"
	"go-rest-api/internal/config"
	"go-rest-api/models"
)

var PermissionRepository Permission

type Permission struct{}

func (p *Permission) CreatePermission(permission *models.Permission) error {
	return config.DB.Create(&permission).Error
}

func (p *Permission) FindByName(name string) (*models.Permission, error) {
	var permission models.Permission
	err := config.DB.Where("name = ?", name).First(&permission).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &permission, nil
}

func (p *Permission) FindByRouteAndMethod(route string, method string) (*models.Permission, error) {
	var permission models.Permission
	err := config.DB.Where("route = ? and method = ?", route, method).First(&permission).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &permission, nil
}
