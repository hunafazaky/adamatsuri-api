package service

import (
	"fmt"
	"time"

	"github.com/hunafazaky/event-booking-api/internal/apperror"
	"github.com/hunafazaky/event-booking-api/internal/dto"
	"github.com/hunafazaky/event-booking-api/internal/model"
	"github.com/hunafazaky/event-booking-api/internal/repository"
)

type CreateBookingInput struct {
	Phone   string `json:"phone" example:"+6281234567890"`
	EventID uint   `json:"event_id" example:"1"`
}

type BookingService interface {
	Create(userID uint, input CreateBookingInput) (*dto.BookingResponse, error)
	GetByUser(userID uint) ([]dto.BookingResponse, error)
	Delete(userID, bookingID uint) error
}

type bookingService struct {
	repo      repository.BookingRepository
	eventRepo repository.EventRepository
}

func NewBookingService(repo repository.BookingRepository, eventRepo repository.EventRepository) BookingService {
	return &bookingService{repo: repo, eventRepo: eventRepo}
}

func (s *bookingService) Create(userID uint, input CreateBookingInput) (*dto.BookingResponse, error) {
	event, err := s.eventRepo.FindByID(input.EventID)
	if err != nil {
		return nil, mapLookupError(err, "event not found", "failed to load event")
	}

	exists, err := s.repo.ExistsByUserAndEvent(userID, input.EventID)
	if err != nil {
		return nil, apperror.Internal("failed to check existing booking", err)
	}
	if exists {
		return nil, apperror.Conflict("you already booked this event")
	}

	bookingCode := fmt.Sprintf("BK-%sE%dU%d", time.Now().Format("20060102"), input.EventID, userID)

	booking := model.Booking{
		BookingCode: bookingCode,
		Phone:       input.Phone,
		UserID:      userID,
		EventID:     input.EventID,
	}

	if err := s.repo.Create(&booking); err != nil {
		return nil, apperror.Internal("failed to create booking", err)
	}

	bookingResponse := dto.BookingResponse{
		ID:          booking.ID,
		BookingCode: bookingCode,
		Phone:       booking.Phone,
		Event:       toEventResponse(*event),
	}

	return &bookingResponse, nil
}

func (s *bookingService) GetByUser(userID uint) ([]dto.BookingResponse, error) {
	bookings, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, apperror.Internal("failed to load bookings", err)
	}

	listBookingResponse := make([]dto.BookingResponse, 0, len(bookings))
	for _, item := range bookings {
		listBookingResponse = append(listBookingResponse, dto.BookingResponse{
			ID:          item.ID,
			BookingCode: item.BookingCode,
			Phone:       item.Phone,
			Event:       toEventResponse(item.Event),
		})
	}
	return listBookingResponse, nil
}

func (s *bookingService) Delete(userID, bookingID uint) error {
	booking, err := s.repo.FindByID(bookingID)
	if err != nil {
		return mapLookupError(err, "booking not found", "failed to load booking")
	}

	if userID != booking.UserID {
		return apperror.Forbidden("you're not authorized")
	}

	if err := s.repo.Delete(booking); err != nil {
		return apperror.Internal("failed to delete booking", err)
	}

	return nil
}
