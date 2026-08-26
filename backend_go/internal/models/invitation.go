package models

import "time"

// InvitationType represents the type of invitation
type InvitationType string

const (
	InvitationTypeOrgAdmin    InvitationType = "ORG_ADMIN"
	InvitationTypeGroupOwner  InvitationType = "GROUP_OWNER"
	InvitationTypeGroupMember InvitationType = "GROUP_MEMBER"
)

// InvitationStatus represents the status of an invitation
type InvitationStatus string

const (
	InvitationStatusPending  InvitationStatus = "PENDING"
	InvitationStatusAccepted InvitationStatus = "ACCEPTED"
	InvitationStatusExpired  InvitationStatus = "EXPIRED"
	InvitationStatusRevoked  InvitationStatus = "REVOKED"
)

// Invitation represents an invitation entity
type Invitation struct {
	ID             string           `json:"id"`
	Token          string           `json:"token"`
	Email          string           `json:"email"`
	Type           InvitationType   `json:"type"`
	OrganizationID *string          `json:"organizationId,omitempty"`
	GroupID        *string          `json:"groupId,omitempty"`
	InviterID      string           `json:"inviterId"`
	Status         InvitationStatus `json:"status"`
	ExpiresAt      time.Time        `json:"expiresAt"`
	AcceptedAt     *time.Time       `json:"acceptedAt,omitempty"`
	AcceptedByID   *string          `json:"acceptedById,omitempty"`
	CreatedAt      time.Time        `json:"createdAt"`
	Inviter        *User            `json:"inviter,omitempty"`
}

// ===== Request DTOs =====

// CreateOrgAdminInvitationRequest is the request for inviting an org admin
type CreateOrgAdminInvitationRequest struct {
	Email          string `json:"email" validate:"required,email"`
	OrganizationID string `json:"organizationId" validate:"required,uuid"`
}

// CreateGroupOwnerInvitationRequest is the request for inviting a group owner
type CreateGroupOwnerInvitationRequest struct {
	Email   string `json:"email" validate:"required,email"`
	GroupID string `json:"groupId" validate:"required,uuid"`
}

// CreateGroupMemberInvitationRequest is the request for inviting a group member
type CreateGroupMemberInvitationRequest struct {
	Email   string `json:"email" validate:"required,email"`
	GroupID string `json:"groupId" validate:"required,uuid"`
}

// AcceptInvitationRequest is the request for accepting an invitation
type AcceptInvitationRequest struct {
	Token       string `json:"token" validate:"required"`
	ExternalID  string `json:"externalId" validate:"required"`
	DisplayName string `json:"displayName" validate:"required,min=1,max=100"`
}

// ===== Response DTOs =====

// InvitationResponse is the response for an invitation
type InvitationResponse struct {
	ID             string  `json:"id"`
	Token          string  `json:"token"`
	Email          string  `json:"email"`
	Type           string  `json:"type"`
	OrganizationID *string `json:"organizationId,omitempty"`
	GroupID        *string `json:"groupId,omitempty"`
	Status         string  `json:"status"`
	ExpiresAt      string  `json:"expiresAt"`
	AcceptedAt     *string `json:"acceptedAt,omitempty"`
	AcceptedByID   *string `json:"acceptedById,omitempty"`
	CreatedAt      string  `json:"createdAt"`
}

// AcceptInvitationResponse is the response after accepting an invitation
type AcceptInvitationResponse struct {
	Invitation *InvitationResponse `json:"invitation"`
	User       *UserResponse       `json:"user"`
}

// ValidateInvitationResponse is the response for validating an invitation token
type ValidateInvitationResponse struct {
	Valid            bool     `json:"valid"`
	Email            *string  `json:"email,omitempty"`
	Type             *string  `json:"type,omitempty"`
	OrganizationName *string  `json:"organizationName,omitempty"`
	GroupName        *string  `json:"groupName,omitempty"`
	InviterName      *string  `json:"inviterName,omitempty"`
	ExpiresAt        *string  `json:"expiresAt,omitempty"`
	Capabilities     []string `json:"capabilities,omitempty"`
}

// InvitationListItemResponse is the response for listing invitations with enriched data
type InvitationListItemResponse struct {
	ID               string  `json:"id"`
	Email            string  `json:"email"`
	Type             string  `json:"type"`
	OrganizationID   *string `json:"organizationId,omitempty"`
	OrganizationName *string `json:"organizationName,omitempty"`
	GroupID          *string `json:"groupId,omitempty"`
	GroupName        *string `json:"groupName,omitempty"`
	InviterName      *string `json:"inviterName,omitempty"`
	Status           string  `json:"status"`
	ExpiresAt        string  `json:"expiresAt"`
	CreatedAt        string  `json:"createdAt"`
}

// ===== Conversion Methods =====

// ToResponse converts an Invitation entity to an InvitationResponse DTO
func (i *Invitation) ToResponse() *InvitationResponse {
	resp := &InvitationResponse{
		ID:             i.ID,
		Token:          i.Token,
		Email:          i.Email,
		Type:           string(i.Type),
		OrganizationID: i.OrganizationID,
		GroupID:        i.GroupID,
		Status:         string(i.Status),
		ExpiresAt:      i.ExpiresAt.Format(time.RFC3339),
		AcceptedByID:   i.AcceptedByID,
		CreatedAt:      i.CreatedAt.Format(time.RFC3339),
	}
	if i.AcceptedAt != nil {
		formatted := i.AcceptedAt.Format(time.RFC3339)
		resp.AcceptedAt = &formatted
	}
	return resp
}
