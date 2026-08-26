package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"strings"
	"time"

	"menunderfire/internal/models"
	"menunderfire/internal/repositories"
)

// InvitationService defines the interface for invitation-related business logic
type InvitationService interface {
	// CreateOrgAdmin creates an invitation for an organization admin (Site Admin only)
	CreateOrgAdmin(ctx context.Context, inviterID string, req *models.CreateOrgAdminInvitationRequest) (*models.Invitation, error)

	// CreateGroupOwner creates an invitation for a group owner (Site Admin or Org Admin)
	CreateGroupOwner(ctx context.Context, inviterID string, req *models.CreateGroupOwnerInvitationRequest) (*models.Invitation, error)

	// CreateGroupMember creates an invitation for a group member (Site Admin, Org Admin, or Group Owner)
	CreateGroupMember(ctx context.Context, inviterID string, req *models.CreateGroupMemberInvitationRequest) (*models.Invitation, error)

	// Delete deletes an invitation
	Delete(ctx context.Context, invitationID string, userID string) error

	// Accept accepts an invitation and creates/updates the user
	Accept(ctx context.Context, req *models.AcceptInvitationRequest) (*models.Invitation, *models.User, error)

	// Validate validates an invitation token and returns details
	Validate(ctx context.Context, token string) (*models.ValidateInvitationResponse, error)

	// GetByID retrieves an invitation by its ID
	GetByID(ctx context.Context, invitationID string) (*models.Invitation, error)

	// GetByToken retrieves an invitation by its token
	GetByToken(ctx context.Context, token string) (*models.Invitation, error)

	// List retrieves invitations visible to the requesting user based on their role
	List(ctx context.Context, userID string) ([]*models.InvitationListItemResponse, error)
}

// invitationService implements the InvitationService interface
type invitationService struct {
	invitationRepo repositories.InvitationRepository
	userRepo       repositories.UserRepository
	orgRepo        repositories.OrganizationRepository
	groupRepo      repositories.GroupRepository
	groupService   GroupService
}

// NewInvitationService creates a new InvitationService implementation
func NewInvitationService(
	invitationRepo repositories.InvitationRepository,
	userRepo repositories.UserRepository,
	orgRepo repositories.OrganizationRepository,
	groupRepo repositories.GroupRepository,
	groupService GroupService,
) InvitationService {
	return &invitationService{
		invitationRepo: invitationRepo,
		userRepo:       userRepo,
		orgRepo:        orgRepo,
		groupRepo:      groupRepo,
		groupService:   groupService,
	}
}

// generateToken generates a cryptographically secure random token
func (s *invitationService) generateToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)
}

// CreateOrgAdmin creates an invitation for an organization admin
// Edge cases:
// - Request is nil -> return ErrValidation
// - Email is empty or whitespace only -> return ErrValidation
// - OrganizationID is empty -> return ErrValidation
// - Inviter is not a site admin -> return ErrForbidden
// - Organization not found -> return ErrOrganizationNotFound
// - Pending invitation already exists for this email/org -> return ErrConflict
// - Repository error during isSiteAdmin check -> return wrapped error
// - Repository error during org lookup -> return wrapped error
// - Repository error during duplicate check -> return wrapped error
// - Repository error during invitation creation -> return wrapped error
func (s *invitationService) CreateOrgAdmin(ctx context.Context, inviterID string, req *models.CreateOrgAdminInvitationRequest) (*models.Invitation, error) {
	if req == nil {
		return nil, ErrValidation
	}

	email := strings.TrimSpace(req.Email)
	if email == "" {
		return nil, ErrValidation
	}

	orgID := strings.TrimSpace(req.OrganizationID)
	if orgID == "" {
		return nil, ErrValidation
	}

	// Check if inviter is a site admin
	isSiteAdmin, err := s.userRepo.IsSiteAdmin(ctx, inviterID)
	if err != nil {
		return nil, err
	}
	if !isSiteAdmin {
		return nil, ErrForbidden
	}

	// Check if organization exists
	_, err = s.orgRepo.FindByID(ctx, orgID)
	if err != nil {
		return nil, ErrOrganizationNotFound
	}

	// Check for existing pending invitation
	existing, err := s.invitationRepo.FindByEmail(ctx, email, models.InvitationTypeOrgAdmin, orgID)
	if err == nil && existing != nil && existing.Status == models.InvitationStatusPending {
		return nil, ErrConflict
	}

	// Create the invitation
	token := s.generateToken()
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 7 days

	invitation, err := s.invitationRepo.Create(ctx, email, models.InvitationTypeOrgAdmin, &orgID, nil, inviterID, token, expiresAt)
	if err != nil {
		return nil, err
	}

	return invitation, nil
}

