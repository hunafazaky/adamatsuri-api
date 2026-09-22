package service

import (
	"errors"

	"github.com/hunafazaky/event-booking-app/internal/apperror"
	"github.com/hunafazaky/event-booking-app/internal/dto"
	"github.com/hunafazaky/event-booking-app/internal/model"
	"gorm.io/gorm"
)

func toUserResponse(user model.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		Interests: toTagResponses(user.Interests),
	}
}

func toTagResponse(tag model.Tag) dto.TagResponse {
	return dto.TagResponse{
		ID:   tag.ID,
		Name: tag.Name,
	}
}

func toTagResponses(tags []model.Tag) []dto.TagResponse {
	responses := make([]dto.TagResponse, 0, len(tags))
	for _, tag := range tags {
		responses = append(responses, toTagResponse(tag))
	}
	return responses
}

func toEventResponse(event model.Event) dto.EventResponse {
	return dto.EventResponse{
		ID:          event.ID,
		Name:        event.Name,
		Description: event.Description,
		Location:    event.Location,
		Image:       event.Image,
		DateTime:    event.DateTime,
		Category:    event.Category,
		Tags:        toTagResponses(event.Tags),
		User:        toUserResponse(event.User),
		CreatedAt:   event.CreatedAt,
	}
}

// mapLookupError turns the error from a SINGLE-RECORD lookup
// (First/Take/Last — e.g. repo.FindByID) into the right AppError:
// NotFound when nothing matched, Internal for any other failure.
//
// Only call this after a lookup that CAN return gorm.ErrRecordNotFound.
// A Find() into a slice never returns that error, even with zero
// results — for those, skip this and just use apperror.Internal(msg, err)
// directly. Calling mapLookupError after a Find() would be harmless
// (the errors.Is check just never matches), but it'd be misleading to
// read, since it implies a "not found" case that can't actually happen.
func mapLookupError(err error, notFoundMsg, internalMsg string) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.NotFound(notFoundMsg)
	}
	return apperror.Internal(internalMsg, err)
}
