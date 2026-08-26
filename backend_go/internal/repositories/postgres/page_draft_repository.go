package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"menunderfire/internal/models"
	"menunderfire/internal/repositories"
)

// pageDraftRepository is the PostgreSQL implementation of PageDraftRepository
type pageDraftRepository struct {
	db *sql.DB
}

// NewPageDraftRepository creates a new instance of PageDraftRepository
func NewPageDraftRepository(db *sql.DB) repositories.PageDraftRepository {
	return &pageDraftRepository{db: db}
}

// Create creates a new page draft with 1-hour expiration
func (r *pageDraftRepository) Create(ctx context.Context, pageID *string, title, content, createdByID string) (*models.PageDraft, error) {
	expiresAt := time.Now().Add(1 * time.Hour)

	query := `
		INSERT INTO page_drafts (page_id, title, content, created_by_id, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, page_id, title, content, created_by_id, created_at, expires_at
	`

	var draft models.PageDraft
	var nullPageID sql.NullString

	err := r.db.QueryRowContext(ctx, query, pageID, title, content, createdByID, expiresAt).Scan(
		&draft.ID,
		&nullPageID,
		&draft.Title,
		&draft.Content,
		&draft.CreatedByID,
		&draft.CreatedAt,
		&draft.ExpiresAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error creating page draft: %w", err)
	}

	if nullPageID.Valid {
		draft.PageID = &nullPageID.String
	}

	return &draft, nil
}

// FindByID retrieves a page draft by its ID (only if not expired)
func (r *pageDraftRepository) FindByID(ctx context.Context, draftID string) (*models.PageDraft, error) {
	query := `
		SELECT id, page_id, title, content, created_by_id, created_at, expires_at
		FROM page_drafts
		WHERE id = $1 AND expires_at > NOW()
	`

	var draft models.PageDraft
	var nullPageID sql.NullString

	err := r.db.QueryRowContext(ctx, query, draftID).Scan(
		&draft.ID,
		&nullPageID,
		&draft.Title,
		&draft.Content,
		&draft.CreatedByID,
		&draft.CreatedAt,
		&draft.ExpiresAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("page draft not found or expired: %s", draftID)
	}
	if err != nil {
		return nil, fmt.Errorf("error finding page draft by ID: %w", err)
	}

	if nullPageID.Valid {
		draft.PageID = &nullPageID.String
	}

	return &draft, nil
}

// Delete deletes a page draft by its ID
func (r *pageDraftRepository) Delete(ctx context.Context, draftID string) error {
	query := `DELETE FROM page_drafts WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, draftID)
	if err != nil {
		return fmt.Errorf("error deleting page draft: %w", err)
	}

	return nil
}

// DeleteExpired deletes all expired drafts
func (r *pageDraftRepository) DeleteExpired(ctx context.Context) (int64, error) {
	query := `DELETE FROM page_drafts WHERE expires_at <= NOW()`

	result, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("error deleting expired drafts: %w", err)
	}

	return result.RowsAffected()
}
