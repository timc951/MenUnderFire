package services

import (
	"context"
	"fmt"

	"menunderfire/internal/models"
	"menunderfire/internal/repositories"
)

// DashboardService defines the interface for dashboard-related business logic
type DashboardService interface {
	// GetStats returns dashboard statistics based on the user's role
	// Site Admin: sees all organizations and groups
	// Org Admin: sees only group count within their organizations
	// Regular users: see minimal stats
	GetStats(ctx context.Context, userID string) (*models.DashboardStatsResponse, error)
}

// dashboardService implements the DashboardService interface
type dashboardService struct {
	userService UserService
	orgRepo     repositories.OrganizationRepository
	groupRepo   repositories.GroupRepository
}

// NewDashboardService creates a new DashboardService implementation
//
// Parameters:
//   - userService: Service to retrieve user permissions
//   - orgRepo: Repository for organization data access
//   - groupRepo: Repository for group data access
func NewDashboardService(
	userService UserService,
	orgRepo repositories.OrganizationRepository,
	groupRepo repositories.GroupRepository,
) DashboardService {
	return &dashboardService{
		userService: userService,
		orgRepo:     orgRepo,
		groupRepo:   groupRepo,
	}
}

// GetStats retrieves dashboard statistics for a user based on their permissions
//
// Edge cases handled:
//   - User not found: Returns ErrUserNotFound
//   - User has no permissions (not admin): Returns null org count, 0 group count
//   - User is site admin (gets full stats): Returns total org count and total group count
//   - User is org admin of one organization: Returns null org count, group count for their org
//   - User is org admin of multiple organizations: Returns null org count, group count for all their orgs
//   - User is org admin with no groups in their orgs: Returns null org count, 0 group count
//   - Database error when getting permissions: Returns wrapped error
//   - Database error when counting orgs: Returns wrapped error
//   - Database error when counting groups: Returns wrapped error
//   - Empty orgIDs list for org admin: Treats as regular user (null org count, 0 group count)
func (s *dashboardService) GetStats(ctx context.Context, userID string) (*models.DashboardStatsResponse, error) {
	// Get user permissions
	permissions, err := s.userService.GetPermissions(ctx, userID)
	if err != nil {
		if err == ErrUserNotFound {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get user permissions: %w", err)
	}

	// Initialize response with default values
	response := &models.DashboardStatsResponse{
		OrganizationCount: nil,
		GroupCount:        0,
	}

	// Site Admin: sees all organizations count and all groups count
	if permissions.IsSiteAdmin {
		orgCount, err := s.orgRepo.Count(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to count organizations: %w", err)
		}

		groupCount, err := s.groupRepo.Count(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to count groups: %w", err)
		}

		response.OrganizationCount = &orgCount
		response.GroupCount = groupCount
		return response, nil
	}

	// Org Admin: sees null for org count + group count ONLY within their organizations
	if len(permissions.AdminOfOrganizationIDs) > 0 {
		groupCount, err := s.groupRepo.CountByOrganizationIDs(ctx, permissions.AdminOfOrganizationIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to count groups for organizations: %w", err)
		}

		response.OrganizationCount = nil
		response.GroupCount = groupCount
		return response, nil
	}

	// Regular users: sees null for org count + 0 for group count
	// This is already the default state of the response
	return response, nil
}
