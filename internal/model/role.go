package model

// Role is a User's permission level. Stored as plain text (not a DB
// enum) so adding a new role later is a Go-side change only — no
// migration needed to alter a Postgres enum type.
type Role string

const (
	// RoleAttendee is the default for every new sign-up: can browse
	// events and create bookings, but cannot create events.
	RoleAttendee Role = "attendee"

	// RoleOrganizer can create, update, and delete their OWN events,
	// in addition to everything an attendee can do.
	RoleOrganizer Role = "organizer"

	// RoleAdmin is reserved for future use (e.g. managing tags,
	// moderating events across organizers). Not assignable via any
	// current API endpoint — there is no signup path or admin-promotion
	// endpoint yet, so this exists only so authorization checks have a
	// name to compare against once one is added.
	RoleAdmin Role = "admin"
)

// IsValidSignupRole reports whether r is a role a new user is allowed to
// request at sign-up. Admin is deliberately excluded.
func IsValidSignupRole(r Role) bool {
	return r == RoleAttendee || r == RoleOrganizer
}