// CreateGroupOwner creates an invitation for a group owner
// Edge cases:
// - Request is nil -> return ErrValidation
// - Email is empty or whitespace only -> return ErrValidation
// - GroupID is empty -> return ErrValidation
// - Group not found -> return ErrGroupNotFound
// - Inviter is not site admin and not org admin -> return ErrForbidden
// - Pending invitation already exists for this email/group -> return ErrConflict
// - Repository error during group lookup -> return wrapped error
// - Repository error during isSiteAdmin check -> return wrapped error
// - Repository error during isOrgAdmin check -> return wrapped error
// - Repository error during duplicate check -> return wrapped error
// - Repository error during invitation creation -> return wrapped error
func (s *invitationService) CreateGroupOwner(ctx context.Context, inviterID string, req *models.CreateGroupOwnerInvitationRequest) (*models.Invitation, error) {
	if req == nil {
		return nil, ErrValidation
	}

	email := strings.TrimSpace(req.Email)
	if email == "" {
		return nil, ErrValidation
	}

	groupID := strings.TrimSpace(req.GroupID)
	if groupID == "" {
		return nil, ErrValidation
	}

	// Check if group exists
	group, err := s.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return nil, ErrGroupNotFound
	}

	// Check authorization: must be site admin or org admin
	isSiteAdmin, err := s.userRepo.IsSiteAdmin(ctx, inviterID)
	if err != nil {
		return nil, err
	}

	if !isSiteAdmin {
		isOrgAdmin, err := s.orgRepo.IsAdmin(ctx, group.OrganizationID, inviterID)
		if err != nil {
			return nil, err
		}
		if !isOrgAdmin {
			return nil, ErrForbidden
		}
	}

	// Check for existing pending invitation
	existing, err := s.invitationRepo.FindByEmail(ctx, email, models.InvitationTypeGroupOwner, groupID)
	if err == nil && existing != nil && existing.Status == models.InvitationStatusPending {
		return nil, ErrConflict
	}

	// Create the invitation
	token := s.generateToken()
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 7 days

	invitation, err := s.invitationRepo.Create(ctx, email, models.InvitationTypeGroupOwner, nil, &groupID, inviterID, token, expiresAt)
	if err != nil {
		return nil, err
	}

	return invitation, nil
}

