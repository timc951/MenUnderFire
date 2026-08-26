package services

import (
	"context"
	"fmt"
	"strings"

	"menunderfire/internal/models"
	"menunderfire/internal/repositories"
)

// FeedbackService defines the interface for feedback-related business logic
type FeedbackService interface {
	// Create creates a new feedback entry
	Create(ctx context.Context, userID string, req *models.CreateFeedbackRequest) (*models.Feedback, error)

	// ListByUser returns all feedback submitted by a user
	ListByUser(ctx context.Context, userID string) ([]*models.Feedback, error)

	// ListAll returns all feedback (site admin only)
	ListAll(ctx context.Context, userID string) ([]*models.Feedback, error)

	// UpdateStatus updates the status of feedback (site admin only)
	UpdateStatus(ctx context.Context, userID, feedbackID string, req *models.UpdateFeedbackStatusRequest) (*models.Feedback, error)
}

type feedbackService struct {
	feedbackRepo repositories.FeedbackRepository
	userRepo     repositories.UserRepository
}

// NewFeedbackService creates a new FeedbackService implementation
func NewFeedbackService(feedbackRepo repositories.FeedbackRepository, userRepo repositories.UserRepository) FeedbackService {
	return &feedbackService{
		feedbackRepo: feedbackRepo,
		userRepo:     userRepo,
	}
}

func (s *feedbackService) Create(ctx context.Context, userID string, req *models.CreateFeedbackRequest) (*models.Feedback, error) {
	// Validate type
	feedbackType := strings.TrimSpace(req.Type)
	if feedbackType != "BUG" && feedbackType != "FEATURE" && feedbackType != "OTHER" {
		return nil, fmt.Errorf("%w: type must be BUG, FEATURE, or OTHER", ErrValidation)
	}

	subject := strings.TrimSpace(req.Subject)
	if subject == "" {
		return nil, fmt.Errorf("%w: subject is required", ErrValidation)
	}
	if len(subject) > 200 {
		return nil, fmt.Errorf("%w: subject must be 200 characters or less", ErrValidation)
	}

	description := strings.TrimSpace(req.Description)
	if description == "" {
		return nil, fmt.Errorf("%w: description is required", ErrValidation)
	}

	feedback, err := s.feedbackRepo.Create(ctx, userID, feedbackType, subject, description)
	if err != nil {
		return nil, fmt.Errorf("failed to create feedback: %w", err)
	}

	return feedback, nil
}

func (s *feedbackService) ListByUser(ctx context.Context, userID string) ([]*models.Feedback, error) {
	items, err := s.feedbackRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list feedback: %w", err)
	}
	if items == nil {
		return []*models.Feedback{}, nil
	}
	return items, nil
}

func (s *feedbackService) ListAll(ctx context.Context, userID string) ([]*models.Feedback, error) {
	// Only site admins can list all feedback
	isSiteAdmin, err := s.userRepo.IsSiteAdmin(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check site admin: %w", err)
	}
	if !isSiteAdmin {
		return nil, ErrForbidden
	}

	items, err := s.feedbackRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list all feedback: %w", err)
	}
	if items == nil {
		return []*models.Feedback{}, nil
	}
	return items, nil
}

func (s *feedbackService) UpdateStatus(ctx context.Context, userID, feedbackID string, req *models.UpdateFeedbackStatusRequest) (*models.Feedback, error) {
	// Only site admins can update feedback status
	isSiteAdmin, err := s.userRepo.IsSiteAdmin(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check site admin: %w", err)
	}
	if !isSiteAdmin {
		return nil, ErrForbidden
	}

	status := strings.TrimSpace(req.Status)
	validStatuses := map[string]bool{
		"OPEN": true, "IN_PROGRESS": true, "RESOLVED": true,
		"CLOSED": true, "ROADMAP": true, "MORE_INFO_REQUIRED": true,
	}
	if !validStatuses[status] {
		return nil, fmt.Errorf("%w: invalid status value", ErrValidation)
	}

	feedback, err := s.feedbackRepo.UpdateStatus(ctx, feedbackID, status)
	if err != nil {
		return nil, fmt.Errorf("failed to update feedback status: %w", err)
	}
	if feedback == nil {
		return nil, ErrNotFound
	}

	return feedback, nil
}
