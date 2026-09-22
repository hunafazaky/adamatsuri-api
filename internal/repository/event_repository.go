package repository

import (
	"github.com/hunafazaky/event-booking-app/internal/model"
	"gorm.io/gorm"
)

type EventRepository interface {
	Create(event *model.Event) error
	Update(event *model.Event) error
	Delete(event *model.Event) error
	FindAll(search string, category model.Category, tag string, page, limit int) (events []model.Event, totalRows, totalPages int64, err error)
	FindByID(id uint) (*model.Event, error)
	FindByUserID(userID uint) ([]model.Event, error)
	// SetTags replaces an event's full set of tag associations with
	// exactly the ones given (an empty slice clears all tags). It only
	// touches the event_tags join table — it never creates/modifies Tag
	// rows themselves, so callers must pass already-persisted tags
	// (e.g. from TagRepository.FindOrCreateByNames).
	SetTags(event *model.Event, tags []model.Tag) error
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
	return r.db.Preload("User", userSummary).Preload("Tags").First(event, event.ID).Error
}

func (r *eventRepository) Update(event *model.Event) error {
	return r.db.Save(event).Error
}

func (r *eventRepository) Delete(event *model.Event) error {
	return r.db.Delete(event).Error
}

func (r *eventRepository) FindAll(search string, category model.Category, tag string, page, limit int) (events []model.Event, totalRows, totalPages int64, err error) {

	// Initialize Query
	query := r.db.Model(&model.Event{})
	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if tag != "" {
		// A subquery (not a JOIN) so filtering by tag can't multiply
		// result rows or break Limit/Offset pagination below.
		query = query.Where(
			"id IN (SELECT event_tags.event_id FROM event_tags JOIN tags ON tags.id = event_tags.tag_id WHERE LOWER(tags.name) = LOWER(?))",
			tag,
		)
	}

	// Total Rows
	query.Count(&totalRows)

	totalPages = (totalRows + int64(limit) - 1) / int64(limit)

	// Event & Error
	offset := (page - 1) * limit
	err = query.
		Preload("User", userSummary).
		Preload("Tags").
		Limit(limit).
		Offset(offset).
		Find(&events).Error
	return
}

func (r *eventRepository) FindByID(id uint) (*model.Event, error) {
	var event model.Event
	err := r.db.
		Preload("User", userSummary).
		Preload("Tags").
		Preload("Booking").
		Preload("Booking.User", userSummary).
		First(&event, id).Error
	return &event, err
}

func (r *eventRepository) FindByUserID(userID uint) ([]model.Event, error) {
	var events []model.Event
	err := r.db.
		Preload("User", userSummary).
		Preload("Tags").
		Where("user_id", userID).
		Find(&events).Error
	return events, err
}

func (r *eventRepository) SetTags(event *model.Event, tags []model.Tag) error {
	if err := r.db.Model(event).Association("Tags").Replace(tags); err != nil {
		return err
	}
	// Association.Replace already updates event.Tags in most GORM
	// versions, but setting it explicitly here removes any doubt for
	// the caller building a response right after this call.
	event.Tags = tags
	return nil
}
