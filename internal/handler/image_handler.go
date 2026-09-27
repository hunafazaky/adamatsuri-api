package handler

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ImageProxyHandler exists for one reason: some visitors' ISPs
// interfere with direct connections to ImageKit's CDN domain (seen in
// the wild as a TLS handshake failure — ERR_SSL_VERSION_OR_CIPHER_MISMATCH
// — that disappears over a VPN, which points to ISP-level DPI
// tampering rather than an ImageKit outage). This server, wherever
// it's actually deployed, isn't behind that same ISP — so it fetches
// the image on the visitor's behalf and streams it back over a
// connection to OUR domain instead, which the visitor's ISP has no
// reason to interfere with.
//
// This intentionally skips the usual handler->service->repository
// layering: there's no business logic and no database involved, just
// a network passthrough, so a service layer would only add
// indirection without adding meaning.
type ImageProxyHandler struct {
	// allowedHosts is the set of hosts this proxy will fetch from —
	// critical for safety, since blindly fetching whatever URL a
	// caller supplies turns this into an open SSRF proxy (a way to
	// make OUR server issue requests to internal/arbitrary hosts on a
	// caller's behalf). ImageKit's shared domain is always allowed;
	// a configured custom domain (IMAGEKIT_URL_ENDPOINT) widens it.
	allowedHosts map[string]bool
	httpClient   *http.Client
}

func NewImageProxyHandler(imageKitURLEndpoint string) *ImageProxyHandler {
	allowed := map[string]bool{"ik.imagekit.io": true}

	if imageKitURLEndpoint != "" {
		if u, err := url.Parse(imageKitURLEndpoint); err == nil && u.Host != "" {
			allowed[u.Host] = true
		}
	}

	return &ImageProxyHandler{
		allowedHosts: allowed,
		httpClient:   &http.Client{Timeout: 15 * time.Second},
	}
}

// GetImage godoc
// @Summary Proxy an event image
// @Description Fetches an image from ImageKit server-side and streams it back — works around ISP-level interference with direct connections to ImageKit's CDN domain. Only ImageKit hosts are allowed as a source, to prevent this becoming an open proxy.
// @Tags Images
// @Param src query string true "The ImageKit image URL to proxy" example(https://ik.imagekit.io/your_id/photo.jpg)
// @Produce image/*
// @Success 200 {file} binary
// @Failure 400 {object} response.Envelope
// @Failure 502 {object} response.Envelope
// @Router /images/proxy [get]
func (h *ImageProxyHandler) GetImage(c *gin.Context) {
	src := c.Query("src")
	if src == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "src is required"})
		return
	}

	parsed, err := url.Parse(src)
	if err != nil || (parsed.Scheme != "https") || !h.allowedHosts[parsed.Host] {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "src must be an https ImageKit URL"})
		return
	}

	req, err := http.NewRequest(http.MethodGet, parsed.String(), nil)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"success": false, "error": "failed to build upstream request"})
		return
	}

	upstream, err := h.httpClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"success": false, "error": "failed to reach image source"})
		return
	}
	defer upstream.Body.Close()

	if upstream.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadGateway, gin.H{"success": false, "error": "image source returned an error"})
		return
	}

	contentType := upstream.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		c.JSON(http.StatusBadGateway, gin.H{"success": false, "error": "source did not return an image"})
		return
	}

	// Long cache lifetime — an uploaded event image is effectively
	// immutable at a given URL (editing an event uploads a NEW file
	// with a new URL rather than overwriting the old one), so there's
	// no correctness cost to caching it hard, and every cache hit is
	// one less request that has to touch this proxy at all.
	c.Header("Cache-Control", "public, max-age=31536000, immutable")
	c.Header("Content-Type", contentType)
	c.Status(http.StatusOK)
	_, _ = io.Copy(c.Writer, upstream.Body)
}
