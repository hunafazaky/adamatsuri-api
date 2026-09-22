package repository

import (
	"github.com/hunafazaky/event-booking-app/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByID(id uint) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
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
	return &user, r.db.Select("id", "name", "email", "role").First(&user, id).Error
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	return &user, r.db.Where("email = ?", email).First(&user).Error
}
