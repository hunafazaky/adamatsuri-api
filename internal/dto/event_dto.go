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
//
// YourBooking is different: it's the VIEWER'S OWN booking on this
// event, if they have one — that's the viewer's own data regardless
// of who they are, so it's populated for any authenticated viewer,
// not just the organizer. This is what lets the frontend show
// "you're booked" instead of a booking form after the fact.
type EventDetailResponse struct {
	EventResponse
	AttendeeCount int                      `json:"attendee_count"`
	Bookings      []BookingSummaryResponse `json:"bookings"`
	YourBooking   *BookingSummaryResponse  `json:"your_booking"`
}

// EventListMeta carries pagination info alongside a list of events.
type EventListMeta struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	TotalRows int64 `json:"total_rows"`
	TotalPage int64 `json:"total_page"`
}
