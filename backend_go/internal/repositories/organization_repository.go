package repositories

import (
	"context"

	"menunderfire/internal/models"
)

// OrganizationRepository defines the interface for organization data access
type OrganizationRepository interface {
	// FindByID retrieves an organization by its ID
	FindByID(ctx context.Context, orgID string) (*models.Organization, error)

	// FindByUserID retrieves all organizations a user belongs to
	FindByUserID(ctx context.Context, userID string) ([]*models.Organization, error)

	// FindAll retrieves all organizations
	FindAll(ctx context.Context) ([]*models.Organization, error)

	// Create creates a new organization
	Create(ctx context.Context, name string, description *string, createdByID string) (*models.Organization, error)

	// Update updates an organization's name and description
	Update(ctx context.Context, orgID string, name string, description *string) (*models.Organization, error)

	// FindAdmins retrieves all admins for an organization
	FindAdmins(ctx context.Context, orgID string) ([]*models.OrganizationAdmin, error)

	// FindAdmin retrieves a specific admin by organization and user ID
	FindAdmin(ctx context.Context, orgID, userID string) (*models.OrganizationAdmin, error)

	// AddAdmin adds a user as an admin of an organization
	AddAdmin(ctx context.Context, orgID, userID string) error

	// IsMember checks if a user is a member of an organization (via any group)
	IsMember(ctx context.Context, orgID, userID string) (bool, error)

	// IsAdmin checks if a user is an admin of an organization
	IsAdmin(ctx context.Context, orgID, userID string) (bool, error)

	// Count returns the total number of organizations
	Count(ctx context.Context) (int64, error)
}
