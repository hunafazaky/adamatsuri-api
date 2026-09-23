package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hunafazaky/adamatsuri-api/internal/apperror"
	"github.com/hunafazaky/adamatsuri-api/internal/dto"
	"github.com/hunafazaky/adamatsuri-api/internal/model"
	"github.com/hunafazaky/adamatsuri-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SignUpInput struct {
	Name     string `json:"name" binding:"required" example:"Jane Doe"`
	Email    string `json:"email" binding:"required,email" example:"jane@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"secret123"`
	// Role is optional and defaults to "attendee" when omitted. Only
	// "attendee" and "organizer" may be requested here — "admin" is
	// deliberately not reachable through sign-up.
	Role model.Role `json:"role" binding:"omitempty,oneof=attendee organizer" example:"attendee"`
}

type SignInInput struct {
	Email    string `json:"email" binding:"required" example:"jane@example.com"`
	Password string `json:"password" binding:"required" example:"secret123"`
}

type UserService interface {
	SignUp(input SignUpInput) (*dto.UserResponse, error)
	SignIn(input SignInInput) (*dto.SignInResponse, error)
	GetAuthUser(userID uint) (*dto.UserResponse, error)
	UpdateInterests(userID uint, names []string) (*dto.UserResponse, error)
	UpdateProfile(userID uint, name string) (*dto.UserResponse, error)
	DeleteAccount(userID uint) error
}

type userService struct {
	repo      repository.UserRepository
	tagRepo   repository.TagRepository
	jwtSecret string
}

func NewUserService(repo repository.UserRepository, tagRepo repository.TagRepository, jwtSecret string) UserService {
	return &userService{repo: repo, tagRepo: tagRepo, jwtSecret: jwtSecret}
}

func (s *userService) SignUp(input SignUpInput) (*dto.UserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperror.Internal("Failed to Hash Password", err)
	}

	role := input.Role
	if role == "" {
		role = model.RoleAttendee
	}

	user := model.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     role,
	}

	if err := s.repo.Create(&user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.Conflict("Email already registered")
		}
		return nil, apperror.Internal("Failed to create User", err)
	}

	response := toUserResponse(user)
	return &response, nil
}

func (s *userService) SignIn(input SignInInput) (*dto.SignInResponse, error) {
	user, err := s.repo.FindByEmail(input.Email)
	if err != nil {
		return nil, apperror.Unauthorized("Invalid Email or Password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, apperror.Unauthorized("Invalid email or password")
	}

	// Role travels inside the JWT so RequireAuth/RequireRole can check
	// it without a DB round trip on every request. Trade-off: if a
	// user's role changes (once there's an admin-promotion path), the
	// change won't take effect until they sign in again and get a new
	// token — acceptable for now, worth revisiting later.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"role": string(user.Role),
		"exp":  time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, apperror.Internal("Failed to sign token", err)
	}

	userResponse := toUserResponse(*user)

	response := dto.SignInResponse{
		Token: tokenString,
		User:  userResponse,
	}

	return &response, nil
}

func (s *userService) GetAuthUser(userID uint) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, mapLookupError(err, "user not found", "failed to load user")
	}

	response := toUserResponse(*user)

	return &response, nil
}

// UpdateInterests replaces the authenticated user's full set of
// favorite-series/interest tags. Unlike Event's tags field, this is a
// dedicated endpoint whose only job is setting interests, so there's no
// "nil means unchanged" ambiguity to worry about — the given names ARE
// the new interest list, and an empty list clears it.
func (s *userService) UpdateInterests(userID uint, names []string) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, mapLookupError(err, "user not found", "failed to load user")
	}

	tags, err := s.tagRepo.FindOrCreateByNames(names)
	if err != nil {
		return nil, apperror.Internal("failed to process interests", err)
	}

	if err := s.repo.SetInterests(user, tags); err != nil {
		return nil, apperror.Internal("failed to update interests", err)
	}

	response := toUserResponse(*user)
	return &response, nil
}

// UpdateProfile currently only lets a user change their name. Email
// isn't editable here — changing it raises questions (re-verification,
// uniqueness checks) this endpoint doesn't attempt to answer.
func (s *userService) UpdateProfile(userID uint, name string) (*dto.UserResponse, error) {
	if name == "" {
		return nil, apperror.BadRequest("name is required")
	}

	// UpdateName touches only the name column (see the repository
	// method's own comment for why) — so re-fetch afterward for a
	// complete, accurate response rather than mutating a struct that
	// FindByID only partially populated in the first place.
	if err := s.repo.UpdateName(userID, name); err != nil {
		return nil, apperror.Internal("failed to update profile", err)
	}

	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, mapLookupError(err, "user not found", "failed to load user")
	}

	response := toUserResponse(*user)
	return &response, nil
}

// DeleteAccount soft-deletes the user row. It does NOT touch their
// events or bookings — those rows stay exactly as they are, just
// pointing at a now-soft-deleted user. That's a real, known
// trade-off (an event's organizer info would show blank once loaded
// through a query that excludes soft-deleted users), not an
// oversight — cascading this properly (reassigning or hard-deleting
// owned events, anonymizing bookings) is a bigger decision than
// "delete my account" implies and deserves its own discussion first.
func (s *userService) DeleteAccount(userID uint) error {
	if err := s.repo.Delete(userID); err != nil {
		return apperror.Internal("failed to delete account", err)
	}
	return nil
}
