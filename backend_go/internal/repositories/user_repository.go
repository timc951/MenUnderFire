package repositories

import (
	"context"

	"menunderfire/internal/models"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// FindByID retrieves a user by their ID
	FindByID(ctx context.Context, userID string) (*models.User, error)

	// FindByExternalID retrieves a user by their external authentication provider ID
	FindByExternalID(ctx context.Context, externalID string) (*models.User, error)

	// FindByEmail retrieves a user by their email address
	FindByEmail(ctx context.Context, email string) (*models.User, error)

	// Count returns the total number of users
	Count(ctx context.Context) (int64, error)

	// Create creates a new user
	Create(ctx context.Context, email, displayName, externalID string) (*models.User, error)

	// CreateAsSiteAdmin creates a new user with site admin privileges
	CreateAsSiteAdmin(ctx context.Context, email, displayName, externalID string) (*models.User, error)

	// Update updates a user's display name
	Update(ctx context.Context, userID, displayName string) (*models.User, error)

	// UpdateExternalID updates a user's external authentication provider ID
	UpdateExternalID(ctx context.Context, userID, externalID string) error

	// UpdateInvitationInfo sets the invited_by_id and invitation_id on a user
	UpdateInvitationInfo(ctx context.Context, userID, invitedByID, invitationID string) error

	// RecordAgreementAcceptance stores the user's agreement acceptance with signature
	RecordAgreementAcceptance(ctx context.Context, userID, version, signature, ipAddress, userAgent string) error

	// IsSiteAdmin checks if a user has site admin privileges
	IsSiteAdmin(ctx context.Context, userID string) (bool, error)

	// FindAdminOrganizationIDs returns the IDs of organizations where the user is an admin
	FindAdminOrganizationIDs(ctx context.Context, userID string) ([]string, error)

	// FindOwnedGroupIDs returns the IDs of groups where the user is an owner
	FindOwnedGroupIDs(ctx context.Context, userID string) ([]string, error)

	// FindMemberGroupIDs returns the IDs of groups where the user is a member
	FindMemberGroupIDs(ctx context.Context, userID string) ([]string, error)
}
