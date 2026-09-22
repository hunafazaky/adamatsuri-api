package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"-"`
	Role     Role   `json:"role" gorm:"type:varchar(20);not null;default:'attendee'"`
	Events   []Event
}
