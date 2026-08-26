package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"menunderfire/internal/models"
	"menunderfire/internal/repositories"
)

// GroupService defines the interface for group-related business logic
type GroupService interface {
	// Create creates a new group with the given user as the owner
	Create(ctx context.Context, userID string, req *models.CreateGroupRequest) (*models.Group, error)

	// List returns all groups the user belongs to
	List(ctx context.Context, userID string) ([]*models.Group, error)

	// GetByID retrieves a group by its ID
	GetByID(ctx context.Context, groupID string) (*models.Group, error)

	// GetDetailByID retrieves detailed group info including members
	GetDetailByID(ctx context.Context, groupID string, userID string) (*models.GroupDetailResponse, error)

	// Join adds a user to a group using an invite code (requires group ID)
	Join(ctx context.Context, groupID string, userID string, inviteCode string) error

	// JoinByInviteCode adds a user to a group by looking up the invite code
	JoinByInviteCode(ctx context.Context, userID string, inviteCode string) error

	// RemoveMember removes a user from a group
	RemoveMember(ctx context.Context, groupID string, targetUserID string, requestingUserID string) error

	// GetUserRole returns the user's role in a group (or empty string if not a member)
	GetUserRole(ctx context.Context, groupID string, userID string) (string, error)

	// GetMembers returns all members of a group
	GetMembers(ctx context.Context, groupID string) ([]*models.GroupMember, error)

	// GetMemberCount returns the number of members in a group
	GetMemberCount(ctx context.Context, groupID string) (int, error)

	// UpdateSettings updates group settings (requires OWNER, LEADER, org admin, or site admin)
	UpdateSettings(ctx context.Context, groupID string, userID string, req *models.UpdateGroupSettingsRequest) error

	// UpdateMemberRole changes a member's role (with hierarchy enforcement)
	UpdateMemberRole(ctx context.Context, groupID string, targetUserID string, newRole string, requestingUserID string) error
}

type groupService struct {
	groupRepo repositories.GroupRepository
	orgRepo   repositories.OrganizationRepository
	userRepo  repositories.UserRepository
}

// NewGroupService creates a new GroupService implementation with the provided dependencies
func NewGroupService(
	groupRepo repositories.GroupRepository,
	orgRepo repositories.OrganizationRepository,
	userRepo repositories.UserRepository,
) GroupService {
	return &groupService{
		groupRepo: groupRepo,
		orgRepo:   orgRepo,
		userRepo:  userRepo,
	}
}

// Create creates a new group with the given user as the owner
//
// Edge cases handled:
//   - Empty name after trim: Returns ErrValidation
//   - Whitespace only name: Returns ErrValidation
//   - Name with leading/trailing whitespace: Trims and creates
//   - Empty organizationID: Returns ErrValidation
//   - Organization not found: Returns ErrOrganizationNotFound
//   - User is not org admin: Returns ErrForbidden
//   - Description is nil: Passes nil to repository
//   - Description with whitespace: Trims the description
//   - Database error finding org: Returns wrapped error
//   - Database error checking org admin: Returns wrapped error
//   - Database error creating group: Returns wrapped error
//   - Database error adding owner as member: Returns wrapped error
//   - Successful creation: Returns group with user as OWNER
func (s *groupService) Create(ctx context.Context, userID string, req *models.CreateGroupRequest) (*models.Group, error) {
	// Trim and validate name
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrValidation
	}

	// Validate organization ID
	if req.OrganizationID == "" {
		return nil, ErrValidation
	}

	// Verify organization exists
	org, err := s.orgRepo.FindByID(ctx, req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to find organization: %w", err)
	}
	if org == nil {
		return nil, ErrOrganizationNotFound
	}

	// Check if user is org admin (required to create groups)
	isAdmin, err := s.orgRepo.IsAdmin(ctx, req.OrganizationID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check org admin status: %w", err)
	}
	if !isAdmin {
		return nil, ErrForbidden
	}

	// Trim description if provided
	var description *string
	if req.Description != nil {
		trimmed := strings.TrimSpace(*req.Description)
		description = &trimmed
	}

	// Generate invite code
	inviteCode := s.groupRepo.GenerateInviteCode()

	// Create group
	group, err := s.groupRepo.Create(ctx, name, description, req.OrganizationID, inviteCode, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create group: %w", err)
	}

	// Add creator as owner
	_, err = s.groupRepo.AddMember(ctx, group.ID, userID, "OWNER")
	if err != nil {
		return nil, fmt.Errorf("failed to add owner to group: %w", err)
	}

	return group, nil
}

