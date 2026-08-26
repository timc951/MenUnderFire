package repositories

import (
	"context"

	"menunderfire/internal/models"
)

// SitePageRepository defines the interface for site page data access
type SitePageRepository interface {
	// FindAll retrieves all site pages
	FindAll(ctx context.Context) ([]*models.SitePage, error)

	// FindByID retrieves a site page by its ID
	FindByID(ctx context.Context, pageID string) (*models.SitePage, error)

	// FindBySlug retrieves a site page by its slug
	FindBySlug(ctx context.Context, slug string) (*models.SitePage, error)

	// Create creates a new site page
	Create(ctx context.Context, slug, title, content string, isPublished bool, createdByID string) (*models.SitePage, error)

	// Update updates an existing site page
	Update(ctx context.Context, pageID, title, content string, isPublished bool, updatedByID string) (*models.SitePage, error)

	// Delete deletes a site page by its ID
	Delete(ctx context.Context, pageID string) error
}
