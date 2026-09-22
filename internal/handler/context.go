package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hunafazaky/event-booking-api/internal/apperror"
	"github.com/hunafazaky/event-booking-api/internal/model"
)

// getUserID reads the authenticated user's ID out of the Gin context.
// RequireAuth (middleware) is what puts it there, always as a uint —
// so no type-switch is needed here. If either check fails, it means
// this handler is mounted on a route that isn't behind RequireAuth,
// or the middleware itself changed shape — either way, that's a real
// bug worth surfacing, not something to silently paper over.
func getUserID(c *gin.Context) (uint, error) {
	v, exists := c.Get("user_id")
	if !exists {
		return 0, apperror.Unauthorized("unauthenticated")
	}

	userID, ok := v.(uint)
	if !ok {
		return 0, apperror.Internal("invalid user_id type in context", nil)
	}

	return userID, nil
}

// getOptionalViewer reads the same context keys as getUserID/RequireAuth,
// but for routes behind OptionalAuth instead — where having no
// authenticated viewer is a normal case, not an error. Returns
// (0, "") for an anonymous viewer rather than failing.
func getOptionalViewer(c *gin.Context) (userID uint, role model.Role) {
	if v, exists := c.Get("user_id"); exists {
		if id, ok := v.(uint); ok {
			userID = id
		}
	}
	if v, exists := c.Get("role"); exists {
		if r, ok := v.(model.Role); ok {
			role = r
		}
	}
	return userID, role
}
