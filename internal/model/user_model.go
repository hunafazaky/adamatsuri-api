package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"-"`
	Role     Role   `json:"role" gorm:"type:varchar(20);not null;default:'attendee'"`
	// Interests reuses the same Tag pool events are tagged with —
	// "favorite series", not a separate vocabulary. A user's
	// interests and an event's tags are just two different many2many
	// relations onto the one Tag table.
	Interests []Tag `json:"interests" gorm:"many2many:user_interests;"`
	Events    []Event
}
