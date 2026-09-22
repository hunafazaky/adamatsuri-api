package service

import (
	"github.com/hunafazaky/event-booking-api/internal/apperror"
	"github.com/hunafazaky/event-booking-api/internal/dto"
	"github.com/hunafazaky/event-booking-api/internal/repository"
)

type TagService interface {
	List() ([]dto.TagResponse, error)
}

type tagService struct {
	repo repository.TagRepository
}

func NewTagService(repo repository.TagRepository) TagService {
	return &tagService{repo: repo}
}

func (s *tagService) List() ([]dto.TagResponse, error) {
	tags, err := s.repo.FindAll()
	if err != nil {
		return nil, apperror.Internal("failed to load tags", err)
	}

	return toTagResponses(tags), nil
}
