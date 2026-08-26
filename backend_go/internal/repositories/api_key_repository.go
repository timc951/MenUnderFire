package repositories

import (
	"context"
	"time"

	"menunderfire/internal/models"
)

// APIKeyRepository defines the interface for API key data access
type APIKeyRepository interface {
	// Create creates a new API key
	Create(ctx context.Context, userID, keyHash, name string, permissions []string, expiresAt *time.Time) (*models.APIKey, error)

	// FindByID retrieves an API key by its ID
	FindByID(ctx context.Context, keyID string) (*models.APIKey, error)

	// FindByUserID retrieves all API keys for a user
	FindByUserID(ctx context.Context, userID string) ([]*models.APIKey, error)

	// FindByKeyHash retrieves an API key by its hash
	FindByKeyHash(ctx context.Context, keyHash string) (*models.APIKey, error)

	// FindAll retrieves all API keys
	FindAll(ctx context.Context) ([]*models.APIKey, error)

	// Delete deletes an API key
	Delete(ctx context.Context, keyID string) error

	// UpdateLastUsed updates the last used timestamp for an API key
	UpdateLastUsed(ctx context.Context, keyID string) error
}
