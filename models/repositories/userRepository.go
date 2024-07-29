package repositories

import (
	"errors"
	"github.com/jinzhu/gorm"
	"go-rest-api/internal/config"
	"go-rest-api/models"
)

var UserRepository User

type User struct{}

func (u *User) CreateUser(user *models.User) error {
	return config.DB.Create(&user).Error
}

func (u *User) GetAllUsers() ([]models.User, error) {
	var user []models.User
	err := config.DB.Find(&user).Error

	for i := range user {
		user[i].Password = ""
	}
	return user, err
}

func (u *User) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := config.DB.Preload("Role.Permissions").Where("username = ?", username).Find(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (u *User) GetUserById(id uint) (*models.User, error) {
	var user models.User
	err := config.DB.Preload("Role.Permissions").Where("id = ?", id).Find(&user).Error
	return &user, err
}
