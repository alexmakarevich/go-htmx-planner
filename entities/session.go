package entities

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	ID        string    `gorm:"primaryKey"`
	CreatedAt time.Time // TODO: IMPORTANT: TTL
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	UserID    uint           `gorm:"not null;<-:create"` // field should not be editable
	User      User
}

func NewSession(ID string, UserID uint) *Session {
	return &Session{ID: ID, UserID: UserID}
}
