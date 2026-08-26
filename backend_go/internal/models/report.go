package models

import "time"

// Report represents an accountability report entity
type Report struct {
	ID                 string    `json:"id"`
	GroupID            string    `json:"groupId"`
	ReporterID         string    `json:"reporterId"`
	Title              string    `json:"title"`
	Content            string    `json:"content"`
	IsAnonymousToGroup bool      `json:"isAnonymousToGroup"`
	CreatedAt          time.Time `json:"createdAt"`
	Reporter           *User     `json:"reporter,omitempty"`
}

// ===== Request DTOs =====

// CreateReportRequest is the request body for creating a report
type CreateReportRequest struct {
	GroupID            string `json:"groupId" validate:"required,uuid"`
	Title              string `json:"title" validate:"required,min=1,max=200"`
	Content            string `json:"content" validate:"required,min=1"`
	IsAnonymousToGroup bool   `json:"isAnonymousToGroup"`
}

// ===== Response DTOs =====

// ReportResponse is the response for a report with visibility-aware author info
type ReportResponse struct {
	ID                 string  `json:"id"`
	GroupID            string  `json:"groupId"`
	Title              string  `json:"title"`
	Content            string  `json:"content"`
	IsAnonymousToGroup bool    `json:"isAnonymousToGroup"`
	ReporterName       *string `json:"reporterName,omitempty"` // null if anonymous and requester is not a leader
	ReporterID         *string `json:"reporterId,omitempty"`   // null if anonymous and requester is not a leader
	IsOwnReport        bool    `json:"isOwnReport"`            // true if the requester authored this report
	CreatedAt          string  `json:"createdAt"`
}

// ===== Conversion Methods =====

// ToResponse converts a Report entity to a ReportResponse DTO
// showAuthor determines if the reporter info should be included (based on requester's role)
// isOwnReport indicates if the requester is the author
func (r *Report) ToResponse(showAuthor bool, isOwnReport bool) *ReportResponse {
	response := &ReportResponse{
		ID:                 r.ID,
		GroupID:            r.GroupID,
		Title:              r.Title,
		Content:            r.Content,
		IsAnonymousToGroup: r.IsAnonymousToGroup,
		IsOwnReport:        isOwnReport,
		CreatedAt:          r.CreatedAt.Format(time.RFC3339),
	}

	// Show author info if:
	// 1. Report is not anonymous, OR
	// 2. Requester is a leader/owner (showAuthor = true), OR
	// 3. Requester is the author (isOwnReport = true)
	if !r.IsAnonymousToGroup || showAuthor || isOwnReport {
		response.ReporterID = &r.ReporterID
		if r.Reporter != nil {
			response.ReporterName = &r.Reporter.DisplayName
		}
	}

	return response
}
