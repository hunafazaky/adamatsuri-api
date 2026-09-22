package service

import (
	"io"
	"time"

	"github.com/hunafazaky/event-booking-api/internal/apperror"
	"github.com/hunafazaky/event-booking-api/internal/dto"
	"github.com/hunafazaky/event-booking-api/internal/model"
	"github.com/hunafazaky/event-booking-api/internal/repository"
)

// CreateEventInput is what the handler builds from the multipart form and
// hands to the service. Image/ImageName are plain io.Reader + string —
// no multipart-specific type crosses into the service layer.
type CreateEventInput struct {
	Name        string
	Description string
	Location    string
	DateTime    time.Time
	Category    model.Category
	Tags        []string // fandom/genre tag names; optional, may be empty
	Image       io.Reader
	ImageName   string
}

type UpdateEventInput struct {
	Name        string
	Description string
	Location    string
	DateTime    *time.Time
	Category    model.Category
	// Tags is nil when the caller didn't send a "tags" field at all —
	// meaning leave the event's existing tags untouched. A non-nil
	// (possibly empty) slice means replace the full tag set with this
	// one; an empty slice clears all tags.
	Tags      []string
	Image     io.Reader
	ImageName string
}

type EventService interface {
	Create(userID uint, input CreateEventInput) (*dto.EventResponse, error)
	List(search string, category model.Category, tag string, page, limit int) ([]dto.EventResponse, dto.EventListMeta, error)
	GetByID(id uint) (*dto.EventDetailResponse, error)
	GetByUser(userID uint) ([]dto.EventResponse, error)
	Update(userID, eventID uint, input UpdateEventInput) (*dto.EventResponse, error)
	Delete(userID, eventID uint) error
}

type eventService struct {
	repo     repository.EventRepository
	tagRepo  repository.TagRepository
	uploader ImageUploader
}

func NewEventService(repo repository.EventRepository, tagRepo repository.TagRepository, uploader ImageUploader) EventService {
	return &eventService{repo: repo, tagRepo: tagRepo, uploader: uploader}
}

func (s *eventService) Create(userID uint, input CreateEventInput) (*dto.EventResponse, error) {
	if !model.IsValidCategory(input.Category) {
		return nil, apperror.BadRequest("invalid category")
	}

	// Upload FIRST, before touching the database. If the upload fails,
	// there's nothing to roll back — we simply never created a DB row
	// with a broken/missing image reference.
	imageURL, imageID, err := s.uploader.Upload(input.Image, input.ImageName)
	if err != nil {
		return nil, apperror.Internal("failed to upload image", err)
	}

	event := model.Event{
		Name:        input.Name,
		Description: input.Description,
		Location:    input.Location,
		DateTime:    input.DateTime,
		Category:    input.Category,
		UserID:      userID,
		Image:       imageURL,
		ImageID:     imageID,
	}

	// Tags are deliberately NOT set on the struct before Create: GORM
	// would try to upsert every associated Tag row too, which risks
	// duplicate-key errors for tags that already exist. Instead the
	// event is created tag-less, then SetTags below only touches the
	// join table against already-persisted Tag rows.
	if err := s.repo.Create(&event); err != nil {
		// The upload already succeeded but the DB write didn't — clean up
		// the now-orphaned file rather than leaving it in ImageKit forever.
		// Best-effort: if THIS also fails, we still report the original
		// DB error to the client, not the cleanup failure. Once Phase 9
		// adds logging, this is exactly the spot that should log instead
		// of silently discarding the cleanup error.
		_ = s.uploader.Delete(imageID)
		return nil, apperror.Internal("failed to create event", err)
	}

	if len(input.Tags) > 0 {
		tags, err := s.tagRepo.FindOrCreateByNames(input.Tags)
		if err != nil {
			return nil, apperror.Internal("failed to process tags", err)
		}
		if err := s.repo.SetTags(&event, tags); err != nil {
			return nil, apperror.Internal("failed to associate tags", err)
		}
	}

	// repo.Create reloads the row with Preload("User") after insert, so
	// event.User is populated here — toEventResponse gets the real name
	// and email, not just the ID. event.Tags is set by SetTags above
	// (or stays nil/empty if no tags were given).
	response := toEventResponse(event)
	return &response, nil
}

