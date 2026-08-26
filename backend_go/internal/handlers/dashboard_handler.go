package handlers

import (
	"errors"
	"net/http"

	"menunderfire/internal/logger"
	"menunderfire/internal/services"
)

// DashboardHandler handles HTTP requests for dashboard endpoints
type DashboardHandler struct {
	dashboardService services.DashboardService
	userService      services.UserService
}

// NewDashboardHandler creates a new DashboardHandler with the given services
func NewDashboardHandler(dashboardService services.DashboardService, userService services.UserService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
		userService:      userService,
	}
}

// GetStats handles GET /api/dashboard/stats
// Returns dashboard statistics (Admin only - Site Admin or Org Admin)
func (h *DashboardHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := GetUserFromContext(ctx)
	if err != nil {
		respondUnauthorized(w, "Authentication required")
		return
	}

	// Check if user is an admin (site admin or org admin)
	permissions, err := h.userService.GetPermissions(ctx, user.ID)
	if err != nil {
		logger.Error().Err(err).Str("user_id", user.ID).Msg("Error getting permissions")
		respondInternalError(w, "Failed to verify permissions")
		return
	}

	isAdmin := permissions.IsSiteAdmin || len(permissions.AdminOfOrganizationIDs) > 0
	if !isAdmin {
		respondForbidden(w, "Admin access required")
		return
	}

	stats, err := h.dashboardService.GetStats(ctx, user.ID)
	if err != nil {
		logger.Error().Err(err).Msg("Error getting dashboard stats")
		if errors.Is(err, services.ErrForbidden) {
			respondForbidden(w, "Admin access required")
			return
		}
		respondInternalError(w, "Failed to get dashboard stats")
		return
	}

	respondOK(w, stats)
}
