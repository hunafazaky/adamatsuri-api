package repository

import (
	"github.com/hunafazaky/event-booking-app/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByID(id uint) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	// SetInterests replaces a user's full set of interest-tag
	// associations with exactly the ones given (an empty slice clears
	// all interests). Like EventRepository.SetTags, it only touches
	// the join table — callers must pass already-persisted Tag rows.
	SetInterests(user *model.User, tags []model.Tag) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	return &user, r.db.Select("id", "name", "email", "role").Preload("Interests").First(&user, id).Error
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	return &user, r.db.Preload("Interests").Where("email = ?", email).First(&user).Error
}

func (r *userRepository) SetInterests(user *model.User, tags []model.Tag) error {
	if err := r.db.Model(user).Association("Interests").Replace(tags); err != nil {
		return err
	}
	user.Interests = tags
	return nil
}
