package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"menunderfire/internal/logger"
	"menunderfire/internal/models"
	"menunderfire/internal/services"

	"github.com/gorilla/mux"
)

// InvitationHandler handles HTTP requests for invitation-related endpoints
type InvitationHandler struct {
	invitationService services.InvitationService
}

// NewInvitationHandler creates a new InvitationHandler with the given service
func NewInvitationHandler(invitationService services.InvitationService) *InvitationHandler {
	return &InvitationHandler{
		invitationService: invitationService,
	}
}

// CreateOrgAdmin handles POST /api/invitations/org-admin
// Creates an invitation for an organization admin (Site Admin only)
func (h *InvitationHandler) CreateOrgAdmin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	var req models.CreateOrgAdminInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondBadRequest(w, "Invalid request body")
		return
	}

	if err := validateCreateOrgAdminInvitationRequest(&req); err != nil {
		respondBadRequest(w, err.Error())
		return
	}

	invitation, err := h.invitationService.CreateOrgAdmin(ctx, user.ID, &req)
	if err != nil {
		logger.Error().Err(err).Msg("Error creating org admin invitation")
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "Site admin access required")
			return
		}
		if errors.Is(err, services.ErrOrganizationNotFound) {
			respondNotFound(w, "Organization not found")
			return
		}
		if errors.Is(err, services.ErrConflict) {
			respondBadRequest(w, "An invitation for this email already exists")
			return
		}
		respondInternalError(w, "Failed to create invitation")
		return
	}

	respondCreated(w, invitation.ToResponse())
}

// CreateGroupOwner handles POST /api/invitations/group-owner
// Creates an invitation for a group owner (Site Admin or Org Admin)
func (h *InvitationHandler) CreateGroupOwner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	var req models.CreateGroupOwnerInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondBadRequest(w, "Invalid request body")
		return
	}

	if err := validateCreateGroupOwnerInvitationRequest(&req); err != nil {
		respondBadRequest(w, err.Error())
		return
	}

	invitation, err := h.invitationService.CreateGroupOwner(ctx, user.ID, &req)
	if err != nil {
		logger.Error().Err(err).Msg("Error creating group owner invitation")
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "Site admin or organization admin access required")
			return
		}
		if errors.Is(err, services.ErrGroupNotFound) {
			respondNotFound(w, "Group not found")
			return
		}
		if errors.Is(err, services.ErrConflict) {
			respondBadRequest(w, "An invitation for this email already exists")
			return
		}
		respondInternalError(w, "Failed to create invitation")
		return
	}

	respondCreated(w, invitation.ToResponse())
}

// CreateGroupMember handles POST /api/invitations/group-member
// Creates an invitation for a group member (Site Admin, Org Admin, or Group Owner)
func (h *InvitationHandler) CreateGroupMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	var req models.CreateGroupMemberInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondBadRequest(w, "Invalid request body")
		return
	}

	if err := validateCreateGroupMemberInvitationRequest(&req); err != nil {
		respondBadRequest(w, err.Error())
		return
	}

	invitation, err := h.invitationService.CreateGroupMember(ctx, user.ID, &req)
	if err != nil {
		logger.Error().Err(err).Msg("Error creating group member invitation")
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "You do not have permission to invite members to this group")
			return
		}
		if errors.Is(err, services.ErrGroupNotFound) {
			respondNotFound(w, "Group not found")
			return
		}
		if errors.Is(err, services.ErrConflict) {
			respondBadRequest(w, "An invitation for this email already exists")
			return
		}
		respondInternalError(w, "Failed to create invitation")
		return
	}

	respondCreated(w, invitation.ToResponse())
}

// List handles GET /api/invitations
// Lists invitations visible to the requesting user based on their role
func (h *InvitationHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	invitations, err := h.invitationService.List(ctx, user.ID)
	if err != nil {
		logger.Error().Err(err).Msg("Error listing invitations")
		respondInternalError(w, "Failed to list invitations")
		return
	}

	respondOK(w, invitations)
}