func (s *eventService) List(search string, category model.Category, tag string, page, limit int) ([]dto.EventResponse, dto.EventListMeta, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 6
	}

	if category != "" && !model.IsValidCategory(category) {
		return nil, dto.EventListMeta{}, apperror.BadRequest("invalid category")
	}

	events, totalRows, totalPage, err := s.repo.FindAll(search, category, tag, page, limit)
	if err != nil {
		return nil, dto.EventListMeta{}, apperror.Internal("failed to load data", err)
	}

	eventResponse := make([]dto.EventResponse, 0, len(events))
	for _, item := range events {
		eventResponse = append(eventResponse, toEventResponse(item))
	}

	eventListMeta := dto.EventListMeta{
		Page:      page,
		Limit:     limit,
		TotalRows: totalRows,
		TotalPage: totalPage,
	}

	return eventResponse, eventListMeta, nil
}

func (s *eventService) GetByID(id uint) (*dto.EventDetailResponse, error) {
	event, err := s.repo.FindByID(id)
	if err != nil {
		return nil, mapLookupError(err, "event not found", "failed to load event")
	}

	eventResponse := toEventResponse(*event)

	bookingSummaryResponse := make([]dto.BookingSummaryResponse, 0, len(event.Booking))
	for _, item := range event.Booking {
		bookingSummaryResponse = append(bookingSummaryResponse, dto.BookingSummaryResponse{
			ID:          item.ID,
			BookingCode: item.BookingCode,
			Phone:       item.Phone,
			User:        toUserResponse(item.User),
		})
	}

	eventDetailResponse := dto.EventDetailResponse{
		EventResponse: eventResponse,
		Bookings:      bookingSummaryResponse,
	}
	return &eventDetailResponse, nil
}

func (s *eventService) GetByUser(userID uint) ([]dto.EventResponse, error) {
	events, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, apperror.Internal("failed to load events", err)
	}

	listEventResponse := make([]dto.EventResponse, 0, len(events))
	for _, item := range events {
		listEventResponse = append(listEventResponse, toEventResponse(item))
	}

	return listEventResponse, nil
}

func (s *eventService) Update(userID, eventID uint, input UpdateEventInput) (*dto.EventResponse, error) {
	// get current event
	event, err := s.repo.FindByID(eventID)
	if err != nil {
		return nil, mapLookupError(err, "event not found", "failed to load event")
	}

	// verify permission
	if event.UserID != userID {
		return nil, apperror.Forbidden("you're not authorized")
	}

	// file handler
	if input.Image != nil {
		// Capture the OLD file ID before we overwrite it on event — once
		// event.ImageID is reassigned below, this is the only reference
		// to the file we're about to replace.
		oldImageID := event.ImageID

		imageURL, imageID, err := s.uploader.Upload(input.Image, input.ImageName)
		if err != nil {
			return nil, apperror.Internal("failed to upload image", err)
		}
		event.Image = imageURL
		event.ImageID = imageID

		// Delete the old file only AFTER the new upload succeeds — if
		// upload had failed above, we return before reaching here, so the
		// old (still valid) image is never touched. Best-effort: a
		// failure to delete the old file doesn't fail this request: the
		// update itself already succeeded from the user's perspective.
		if oldImageID != "" {
			_ = s.uploader.Delete(oldImageID)
		}
	}

	// replace data
	if input.Name != "" {
		event.Name = input.Name
	}
	if input.Description != "" {
		event.Description = input.Description
	}
	if input.Location != "" {
		event.Location = input.Location
	}
	if input.DateTime != nil {
		event.DateTime = *input.DateTime
	}
	if input.Category != "" {
		if !model.IsValidCategory(input.Category) {
			return nil, apperror.BadRequest("invalid category")
		}
		event.Category = input.Category
	}

	if err := s.repo.Update(event); err != nil {
		return nil, apperror.Internal("failed to update event", err)
	}

	// nil means the request didn't send a "tags" field at all — leave
	// the existing tags alone. A non-nil (possibly empty) slice means
	// replace the full set, which SetTags does even for an empty slice
	// (clearing every tag).
	if input.Tags != nil {
		tags, err := s.tagRepo.FindOrCreateByNames(input.Tags)
		if err != nil {
			return nil, apperror.Internal("failed to process tags", err)
		}
		if err := s.repo.SetTags(event, tags); err != nil {
			return nil, apperror.Internal("failed to associate tags", err)
		}
	}

	// event.User and event.Tags were already preloaded by the FindByID
	// call above (Tags is refreshed by SetTags when tags changed) —
	// toEventResponse gets the real data here, not an ID-only stub.
	eventResponse := toEventResponse(*event)

	return &eventResponse, nil
}

func (s *eventService) Delete(userID, eventID uint) error {
	event, err := s.repo.FindByID(eventID)
	if err != nil {
		return mapLookupError(err, "event not found", "failed to load event")
	}

	// verify permission
	if event.UserID != userID {
		return apperror.Forbidden("you're not authorized")
	}

	if err := s.repo.Delete(event); err != nil {
		return apperror.Internal("failed to delete event", err)
	}

	return nil
}
