package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"menunderfire/internal/models"
	"menunderfire/internal/repositories"
)

// sitePageRepository is the PostgreSQL implementation of SitePageRepository
type sitePageRepository struct {
	db *sql.DB
}

// NewSitePageRepository creates a new instance of SitePageRepository
func NewSitePageRepository(db *sql.DB) repositories.SitePageRepository {
	return &sitePageRepository{db: db}
}

// FindAll retrieves all site pages
func (r *sitePageRepository) FindAll(ctx context.Context) ([]*models.SitePage, error) {
	query := `
		SELECT id, slug, title, content, is_published, created_at, updated_at
		FROM site_pages
		ORDER BY updated_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error finding all site pages: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var pages []*models.SitePage
	for rows.Next() {
		var page models.SitePage
		if err := rows.Scan(
			&page.ID,
			&page.Slug,
			&page.Title,
			&page.Content,
			&page.IsPublished,
			&page.CreatedAt,
			&page.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning site page: %w", err)
		}

		pages = append(pages, &page)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating site pages: %w", err)
	}

	return pages, nil
}

// FindByID retrieves a site page by its ID
func (r *sitePageRepository) FindByID(ctx context.Context, pageID string) (*models.SitePage, error) {
	query := `
		SELECT id, slug, title, content, is_published, created_at, updated_at
		FROM site_pages
		WHERE id = $1
	`

	var page models.SitePage
	err := r.db.QueryRowContext(ctx, query, pageID).Scan(
		&page.ID,
		&page.Slug,
		&page.Title,
		&page.Content,
		&page.IsPublished,
		&page.CreatedAt,
		&page.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("site page not found: %s", pageID)
	}
	if err != nil {
		return nil, fmt.Errorf("error finding site page by ID: %w", err)
	}

	return &page, nil
}

// FindBySlug retrieves a site page by its slug
func (r *sitePageRepository) FindBySlug(ctx context.Context, slug string) (*models.SitePage, error) {
	query := `
		SELECT id, slug, title, content, is_published, created_at, updated_at
		FROM site_pages
		WHERE slug = $1
	`

	var page models.SitePage
	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&page.ID,
		&page.Slug,
		&page.Title,
		&page.Content,
		&page.IsPublished,
		&page.CreatedAt,
		&page.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("site page not found with slug: %s", slug)
	}
	if err != nil {
		return nil, fmt.Errorf("error finding site page by slug: %w", err)
	}

	return &page, nil
}

// Create creates a new site page
func (r *sitePageRepository) Create(ctx context.Context, slug, title, content string, isPublished bool, createdByID string) (*models.SitePage, error) {
	query := `
		INSERT INTO site_pages (slug, title, content, is_published, created_by_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, slug, title, content, is_published, created_at, updated_at
	`

	var page models.SitePage
	err := r.db.QueryRowContext(ctx, query, slug, title, content, isPublished, createdByID).Scan(
		&page.ID,
		&page.Slug,
		&page.Title,
		&page.Content,
		&page.IsPublished,
		&page.CreatedAt,
		&page.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error creating site page: %w", err)
	}

	return &page, nil
}

// Update updates an existing site page
func (r *sitePageRepository) Update(ctx context.Context, pageID, title, content string, isPublished bool, updatedByID string) (*models.SitePage, error) {
	query := `
		UPDATE site_pages
		SET title = $1, content = $2, is_published = $3, updated_at = NOW(), updated_by_id = $4
		WHERE id = $5
		RETURNING id, slug, title, content, is_published, created_at, updated_at
	`

	var page models.SitePage
	err := r.db.QueryRowContext(ctx, query, title, content, isPublished, updatedByID, pageID).Scan(
		&page.ID,
		&page.Slug,
		&page.Title,
		&page.Content,
		&page.IsPublished,
		&page.CreatedAt,
		&page.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("site page not found: %s", pageID)
	}
	if err != nil {
		return nil, fmt.Errorf("error updating site page: %w", err)
	}

	return &page, nil
}

// Delete deletes a site page by its ID
func (r *sitePageRepository) Delete(ctx context.Context, pageID string) error {
	query := `
		DELETE FROM site_pages
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, pageID)
	if err != nil {
		return fmt.Errorf("error deleting site page: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("site page not found: %s", pageID)
	}

	return nil
}
