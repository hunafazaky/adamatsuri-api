package dto

import (
	"time"

	"github.com/hunafazaky/adamatsuri-api/internal/model"
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

// EventDetailResponse extends EventResponse with attendee info. Used
// only for GET /events/:id, the one query that actually preloads
// Booking + Booking.User.
//
// AttendeeCount is always accurate and visible to anyone. Bookings —
// the actual list, with each attendee's phone number and booking
// code — is only populated for the event's organizer (or an admin);
// EventService.GetByID leaves it nil/empty for every other viewer,
// including an anonymous one. Getting this backwards would leak every
// attendee's phone number to the public on a route that requires no
// authentication at all.
type EventDetailResponse struct {
	EventResponse
	AttendeeCount int                      `json:"attendee_count"`
	Bookings      []BookingSummaryResponse `json:"bookings"`
}

// EventListMeta carries pagination info alongside a list of events.
type EventListMeta struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	TotalRows int64 `json:"total_rows"`
	TotalPage int64 `json:"total_page"`
}
