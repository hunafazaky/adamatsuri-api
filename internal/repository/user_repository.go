package repository

import (
	"github.com/hunafazaky/adamatsuri-api/internal/model"
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
	// UpdateName updates only the name column, by ID — deliberately
	// NOT a full Save(user) on a struct fetched via FindByID, whose
	// Select() only loads id/name/email/role. Saving that struct back
	// would write empty strings over Password and every other column
	// FindByID didn't select. Update() with an explicit column list
	// has no such risk.
	UpdateName(userID uint, name string) error
	// Delete soft-deletes the user by ID (GORM sets DeletedAt on any
	// model.User — no need to fetch the row first).
	Delete(userID uint) error
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

func (r *userRepository) UpdateName(userID uint, name string) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).Update("name", name).Error
}

func (r *userRepository) Delete(userID uint) error {
	return r.db.Delete(&model.User{}, userID).Error
}
