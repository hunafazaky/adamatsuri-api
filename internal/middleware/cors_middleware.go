package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS lets a browser-based frontend on a different origin (e.g. the
// Vite dev server on http://localhost:5173) call this API at all —
// without it, the browser blocks the response from ever reaching the
// frontend's JavaScript, regardless of what status code the server
// actually returned.
//
// allowedOrigins is comma-separated (CLIENT_ORIGIN in .env), so more
// than one frontend — e.g. local dev plus a deployed one — can be
// allowed at once.
func CORS(allowedOrigins string) gin.HandlerFunc {
	origins := make(map[string]bool)
	for _, o := range strings.Split(allowedOrigins, ",") {
		if o = strings.TrimSpace(o); o != "" {
			origins[o] = true
		}
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origins[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		}

		// The browser sends a preflight OPTIONS request before certain
		// cross-origin requests (anything carrying our Authorization
		// header, or a multipart POST/PUT) to check these headers BEFORE
		// sending the real request. Answer it here directly — it should
		// never reach a route handler, none of which register OPTIONS.
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
