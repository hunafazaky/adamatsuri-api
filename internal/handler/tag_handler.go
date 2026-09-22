package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hunafazaky/event-booking-api/internal/response"
	"github.com/hunafazaky/event-booking-api/internal/service"
)

type TagHandler struct {
	service service.TagService
}

func NewTagHandler(service service.TagService) *TagHandler {
	return &TagHandler{service: service}
}

// GetTags godoc
// @Summary Get all tags
// @Description Returns every fandom/genre tag currently in use, for building a filter UI or a tag picker.
// @Tags Tags
// @Produce json
// @Success 200 {object} response.Envelope{data=[]dto.TagResponse}
// @Failure 500 {object} response.Envelope
// @Router /tags [get]
func (h *TagHandler) GetTags(c *gin.Context) {
	tags, err := h.service.List()
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "data tags retrieved", tags)
}