// Delete handles DELETE /api/invitations/{id}
// Deletes an invitation
func (h *InvitationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	invitationID := vars["id"]

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	err = h.invitationService.Delete(ctx, invitationID, user.ID)
	if err != nil {
		logger.Error().Err(err).Str("invitation_id", invitationID).Msg("Error deleting invitation")
		if errors.Is(err, services.ErrInvitationNotFound) {
			respondNotFound(w, "Invitation not found")
			return
		}
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "You do not have permission to delete this invitation")
			return
		}
		respondInternalError(w, "Failed to delete invitation")
		return
	}

	respondNoContent(w)
}

// Accept handles POST /api/invitations/accept
// Accepts an invitation (Public - no auth required)
func (h *InvitationHandler) Accept(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req models.AcceptInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondBadRequest(w, "Invalid request body")
		return
	}

	if err := validateAcceptInvitationRequest(&req); err != nil {
		respondBadRequest(w, err.Error())
		return
	}

	invitation, user, err := h.invitationService.Accept(ctx, &req)
	if err != nil {
		logger.Error().Err(err).Msg("Error accepting invitation")
		if errors.Is(err, services.ErrInvitationNotFound) {
			respondNotFound(w, "Invitation not found")
			return
		}
		if errors.Is(err, services.ErrInvalidToken) {
			respondBadRequest(w, "Invalid or expired invitation token")
			return
		}
		respondInternalError(w, "Failed to accept invitation")
		return
	}

	respondOK(w, &models.AcceptInvitationResponse{
		Invitation: invitation.ToResponse(),
		User:       user.ToResponse(),
	})
}

// Validate handles GET /api/invitations/validate/{token}
// Validates an invitation token (Public - no auth required)
func (h *InvitationHandler) Validate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	token := vars["token"]

	if token == "" {
		respondBadRequest(w, "Token is required")
		return
	}

	response, err := h.invitationService.Validate(ctx, token)
	if err != nil {
		logger.Error().Err(err).Msg("Error validating invitation token")
		// Return invalid response instead of error for invalid tokens
		respondOK(w, &models.ValidateInvitationResponse{Valid: false})
		return
	}

	respondOK(w, response)
}

// validateCreateOrgAdminInvitationRequest validates the request
func validateCreateOrgAdminInvitationRequest(req *models.CreateOrgAdminInvitationRequest) error {
	if req.Email == "" {
		return errors.New("email is required")
	}
	if !isValidEmail(req.Email) {
		return errors.New("invalid email format")
	}
	if req.OrganizationID == "" {
		return errors.New("organizationId is required")
	}
	return nil
}

// validateCreateGroupOwnerInvitationRequest validates the request
func validateCreateGroupOwnerInvitationRequest(req *models.CreateGroupOwnerInvitationRequest) error {
	if req.Email == "" {
		return errors.New("email is required")
	}
	if !isValidEmail(req.Email) {
		return errors.New("invalid email format")
	}
	if req.GroupID == "" {
		return errors.New("groupId is required")
	}
	return nil
}

// validateCreateGroupMemberInvitationRequest validates the request
func validateCreateGroupMemberInvitationRequest(req *models.CreateGroupMemberInvitationRequest) error {
	if req.Email == "" {
		return errors.New("email is required")
	}
	if !isValidEmail(req.Email) {
		return errors.New("invalid email format")
	}
	if req.GroupID == "" {
		return errors.New("groupId is required")
	}
	return nil
}

// validateAcceptInvitationRequest validates the request
func validateAcceptInvitationRequest(req *models.AcceptInvitationRequest) error {
	if req.Token == "" {
		return errors.New("token is required")
	}
	if req.ExternalID == "" {
		return errors.New("externalId is required")
	}
	if req.DisplayName == "" {
		return errors.New("displayName is required")
	}
	if len(req.DisplayName) > 100 {
		return errors.New("displayName must be 100 characters or less")
	}
	return nil
}

// isValidEmail performs basic email validation
func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}
