package model

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Name        string    `json:"name" binding:"required" `
	Description string    `json:"description" binding:"required" `
	Location    string    `json:"location" binding:"required" `
	Image       string    `json:"image"`
	ImageID     string    `json:"image_id"`
	UserID      uint      `json:"user_id"`
	User        User      `json:"user" gorm:"foreignKey:user_id"`
	DateTime    time.Time `json:"datetime" binding:"required" `
	Category    Category  `json:"category" gorm:"type:varchar(30);not null;default:'meetup';index"`
	Booking     []Booking `json:"booking_list" gorm:"foreignKey:event_id"`
}

var events []Event = []Event{}

func (e Event) Save() {
	events = append(events, e)
}

func GetAllEvents() []Event {
	return events
}
