package models

import "time"

// Organization represents an organization entity
type Organization struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	CreatedByID string    `json:"createdById"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// OrganizationAdmin represents an admin of an organization
type OrganizationAdmin struct {
	ID             string    `json:"id"`
	UserID         string    `json:"userId"`
	OrganizationID string    `json:"organizationId"`
	InvitedByID    string    `json:"invitedById"`
	JoinedAt       time.Time `json:"joinedAt"`
	User           *User     `json:"user,omitempty"`
}

// ===== Request DTOs =====

// CreateOrganizationRequest is the request body for creating an organization
type CreateOrganizationRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=100"`
	Description *string `json:"description,omitempty"`
}

// UpdateOrganizationRequest is the request body for updating an organization
type UpdateOrganizationRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=100"`
	Description *string `json:"description,omitempty"`
}

// ===== Response DTOs =====

// OrganizationResponse is a summary of an organization for list views
type OrganizationResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	CreatedAt   string  `json:"createdAt"`
}

// OrganizationDetailResponse is detailed organization info
type OrganizationDetailResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	CreatedByID string  `json:"createdById"`
	CanEdit     bool    `json:"canEdit"`
	IsSiteAdmin bool    `json:"isSiteAdmin"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

// OrganizationAdminResponse is the response for an organization admin
type OrganizationAdminResponse struct {
	ID             string `json:"id"`
	UserID         string `json:"userId"`
	OrganizationID string `json:"organizationId"`
	DisplayName    string `json:"displayName"`
	Email          string `json:"email"`
	JoinedAt       string `json:"joinedAt"`
}

// ===== Conversion Methods =====

// ToResponse converts an Organization entity to an OrganizationResponse DTO
func (o *Organization) ToResponse() *OrganizationResponse {
	return &OrganizationResponse{
		ID:          o.ID,
		Name:        o.Name,
		Description: o.Description,
		CreatedAt:   o.CreatedAt.Format(time.RFC3339),
	}
}

// ToDetailResponse converts an Organization entity to an OrganizationDetailResponse DTO
func (o *Organization) ToDetailResponse(canEdit bool, isSiteAdmin bool) *OrganizationDetailResponse {
	return &OrganizationDetailResponse{
		ID:          o.ID,
		Name:        o.Name,
		Description: o.Description,
		CreatedByID: o.CreatedByID,
		CanEdit:     canEdit,
		IsSiteAdmin: isSiteAdmin,
		CreatedAt:   o.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   o.UpdatedAt.Format(time.RFC3339),
	}
}

// ToAdminResponse converts an OrganizationAdmin to an OrganizationAdminResponse DTO
func (a *OrganizationAdmin) ToAdminResponse() *OrganizationAdminResponse {
	displayName := ""
	email := ""
	if a.User != nil {
		displayName = a.User.DisplayName
		email = a.User.Email
	}
	return &OrganizationAdminResponse{
		ID:             a.ID,
		UserID:         a.UserID,
		OrganizationID: a.OrganizationID,
		DisplayName:    displayName,
		Email:          email,
		JoinedAt:       a.JoinedAt.Format(time.RFC3339),
	}
}
