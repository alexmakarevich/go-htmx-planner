package entities

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName string `gorm:"not null;unique"`
	Password string `gorm:"not null"` // TODO: IMPORTANT: OBVIOUSLY, HASH & SALT THIS
}

func NewUser(AccountName string, Password string) *User {
	return &User{UserName: AccountName, Password: Password}
}
