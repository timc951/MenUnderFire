package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"menunderfire/internal/logger"
	"menunderfire/internal/models"
	"menunderfire/internal/services"
)

// UserHandler handles HTTP requests for user-related endpoints
type UserHandler struct {
	userService services.UserService
}

// NewUserHandler creates a new UserHandler with the given service
func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetMe handles GET /api/users/me
// Returns the current authenticated user's profile
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger.Debug().Str("remote_addr", r.RemoteAddr).Msg("GetMe request received")

	// Get authenticated user from context (set by auth middleware)
	user, err := GetUserFromContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get user from context - authentication middleware may not have set user")
		respondUnauthorized(w, "Authentication required")
		return
	}

	logger.Debug().Str("user_id", user.ID).Str("email", user.Email).Msg("User found in context")

	// Fetch fresh user data from service
	freshUser, err := h.userService.GetByID(ctx, user.ID)
	if err != nil {
		logger.Error().Err(err).Str("user_id", user.ID).Msg("Error fetching user")
		if errors.Is(err, services.ErrUserNotFound) {
			respondNotFound(w, "User not found")
			return
		}
		respondInternalError(w, "Failed to fetch user")
		return
	}

	respondOK(w, freshUser.ToResponse())
}

// GetPermissions handles GET /api/users/me/permissions
// Returns the current user's roles and permissions
func (h *UserHandler) GetPermissions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get authenticated user from context
	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	logger.Info().
		Str("user_id", user.ID).
		Str("email", user.Email).
		Str("external_id", user.ExternalID).
		Bool("is_site_admin", user.IsSiteAdmin).
		Str("display_name", user.DisplayName).
		Msg("GetPermissions called")

	// Fetch permissions from service
	permissions, err := h.userService.GetPermissions(ctx, user.ID)
	if err != nil {
		logger.Error().Err(err).Str("user_id", user.ID).Msg("Error fetching permissions")
		if errors.Is(err, services.ErrUserNotFound) {
			respondNotFound(w, "User not found")
			return
		}
		respondInternalError(w, "Failed to fetch permissions")
		return
	}

	logger.Info().
		Str("user_id", user.ID).
		Bool("is_site_admin", permissions.IsSiteAdmin).
		Int("admin_orgs", len(permissions.AdminOfOrganizationIDs)).
		Int("owned_groups", len(permissions.OwnedGroupIDs)).
		Int("member_groups", len(permissions.MemberGroupIDs)).
		Msg("GetPermissions result")

	respondOK(w, permissions)
}

// UpdateMe handles PUT /api/users/me
// Updates the current user's profile
func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get authenticated user from context
	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	// Parse request body
	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondBadRequest(w, "Invalid request body")
		return
	}

	// Validate request
	if err := validateUpdateUserRequest(&req); err != nil {
		respondBadRequest(w, err.Error())
		return
	}

	// Update user via service
	updatedUser, err := h.userService.Update(ctx, user.ID, &req)
	if err != nil {
		logger.Error().Err(err).Str("user_id", user.ID).Msg("Error updating user")
		if errors.Is(err, services.ErrUserNotFound) {
			respondNotFound(w, "User not found")
			return
		}
		if errors.Is(err, services.ErrValidation) {
			respondBadRequest(w, err.Error())
			return
		}
		respondInternalError(w, "Failed to update user")
		return
	}

	respondOK(w, updatedUser.ToResponse())
}

// AcceptAgreement handles POST /api/users/me/accept-agreement
// Records the user's acceptance of the testing agreement
func (h *UserHandler) AcceptAgreement(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	var req models.AcceptAgreementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondBadRequest(w, "Invalid request body")
		return
	}

	if req.AgreementVersion == "" {
		respondBadRequest(w, "agreementVersion is required")
		return
	}

	// Extract IP and user agent for audit trail
	ipAddress := r.Header.Get("X-Forwarded-For")
	if ipAddress == "" {
		ipAddress = r.RemoteAddr
	}
	userAgent := r.Header.Get("User-Agent")

	response, err := h.userService.AcceptAgreement(ctx, user.ID, req.AgreementVersion, ipAddress, userAgent)
	if err != nil {
		logger.Error().Err(err).Str("user_id", user.ID).Msg("Error accepting agreement")
		if errors.Is(err, services.ErrValidation) {
			respondBadRequest(w, "Invalid request")
			return
		}
		if errors.Is(err, services.ErrConflict) {
			respondBadRequest(w, "Agreement already accepted")
			return
		}
		respondInternalError(w, "Failed to record agreement acceptance")
		return
	}

	respondOK(w, response)
}

// validateUpdateUserRequest validates the UpdateUserRequest
func validateUpdateUserRequest(req *models.UpdateUserRequest) error {
	if req.DisplayName == "" {
		return errors.New("displayName is required")
	}
	if len(req.DisplayName) > 100 {
		return errors.New("displayName must be 100 characters or less")
	}
	return nil
}
