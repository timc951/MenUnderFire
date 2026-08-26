package repositories

import (
	"context"

	"menunderfire/internal/models"
)

// ReportRepository defines the interface for report data access
type ReportRepository interface {
	// FindByID retrieves a report by its ID
	FindByID(ctx context.Context, reportID string) (*models.Report, error)

	// FindByGroupID retrieves all reports for a group
	FindByGroupID(ctx context.Context, groupID string) ([]*models.Report, error)

	// Create creates a new accountability report
	Create(ctx context.Context, groupID, reporterID, title, content string, isAnonymous bool) (*models.Report, error)
}
