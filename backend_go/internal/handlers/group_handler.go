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

// GroupHandler handles HTTP requests for group-related endpoints
type GroupHandler struct {
	groupService services.GroupService
}

// NewGroupHandler creates a new GroupHandler with the given service
func NewGroupHandler(groupService services.GroupService) *GroupHandler {
	return &GroupHandler{
		groupService: groupService,
	}
}

// Create handles POST /api/groups
// Creates a new group with the current user as owner
func (h *GroupHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	var req models.CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondBadRequest(w, "Invalid request body")
		return
	}

	if err := validateCreateGroupRequest(&req); err != nil {
		respondBadRequest(w, err.Error())
		return
	}

	group, err := h.groupService.Create(ctx, user.ID, &req)
	if err != nil {
		logger.Error().Err(err).Msg("Error creating group")
		if errors.Is(err, services.ErrOrganizationNotFound) {
			respondNotFound(w, "Organization not found")
			return
		}
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "You do not have permission to create groups in this organization")
			return
		}
		respondInternalError(w, "Failed to create group")
		return
	}

	role := "OWNER"
	memberCount := 1
	respondCreated(w, group.ToResponse(&role, &memberCount, true))
}

// List handles GET /api/groups
// Returns all groups the current user belongs to
func (h *GroupHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	groups, err := h.groupService.List(ctx, user.ID)
	if err != nil {
		logger.Error().Err(err).Str("user_id", user.ID).Msg("Error listing groups")
		respondInternalError(w, "Failed to list groups")
		return
	}

	responses := make([]*models.GroupResponse, len(groups))
	for i, group := range groups {
		role, _ := h.groupService.GetUserRole(ctx, group.ID, user.ID)
		memberCount, _ := h.groupService.GetMemberCount(ctx, group.ID)
		// If user has no direct role but the group was returned (org admin access), show ADMIN role
		if role == "" {
			role = "ADMIN"
		}
		showInviteCode := role == "OWNER" || role == "LEADER" || role == "ADMIN"
		responses[i] = group.ToResponse(&role, &memberCount, showInviteCode)
	}

	respondOK(w, responses)
}

// Get handles GET /api/groups/{id}
// Returns detailed information about a specific group
func (h *GroupHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	groupID := vars["id"]

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	detail, err := h.groupService.GetDetailByID(ctx, groupID, user.ID)
	if err != nil {
		logger.Error().Err(err).Str("group_id", groupID).Msg("Error getting group")
		if errors.Is(err, services.ErrGroupNotFound) {
			respondNotFound(w, "Group not found")
			return
		}
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "You are not a member of this group")
			return
		}
		respondInternalError(w, "Failed to get group")
		return
	}

	respondOK(w, detail)
}

// Join handles POST /api/groups/{id}/join
// Joins a group using an invite code
func (h *GroupHandler) Join(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	groupID := vars["id"]

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	var req models.JoinGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondBadRequest(w, "Invalid request body")
		return
	}

	if req.InviteCode == "" {
		respondBadRequest(w, "inviteCode is required")
		return
	}

	err = h.groupService.Join(ctx, groupID, user.ID, req.InviteCode)
	if err != nil {
		logger.Error().Err(err).Str("group_id", groupID).Msg("Error joining group")
		if errors.Is(err, services.ErrGroupNotFound) {
			respondNotFound(w, "Group not found")
			return
		}
		if errors.Is(err, services.ErrInvalidInviteCode) {
			respondBadRequest(w, "Invalid invite code")
			return
		}
		if errors.Is(err, services.ErrConflict) {
			respondBadRequest(w, "You are already a member of this group")
			return
		}
		respondInternalError(w, "Failed to join group")
		return
	}

	respondOK(w, MessageResponse{Message: "Successfully joined the group"})
}

