package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"html"
	"strings"
)

type User struct {
	gorm.Model
	Username string `gorm:"index:idx_username;unique;not null" json:"username"`
	Password string `json:"password"`
	RoleID   uint   `json:"role_id"`
	Role     Role   `gorm:"foreignKey:roleID" json:"role"`
}

func (u *User) BeforeCreate() error {

	//turn password into hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	//remove spaces in username
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

	return nil

}

const UserToken = "token:user:"
