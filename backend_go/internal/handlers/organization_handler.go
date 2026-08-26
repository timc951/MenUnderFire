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

// OrganizationHandler handles HTTP requests for organization-related endpoints
type OrganizationHandler struct {
	orgService  services.OrganizationService
	userService services.UserService
}

// NewOrganizationHandler creates a new OrganizationHandler with the given services
func NewOrganizationHandler(orgService services.OrganizationService, userService services.UserService) *OrganizationHandler {
	return &OrganizationHandler{
		orgService:  orgService,
		userService: userService,
	}
}

// List handles GET /api/organizations
// Returns all organizations the current user belongs to (or all organizations if site admin)
func (h *OrganizationHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := GetUserFromContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get user from context")
		respondUnauthorized(w, "Authentication required")
		return
	}

	logger.Debug().Str("user_id", user.ID).Str("email", user.Email).Msg("Listing organizations for user")

	// Check if user is site admin - if so, return all organizations
	permissions, err := h.userService.GetPermissions(ctx, user.ID)
	if err != nil {
		logger.Error().Err(err).Str("user_id", user.ID).Msg("Error getting permissions")
		respondInternalError(w, "Failed to verify permissions")
		return
	}

	var orgs []*models.Organization
	if permissions.IsSiteAdmin {
		logger.Debug().Msg("User is site admin, returning all organizations")
		orgs, err = h.orgService.ListAll(ctx)
	} else {
		logger.Debug().Msg("User is not site admin, returning organizations user is admin of")
		orgs, err = h.orgService.List(ctx, user.ID)
	}

	if err != nil {
		logger.Error().Err(err).Str("user_id", user.ID).Msg("Error listing organizations")
		respondInternalError(w, "Failed to list organizations")
		return
	}

	logger.Debug().Int("count", len(orgs)).Msg("Found organizations")

	responses := make([]*models.OrganizationResponse, len(orgs))
	for i, org := range orgs {
		responses[i] = org.ToResponse()
	}

	respondOK(w, responses)
}

// ListAll handles GET /api/organizations/all
// Returns all organizations (Site Admin only)
func (h *OrganizationHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	// Check if user is site admin
	permissions, err := h.userService.GetPermissions(ctx, user.ID)
	if err != nil {
		logger.Error().Err(err).Str("user_id", user.ID).Msg("Error getting permissions")
		respondInternalError(w, "Failed to verify permissions")
		return
	}

	if !permissions.IsSiteAdmin {
		respondForbidden(w, "Site admin access required")
		return
	}

	orgs, err := h.orgService.ListAll(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Error listing all organizations")
		respondInternalError(w, "Failed to list organizations")
		return
	}

	responses := make([]*models.OrganizationResponse, len(orgs))
	for i, org := range orgs {
		responses[i] = org.ToResponse()
	}

	respondOK(w, responses)
}

// Create handles POST /api/organizations
// Creates a new organization (Site Admin only)
func (h *OrganizationHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	// Check if user is site admin
	permissions, err := h.userService.GetPermissions(ctx, user.ID)
	if err != nil {
		logger.Error().Err(err).Str("user_id", user.ID).Msg("Error getting permissions")
		respondInternalError(w, "Failed to verify permissions")
		return
	}

	if !permissions.IsSiteAdmin {
		respondForbidden(w, "Site admin access required")
		return
	}

	var req models.CreateOrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondBadRequest(w, "Invalid request body")
		return
	}

	if err := validateCreateOrganizationRequest(&req); err != nil {
		respondBadRequest(w, err.Error())
		return
	}

	org, err := h.orgService.Create(ctx, user.ID, &req)
	if err != nil {
		logger.Error().Err(err).Msg("Error creating organization")
		if errors.Is(err, services.ErrConflict) {
			respondBadRequest(w, "An organization with this name already exists")
			return
		}
		respondInternalError(w, "Failed to create organization")
		return
	}

	respondCreated(w, org.ToResponse())
}

// Get handles GET /api/organizations/{id}
// Returns detailed information about an organization
func (h *OrganizationHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	orgID := vars["id"]

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	org, err := h.orgService.GetByID(ctx, orgID)
	if err != nil {
		logger.Error().Err(err).Str("org_id", orgID).Msg("Error getting organization")
		if errors.Is(err, services.ErrOrganizationNotFound) {
			respondNotFound(w, "Organization not found")
			return
		}
		respondInternalError(w, "Failed to get organization")
		return
	}

	// Check user's permissions
	permissions, err := h.userService.GetPermissions(ctx, user.ID)
	if err != nil {
		logger.Error().Err(err).Str("user_id", user.ID).Msg("Error getting permissions")
		respondInternalError(w, "Failed to verify permissions")
		return
	}

	isSiteAdmin := permissions.IsSiteAdmin
	isOrgAdmin, _ := h.orgService.IsAdmin(ctx, orgID, user.ID)
	isMember, _ := h.orgService.IsMember(ctx, orgID, user.ID)

	// User must be a member, org admin, or site admin to view
	if !isSiteAdmin && !isOrgAdmin && !isMember {
		respondForbidden(w, "You do not have access to this organization")
		return
	}

	canEdit := isSiteAdmin || isOrgAdmin
	respondOK(w, org.ToDetailResponse(canEdit, isSiteAdmin))
}

