package models

import "time"

// User represents a user entity
type User struct {
	ID                  string     `json:"id"`
	Email               string     `json:"email"`
	DisplayName         string     `json:"displayName"`
	ExternalID          string     `json:"externalId"`
	IsSiteAdmin         bool       `json:"isSiteAdmin"`
	InvitedByID         *string    `json:"invitedById,omitempty"`
	InvitationID        *string    `json:"invitationId,omitempty"`
	AgreementAcceptedAt *time.Time `json:"agreementAcceptedAt,omitempty"`
	AgreementVersion    *string    `json:"agreementVersion,omitempty"`
	CreatedAt           time.Time  `json:"createdAt"`
	UpdatedAt           time.Time  `json:"updatedAt"`
}

// ===== Request DTOs =====

// UpdateUserRequest is the request body for updating a user's profile
type UpdateUserRequest struct {
	DisplayName string `json:"displayName" validate:"required,min=1,max=100"`
}

// ===== Response DTOs =====

// UserResponse is the public user profile response
type UserResponse struct {
	ID                  string  `json:"id"`
	Email               string  `json:"email"`
	DisplayName         string  `json:"displayName"`
	AgreementAcceptedAt *string `json:"agreementAcceptedAt"`
	CreatedAt           string  `json:"createdAt"`
}

// UserPermissionsResponse contains the user's roles and permissions
type UserPermissionsResponse struct {
	IsSiteAdmin            bool     `json:"isSiteAdmin"`
	AdminOfOrganizationIDs []string `json:"adminOfOrganizationIds"`
	OwnedGroupIDs          []string `json:"ownedGroupIds"`
	MemberGroupIDs         []string `json:"memberGroupIds"`
}

// ===== Conversion Methods =====

// ToResponse converts a User entity to a UserResponse DTO
func (u *User) ToResponse() *UserResponse {
	resp := &UserResponse{
		ID:          u.ID,
		Email:       u.Email,
		DisplayName: u.DisplayName,
		CreatedAt:   u.CreatedAt.Format(time.RFC3339),
	}
	if u.AgreementAcceptedAt != nil {
		formatted := u.AgreementAcceptedAt.Format(time.RFC3339)
		resp.AgreementAcceptedAt = &formatted
	}
	return resp
}

// ===== Agreement DTOs =====

// AcceptAgreementRequest is the request body for accepting the testing agreement
type AcceptAgreementRequest struct {
	AgreementVersion string `json:"agreementVersion"`
}

// AcceptAgreementResponse is the response after accepting the testing agreement
type AcceptAgreementResponse struct {
	AcceptedAt string `json:"acceptedAt"`
	Version    string `json:"version"`
	Signature  string `json:"signature"`
}
