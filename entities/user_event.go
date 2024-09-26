package entities

import (
	"gorm.io/gorm"
)

type UserEvent struct {
	gorm.Model
	UserID          uint `gorm:"not null;<-:create"`
	CalendarEventID uint `gorm:"not null;<-:create"`
	User            User
	CalendarEvent   CalendarEvent
}

func NewUserEvent(UserID uint, CalendarEventID uint) *UserEvent {
	return &UserEvent{UserID: UserID, CalendarEventID: CalendarEventID}
}
