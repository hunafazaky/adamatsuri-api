package repository

import (
	"errors"

	"github.com/hunafazaky/event-booking-app/internal/model"
	"gorm.io/gorm"
)

type BookingRepository interface {
	Create(booking *model.Booking) error
	Delete(booking *model.Booking) error
	FindByID(id uint) (*model.Booking, error)
	FindByUserID(userID uint) ([]model.Booking, error)
	ExistsByUserAndEvent(userID, eventID uint) (bool, error)
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) Create(booking *model.Booking) error {
	return r.db.Create(booking).Error
}

func (r *bookingRepository) Delete(booking *model.Booking) error {
	return r.db.Delete(booking).Error
}

func (r *bookingRepository) FindByID(id uint) (*model.Booking, error) {
	var booking model.Booking
	return &booking, r.db.First(&booking, id).Error
}

func (r *bookingRepository) FindByUserID(userID uint) ([]model.Booking, error) {
	var bookings []model.Booking

	err := r.db.
		Preload("Event").
		Preload("Event.User", userSummary).
		Preload("Event.Tags").
		Where("user_id = ?", userID).
		Find(&bookings).Error
	return bookings, err
}

func (r *bookingRepository) ExistsByUserAndEvent(userID, eventID uint) (bool, error) {
	err := r.db.
		Where("user_id = ? AND event_id = ?", userID, eventID).
		First(&model.Booking{}).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
