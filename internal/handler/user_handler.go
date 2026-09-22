package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hunafazaky/event-booking-app/internal/response"
	"github.com/hunafazaky/event-booking-app/internal/service"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// SignUp godoc
// @Summary Sign Up
// @Description Creates a new user account.
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body service.SignUpInput true "Sign up payload"
// @Success 201 {object} response.Envelope{data=dto.UserResponse}
// @Failure 400 {object} response.Envelope
// @Failure 409 {object} response.Envelope
// @Router /auth/signup [post]
func (h *UserHandler) SignUp(c *gin.Context) {
	var input service.SignUpInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.SignUp(input)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "user created successfully", user)
}

// SignIn godoc
// @Summary Sign In
// @Description Authenticates a user and returns a JWT with the user's profile.
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body service.SignInInput true "Sign in payload"
// @Success 200 {object} response.Envelope{data=dto.SignInResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Router /auth/signin [post]
func (h *UserHandler) SignIn(c *gin.Context) {
	var input service.SignInInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.SignIn(input)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "user signed successfully", user)
}

// UpdateInterests godoc
// @Summary Update favorite series
// @Description Replaces the authenticated user's full list of favorite-series/interest tags. Unknown tag names are created on the fly, reusing the same tag pool events are tagged with.
// @Security BearerAuth
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body UpdateInterestsInput true "New interest list (replaces the existing one; empty array clears it)"
// @Success 200 {object} response.Envelope{data=dto.UserResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Router /auth/me/interests [put]
func (h *UserHandler) UpdateInterests(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.FromError(c, err)
		return
	}

	var input UpdateInterestsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.UpdateInterests(userID, input.Interests)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "interests updated successfully", user)
}

// UpdateInterestsInput is the request body for PUT /auth/me/interests.
// Interests has no "required" binding on purpose — an empty array (or
// an omitted field) is valid and means "clear all interests", not a
// validation error.
type UpdateInterestsInput struct {
	Interests []string `json:"interests" example:"Jujutsu Kaisen,shounen"`
}

// @Summary Get user's data
// @Description Returns the profile of the authenticated user.
// @Security BearerAuth
// @Tags Auth
// @Produce json
// @Success 200 {object} response.Envelope{data=dto.UserResponse}
// @Failure 401 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Router /auth/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.FromError(c, err)
		return
	}

	user, err := h.service.GetAuthUser(userID)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "data user retrieved", user)
}