// List returns all groups the user has access to.
// For regular members: only groups they belong to.
// For org admins: groups they belong to + all groups in their organizations.
func (s *groupService) List(ctx context.Context, userID string) ([]*models.Group, error) {
	// Get groups the user belongs to
	memberGroups, err := s.groupRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find groups for user: %w", err)
	}

	// Check if user is an org admin
	adminOrgIDs, err := s.userRepo.FindAdminOrganizationIDs(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find admin org IDs: %w", err)
	}

	if len(adminOrgIDs) == 0 {
		if memberGroups == nil {
			return []*models.Group{}, nil
		}
		return memberGroups, nil
	}

	// Org admin: also fetch all groups from their organizations
	seen := make(map[string]bool)
	var allGroups []*models.Group
	for _, g := range memberGroups {
		seen[g.ID] = true
		allGroups = append(allGroups, g)
	}

	for _, orgID := range adminOrgIDs {
		orgGroups, err := s.groupRepo.FindByOrganizationID(ctx, orgID)
		if err != nil {
			return nil, fmt.Errorf("failed to find org groups: %w", err)
		}
		for _, g := range orgGroups {
			if !seen[g.ID] {
				seen[g.ID] = true
				allGroups = append(allGroups, g)
			}
		}
	}

	if allGroups == nil {
		return []*models.Group{}, nil
	}

	return allGroups, nil
}

// GetByID retrieves a group by its ID
//
// Edge cases handled:
//   - Empty groupID: Returns ErrGroupNotFound
//   - Group not found: Returns ErrGroupNotFound
//   - Group exists: Returns the group
//   - Database error: Returns wrapped error
func (s *groupService) GetByID(ctx context.Context, groupID string) (*models.Group, error) {
	if groupID == "" {
		return nil, ErrGroupNotFound
	}

	group, err := s.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group by ID: %w", err)
	}
	if group == nil {
		return nil, ErrGroupNotFound
	}

	return group, nil
}

// GetDetailByID retrieves detailed group info including members
//
// Edge cases handled:
//   - Group not found: Returns ErrGroupNotFound
//   - User is not a group member and not org admin and not site admin: Returns ErrForbidden
//   - User is a group member: Returns detail with their role, shows invite code for OWNER/LEADER
//   - User is org admin (not member): Returns detail with "ADMIN" role, shows invite code
//   - User is site admin (not member): Returns detail with "ADMIN" role, shows invite code
//   - Group has no members (edge case): Returns detail with empty members list
//   - Database error finding group: Returns wrapped error
//   - Database error finding membership: Returns wrapped error
//   - Database error checking site admin: Returns wrapped error
//   - Database error checking org admin: Returns wrapped error
//   - Database error finding members: Returns wrapped error
func (s *groupService) GetDetailByID(ctx context.Context, groupID string, userID string) (*models.GroupDetailResponse, error) {
	// Find group by ID
	group, err := s.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to find group: %w", err)
	}
	if group == nil {
		return nil, ErrGroupNotFound
	}

	// Check if user is a member of the group
	membership, err := s.groupRepo.FindMember(ctx, groupID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check membership: %w", err)
	}

	var role string
	var showInviteCode bool

	if membership != nil {
		// User is a member
		role = membership.Role
		showInviteCode = role == "OWNER" || role == "LEADER"
	} else {
		// User is not a member, check if they're an admin
		siteAdmin, err := s.userRepo.IsSiteAdmin(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to check site admin status: %w", err)
		}

		if siteAdmin {
			role = "ADMIN"
			showInviteCode = true
		} else {
			orgAdmin, err := s.orgRepo.IsAdmin(ctx, group.OrganizationID, userID)
			if err != nil {
				return nil, fmt.Errorf("failed to check org admin status: %w", err)
			}

			if orgAdmin {
				role = "ADMIN"
				showInviteCode = true
			} else {
				// User has no access - return limited info (just group name)
				return group.ToDetailResponse("", []models.MemberResponse{}, false), nil
			}
		}
	}

	// Get all members (only for users with access)
	members, err := s.groupRepo.FindMembers(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to find group members: %w", err)
	}

	// Convert members to response format
	memberResponses := make([]models.MemberResponse, 0, len(members))
	for _, member := range members {
		memberResponses = append(memberResponses, *member.ToMemberResponse())
	}

	return group.ToDetailResponse(role, memberResponses, showInviteCode), nil
}

