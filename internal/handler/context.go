package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hunafazaky/event-booking-api/internal/apperror"
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
