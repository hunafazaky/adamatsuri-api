package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hunafazaky/adamatsuri-api/internal/model"
	"github.com/hunafazaky/adamatsuri-api/internal/response"
)

func RequireAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			response.Fail(c, http.StatusUnauthorized, "invalid or missing token")
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || token == nil || !token.Valid {
			response.Fail(c, http.StatusUnauthorized, "invalid or missing token")
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if sub, ok := claims["sub"].(float64); ok {
				c.Set("user_id", uint(sub))

				// role defaults to attendee if the claim is missing —
				// e.g. a token signed before role support existed.
				// Such a user simply can't hit organizer-only routes
				// until they sign in again for a fresh token.
				role := model.RoleAttendee
				if r, ok := claims["role"].(string); ok && r != "" {
					role = model.Role(r)
				}
				c.Set("role", role)

				c.Next()
				return
			}
		}

		response.Fail(c, http.StatusUnauthorized, "invalid or missing token")
		c.Abort()
	}
}

// OptionalAuth is for routes that are public but behave DIFFERENTLY
// for a signed-in viewer — e.g. GET /events/:id, which shows full
// attendee info to the event's organizer but not to anyone else. If
// there's no token, or it's invalid/expired, this simply doesn't set
// "user_id"/"role" and moves on — unlike RequireAuth, it never aborts
// the request. Handlers on this route should check for "user_id"
// with the two-value c.Get form (or a dedicated helper), never the
// version that errors when it's absent.
func OptionalAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.Next()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			c.Next()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if sub, ok := claims["sub"].(float64); ok {
				c.Set("user_id", uint(sub))

				role := model.RoleAttendee
				if r, ok := claims["role"].(string); ok && r != "" {
					role = model.Role(r)
				}
				c.Set("role", role)
			}
		}

		c.Next()
	}
}

// RequireRole guards a route to only the listed roles. It must be
// mounted AFTER RequireAuth, which is what puts "role" into the
// context in the first place — using it standalone always denies,
// since there'd be no role to check.
func RequireRole(allowed ...model.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		v, exists := c.Get("role")
		if !exists {
			response.Fail(c, http.StatusForbidden, "you're not authorized")
			c.Abort()
			return
		}

		role, ok := v.(model.Role)
		if !ok {
			response.Fail(c, http.StatusForbidden, "you're not authorized")
			c.Abort()
			return
		}

		for _, a := range allowed {
			if role == a {
				c.Next()
				return
			}
		}

		response.Fail(c, http.StatusForbidden, "you're not authorized")
		c.Abort()
	}
}
