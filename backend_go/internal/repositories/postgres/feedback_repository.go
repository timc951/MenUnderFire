package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"menunderfire/internal/models"
	"menunderfire/internal/repositories"
)

type feedbackRepository struct {
	db *sql.DB
}

// NewFeedbackRepository creates a new instance of FeedbackRepository
func NewFeedbackRepository(db *sql.DB) repositories.FeedbackRepository {
	return &feedbackRepository{db: db}
}

func (r *feedbackRepository) Create(ctx context.Context, userID, feedbackType, subject, description string) (*models.Feedback, error) {
	query := `
		INSERT INTO feedback (user_id, type, subject, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, type, subject, description, status, created_at, updated_at
	`

	var f models.Feedback
	err := r.db.QueryRowContext(ctx, query, userID, feedbackType, subject, description).Scan(
		&f.ID, &f.UserID, &f.Type, &f.Subject, &f.Description, &f.Status, &f.CreatedAt, &f.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating feedback: %w", err)
	}

	return &f, nil
}

func (r *feedbackRepository) FindByID(ctx context.Context, feedbackID string) (*models.Feedback, error) {
	query := `
		SELECT id, user_id, type, subject, description, status, created_at, updated_at
		FROM feedback
		WHERE id = $1
	`

	var f models.Feedback
	err := r.db.QueryRowContext(ctx, query, feedbackID).Scan(
		&f.ID, &f.UserID, &f.Type, &f.Subject, &f.Description, &f.Status, &f.CreatedAt, &f.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error finding feedback by ID: %w", err)
	}

	return &f, nil
}

func (r *feedbackRepository) FindByUserID(ctx context.Context, userID string) ([]*models.Feedback, error) {
	query := `
		SELECT id, user_id, type, subject, description, status, created_at, updated_at
		FROM feedback
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error finding feedback by user ID: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var items []*models.Feedback
	for rows.Next() {
		var f models.Feedback
		if err := rows.Scan(
			&f.ID, &f.UserID, &f.Type, &f.Subject, &f.Description, &f.Status, &f.CreatedAt, &f.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning feedback: %w", err)
		}
		items = append(items, &f)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating feedback: %w", err)
	}

	return items, nil
}

func (r *feedbackRepository) FindAll(ctx context.Context) ([]*models.Feedback, error) {
	query := `
		SELECT id, user_id, type, subject, description, status, created_at, updated_at
		FROM feedback
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error finding all feedback: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var items []*models.Feedback
	for rows.Next() {
		var f models.Feedback
		if err := rows.Scan(
			&f.ID, &f.UserID, &f.Type, &f.Subject, &f.Description, &f.Status, &f.CreatedAt, &f.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning feedback: %w", err)
		}
		items = append(items, &f)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating feedback: %w", err)
	}

	return items, nil
}

func (r *feedbackRepository) UpdateStatus(ctx context.Context, feedbackID, status string) (*models.Feedback, error) {
	query := `
		UPDATE feedback
		SET status = $2, updated_at = NOW()
		WHERE id = $1
		RETURNING id, user_id, type, subject, description, status, created_at, updated_at
	`

	var f models.Feedback
	err := r.db.QueryRowContext(ctx, query, feedbackID, status).Scan(
		&f.ID, &f.UserID, &f.Type, &f.Subject, &f.Description, &f.Status, &f.CreatedAt, &f.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error updating feedback status: %w", err)
	}

	return &f, nil
}