// CreateGroupMember creates an invitation for a group member
// Edge cases:
// - Request is nil -> return ErrValidation
// - Email is empty or whitespace only -> return ErrValidation
// - GroupID is empty -> return ErrValidation
// - Group not found -> return ErrGroupNotFound
// - Inviter is not site admin, not org admin, and not group owner -> return ErrForbidden
// - Pending invitation already exists for this email/group -> return ErrConflict
// - Repository error during group lookup -> return wrapped error
// - Repository error during isSiteAdmin check -> return wrapped error
// - Repository error during isOrgAdmin check -> return wrapped error
// - Repository error during getUserRole check -> return wrapped error
// - Repository error during duplicate check -> return wrapped error
// - Repository error during invitation creation -> return wrapped error
func (s *invitationService) CreateGroupMember(ctx context.Context, inviterID string, req *models.CreateGroupMemberInvitationRequest) (*models.Invitation, error) {
	if req == nil {
		return nil, ErrValidation
	}

	email := strings.TrimSpace(req.Email)
	if email == "" {
		return nil, ErrValidation
	}

	groupID := strings.TrimSpace(req.GroupID)
	if groupID == "" {
		return nil, ErrValidation
	}

	// Check if group exists
	group, err := s.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return nil, ErrGroupNotFound
	}

	// Check authorization: must be site admin, org admin, or group owner
	isSiteAdmin, err := s.userRepo.IsSiteAdmin(ctx, inviterID)
	if err != nil {
		return nil, err
	}

	if !isSiteAdmin {
		isOrgAdmin, err := s.orgRepo.IsAdmin(ctx, group.OrganizationID, inviterID)
		if err != nil {
			return nil, err
		}

		if !isOrgAdmin {
			role, err := s.groupService.GetUserRole(ctx, groupID, inviterID)
			if err != nil {
				return nil, err
			}
			if role != "OWNER" {
				return nil, ErrForbidden
			}
		}
	}

	// Check for existing pending invitation
	existing, err := s.invitationRepo.FindByEmail(ctx, email, models.InvitationTypeGroupMember, groupID)
	if err == nil && existing != nil && existing.Status == models.InvitationStatusPending {
		return nil, ErrConflict
	}

	// Create the invitation
	token := s.generateToken()
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 7 days

	invitation, err := s.invitationRepo.Create(ctx, email, models.InvitationTypeGroupMember, nil, &groupID, inviterID, token, expiresAt)
	if err != nil {
		return nil, err
	}

	return invitation, nil
}

// Delete deletes an invitation
// Edge cases:
// - InvitationID is empty -> return ErrValidation
// - Invitation not found -> return ErrInvitationNotFound
// - For ORG_ADMIN invitation: user is not site admin -> return ErrForbidden
// - For GROUP_OWNER invitation: user is not site admin and not org admin -> return ErrForbidden
// - For GROUP_MEMBER invitation: user is not site admin, not org admin, and not group owner -> return ErrForbidden
// - Invitation is already accepted -> return ErrConflict
// - Repository error during invitation lookup -> return wrapped error
// - Repository error during authorization checks -> return wrapped error
// - Repository error during deletion -> return wrapped error
func (s *invitationService) Delete(ctx context.Context, invitationID string, userID string) error {
	invitationID = strings.TrimSpace(invitationID)
	if invitationID == "" {
		return ErrValidation
	}

	// Get the invitation
	invitation, err := s.invitationRepo.FindByID(ctx, invitationID)
	if err != nil {
		return ErrInvitationNotFound
	}

	// Cannot delete already accepted invitations
	if invitation.Status == models.InvitationStatusAccepted {
		return ErrConflict
	}

	// Check authorization based on invitation type
	isSiteAdmin, err := s.userRepo.IsSiteAdmin(ctx, userID)
	if err != nil {
		return err
	}

	if !isSiteAdmin {
		switch invitation.Type {
		case models.InvitationTypeOrgAdmin:
			// Only site admin can delete org admin invitations
			return ErrForbidden

		case models.InvitationTypeGroupOwner:
			// Site admin or org admin can delete
			if invitation.GroupID == nil {
				return ErrForbidden
			}
			group, err := s.groupRepo.FindByID(ctx, *invitation.GroupID)
			if err != nil {
				return err
			}
			isOrgAdmin, err := s.orgRepo.IsAdmin(ctx, group.OrganizationID, userID)
			if err != nil {
				return err
			}
			if !isOrgAdmin {
				return ErrForbidden
			}

		case models.InvitationTypeGroupMember:
			// Site admin, org admin, or group owner can delete
			if invitation.GroupID == nil {
				return ErrForbidden
			}
			group, err := s.groupRepo.FindByID(ctx, *invitation.GroupID)
			if err != nil {
				return err
			}
			isOrgAdmin, err := s.orgRepo.IsAdmin(ctx, group.OrganizationID, userID)
			if err != nil {
				return err
			}
			if !isOrgAdmin {
				role, err := s.groupService.GetUserRole(ctx, *invitation.GroupID, userID)
				if err != nil {
					return err
				}
				if role != "OWNER" {
					return ErrForbidden
				}
			}
		}
	}

	return s.invitationRepo.Delete(ctx, invitationID)
}

