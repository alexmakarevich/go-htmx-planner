package entities

import (
	"gorm.io/gorm"
)

type UserEvent struct {
	gorm.Model
	UserID          uint `gorm:"not null;<-:create"` // field should not be editable
	CalendarEventID uint `gorm:"not null;<-:create"` // field should not be editable
	User            User
	CalendarEvent   CalendarEvent
}

func NewUserEvent(UserID uint, CalendarEventID uint) *UserEvent {
	return &UserEvent{UserID: UserID, CalendarEventID: CalendarEventID}
}
