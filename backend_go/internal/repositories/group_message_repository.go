package repositories

import (
	"context"

	"menunderfire/internal/models"
)

// GroupMessageRepository defines the interface for group message data access
type GroupMessageRepository interface {
	// FindByID retrieves a message by its ID
	FindByID(ctx context.Context, messageID string) (*models.GroupMessage, error)

	// FindByGroupID retrieves all messages for a group
	FindByGroupID(ctx context.Context, groupID string) ([]*models.GroupMessage, error)

	// FindFormMessagesByGroupID retrieves all messages with a form_id for a group
	FindFormMessagesByGroupID(ctx context.Context, groupID string) ([]*models.GroupMessage, error)

	// Create creates a new group message
	Create(ctx context.Context, groupID, senderID, content string, notifyMembers bool, formID *string) (*models.GroupMessage, error)

	// Delete deletes a group message
	Delete(ctx context.Context, messageID string) error
}
