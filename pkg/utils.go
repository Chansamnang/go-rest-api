package pkg

import (
	"fmt"
	"go-rest-api/internal/config"
	"golang.org/x/crypto/bcrypt"
)

func CompareHashAndPassword(hashedPassword []byte, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	if err != nil {
		config.Logger.Error(fmt.Sprintf("incorrect password %s", err.Error()))
		return false
	}
	return true
}
