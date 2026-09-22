package repository

import (
	"github.com/hunafazaky/event-booking-app/internal/model"
	"gorm.io/gorm"
)

type EventRepository interface {
	Create(event *model.Event) error
	Update(event *model.Event) error
	Delete(event *model.Event) error
	FindAll(search string, category model.Category, page, limit int) (events []model.Event, totalRows, totalPages int64, err error)
	FindByID(id uint) (*model.Event, error)
	FindByUserID(userID uint) ([]model.Event, error)
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) Create(event *model.Event) error {
	if err := r.db.Create(event).Error; err != nil {
		return err
	}
	return r.db.Preload("User", userSummary).First(event, event.ID).Error
}

func (r *eventRepository) Update(event *model.Event) error {
	return r.db.Save(event).Error
}

func (r *eventRepository) Delete(event *model.Event) error {
	return r.db.Delete(event).Error
}

func (r *eventRepository) FindAll(search string, category model.Category, page, limit int) (events []model.Event, totalRows, totalPages int64, err error) {

	// Initialize Query
	query := r.db.Model(&model.Event{})
	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// Total Rows
	query.Count(&totalRows)

	totalPages = (totalRows + int64(limit) - 1) / int64(limit)

	// Event & Error
	offset := (page - 1) * limit
	err = query.
		Preload("User", userSummary).
		Limit(limit).
		Offset(offset).
		Find(&events).Error
	return
}

func (r *eventRepository) FindByID(id uint) (*model.Event, error) {
	var event model.Event
	err := r.db.
		Preload("User", userSummary).
		Preload("Booking").
		Preload("Booking.User", userSummary).
		First(&event, id).Error
	return &event, err
}

func (r *eventRepository) FindByUserID(userID uint) ([]model.Event, error) {
	var events []model.Event
	err := r.db.
		Preload("User", userSummary).
		Where("user_id", userID).
		Find(&events).Error
	return events, err
}
