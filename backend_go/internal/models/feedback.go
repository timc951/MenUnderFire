package models

import "time"

// Feedback represents a user-submitted feedback item (bug report, feature request, etc.)
type Feedback struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	Type        string    `json:"type"` // BUG, FEATURE, OTHER
	Subject     string    `json:"subject"`
	Description string    `json:"description"`
	Status      string    `json:"status"` // OPEN, IN_PROGRESS, RESOLVED, CLOSED
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// ===== Request DTOs =====

// CreateFeedbackRequest is the request body for creating feedback
type CreateFeedbackRequest struct {
	Type        string `json:"type" validate:"required,oneof=BUG FEATURE OTHER"`
	Subject     string `json:"subject" validate:"required,min=1,max=200"`
	Description string `json:"description" validate:"required,min=1"`
}

// UpdateFeedbackStatusRequest is the request body for updating feedback status (site admin only)
type UpdateFeedbackStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=OPEN IN_PROGRESS RESOLVED CLOSED ROADMAP MORE_INFO_REQUIRED"`
}

// ===== Response DTOs =====

// FeedbackResponse is the public representation of a feedback item
type FeedbackResponse struct {
	ID              string `json:"id"`
	UserID          string `json:"userId"`
	UserDisplayName string `json:"userDisplayName"`
	UserEmail       string `json:"userEmail"`
	Type            string `json:"type"`
	Subject         string `json:"subject"`
	Description     string `json:"description"`
	Status          string `json:"status"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}

// ===== Conversion Methods =====

// ToResponse converts a Feedback entity to a FeedbackResponse DTO
func (f *Feedback) ToResponse(displayName, email string) *FeedbackResponse {
	return &FeedbackResponse{
		ID:              f.ID,
		UserID:          f.UserID,
		UserDisplayName: displayName,
		UserEmail:       email,
		Type:            f.Type,
		Subject:         f.Subject,
		Description:     f.Description,
		Status:          f.Status,
		CreatedAt:       f.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       f.UpdatedAt.Format(time.RFC3339),
	}
}
