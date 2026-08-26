package repositories

import (
	"context"

	"menunderfire/internal/models"
)

// GroupRepository defines the interface for group data access
type GroupRepository interface {
	// FindByID retrieves a group by its ID
	FindByID(ctx context.Context, groupID string) (*models.Group, error)

	// FindByInviteCode retrieves a group by its invite code
	FindByInviteCode(ctx context.Context, inviteCode string) (*models.Group, error)

	// FindByUserID retrieves all groups a user belongs to
	FindByUserID(ctx context.Context, userID string) ([]*models.Group, error)

	// FindByOrganizationID retrieves all groups for an organization
	FindByOrganizationID(ctx context.Context, orgID string) ([]*models.Group, error)

	// Create creates a new group
	Create(ctx context.Context, name string, description *string, orgID string, inviteCode string, createdBy string) (*models.Group, error)

	// GenerateInviteCode generates a unique invite code for a group
	GenerateInviteCode() string

	// FindMember retrieves a group member by group and user ID
	FindMember(ctx context.Context, groupID, userID string) (*models.GroupMember, error)

	// FindMembers retrieves all members of a group
	FindMembers(ctx context.Context, groupID string) ([]*models.GroupMember, error)

	// CountMembers returns the number of members in a group
	CountMembers(ctx context.Context, groupID string) (int, error)

	// AddMember adds a member to a group with a specific role
	AddMember(ctx context.Context, groupID, userID, role string) (*models.GroupMember, error)

	// RemoveMember removes a member from a group
	RemoveMember(ctx context.Context, groupID, userID string) error

	// UpdateSettings updates the group's settings
	UpdateSettings(ctx context.Context, groupID string, requirePostApproval, allowAnonymousPosts bool) error

	// UpdateMemberRole updates a member's role within a group
	UpdateMemberRole(ctx context.Context, groupID, userID, role string) error

	// Count returns the total number of groups
	Count(ctx context.Context) (int64, error)

	// CountByOrganizationIDs returns the number of groups within specific organizations
	CountByOrganizationIDs(ctx context.Context, orgIDs []string) (int64, error)
}