// Accept accepts an invitation and creates/updates the user
// Edge cases:
// - Request is nil -> return ErrValidation
// - Token is empty -> return ErrValidation
// - ExternalID is empty -> return ErrValidation
// - DisplayName is empty or whitespace only -> return ErrValidation
// - Invitation not found by token -> return ErrInvitationNotFound
// - Invitation is expired -> return ErrInvalidToken
// - Invitation is already accepted -> return ErrConflict
// - Invitation is revoked -> return ErrInvalidToken
// - For ORG_ADMIN: OrganizationID is nil -> return ErrValidation
// - For GROUP_OWNER/GROUP_MEMBER: GroupID is nil -> return ErrValidation
// - Repository error during invitation lookup -> return wrapped error
// - Repository error during user creation -> return wrapped error
// - Repository error during role assignment -> return wrapped error
// - Repository error during status update -> return wrapped error
func (s *invitationService) Accept(ctx context.Context, req *models.AcceptInvitationRequest) (*models.Invitation, *models.User, error) {
	if req == nil {
		return nil, nil, ErrValidation
	}

	token := strings.TrimSpace(req.Token)
	if token == "" {
		return nil, nil, ErrValidation
	}

	externalID := strings.TrimSpace(req.ExternalID)
	if externalID == "" {
		return nil, nil, ErrValidation
	}

	displayName := strings.TrimSpace(req.DisplayName)
	if displayName == "" {
		return nil, nil, ErrValidation
	}

	// Find invitation by token
	invitation, err := s.invitationRepo.FindByToken(ctx, token)
	if err != nil {
		return nil, nil, ErrInvitationNotFound
	}

	// Check invitation status
	switch invitation.Status {
	case models.InvitationStatusAccepted:
		return nil, nil, ErrConflict
	case models.InvitationStatusExpired, models.InvitationStatusRevoked:
		return nil, nil, ErrInvalidToken
	}

	// Check if expired
	if time.Now().After(invitation.ExpiresAt) {
		return nil, nil, ErrInvalidToken
	}

	// Create or find user
	user, err := s.userRepo.FindByEmail(ctx, invitation.Email)
	if err != nil {
		// User doesn't exist, create new user
		user, err = s.userRepo.Create(ctx, invitation.Email, displayName, externalID)
		if err != nil {
			return nil, nil, err
		}
	}

	// Link the user to the invitation that brought them in
	if err := s.userRepo.UpdateInvitationInfo(ctx, user.ID, invitation.InviterID, invitation.ID); err != nil {
		return nil, nil, err
	}
	user.InvitedByID = &invitation.InviterID
	user.InvitationID = &invitation.ID

	// Assign role based on invitation type
	switch invitation.Type {
	case models.InvitationTypeOrgAdmin:
		if invitation.OrganizationID == nil {
			return nil, nil, ErrValidation
		}
		err = s.orgRepo.AddAdmin(ctx, *invitation.OrganizationID, user.ID)
		if err != nil {
			return nil, nil, err
		}

	case models.InvitationTypeGroupOwner:
		if invitation.GroupID == nil {
			return nil, nil, ErrValidation
		}
		_, err = s.groupRepo.AddMember(ctx, *invitation.GroupID, user.ID, "OWNER")
		if err != nil {
			return nil, nil, err
		}

	case models.InvitationTypeGroupMember:
		if invitation.GroupID == nil {
			return nil, nil, ErrValidation
		}
		_, err = s.groupRepo.AddMember(ctx, *invitation.GroupID, user.ID, "MEMBER")
		if err != nil {
			return nil, nil, err
		}
	}

	// Update invitation status, accepted_at, and accepted_by_id in one call
	err = s.invitationRepo.UpdateStatus(ctx, invitation.ID, models.InvitationStatusAccepted, &user.ID)
	if err != nil {
		return nil, nil, err
	}

	now := time.Now()
	invitation.Status = models.InvitationStatusAccepted
	invitation.AcceptedAt = &now
	invitation.AcceptedByID = &user.ID
	return invitation, user, nil
}