// JoinByCode handles POST /api/groups/join
// Joins a group using only an invite code (no group ID required)
func (h *GroupHandler) JoinByCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	var req models.JoinGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondBadRequest(w, "Invalid request body")
		return
	}

	if req.InviteCode == "" {
		respondBadRequest(w, "inviteCode is required")
		return
	}

	err = h.groupService.JoinByInviteCode(ctx, user.ID, req.InviteCode)
	if err != nil {
		logger.Error().Err(err).Str("invite_code", req.InviteCode).Msg("Error joining group by invite code")
		if errors.Is(err, services.ErrGroupNotFound) {
			respondNotFound(w, "Group not found")
			return
		}
		if errors.Is(err, services.ErrInvalidInviteCode) {
			respondBadRequest(w, "Invalid invite code")
			return
		}
		if errors.Is(err, services.ErrConflict) {
			respondBadRequest(w, "You are already a member of this group")
			return
		}
		respondInternalError(w, "Failed to join group")
		return
	}

	respondOK(w, MessageResponse{Message: "Successfully joined the group"})
}

// RemoveMember handles DELETE /api/groups/{id}/members/{userId}
// Removes a member from a group
func (h *GroupHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	groupID := vars["id"]
	targetUserID := vars["userId"]

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	err = h.groupService.RemoveMember(ctx, groupID, targetUserID, user.ID)
	if err != nil {
		logger.Error().Err(err).Str("target_user_id", targetUserID).Str("group_id", groupID).Msg("Error removing member")
		if errors.Is(err, services.ErrGroupNotFound) {
			respondNotFound(w, "Group not found")
			return
		}
		if errors.Is(err, services.ErrUserNotFound) {
			respondNotFound(w, "User not found in group")
			return
		}
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "You do not have permission to remove this member")
			return
		}
		respondInternalError(w, "Failed to remove member")
		return
	}

	respondNoContent(w)
}

// UpdateSettings handles PUT /api/groups/{id}/settings
// Updates group settings (requirePostApproval, allowAnonymousPosts)
func (h *GroupHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	groupID := vars["id"]

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	var req models.UpdateGroupSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondBadRequest(w, "Invalid request body")
		return
	}

	err = h.groupService.UpdateSettings(ctx, groupID, user.ID, &req)
	if err != nil {
		logger.Error().Err(err).Str("group_id", groupID).Msg("Error updating group settings")
		if errors.Is(err, services.ErrGroupNotFound) {
			respondNotFound(w, "Group not found")
			return
		}
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "You do not have permission to update group settings")
			return
		}
		respondInternalError(w, "Failed to update group settings")
		return
	}

	respondOK(w, MessageResponse{Message: "Settings updated"})
}

// UpdateMemberRole handles PUT /api/groups/{id}/members/{userId}/role
// Changes a member's role within a group
func (h *GroupHandler) UpdateMemberRole(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	groupID := vars["id"]
	targetUserID := vars["userId"]

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	var body struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondBadRequest(w, "Invalid request body")
		return
	}

	if body.Role == "" {
		respondBadRequest(w, "role is required")
		return
	}

	err = h.groupService.UpdateMemberRole(ctx, groupID, targetUserID, body.Role, user.ID)
	if err != nil {
		logger.Error().Err(err).Str("group_id", groupID).Str("target_user_id", targetUserID).Msg("Error updating member role")
		if errors.Is(err, services.ErrGroupNotFound) {
			respondNotFound(w, "Group not found")
			return
		}
		if errors.Is(err, services.ErrNotFound) {
			respondNotFound(w, "User not found in group")
			return
		}
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "You do not have permission to change this member's role")
			return
		}
		if errors.Is(err, services.ErrValidation) {
			respondBadRequest(w, err.Error())
			return
		}
		respondInternalError(w, "Failed to update member role")
		return
	}

	respondOK(w, MessageResponse{Message: "Role updated"})
}

// validateCreateGroupRequest validates the CreateGroupRequest
func validateCreateGroupRequest(req *models.CreateGroupRequest) error {
	if req.Name == "" {
		return errors.New("name is required")
	}
	if len(req.Name) > 100 {
		return errors.New("name must be 100 characters or less")
	}
	if req.OrganizationID == "" {
		return errors.New("organizationId is required")
	}
	return nil
}
