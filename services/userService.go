package services

import (
	"go-rest-api/models"
	"go-rest-api/models/repositories"
)

func CreateUser(user *models.User) error {
	return repositories.UserRepository.CreateUser(user)
}

func GetAllUser() ([]models.User, error) {
	return repositories.UserRepository.GetAllUsers()
}

func GetUserByUsername(username string) (*models.User, error) {
	return repositories.UserRepository.GetUserByUsername(username)
}

func GetUserById(id uint) (*models.User, error) {
	return repositories.UserRepository.GetUserById(id)
}
