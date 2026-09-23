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
<svg
   width="300"
   height="300"
   viewBox="0 0 300 300"
   version="1.1"
   id="svg1"
   inkscape:version="1.4 (86a8ad7, 2024-10-11)"
   sodipodi:docname="hz-negative.svg"
   inkscape:export-filename="hz-dark.svg"
   inkscape:export-xdpi="57.599998"
   inkscape:export-ydpi="57.599998"
   xml:space="preserve"
   xmlns:inkscape="http://www.inkscape.org/namespaces/inkscape"
   xmlns:sodipodi="http://sodipodi.sourceforge.net/DTD/sodipodi-0.dtd"
   xmlns="http://www.w3.org/2000/svg"
   xmlns:svg="http://www.w3.org/2000/svg"><sodipodi:namedview
     id="namedview1"
     pagecolor="#ffffff"
     bordercolor="#000000"
     borderopacity="0.25"
     inkscape:showpageshadow="2"
     inkscape:pageopacity="0.0"
     inkscape:pagecheckerboard="0"
     inkscape:deskcolor="#d1d1d1"
     inkscape:document-units="px"
     inkscape:zoom="2"
     inkscape:cx="77.5"
     inkscape:cy="145.25"
     inkscape:window-width="1920"
     inkscape:window-height="974"
     inkscape:window-x="-11"
     inkscape:window-y="-11"
     inkscape:window-maximized="1"
     inkscape:current-layer="svg1"
     showguides="false"><inkscape:page
       x="0"
       y="0"
       width="300"
       height="300"
       id="page2"
       margin="0"
       bleed="0" /></sodipodi:namedview><defs
     id="defs1"><mask
       maskUnits="userSpaceOnUse"
       id="mask19"><circle
         style="font-variation-settings:normal;opacity:1;fill:#ffffff;fill-opacity:1;stroke:none;stroke-width:10.6762;stroke-linecap:butt;stroke-linejoin:round;stroke-miterlimit:4;stroke-dasharray:none;stroke-dashoffset:0;stroke-opacity:1;paint-order:normal"
         id="circle19"
         cx="150"
         cy="150"
         r="150" /></mask></defs><g
     id="g20"
     inkscape:label="full"
     style="display:inline"
     transform="matrix(0.96721311,0,0,0.96721311,4.9180335,4.9180335)"><rect
       style="display:inline;fill:#0b192c;fill-opacity:1;stroke:none;stroke-width:4.99717;stroke-dasharray:none;stroke-opacity:1"
       id="rect3"
       width="299.83051"
       height="299.83051"
       x="0.084747314"
       y="0.084747314"
       inkscape:label="bg"
       ry="46.525425"
       inkscape:export-filename="rect3.svg"
       inkscape:export-xdpi="57.599998"
       inkscape:export-ydpi="57.599998" /><g
       id="g17"
       mask="url(#mask19)"
       inkscape:export-filename="rect4.svg"
       inkscape:export-xdpi="57.599998"
       inkscape:export-ydpi="57.599998"><path
         id="path5"
         style="display:inline;fill:#ff6500;fill-opacity:1;stroke-width:1.31457"
         d="m 236.60156,30 -17.32031,30 -43.30078,75 h -19.64063 -15 l -17.32031,30 h 15 19.64063 l -43.30078,75 -17.320318,30 H 270 V 240 H 150.00195 L 253.92383,60 H 300 V 30 Z" /><path
         id="path8"
         style="fill:#f2f2f2;fill-opacity:1;stroke-width:1.31457"
         d="m 30,30 v 30 h 120 l -43.30078,75 H 72.058594 L 54.738281,165 H 89.376953 L 46.076172,240 H 0 v 30 H 63.398438 L 80.716797,240 184.64062,60 201.96094,30 Z" /></g></g></svg>
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
) {
	server.Use(middleware.CORS(cfg.ClientOrigin))

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
