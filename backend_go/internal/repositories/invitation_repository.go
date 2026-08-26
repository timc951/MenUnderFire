package repositories

import (
	"context"
	"time"

	"menunderfire/internal/models"
)

// InvitationRepository defines the interface for invitation data access
type InvitationRepository interface {
	// FindByID retrieves an invitation by its ID
	FindByID(ctx context.Context, invitationID string) (*models.Invitation, error)

	// FindByToken retrieves an invitation by its token
	FindByToken(ctx context.Context, token string) (*models.Invitation, error)

	// FindByEmail retrieves an invitation by email, type, and target ID
	FindByEmail(ctx context.Context, email string, invType models.InvitationType, targetID string) (*models.Invitation, error)

	// Create creates a new invitation
	Create(ctx context.Context, email string, invType models.InvitationType, orgID, groupID *string, inviterID string, token string, expiresAt time.Time) (*models.Invitation, error)

	// UpdateStatus updates the status of an invitation.
	// When accepting, pass the accepting user's ID to populate accepted_at and accepted_by_id.
	UpdateStatus(ctx context.Context, invitationID string, status models.InvitationStatus, acceptedByID *string) error

	// Delete deletes an invitation
	Delete(ctx context.Context, invitationID string) error

	// ListAll retrieves all invitations (for site admin)
	ListAll(ctx context.Context) ([]*models.Invitation, error)

	// ListByOrganizationIDs retrieves invitations for given organization IDs (for org admin)
	ListByOrganizationIDs(ctx context.Context, orgIDs []string) ([]*models.Invitation, error)

	// ListByGroupIDs retrieves invitations for given group IDs (for group admin)
	ListByGroupIDs(ctx context.Context, groupIDs []string) ([]*models.Invitation, error)
}
