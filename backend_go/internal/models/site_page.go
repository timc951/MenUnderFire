package models

import "time"

// SitePage represents a site-managed page entity
type SitePage struct {
	ID          string    `json:"id"`
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	IsPublished bool      `json:"isPublished"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// ===== Request DTOs =====

// CreateSitePageRequest is the request for creating a site page
type CreateSitePageRequest struct {
	Slug        string `json:"slug" validate:"required,min=1,max=100"`
	Title       string `json:"title" validate:"required,min=1,max=200"`
	Content     string `json:"content"`
	IsPublished bool   `json:"isPublished"`
}

// UpdateSitePageRequest is the request for updating a site page
type UpdateSitePageRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=200"`
	Content     string `json:"content"`
	IsPublished bool   `json:"isPublished"`
}

// ===== Response DTOs =====

// SitePageSummaryResponse is a summary of a site page for list views
type SitePageSummaryResponse struct {
	ID          string `json:"id"`
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	IsPublished bool   `json:"isPublished"`
	UpdatedAt   string `json:"updatedAt"`
}

// SitePageResponse is the full site page response
type SitePageResponse struct {
	ID          string `json:"id"`
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	IsPublished bool   `json:"isPublished"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// ===== Conversion Methods =====

// ToSummaryResponse converts a SitePage entity to a SitePageSummaryResponse DTO
func (p *SitePage) ToSummaryResponse() *SitePageSummaryResponse {
	return &SitePageSummaryResponse{
		ID:          p.ID,
		Slug:        p.Slug,
		Title:       p.Title,
		IsPublished: p.IsPublished,
		UpdatedAt:   p.UpdatedAt.Format(time.RFC3339),
	}
}

// ToResponse converts a SitePage entity to a SitePageResponse DTO
func (p *SitePage) ToResponse() *SitePageResponse {
	return &SitePageResponse{
		ID:          p.ID,
		Slug:        p.Slug,
		Title:       p.Title,
		Content:     p.Content,
		IsPublished: p.IsPublished,
		CreatedAt:   p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   p.UpdatedAt.Format(time.RFC3339),
	}
}