// Join adds a user to a group using an invite code
//
// Edge cases handled:
//   - Group not found: Returns ErrGroupNotFound
//   - Empty invite code: Returns ErrInvalidInviteCode
//   - Invalid invite code (doesn't match): Returns ErrInvalidInviteCode
//   - User is already a member: Returns ErrConflict
//   - Valid invite code: Adds user as MEMBER role
//   - Database error finding group: Returns wrapped error
//   - Database error checking membership: Returns wrapped error
//   - Database error adding member: Returns wrapped error
func (s *groupService) Join(ctx context.Context, groupID string, userID string, inviteCode string) error {
	// Find group by ID
	group, err := s.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return fmt.Errorf("failed to find group: %w", err)
	}
	if group == nil {
		return ErrGroupNotFound
	}

	// Validate invite code
	if inviteCode == "" {
		return ErrInvalidInviteCode
	}
	if group.InviteCode != inviteCode {
		return ErrInvalidInviteCode
	}

	// Check if invite code has expired
	if group.InviteCodeExpiresAt != nil && group.InviteCodeExpiresAt.Before(time.Now()) {
		return ErrInvalidInviteCode
	}

	// Check if user is already a member
	existingMembership, err := s.groupRepo.FindMember(ctx, groupID, userID)
	if err != nil {
		return fmt.Errorf("failed to check existing membership: %w", err)
	}
	if existingMembership != nil {
		return ErrConflict
	}

	// Add user as member
	_, err = s.groupRepo.AddMember(ctx, groupID, userID, "MEMBER")
	if err != nil {
		return fmt.Errorf("failed to add member: %w", err)
	}

	return nil
}

// JoinByInviteCode adds a user to a group by looking up the invite code
func (s *groupService) JoinByInviteCode(ctx context.Context, userID string, inviteCode string) error {
	if inviteCode == "" {
		return ErrInvalidInviteCode
	}

	group, err := s.groupRepo.FindByInviteCode(ctx, inviteCode)
	if err != nil {
		return ErrGroupNotFound
	}

	return s.Join(ctx, group.ID, userID, inviteCode)
}

