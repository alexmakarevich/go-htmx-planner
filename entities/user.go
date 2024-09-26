package entities

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName string `gorm:"not null;unique"`
}

func NewUser(AccountName string) *User {
	return &User{UserName: AccountName}
}
