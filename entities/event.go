package entities

import (
	"time"

	"gorm.io/gorm"
)

type CalendarEvent struct {
	gorm.Model
	Title    string    `gorm:"not null"`
	DateTime time.Time `gorm:"not null"`
}

func NewCalendarEvent(Title string, DateTime time.Time) *CalendarEvent {
	return &CalendarEvent{Title: Title, DateTime: DateTime}
}
