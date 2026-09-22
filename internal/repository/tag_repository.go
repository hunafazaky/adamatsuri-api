package repository

import (
	"errors"
	"strings"

	"github.com/hunafazaky/event-booking-api/internal/model"
	"gorm.io/gorm"
)

type TagRepository interface {
	// FindOrCreateByNames resolves a list of tag names to persisted
	// Tag records — reusing an existing tag (matched case-insensitively)
	// where one exists, creating a new one otherwise. Blank names are
	// skipped and duplicate names (after trimming/case-folding) are
	// collapsed to a single tag.
	FindOrCreateByNames(names []string) ([]model.Tag, error)
	FindAll() ([]model.Tag, error)
}

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) FindOrCreateByNames(names []string) ([]model.Tag, error) {
	tags := make([]model.Tag, 0, len(names))
	seen := make(map[string]bool, len(names))

	for _, raw := range names {
		name := strings.TrimSpace(raw)
		if name == "" {
			continue
		}

		key := strings.ToLower(name)
		if seen[key] {
			continue
		}
		seen[key] = true

		var tag model.Tag
		err := r.db.Where("LOWER(name) = LOWER(?)", name).First(&tag).Error
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			tag = model.Tag{Name: name}
			if err := r.db.Create(&tag).Error; err != nil {
				return nil, err
			}
		case err != nil:
			return nil, err
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func (r *tagRepository) FindAll() ([]model.Tag, error) {
	var tags []model.Tag
	err := r.db.Order("name asc").Find(&tags).Error
	return tags, err
}
