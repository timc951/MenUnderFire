package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"menunderfire/internal/logger"
	"menunderfire/internal/models"
	"menunderfire/internal/services"
)

// ReportHandler handles HTTP requests for report-related endpoints
type ReportHandler struct {
	reportService services.ReportService
	groupService  services.GroupService
}

// NewReportHandler creates a new ReportHandler with the given services
func NewReportHandler(reportService services.ReportService, groupService services.GroupService) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
		groupService:  groupService,
	}
}

// Create handles POST /api/reports
// Creates a new accountability report
func (h *ReportHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	var req models.CreateReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondBadRequest(w, "Invalid request body")
		return
	}

	if err := validateCreateReportRequest(&req); err != nil {
		respondBadRequest(w, err.Error())
		return
	}

	report, err := h.reportService.Create(ctx, user.ID, &req)
	if err != nil {
		logger.Error().Err(err).Msg("Error creating report")
		if errors.Is(err, services.ErrGroupNotFound) {
			respondNotFound(w, "Group not found")
			return
		}
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "You are not a member of this group")
			return
		}
		respondInternalError(w, "Failed to create report")
		return
	}

	// The author can always see their own report info
	respondCreated(w, report.ToResponse(true, true))
}

// List handles GET /api/reports
// Returns all reports for a group (requires groupId query parameter)
func (h *ReportHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	groupID := r.URL.Query().Get("groupId")
	if groupID == "" {
		respondBadRequest(w, "groupId query parameter is required")
		return
	}

	// Check user's role in the group to determine visibility
	role, err := h.groupService.GetUserRole(ctx, groupID, user.ID)
	if err != nil {
		logger.Error().Err(err).Str("group_id", groupID).Msg("Error getting user role in group")
		if errors.Is(err, services.ErrGroupNotFound) {
			respondNotFound(w, "Group not found")
			return
		}
		respondInternalError(w, "Failed to verify group membership")
		return
	}

	if role == "" {
		respondForbidden(w, "You are not a member of this group")
		return
	}

	reports, err := h.reportService.List(ctx, groupID, user.ID)
	if err != nil {
		logger.Error().Err(err).Str("group_id", groupID).Msg("Error listing reports for group")
		respondInternalError(w, "Failed to list reports")
		return
	}

	// Leaders and owners can see all author info
	isLeaderOrOwner := role == "OWNER" || role == "LEADER"

	responses := make([]*models.ReportResponse, len(reports))
	for i, report := range reports {
		isOwnReport := report.ReporterID == user.ID
		responses[i] = report.ToResponse(isLeaderOrOwner, isOwnReport)
	}

	respondOK(w, responses)
}

// validateCreateReportRequest validates the CreateReportRequest
func validateCreateReportRequest(req *models.CreateReportRequest) error {
	if req.GroupID == "" {
		return errors.New("groupId is required")
	}
	if req.Title == "" {
		return errors.New("title is required")
	}
	if len(req.Title) > 200 {
		return errors.New("title must be 200 characters or less")
	}
	if req.Content == "" {
		return errors.New("content is required")
	}
	return nil
}
