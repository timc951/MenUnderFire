package models

import "time"

// PageDraft represents a temporary preview draft of a site page
type PageDraft struct {
	ID          string    `json:"id"`
	PageID      *string   `json:"pageId,omitempty"` // nil for new pages
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	CreatedByID string    `json:"createdById"`
	CreatedAt   time.Time `json:"createdAt"`
	ExpiresAt   time.Time `json:"expiresAt"`
}

// ===== Request DTOs =====

// CreatePageDraftRequest is the request for creating a preview draft
type CreatePageDraftRequest struct {
	Title   string `json:"title" validate:"required,min=1,max=200"`
	Content string `json:"content"`
}

// ===== Response DTOs =====

// PageDraftResponse is the response when creating a draft
type PageDraftResponse struct {
	ID        string `json:"id"`
	ExpiresAt string `json:"expiresAt"`
}

// PageDraftContentResponse is the response when fetching draft content for preview
type PageDraftContentResponse struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	ExpiresAt string `json:"expiresAt"`
}

// ===== Conversion Methods =====

// ToResponse converts a PageDraft to a PageDraftResponse
func (d *PageDraft) ToResponse() *PageDraftResponse {
	return &PageDraftResponse{
		ID:        d.ID,
		ExpiresAt: d.ExpiresAt.Format(time.RFC3339),
	}
}

// ToContentResponse converts a PageDraft to a PageDraftContentResponse
func (d *PageDraft) ToContentResponse() *PageDraftContentResponse {
	return &PageDraftContentResponse{
		Title:     d.Title,
		Content:   d.Content,
		ExpiresAt: d.ExpiresAt.Format(time.RFC3339),
	}
}
