// internal/router/router.go
package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hunafazaky/adamatsuri-api/internal/config"
	"github.com/hunafazaky/adamatsuri-api/internal/handler"
	"github.com/hunafazaky/adamatsuri-api/internal/middleware"
	"github.com/hunafazaky/adamatsuri-api/internal/model"
	"github.com/hunafazaky/adamatsuri-api/internal/response"

	"github.com/hunafazaky/adamatsuri-api/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// faviconSVG is a minimal, dependency-free favicon — no image file to
// manage, no extra static-serving route setup. Modern browsers accept an
// SVG returned from /favicon.ico as long as the Content-Type is correct,
// which is what the /favicon.ico route below sets. Swap the fill color or
// letter to reuse this in another project.
const faviconSVG = `
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32">
  <rect x="1.5" y="1.5" width="29" height="29" rx="6" fill="#A8D8FF" stroke="#1A1A1A" stroke-width="2.5"/>
  <text x="16" y="23" text-anchor="middle" font-family="Arial, sans-serif" font-weight="800" font-size="18" fill="#1A1A1A">A</text>
</svg>
`

// scalarReferenceHTML is a static page embedding Scalar's API Reference
// component via CDN. It only needs a URL to fetch the OpenAPI spec from —
// /openapi.json below, which is the same generated spec Swagger UI reads,
// just served as raw JSON instead of wrapped in a UI.
//
// data-configuration is Scalar's own settings object — theme:"default" and
// layout:"classic" together give the plainest, least "app-like" look
// Scalar offers: a single scrolling document instead of a three-pane
// dashboard. This block is intentionally generic — copy it as-is into any
// other project's docs page for the same clean baseline.
const scalarReferenceHTML = `<!doctype html>
<html>
<head>
	<title>AdaMatsuri API Reference</title>
	<meta charset="utf-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1" />
	<link rel="icon" type="image/svg+xml" href="/favicon.ico" />
</head>
<body>
	<script
		id="api-reference"
		data-url="/openapi.json"
		data-configuration='{"theme":"elysiajs","darkMode":true}'
	></script>
	<script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>`

// Setup registers every route on server. It takes the already-constructed
// handlers and config as parameters — it builds NOTHING itself, it only
// wires paths to methods. Construction (repos → services → handlers)
// stays entirely in main.go; this file's only job is routing.
func Setup(
	server *gin.Engine,
	cfg *config.Config,
	userHandler *handler.UserHandler,
	eventHandler *handler.EventHandler,
	bookingHandler *handler.BookingHandler,
	tagHandler *handler.TagHandler,
	imageHandler *handler.ImageProxyHandler,
) {
	server.Use(middleware.CORS(cfg.ClientOrigin))

	{
		api := server.Group("/api/images")
		api.GET("/proxy", imageHandler.GetImage)
	}

	{
		api := server.Group("/api/events")
		api.GET("", eventHandler.GetEvents)
		api.GET("/:id", middleware.OptionalAuth(cfg.JWTSecret), eventHandler.GetEventByID)
	}

	{
		api := server.Group("/api/tags")
		api.GET("", tagHandler.GetTags)
	}

	{
		api := server.Group("/api/auth")
		api.POST("/signup", userHandler.SignUp)
		api.POST("/signin", userHandler.SignIn)
	}

	{
		protectedApi := server.Group("/api")
		protectedApi.Use(middleware.RequireAuth(cfg.JWTSecret))
		protectedApi.GET("/auth/me", userHandler.GetMe)
		protectedApi.PATCH("/auth/me", userHandler.UpdateProfile)
		protectedApi.DELETE("/auth/me", userHandler.DeleteAccount)
		protectedApi.PUT("/auth/me/interests", userHandler.UpdateInterests)
		protectedApi.POST("/bookings", bookingHandler.CreateBooking)
		protectedApi.GET("/bookings", bookingHandler.GetBooks)
		protectedApi.DELETE("/bookings/:id", bookingHandler.DeleteBooking)
	}

	{
		// Event management (create/update/delete/list-mine) is
		// organizer-and-admin only. Ownership (a user can only touch
		// THEIR OWN event) is still checked separately inside
		// EventService — this middleware only gates "can create/manage
		// events at all", not "can manage this specific event".
		organizerApi := server.Group("/api")
		organizerApi.Use(middleware.RequireAuth(cfg.JWTSecret))
		organizerApi.Use(middleware.RequireRole(model.RoleOrganizer, model.RoleAdmin))
		organizerApi.POST("/events", eventHandler.CreateEvent)
		organizerApi.PUT("/events/:id", eventHandler.UpdateEvent)
		organizerApi.DELETE("/events/:id", eventHandler.DeleteEvent)
		organizerApi.GET("/events/mine", eventHandler.GetEventsMine)
	}

	// --- API documentation ---
	// Both UIs read the SAME generated spec (docs.SwaggerInfo) — main.go
	// overrides its Host once at startup if cfg.PublicHost is set, so
	// neither UI needs to know about deployment vs local on its own.
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server.GET("/openapi.json", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", []byte(docs.SwaggerInfo.ReadDoc()))
	})

	server.GET("/docs", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(scalarReferenceHTML))
	})

	server.GET("/favicon.ico", func(c *gin.Context) {
		c.Data(http.StatusOK, "image/svg+xml", []byte(faviconSVG))
	})

	// A friendly root instead of a bare 404 — mostly useful for anyone
	// (or any uptime monitor) hitting the bare domain and wondering what
	// they landed on.
	server.GET("/", func(c *gin.Context) {
		response.Success(c, http.StatusOK, "AdaMatsuri API", gin.H{
			"docs":    "/docs",
			"swagger": "/swagger/index.html",
		})
	})

	// Every OTHER response in this API — success or failure — uses
	// response.Envelope. Gin's default 404 (plain text "404 page not
	// found") was the one exception; this makes an unmatched route
	// consistent with everything else, including /api itself, which was
	// never registered as its own route and falls through to here too.
	server.NoRoute(func(c *gin.Context) {
		response.Fail(c, http.StatusNotFound, "route not found")
	})
}
