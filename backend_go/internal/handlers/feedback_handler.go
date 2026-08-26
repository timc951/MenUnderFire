package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"menunderfire/internal/logger"
	"menunderfire/internal/models"
	"menunderfire/internal/services"

	"github.com/gorilla/mux"
)

// FeedbackHandler handles HTTP requests for feedback endpoints
type FeedbackHandler struct {
	feedbackService services.FeedbackService
	userService     services.UserService
}

// NewFeedbackHandler creates a new FeedbackHandler
func NewFeedbackHandler(feedbackService services.FeedbackService, userService services.UserService) *FeedbackHandler {
	return &FeedbackHandler{
		feedbackService: feedbackService,
		userService:     userService,
	}
}

// Create handles POST /api/feedback
func (h *FeedbackHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	var req models.CreateFeedbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondBadRequest(w, "Invalid request body")
		return
	}

	feedback, err := h.feedbackService.Create(ctx, user.ID, &req)
	if err != nil {
		logger.Error().Err(err).Msg("Error creating feedback")
		if errors.Is(err, services.ErrValidation) {
			respondBadRequest(w, err.Error())
			return
		}
		respondInternalError(w, "Failed to create feedback")
		return
	}

	respondCreated(w, feedback.ToResponse(user.DisplayName, user.Email))
}

// ListMine handles GET /api/feedback/me
func (h *FeedbackHandler) ListMine(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	items, err := h.feedbackService.ListByUser(ctx, user.ID)
	if err != nil {
		logger.Error().Err(err).Msg("Error listing user feedback")
		respondInternalError(w, "Failed to list feedback")
		return
	}

	responses := make([]*models.FeedbackResponse, len(items))
	for i, f := range items {
		responses[i] = f.ToResponse(user.DisplayName, user.Email)
	}

	respondOK(w, responses)
}

// ListAll handles GET /api/feedback (site admin only)
func (h *FeedbackHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	items, err := h.feedbackService.ListAll(ctx, user.ID)
	if err != nil {
		logger.Error().Err(err).Msg("Error listing all feedback")
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "Only site admins can view all feedback")
			return
		}
		respondInternalError(w, "Failed to list feedback")
		return
	}

	// Build responses with user info
	responses := make([]*models.FeedbackResponse, len(items))
	for i, f := range items {
		u, _ := h.userService.GetByID(ctx, f.UserID)
		displayName := ""
		email := ""
		if u != nil {
			displayName = u.DisplayName
			email = u.Email
		}
		responses[i] = f.ToResponse(displayName, email)
	}

	respondOK(w, responses)
}

// UpdateStatus handles PATCH /api/feedback/{id}/status (site admin only)
func (h *FeedbackHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	feedbackID := vars["id"]

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	var req models.UpdateFeedbackStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondBadRequest(w, "Invalid request body")
		return
	}

	feedback, err := h.feedbackService.UpdateStatus(ctx, user.ID, feedbackID, &req)
	if err != nil {
		logger.Error().Err(err).Msg("Error updating feedback status")
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "Only site admins can update feedback status")
			return
		}
		if errors.Is(err, services.ErrNotFound) {
			respondNotFound(w, "Feedback not found")
			return
		}
		if errors.Is(err, services.ErrValidation) {
			respondBadRequest(w, err.Error())
			return
		}
		respondInternalError(w, "Failed to update feedback status")
		return
	}

	u, _ := h.userService.GetByID(ctx, feedback.UserID)
	displayName := ""
	email := ""
	if u != nil {
		displayName = u.DisplayName
		email = u.Email
	}

	respondOK(w, feedback.ToResponse(displayName, email))
}
