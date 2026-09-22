package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hunafazaky/event-booking-api/internal/apperror"
)

// Envelope is the standard shape every API response uses, success or failure.
// Keeping one shape means clients never have to guess which key holds the
// error message, and OpenAPI docs only need to describe it once.
type Envelope struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Success sends a successful response with the given status code and data.
func Success(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, Envelope{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SuccessWithMeta sends a successful response that also carries pagination
// or other metadata (page, limit, total_rows, etc).
func SuccessWithMeta(c *gin.Context, status int, message string, data interface{}, meta interface{}) {
	c.JSON(status, Envelope{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// Fail sends a plain error response with an explicit status code.
// Prefer FromError below once handlers call into the service layer —
// this is for validation failures caught right in the handler (e.g. bad JSON).
func Fail(c *gin.Context, status int, message string) {
	c.JSON(status, Envelope{
		Success: false,
		Error:   message,
	})
}

// FromError inspects err: if it's an *apperror.AppError, its Code and Message
// drive the response. Otherwise it's treated as an unexpected failure and
// mapped to 500, without leaking internal error details to the client.
func FromError(c *gin.Context, err error) {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.Code, Envelope{
			Success: false,
			Error:   appErr.Message,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, Envelope{
		Success: false,
		Error:   "internal server error",
	})
}
