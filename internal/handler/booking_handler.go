package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hunafazaky/event-booking-api/internal/response"
	"github.com/hunafazaky/event-booking-api/internal/service"
)

type BookingInput struct {
	Phone   string `json:"phone"`
	EventID uint   `json:"event_id"`
}

type BookingHandler struct {
	service service.BookingService
}

func NewBookingHandler(service service.BookingService) *BookingHandler {
	return &BookingHandler{service: service}
}

// CreateBooking godoc
// @Summary Create a booking
// @Description Creates a booking for an event for the authenticated user.
// @Security BearerAuth
// @Tags Bookings
// @Accept json
// @Produce json
// @Param input body service.CreateBookingInput true "Booking payload"
// @Success 201 {object} response.Envelope{data=dto.BookingResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 409 {object} response.Envelope
// @Router /bookings [post]
func (h *BookingHandler) CreateBooking(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.FromError(c, err)
		return
	}

	var input service.CreateBookingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	booking, err := h.service.Create(userID, input)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "booking created successfully", booking)
}

// GetBooks godoc
// @Summary Get user's booking list
// @Description Returns all bookings belonging to the authenticated user.
// @Security BearerAuth
// @Tags Bookings
// @Produce json
// @Success 200 {object} response.Envelope{data=[]dto.BookingResponse}
// @Failure 401 {object} response.Envelope
// @Router /bookings [get]
func (h *BookingHandler) GetBooks(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.FromError(c, err)
		return
	}

	bookings, err := h.service.GetByUser(userID)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "booking retrieved successfully", bookings)
}

// DeleteBooking godoc
// @Summary Delete a booking
// @Description Deletes a booking owned by the authenticated user.
// @Security BearerAuth
// @Tags Bookings
// @Produce json
// @Param id path int true "Booking ID" example(1)
// @Success 200 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Router /bookings/{id} [delete]
func (h *BookingHandler) DeleteBooking(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.FromError(c, err)
		return
	}

	bookingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid query parameter value")
		return
	}

	if err := h.service.Delete(userID, uint(bookingID)); err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "data deleted successfully", nil)
}
