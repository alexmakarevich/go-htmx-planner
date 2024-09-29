package entities

import (
	"time"

	"gorm.io/gorm"
)

type CalendarEvent struct {
	gorm.Model
	Title    string    `gorm:"not null"`
	DateTime time.Time `gorm:"not null"`
	OwnerId  uint      `gorm:"not null"`
	Owner    User      `gorm:"foreignKey:OwnerId"`
}

func NewCalendarEvent(Title string, DateTime time.Time, OwnerId uint) *CalendarEvent {
	return &CalendarEvent{Title: Title, DateTime: DateTime, OwnerId: OwnerId}
}