// Validate validates an invitation token and returns details
// Edge cases:
// - Token is empty -> return invalid response with Valid=false
// - Invitation not found -> return invalid response with Valid=false
// - Invitation is expired -> return invalid response with Valid=false
// - Invitation is already accepted -> return invalid response with Valid=false
// - Invitation is revoked -> return invalid response with Valid=false
// - For ORG_ADMIN: org lookup fails -> return valid response without org name
// - For GROUP_OWNER/GROUP_MEMBER: group lookup fails -> return valid response without group name
// - Inviter lookup fails -> return valid response without inviter name
// - Repository error during invitation lookup -> return invalid response with Valid=false
func (s *invitationService) Validate(ctx context.Context, token string) (*models.ValidateInvitationResponse, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return &models.ValidateInvitationResponse{Valid: false}, nil
	}

	invitation, err := s.invitationRepo.FindByToken(ctx, token)
	if err != nil {
		return &models.ValidateInvitationResponse{Valid: false}, nil
	}

	// Check if expired
	if time.Now().After(invitation.ExpiresAt) {
		return &models.ValidateInvitationResponse{Valid: false}, nil
	}

	// Check status
	if invitation.Status != models.InvitationStatusPending {
		return &models.ValidateInvitationResponse{Valid: false}, nil
	}

	// Build response
	invType := string(invitation.Type)
	expiresAt := invitation.ExpiresAt.Format(time.RFC3339)

	response := &models.ValidateInvitationResponse{
		Valid:     true,
		Email:     &invitation.Email,
		Type:      &invType,
		ExpiresAt: &expiresAt,
	}

	// Get organization name for ORG_ADMIN invitations
	if invitation.Type == models.InvitationTypeOrgAdmin && invitation.OrganizationID != nil {
		org, err := s.orgRepo.FindByID(ctx, *invitation.OrganizationID)
		if err == nil {
			response.OrganizationName = &org.Name
		}
	}

	// Get group name for GROUP_OWNER/GROUP_MEMBER invitations
	if (invitation.Type == models.InvitationTypeGroupOwner || invitation.Type == models.InvitationTypeGroupMember) && invitation.GroupID != nil {
		group, err := s.groupRepo.FindByID(ctx, *invitation.GroupID)
		if err == nil {
			response.GroupName = &group.Name
		}
	}

	// Get inviter name
	inviter, err := s.userRepo.FindByID(ctx, invitation.InviterID)
	if err == nil {
		response.InviterName = &inviter.DisplayName
	}

	return response, nil
}

// GetByID retrieves an invitation by its ID
// Edge cases:
// - InvitationID is empty -> return ErrValidation
// - Invitation not found -> return ErrInvitationNotFound
// - Repository error -> return wrapped error
func (s *invitationService) GetByID(ctx context.Context, invitationID string) (*models.Invitation, error) {
	invitationID = strings.TrimSpace(invitationID)
	if invitationID == "" {
		return nil, ErrValidation
	}

	invitation, err := s.invitationRepo.FindByID(ctx, invitationID)
	if err != nil {
		return nil, ErrInvitationNotFound
	}

	return invitation, nil
}

// GetByToken retrieves an invitation by its token
// Edge cases:
// - Token is empty -> return ErrValidation
// - Invitation not found -> return ErrInvitationNotFound
// - Repository error -> return wrapped error
func (s *invitationService) GetByToken(ctx context.Context, token string) (*models.Invitation, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, ErrValidation
	}

	invitation, err := s.invitationRepo.FindByToken(ctx, token)
	if err != nil {
		return nil, ErrInvitationNotFound
	}

	return invitation, nil
}