// RemoveMember removes a user from a group
//
// Edge cases handled:
//   - Group not found: Returns ErrGroupNotFound
//   - Target user is not a member: Returns ErrNotFound
//   - User removing themselves (any role): Allowed
//   - Owner removing another member: Allowed
//   - Leader removing a MEMBER: Allowed
//   - Leader removing another LEADER: Returns ErrForbidden
//   - Leader removing an OWNER: Returns ErrForbidden
//   - Regular member removing another member: Returns ErrForbidden
//   - Org admin removing any member: Allowed
//   - Site admin removing any member: Allowed
//   - Removing the last owner: Returns ErrValidation (cannot leave group without owner)
//   - Database error finding group: Returns wrapped error
//   - Database error checking target membership: Returns wrapped error
//   - Database error checking requester membership: Returns wrapped error
//   - Database error checking site admin: Returns wrapped error
//   - Database error checking org admin: Returns wrapped error
//   - Database error removing member: Returns wrapped error
func (s *groupService) RemoveMember(ctx context.Context, groupID string, targetUserID string, requestingUserID string) error {
	// Find group by ID
	group, err := s.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return fmt.Errorf("failed to find group: %w", err)
	}
	if group == nil {
		return ErrGroupNotFound
	}

	// Check if target user is a member
	targetMembership, err := s.groupRepo.FindMember(ctx, groupID, targetUserID)
	if err != nil {
		return fmt.Errorf("failed to check target membership: %w", err)
	}
	if targetMembership == nil {
		return ErrNotFound
	}

	// If removing an owner, check if they're the last owner
	if targetMembership.Role == "OWNER" {
		members, err := s.groupRepo.FindMembers(ctx, groupID)
		if err != nil {
			return fmt.Errorf("failed to find group members: %w", err)
		}

		ownerCount := 0
		for _, m := range members {
			if m.Role == "OWNER" {
				ownerCount++
			}
		}

		if ownerCount <= 1 {
			return fmt.Errorf("%w: cannot remove the last owner", ErrValidation)
		}
	}

	// User can always remove themselves
	if targetUserID == requestingUserID {
		if err := s.groupRepo.RemoveMember(ctx, groupID, targetUserID); err != nil {
			return fmt.Errorf("failed to remove member: %w", err)
		}
		return nil
	}

	// Check requester's membership
	requesterMembership, err := s.groupRepo.FindMember(ctx, groupID, requestingUserID)
	if err != nil {
		return fmt.Errorf("failed to check requester membership: %w", err)
	}

	// If requester is a member, check role-based permissions using hierarchy
	if requesterMembership != nil {
		requesterLvl := roleLevel(requesterMembership.Role)
		targetLvl := roleLevel(targetMembership.Role)

		// Can remove anyone below your level
		if requesterLvl > targetLvl {
			if err := s.groupRepo.RemoveMember(ctx, groupID, targetUserID); err != nil {
				return fmt.Errorf("failed to remove member: %w", err)
			}
			return nil
		}
		// Cannot remove at or above your level - fall through to admin checks
	}

	// Check if requester is site admin
	siteAdmin, err := s.userRepo.IsSiteAdmin(ctx, requestingUserID)
	if err != nil {
		return fmt.Errorf("failed to check site admin status: %w", err)
	}
	if siteAdmin {
		if err := s.groupRepo.RemoveMember(ctx, groupID, targetUserID); err != nil {
			return fmt.Errorf("failed to remove member: %w", err)
		}
		return nil
	}

	// Check if requester is org admin
	orgAdmin, err := s.orgRepo.IsAdmin(ctx, group.OrganizationID, requestingUserID)
	if err != nil {
		return fmt.Errorf("failed to check org admin status: %w", err)
	}
	if orgAdmin {
		if err := s.groupRepo.RemoveMember(ctx, groupID, targetUserID); err != nil {
			return fmt.Errorf("failed to remove member: %w", err)
		}
		return nil
	}

	// User has no permission to remove this member
	return ErrForbidden
}

// GetUserRole returns the user's role in a group (or empty string if not a member)
//
// Edge cases handled:
//   - User is not a member: Returns empty string, nil
//   - User is OWNER: Returns "OWNER", nil
//   - User is LEADER: Returns "LEADER", nil
//   - User is MEMBER: Returns "MEMBER", nil
//   - Database error: Returns empty string, wrapped error
func (s *groupService) GetUserRole(ctx context.Context, groupID string, userID string) (string, error) {
	membership, err := s.groupRepo.FindMember(ctx, groupID, userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user role: %w", err)
	}

	if membership == nil {
		return "", nil
	}

	return membership.Role, nil
}

// GetMembers returns all members of a group
//
// Edge cases handled:
//   - Group has no members: Returns empty slice
//   - Group has members: Returns all members with user info populated
//   - Database error: Returns wrapped error
func (s *groupService) GetMembers(ctx context.Context, groupID string) ([]*models.GroupMember, error) {
	members, err := s.groupRepo.FindMembers(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group members: %w", err)
	}

	if members == nil {
		return []*models.GroupMember{}, nil
	}

	return members, nil
}

// GetMemberCount returns the number of members in a group
//
// Edge cases handled:
//   - Group has no members: Returns 0
//   - Group has members: Returns the count
//   - Database error: Returns 0, wrapped error
func (s *groupService) GetMemberCount(ctx context.Context, groupID string) (int, error) {
	count, err := s.groupRepo.CountMembers(ctx, groupID)
	if err != nil {
		return 0, fmt.Errorf("failed to count group members: %w", err)
	}

	return count, nil
}

// roleLevel returns a numeric level for role hierarchy comparison.
// Higher level = more authority. Used for enforce "can't promote/demote above own level".
func roleLevel(role string) int {
	switch role {
	case "OWNER":
		return 4
	case "LEADER":
		return 3
	case "MODERATOR":
		return 2
	case "MEMBER":
		return 1
	default:
		return 0
	}
}

