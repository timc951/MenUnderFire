package repositories

import (
	"context"

	"menunderfire/internal/models"
)

// PageDraftRepository defines the interface for page draft data access
type PageDraftRepository interface {
	// Create creates a new page draft
	Create(ctx context.Context, pageID *string, title, content, createdByID string) (*models.PageDraft, error)

	// FindByID retrieves a page draft by its ID
	FindByID(ctx context.Context, draftID string) (*models.PageDraft, error)

	// Delete deletes a page draft by its ID
	Delete(ctx context.Context, draftID string) error

	// DeleteExpired deletes all expired drafts
	DeleteExpired(ctx context.Context) (int64, error)
}