// List retrieves invitations visible to the requesting user based on their role:
// - Site admin: sees all invitations
// - Org admin: sees invitations for their organizations and groups within them
// - Group owner: sees invitations for their groups
func (s *invitationService) List(ctx context.Context, userID string) ([]*models.InvitationListItemResponse, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, ErrValidation
	}

	isSiteAdmin, err := s.userRepo.IsSiteAdmin(ctx, userID)
	if err != nil {
		return nil, err
	}

	var invitations []*models.Invitation

	if isSiteAdmin {
		invitations, err = s.invitationRepo.ListAll(ctx)
		if err != nil {
			return nil, err
		}
	} else {
		// Get org admin IDs
		adminOrgIDs, err := s.userRepo.FindAdminOrganizationIDs(ctx, userID)
		if err != nil {
			return nil, err
		}
		// Get owned group IDs
		ownedGroupIDs, err := s.userRepo.FindOwnedGroupIDs(ctx, userID)
		if err != nil {
			return nil, err
		}

		if len(adminOrgIDs) == 0 && len(ownedGroupIDs) == 0 {
			return []*models.InvitationListItemResponse{}, nil
		}

		// Combine results from org-level and group-level queries
		seen := make(map[string]bool)

		if len(adminOrgIDs) > 0 {
			orgInvitations, err := s.invitationRepo.ListByOrganizationIDs(ctx, adminOrgIDs)
			if err != nil {
				return nil, err
			}
			for _, inv := range orgInvitations {
				if !seen[inv.ID] {
					invitations = append(invitations, inv)
					seen[inv.ID] = true
				}
			}
		}

		if len(ownedGroupIDs) > 0 {
			groupInvitations, err := s.invitationRepo.ListByGroupIDs(ctx, ownedGroupIDs)
			if err != nil {
				return nil, err
			}
			for _, inv := range groupInvitations {
				if !seen[inv.ID] {
					invitations = append(invitations, inv)
					seen[inv.ID] = true
				}
			}
		}
	}

	// Enrich with org names, group names, and inviter names
	// Cache lookups to avoid repeated queries
	orgNameCache := make(map[string]string)
	groupNameCache := make(map[string]string)
	userNameCache := make(map[string]string)

	result := make([]*models.InvitationListItemResponse, len(invitations))
	for i, inv := range invitations {
		item := &models.InvitationListItemResponse{
			ID:             inv.ID,
			Email:          inv.Email,
			Type:           string(inv.Type),
			OrganizationID: inv.OrganizationID,
			GroupID:        inv.GroupID,
			Status:         string(inv.Status),
			ExpiresAt:      inv.ExpiresAt.Format(time.RFC3339),
			CreatedAt:      inv.CreatedAt.Format(time.RFC3339),
		}

		// Resolve org name
		if inv.OrganizationID != nil {
			if name, ok := orgNameCache[*inv.OrganizationID]; ok {
				item.OrganizationName = &name
			} else {
				org, err := s.orgRepo.FindByID(ctx, *inv.OrganizationID)
				if err == nil {
					orgNameCache[*inv.OrganizationID] = org.Name
					item.OrganizationName = &org.Name
				}
			}
		}

		// Resolve group name (and org name if not already set)
		if inv.GroupID != nil {
			if name, ok := groupNameCache[*inv.GroupID]; ok {
				item.GroupName = &name
			} else {
				group, err := s.groupRepo.FindByID(ctx, *inv.GroupID)
				if err == nil {
					groupNameCache[*inv.GroupID] = group.Name
					item.GroupName = &group.Name
					// Also resolve the org for group invitations
					if item.OrganizationName == nil {
						if name, ok := orgNameCache[group.OrganizationID]; ok {
							item.OrganizationName = &name
						} else {
							org, err := s.orgRepo.FindByID(ctx, group.OrganizationID)
							if err == nil {
								orgNameCache[group.OrganizationID] = org.Name
								item.OrganizationName = &org.Name
								item.OrganizationID = &group.OrganizationID
							}
						}
					}
				}
			}
		}

		// Resolve inviter name
		if name, ok := userNameCache[inv.InviterID]; ok {
			item.InviterName = &name
		} else {
			inviter, err := s.userRepo.FindByID(ctx, inv.InviterID)
			if err == nil {
				userNameCache[inv.InviterID] = inviter.DisplayName
				item.InviterName = &inviter.DisplayName
			}
		}

		result[i] = item
	}

	return result, nil
}