// Update handles PUT /api/organizations/{id}
// Updates an organization (Admin only)
func (h *OrganizationHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	orgID := vars["id"]

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	var req models.UpdateOrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondBadRequest(w, "Invalid request body")
		return
	}

	if err := validateUpdateOrganizationRequest(&req); err != nil {
		respondBadRequest(w, err.Error())
		return
	}

	org, err := h.orgService.Update(ctx, orgID, user.ID, &req)
	if err != nil {
		logger.Error().Err(err).Str("org_id", orgID).Msg("Error updating organization")
		if errors.Is(err, services.ErrOrganizationNotFound) {
			respondNotFound(w, "Organization not found")
			return
		}
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "You do not have permission to update this organization")
			return
		}
		respondInternalError(w, "Failed to update organization")
		return
	}

	respondOK(w, org.ToResponse())
}

// ListAdmins handles GET /api/organizations/{id}/admins
// Returns all admins of an organization
func (h *OrganizationHandler) ListAdmins(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	orgID := vars["id"]

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	// Verify user has access to view admins
	permissions, err := h.userService.GetPermissions(ctx, user.ID)
	if err != nil {
		logger.Error().Err(err).Str("user_id", user.ID).Msg("Error getting permissions")
		respondInternalError(w, "Failed to verify permissions")
		return
	}

	isSiteAdmin := permissions.IsSiteAdmin
	isOrgAdmin, _ := h.orgService.IsAdmin(ctx, orgID, user.ID)
	isMember, _ := h.orgService.IsMember(ctx, orgID, user.ID)

	if !isSiteAdmin && !isOrgAdmin && !isMember {
		respondForbidden(w, "You do not have access to this organization")
		return
	}

	admins, err := h.orgService.ListAdmins(ctx, orgID)
	if err != nil {
		logger.Error().Err(err).Str("org_id", orgID).Msg("Error listing admins")
		if errors.Is(err, services.ErrOrganizationNotFound) {
			respondNotFound(w, "Organization not found")
			return
		}
		respondInternalError(w, "Failed to list admins")
		return
	}

	responses := make([]*models.OrganizationAdminResponse, len(admins))
	for i, admin := range admins {
		responses[i] = admin.ToAdminResponse()
	}

	respondOK(w, responses)
}

// ListGroups handles GET /api/organizations/{id}/groups
// Returns all groups within an organization
func (h *OrganizationHandler) ListGroups(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	orgID := vars["id"]

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	// Verify user has access
	permissions, err := h.userService.GetPermissions(ctx, user.ID)
	if err != nil {
		logger.Error().Err(err).Str("user_id", user.ID).Msg("Error getting permissions")
		respondInternalError(w, "Failed to verify permissions")
		return
	}

	isSiteAdmin := permissions.IsSiteAdmin
	isOrgAdmin, _ := h.orgService.IsAdmin(ctx, orgID, user.ID)
	isMember, _ := h.orgService.IsMember(ctx, orgID, user.ID)

	if !isSiteAdmin && !isOrgAdmin && !isMember {
		respondForbidden(w, "You do not have access to this organization")
		return
	}

	groups, err := h.orgService.ListGroups(ctx, orgID)
	if err != nil {
		logger.Error().Err(err).Str("org_id", orgID).Msg("Error listing groups")
		if errors.Is(err, services.ErrOrganizationNotFound) {
			respondNotFound(w, "Organization not found")
			return
		}
		respondInternalError(w, "Failed to list groups")
		return
	}

	responses := make([]*models.GroupResponse, len(groups))
	for i, group := range groups {
		// Don't show invite codes or member counts in org group list
		responses[i] = group.ToResponse(nil, nil, false)
	}

	respondOK(w, responses)
}

// validateCreateOrganizationRequest validates the CreateOrganizationRequest
func validateCreateOrganizationRequest(req *models.CreateOrganizationRequest) error {
	if req.Name == "" {
		return errors.New("name is required")
	}
	if len(req.Name) > 100 {
		return errors.New("name must be 100 characters or less")
	}
	return nil
}

// validateUpdateOrganizationRequest validates the UpdateOrganizationRequest
func validateUpdateOrganizationRequest(req *models.UpdateOrganizationRequest) error {
	if req.Name == "" {
		return errors.New("name is required")
	}
	if len(req.Name) > 100 {
		return errors.New("name must be 100 characters or less")
	}
	return nil
}
