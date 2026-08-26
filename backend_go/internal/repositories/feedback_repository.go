package repositories

import (
	"context"

	"menunderfire/internal/models"
)

// FeedbackRepository defines the interface for feedback data access
type FeedbackRepository interface {
	// Create creates a new feedback entry
	Create(ctx context.Context, userID, feedbackType, subject, description string) (*models.Feedback, error)

	// FindByID retrieves feedback by its ID
	FindByID(ctx context.Context, feedbackID string) (*models.Feedback, error)

	// FindByUserID retrieves all feedback submitted by a user
	FindByUserID(ctx context.Context, userID string) ([]*models.Feedback, error)

	// FindAll retrieves all feedback (for site admin)
	FindAll(ctx context.Context) ([]*models.Feedback, error)

	// UpdateStatus updates the status of a feedback entry
	UpdateStatus(ctx context.Context, feedbackID, status string) (*models.Feedback, error)
}
