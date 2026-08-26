package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"menunderfire/internal/logger"
	"menunderfire/internal/models"
	"menunderfire/internal/services"

	"github.com/gorilla/mux"
)

// SitePageHandler handles HTTP requests for site page endpoints
type SitePageHandler struct {
	pageService services.SitePageService
	userService services.UserService
}

// NewSitePageHandler creates a new SitePageHandler with the given services
func NewSitePageHandler(pageService services.SitePageService, userService services.UserService) *SitePageHandler {
	return &SitePageHandler{
		pageService: pageService,
		userService: userService,
	}
}

// List handles GET /api/pages
// Returns all site pages (Public - no auth required)
func (h *SitePageHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pages, err := h.pageService.List(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Error listing site pages")
		respondInternalError(w, "Failed to list pages")
		return
	}

	responses := make([]*models.SitePageSummaryResponse, len(pages))
	for i, page := range pages {
		responses[i] = page.ToSummaryResponse()
	}

	respondOK(w, responses)
}

// GetBySlug handles GET /api/pages/{slug}
// Returns a site page by slug (Public - no auth required)
func (h *SitePageHandler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	slug := vars["slug"]

	page, err := h.pageService.GetBySlug(ctx, slug)
	if err != nil {
		logger.Error().Err(err).Str("slug", slug).Msg("Error getting page by slug")
		if errors.Is(err, services.ErrPageNotFound) {
			respondNotFound(w, "Page not found")
			return
		}
		respondInternalError(w, "Failed to get page")
		return
	}

	respondOK(w, page.ToResponse())
}

// Create handles POST /api/pages
// Creates a new site page (Site Admin only)
func (h *SitePageHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	// Check if user is site admin
	permissions, err := h.userService.GetPermissions(ctx, user.ID)
	if err != nil {
		logger.Error().Err(err).Str("user_id", user.ID).Msg("Error getting permissions for user")
		respondInternalError(w, "Failed to verify permissions")
		return
	}

	if !permissions.IsSiteAdmin {
		respondForbidden(w, "Site admin access required")
		return
	}

	var req models.CreateSitePageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondBadRequest(w, "Invalid request body")
		return
	}

	if err := validateCreateSitePageRequest(&req); err != nil {
		respondBadRequest(w, err.Error())
		return
	}

	page, err := h.pageService.Create(ctx, user.ID, &req)
	if err != nil {
		logger.Error().Err(err).Msg("Error creating site page")
		if errors.Is(err, services.ErrConflict) {
			respondBadRequest(w, "A page with this slug already exists")
			return
		}
		respondInternalError(w, "Failed to create page")
		return
	}

	respondCreated(w, page.ToResponse())
}

// Update handles PUT /api/pages/{id}
// Updates a site page (Site Admin only)
func (h *SitePageHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	pageID := vars["id"]

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	// Check if user is site admin
	permissions, err := h.userService.GetPermissions(ctx, user.ID)
	if err != nil {
		logger.Error().Err(err).Str("user_id", user.ID).Msg("Error getting permissions for user")
		respondInternalError(w, "Failed to verify permissions")
		return
	}

	if !permissions.IsSiteAdmin {
		respondForbidden(w, "Site admin access required")
		return
	}

	var req models.UpdateSitePageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondBadRequest(w, "Invalid request body")
		return
	}

	if err := validateUpdateSitePageRequest(&req); err != nil {
		respondBadRequest(w, err.Error())
		return
	}

	page, err := h.pageService.Update(ctx, pageID, user.ID, &req)
	if err != nil {
		logger.Error().Err(err).Str("page_id", pageID).Msg("Error updating site page")
		if errors.Is(err, services.ErrPageNotFound) {
			respondNotFound(w, "Page not found")
			return
		}
		respondInternalError(w, "Failed to update page")
		return
	}

	respondOK(w, page.ToResponse())
}

// Delete handles DELETE /api/pages/{id}
// Deletes a site page (Site Admin only)
func (h *SitePageHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	pageID := vars["id"]

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	// Check if user is site admin
	permissions, err := h.userService.GetPermissions(ctx, user.ID)
	if err != nil {
		logger.Error().Err(err).Str("user_id", user.ID).Msg("Error getting permissions for user")
		respondInternalError(w, "Failed to verify permissions")
		return
	}

	if !permissions.IsSiteAdmin {
		respondForbidden(w, "Site admin access required")
		return
	}

	err = h.pageService.Delete(ctx, pageID, user.ID)
	if err != nil {
		logger.Error().Err(err).Str("page_id", pageID).Msg("Error deleting site page")
		if errors.Is(err, services.ErrPageNotFound) {
			respondNotFound(w, "Page not found")
			return
		}
		respondInternalError(w, "Failed to delete page")
		return
	}

	respondNoContent(w)
}

// Validation helpers

var slugRegex = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

func validateCreateSitePageRequest(req *models.CreateSitePageRequest) error {
	if req.Slug == "" {
		return errors.New("slug is required")
	}
	if len(req.Slug) > 100 {
		return errors.New("slug must be 100 characters or less")
	}
	if !slugRegex.MatchString(req.Slug) {
		return errors.New("slug must contain only lowercase letters, numbers, and hyphens")
	}
	if req.Title == "" {
		return errors.New("title is required")
	}
	if len(req.Title) > 200 {
		return errors.New("title must be 200 characters or less")
	}
	return nil
}

func validateUpdateSitePageRequest(req *models.UpdateSitePageRequest) error {
	if req.Title == "" {
		return errors.New("title is required")
	}
	if len(req.Title) > 200 {
		return errors.New("title must be 200 characters or less")
	}
	return nil
}

// CreateDraft handles POST /api/pages/{id}/preview
// Creates a temporary preview draft (Site Admin only)
func (h *SitePageHandler) CreateDraft(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	pageID := vars["id"]

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	// Check if user is site admin
	permissions, err := h.userService.GetPermissions(ctx, user.ID)
	if err != nil {
		logger.Error().Err(err).Str("user_id", user.ID).Msg("Error getting permissions for user")
		respondInternalError(w, "Failed to verify permissions")
		return
	}

	if !permissions.IsSiteAdmin {
		respondForbidden(w, "Site admin access required")
		return
	}

	var req models.CreatePageDraftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondBadRequest(w, "Invalid request body")
		return
	}

	if req.Title == "" {
		respondBadRequest(w, "Title is required")
		return
	}

	// Pass pageID as pointer (nil for new pages)
	var pageIDPtr *string
	if pageID != "" && pageID != "new" {
		pageIDPtr = &pageID
	}

	draft, err := h.pageService.CreateDraft(ctx, pageIDPtr, user.ID, &req)
	if err != nil {
		logger.Error().Err(err).Msg("Error creating page draft")
		respondInternalError(w, "Failed to create preview draft")
		return
	}

	respondCreated(w, draft.ToResponse())
}

// GetDraft handles GET /api/pages/preview/{draftId}
// Returns a draft for preview (Public - no auth required)
func (h *SitePageHandler) GetDraft(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	draftID := vars["draftId"]

	draft, err := h.pageService.GetDraft(ctx, draftID)
	if err != nil {
		logger.Error().Err(err).Str("draft_id", draftID).Msg("Error getting draft")
		if errors.Is(err, services.ErrDraftNotFound) {
			respondNotFound(w, "Preview not found or expired")
			return
		}
		respondInternalError(w, "Failed to get preview")
		return
	}

	respondOK(w, draft.ToContentResponse())
}
