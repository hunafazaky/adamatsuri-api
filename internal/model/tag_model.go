package model

import "gorm.io/gorm"

// Tag is a free-form fandom/genre label attached to events — e.g.
// "Jujutsu Kaisen", "shounen". Unlike Category (a fixed, small enum
// of event KINDS), tags are open-ended and user-created: organizers
// type a tag name when creating an event, and TagRepository finds an
// existing tag with that name or creates a new one on the fly.
type Tag struct {
	gorm.Model
	Name   string  `json:"name" gorm:"unique;not null"`
	Events []Event `json:"-" gorm:"many2many:event_tags;"`
}
