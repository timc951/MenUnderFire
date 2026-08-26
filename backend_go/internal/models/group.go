package models

import "time"

// Group represents a group entity
type Group struct {
	ID                  string     `json:"id"`
	Name                string     `json:"name"`
	Description         *string    `json:"description"`
	OrganizationID      string     `json:"organizationId"`
	OwnerID             *string    `json:"ownerId,omitempty"`
	CreatedBy           string     `json:"createdBy"`
	InviteCode          string     `json:"inviteCode"`
	InviteCodeExpiresAt *time.Time `json:"inviteCodeExpiresAt,omitempty"`
	RequirePostApproval bool       `json:"requirePostApproval"`
	AllowAnonymousPosts bool       `json:"allowAnonymousPosts"`
	CreatedAt           time.Time  `json:"createdAt"`
	UpdatedAt           time.Time  `json:"updatedAt"`
}

// GroupMember represents a user's membership in a group
type GroupMember struct {
	ID       string    `json:"id"`
	GroupID  string    `json:"groupId"`
	UserID   string    `json:"userId"`
	Role     string    `json:"role"` // OWNER, LEADER, MODERATOR, MEMBER
	JoinedAt time.Time `json:"joinedAt"`
	User     *User     `json:"user,omitempty"`
}

// ===== Request DTOs =====

// CreateGroupRequest is the request body for creating a group
type CreateGroupRequest struct {
	Name           string  `json:"name" validate:"required,min=1,max=100"`
	Description    *string `json:"description,omitempty"`
	OrganizationID string  `json:"organizationId" validate:"required,uuid"`
}

// JoinGroupRequest is the request body for joining a group
type JoinGroupRequest struct {
	InviteCode string `json:"inviteCode" validate:"required"`
}

// UpdateGroupSettingsRequest is the request body for updating group settings
type UpdateGroupSettingsRequest struct {
	RequirePostApproval bool `json:"requirePostApproval"`
	AllowAnonymousPosts bool `json:"allowAnonymousPosts"`
}

// ===== Response DTOs =====

// GroupResponse is a summary of a group for list views
type GroupResponse struct {
	ID                  string  `json:"id"`
	Name                string  `json:"name"`
	Description         *string `json:"description,omitempty"`
	InviteCode          *string `json:"inviteCode,omitempty"`          // Only visible to leaders/owners
	InviteCodeExpiresAt *string `json:"inviteCodeExpiresAt,omitempty"` // Only visible to leaders/owners
	Role                *string `json:"role,omitempty"`                // User's role in the group
	MemberCount         *int    `json:"memberCount,omitempty"`
	CreatedAt           string  `json:"createdAt"`
}

// GroupDetailResponse is detailed group info including members
type GroupDetailResponse struct {
	ID                  string           `json:"id"`
	Name                string           `json:"name"`
	Description         *string          `json:"description,omitempty"`
	OrganizationID      string           `json:"organizationId"`
	InviteCode          *string          `json:"inviteCode,omitempty"`          // Only visible to leaders/owners
	InviteCodeExpiresAt *string          `json:"inviteCodeExpiresAt,omitempty"` // Only visible to leaders/owners
	RequirePostApproval bool             `json:"requirePostApproval"`
	AllowAnonymousPosts bool             `json:"allowAnonymousPosts"`
	Role                string           `json:"role"`
	Members             []MemberResponse `json:"members"`
	CreatedAt           string           `json:"createdAt"`
}

// MemberResponse is a group member's public info
type MemberResponse struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	JoinedAt    string `json:"joinedAt"`
}

// ===== Conversion Methods =====

// ToResponse converts a Group entity to a GroupResponse DTO
func (g *Group) ToResponse(role *string, memberCount *int, showInviteCode bool) *GroupResponse {
	var inviteCode *string
	var inviteCodeExpiresAt *string
	if showInviteCode {
		inviteCode = &g.InviteCode
		if g.InviteCodeExpiresAt != nil {
			exp := g.InviteCodeExpiresAt.Format(time.RFC3339)
			inviteCodeExpiresAt = &exp
		}
	}
	return &GroupResponse{
		ID:                  g.ID,
		Name:                g.Name,
		Description:         g.Description,
		InviteCode:          inviteCode,
		InviteCodeExpiresAt: inviteCodeExpiresAt,
		Role:                role,
		MemberCount:         memberCount,
		CreatedAt:           g.CreatedAt.Format(time.RFC3339),
	}
}

// ToDetailResponse converts a Group entity to a GroupDetailResponse DTO
func (g *Group) ToDetailResponse(role string, members []MemberResponse, showInviteCode bool) *GroupDetailResponse {
	var inviteCode *string
	var inviteCodeExpiresAt *string
	if showInviteCode {
		inviteCode = &g.InviteCode
		if g.InviteCodeExpiresAt != nil {
			exp := g.InviteCodeExpiresAt.Format(time.RFC3339)
			inviteCodeExpiresAt = &exp
		}
	}
	return &GroupDetailResponse{
		ID:                  g.ID,
		Name:                g.Name,
		Description:         g.Description,
		OrganizationID:      g.OrganizationID,
		InviteCode:          inviteCode,
		InviteCodeExpiresAt: inviteCodeExpiresAt,
		RequirePostApproval: g.RequirePostApproval,
		AllowAnonymousPosts: g.AllowAnonymousPosts,
		Role:                role,
		Members:             members,
		CreatedAt:           g.CreatedAt.Format(time.RFC3339),
	}
}

// ToMemberResponse converts a GroupMember to a MemberResponse DTO
func (m *GroupMember) ToMemberResponse() *MemberResponse {
	displayName := ""
	email := ""
	if m.User != nil {
		displayName = m.User.DisplayName
		email = m.User.Email
	}
	return &MemberResponse{
		ID:          m.UserID,
		DisplayName: displayName,
		Email:       email,
		Role:        m.Role,
		JoinedAt:    m.JoinedAt.Format(time.RFC3339),
	}
}