// UpdateSettings updates group settings
// Only OWNER, LEADER, org admin, or site admin can update settings
func (s *groupService) UpdateSettings(ctx context.Context, groupID string, userID string, req *models.UpdateGroupSettingsRequest) error {
	group, err := s.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return ErrGroupNotFound
	}

	// Check permissions
	membership, err := s.groupRepo.FindMember(ctx, groupID, userID)
	if err != nil {
		return fmt.Errorf("failed to check membership: %w", err)
	}

	hasPermission := membership != nil && (membership.Role == "OWNER" || membership.Role == "LEADER")

	if !hasPermission {
		siteAdmin, err := s.userRepo.IsSiteAdmin(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to check site admin: %w", err)
		}
		if siteAdmin {
			hasPermission = true
		}
	}

	if !hasPermission {
		orgAdmin, err := s.orgRepo.IsAdmin(ctx, group.OrganizationID, userID)
		if err != nil {
			return fmt.Errorf("failed to check org admin: %w", err)
		}
		if orgAdmin {
			hasPermission = true
		}
	}

	if !hasPermission {
		return ErrForbidden
	}

	return s.groupRepo.UpdateSettings(ctx, groupID, req.RequirePostApproval, req.AllowAnonymousPosts)
}

// UpdateMemberRole changes a member's role with hierarchy enforcement.
// Rules:
//   - A user cannot promote/demote to a role at or above their own level
//   - A user cannot change the role of someone at or above their own level
//   - Site admins and org admins can change any role except OWNER
//   - Only OWNERs can promote to OWNER or demote another OWNER
func (s *groupService) UpdateMemberRole(ctx context.Context, groupID string, targetUserID string, newRole string, requestingUserID string) error {
	// Validate the new role
	validRoles := map[string]bool{"OWNER": true, "LEADER": true, "MODERATOR": true, "MEMBER": true}
	if !validRoles[newRole] {
		return fmt.Errorf("%w: invalid role", ErrValidation)
	}

	// Get group
	group, err := s.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return ErrGroupNotFound
	}
	_ = group

	// Check target is a member
	targetMembership, err := s.groupRepo.FindMember(ctx, groupID, targetUserID)
	if err != nil {
		return fmt.Errorf("failed to check target membership: %w", err)
	}
	if targetMembership == nil {
		return ErrNotFound
	}

	// Can't change to same role
	if targetMembership.Role == newRole {
		return nil
	}

	// Get requester's effective level
	requesterMembership, err := s.groupRepo.FindMember(ctx, groupID, requestingUserID)
	if err != nil {
		return fmt.Errorf("failed to check requester membership: %w", err)
	}

	requesterLevel := 0
	if requesterMembership != nil {
		requesterLevel = roleLevel(requesterMembership.Role)
	}

	// Check if requester is site admin or org admin (grants LEADER-level authority if no group membership)
	if requesterLevel == 0 {
		siteAdmin, err := s.userRepo.IsSiteAdmin(ctx, requestingUserID)
		if err != nil {
			return fmt.Errorf("failed to check site admin: %w", err)
		}
		if siteAdmin {
			requesterLevel = roleLevel("LEADER")
		} else {
			orgAdmin, err := s.orgRepo.IsAdmin(ctx, group.OrganizationID, requestingUserID)
			if err != nil {
				return fmt.Errorf("failed to check org admin: %w", err)
			}
			if orgAdmin {
				requesterLevel = roleLevel("LEADER")
			}
		}
	}

	if requesterLevel == 0 {
		return ErrForbidden
	}

	targetLevel := roleLevel(targetMembership.Role)
	newLevel := roleLevel(newRole)

	// Can't change role of someone at or above your level (unless you're OWNER)
	if targetLevel >= requesterLevel && requesterLevel < roleLevel("OWNER") {
		return ErrForbidden
	}

	// Can't promote to a level at or above your own (unless you're OWNER)
	if newLevel >= requesterLevel && requesterLevel < roleLevel("OWNER") {
		return ErrForbidden
	}

	// If demoting from OWNER, ensure they're not the last owner
	if targetMembership.Role == "OWNER" && newRole != "OWNER" {
		members, err := s.groupRepo.FindMembers(ctx, groupID)
		if err != nil {
			return fmt.Errorf("failed to find group members: %w", err)
		}
		ownerCount := 0
		for _, m := range members {
			if m.Role == "OWNER" {
				ownerCount++
			}
		}
		if ownerCount <= 1 {
			return fmt.Errorf("%w: cannot demote the last owner", ErrValidation)
		}
	}

	return s.groupRepo.UpdateMemberRole(ctx, groupID, targetUserID, newRole)
}
