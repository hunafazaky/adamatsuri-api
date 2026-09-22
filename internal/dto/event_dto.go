package dto

import (
	"time"

	"github.com/hunafazaky/event-booking-app/internal/model"
)

// EventResponse is the shape of an event in LIST results
// (GET /events, GET /events/mine). No bookings — those queries never
// preload them, so this DTO doesn't promise data that isn't there.
type EventResponse struct {
	ID          uint           `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Location    string         `json:"location"`
	Image       string         `json:"image"`
	DateTime    time.Time      `json:"datetime"`
	Category    model.Category `json:"category"`
	Tags        []TagResponse  `json:"tags"`
	User        UserResponse   `json:"user"`
	CreatedAt   time.Time      `json:"created_at"`
}

// EventDetailResponse extends EventResponse with the list of bookings.
// Used only for GET /events/:id, the one query that actually preloads
// Booking + Booking.User.
type EventDetailResponse struct {
	EventResponse
	Bookings []BookingSummaryResponse `json:"bookings"`
}

// EventListMeta carries pagination info alongside a list of events.
type EventListMeta struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	TotalRows int64 `json:"total_rows"`
	TotalPage int64 `json:"total_page"`
}
