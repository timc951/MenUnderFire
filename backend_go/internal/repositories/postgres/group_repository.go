package postgres

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base32"
	"fmt"
	"strings"

	"menunderfire/internal/models"
	"menunderfire/internal/repositories"

	"github.com/lib/pq"
)

// groupRepository is the PostgreSQL implementation of GroupRepository
type groupRepository struct {
	db *sql.DB
}

// NewGroupRepository creates a new instance of GroupRepository
func NewGroupRepository(db *sql.DB) repositories.GroupRepository {
	return &groupRepository{db: db}
}

// FindByID retrieves a group by its ID
func (r *groupRepository) FindByID(ctx context.Context, groupID string) (*models.Group, error) {
	query := `
		SELECT id, name, description, organization_id, owner_id, created_by, invite_code, invite_code_expires_at, require_post_approval, allow_anonymous_posts, created_at, updated_at
		FROM groups
		WHERE id = $1
	`

	var group models.Group
	var description, ownerID sql.NullString
	var inviteCodeExpiresAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, groupID).Scan(
		&group.ID,
		&group.Name,
		&description,
		&group.OrganizationID,
		&ownerID,
		&group.CreatedBy,
		&group.InviteCode,
		&inviteCodeExpiresAt,
		&group.RequirePostApproval,
		&group.AllowAnonymousPosts,
		&group.CreatedAt,
		&group.UpdatedAt,
	)

	if description.Valid {
		group.Description = &description.String
	}
	if ownerID.Valid {
		group.OwnerID = &ownerID.String
	}
	if inviteCodeExpiresAt.Valid {
		group.InviteCodeExpiresAt = &inviteCodeExpiresAt.Time
	}

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("group not found: %s", groupID)
	}
	if err != nil {
		return nil, fmt.Errorf("error finding group by ID: %w", err)
	}

	return &group, nil
}

// FindByInviteCode retrieves a group by its invite code
func (r *groupRepository) FindByInviteCode(ctx context.Context, inviteCode string) (*models.Group, error) {
	query := `
		SELECT id, name, description, organization_id, owner_id, created_by, invite_code, invite_code_expires_at, require_post_approval, allow_anonymous_posts, created_at, updated_at
		FROM groups
		WHERE invite_code = $1
	`

	var group models.Group
	var description, ownerID sql.NullString
	var inviteCodeExpiresAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, inviteCode).Scan(
		&group.ID,
		&group.Name,
		&description,
		&group.OrganizationID,
		&ownerID,
		&group.CreatedBy,
		&group.InviteCode,
		&inviteCodeExpiresAt,
		&group.RequirePostApproval,
		&group.AllowAnonymousPosts,
		&group.CreatedAt,
		&group.UpdatedAt,
	)

	if description.Valid {
		group.Description = &description.String
	}
	if ownerID.Valid {
		group.OwnerID = &ownerID.String
	}
	if inviteCodeExpiresAt.Valid {
		group.InviteCodeExpiresAt = &inviteCodeExpiresAt.Time
	}

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("group not found with invite code: %s", inviteCode)
	}
	if err != nil {
		return nil, fmt.Errorf("error finding group by invite code: %w", err)
	}

	return &group, nil
}

// FindByUserID retrieves all groups a user belongs to
func (r *groupRepository) FindByUserID(ctx context.Context, userID string) ([]*models.Group, error) {
	query := `
		SELECT DISTINCT g.id, g.name, g.description, g.organization_id, g.owner_id, g.created_by, g.invite_code, g.invite_code_expires_at, g.require_post_approval, g.allow_anonymous_posts, g.created_at, g.updated_at
		FROM groups g
		INNER JOIN group_memberships gm ON g.id = gm.group_id
		WHERE gm.user_id = $1
		ORDER BY g.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error finding groups by user ID: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var groups []*models.Group
	for rows.Next() {
		var group models.Group
		var description, ownerID sql.NullString
		var inviteCodeExpiresAt sql.NullTime
		if err := rows.Scan(
			&group.ID,
			&group.Name,
			&description,
			&group.OrganizationID,
			&ownerID,
			&group.CreatedBy,
			&group.InviteCode,
			&inviteCodeExpiresAt,
			&group.RequirePostApproval,
			&group.AllowAnonymousPosts,
			&group.CreatedAt,
			&group.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning group: %w", err)
		}

		if description.Valid {
			group.Description = &description.String
		}
		if ownerID.Valid {
			group.OwnerID = &ownerID.String
		}
		if inviteCodeExpiresAt.Valid {
			group.InviteCodeExpiresAt = &inviteCodeExpiresAt.Time
		}

		groups = append(groups, &group)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating groups: %w", err)
	}

	return groups, nil
}

// FindByOrganizationID retrieves all groups for an organization
func (r *groupRepository) FindByOrganizationID(ctx context.Context, orgID string) ([]*models.Group, error) {
	query := `
		SELECT id, name, description, organization_id, owner_id, created_by, invite_code, invite_code_expires_at, require_post_approval, allow_anonymous_posts, created_at, updated_at
		FROM groups
		WHERE organization_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, orgID)
	if err != nil {
		return nil, fmt.Errorf("error finding groups by organization ID: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var groups []*models.Group
	for rows.Next() {
		var group models.Group
		var description, ownerID sql.NullString
		var inviteCodeExpiresAt sql.NullTime
		if err := rows.Scan(
			&group.ID,
			&group.Name,
			&description,
			&group.OrganizationID,
			&ownerID,
			&group.CreatedBy,
			&group.InviteCode,
			&inviteCodeExpiresAt,
			&group.RequirePostApproval,
			&group.AllowAnonymousPosts,
			&group.CreatedAt,
			&group.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning group: %w", err)
		}

		if description.Valid {
			group.Description = &description.String
		}
		if ownerID.Valid {
			group.OwnerID = &ownerID.String
		}
		if inviteCodeExpiresAt.Valid {
			group.InviteCodeExpiresAt = &inviteCodeExpiresAt.Time
		}

		groups = append(groups, &group)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating groups: %w", err)
	}

	return groups, nil
}

// Create creates a new group
func (r *groupRepository) Create(ctx context.Context, name string, description *string, orgID string, inviteCode string, createdBy string) (*models.Group, error) {
	query := `
		INSERT INTO groups (name, description, organization_id, invite_code, invite_code_expires_at, created_by, owner_id)
		VALUES ($1, $2, $3, $4, NOW() + INTERVAL '30 days', $5, $5)
		RETURNING id, name, description, organization_id, owner_id, created_by, invite_code, invite_code_expires_at, require_post_approval, allow_anonymous_posts, created_at, updated_at
	`

	var group models.Group
	var nullDesc, ownerID sql.NullString
	var inviteCodeExpiresAt sql.NullTime
	if description != nil {
		nullDesc = sql.NullString{String: *description, Valid: true}
	}

	err := r.db.QueryRowContext(ctx, query, name, nullDesc, orgID, inviteCode, createdBy).Scan(
		&group.ID,
		&group.Name,
		&nullDesc,
		&group.OrganizationID,
		&ownerID,
		&group.CreatedBy,
		&group.InviteCode,
		&inviteCodeExpiresAt,
		&group.RequirePostApproval,
		&group.AllowAnonymousPosts,
		&group.CreatedAt,
		&group.UpdatedAt,
	)

	if nullDesc.Valid {
		group.Description = &nullDesc.String
	}
	if ownerID.Valid {
		group.OwnerID = &ownerID.String
	}
	if inviteCodeExpiresAt.Valid {
		group.InviteCodeExpiresAt = &inviteCodeExpiresAt.Time
	}

	if err != nil {
		return nil, fmt.Errorf("error creating group: %w", err)
	}

	return &group, nil
}

// GenerateInviteCode generates a unique invite code for a group
func (r *groupRepository) GenerateInviteCode() string {
	bytes := make([]byte, 10)
	rand.Read(bytes)
	// Use base32 encoding and take first 12 characters for readability
	encoded := base32.StdEncoding.EncodeToString(bytes)
	return strings.ToUpper(encoded[:12])
}

// FindMember retrieves a group member by group and user ID
func (r *groupRepository) FindMember(ctx context.Context, groupID, userID string) (*models.GroupMember, error) {
	query := `
		SELECT id, group_id, user_id, role, joined_at
		FROM group_memberships
		WHERE group_id = $1 AND user_id = $2
	`

	var member models.GroupMember
	err := r.db.QueryRowContext(ctx, query, groupID, userID).Scan(
		&member.ID,
		&member.GroupID,
		&member.UserID,
		&member.Role,
		&member.JoinedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Not a member - this is not an error condition
	}
	if err != nil {
		return nil, fmt.Errorf("error finding group member: %w", err)
	}

	return &member, nil
}

// FindMembers retrieves all members of a group with user details
func (r *groupRepository) FindMembers(ctx context.Context, groupID string) ([]*models.GroupMember, error) {
	query := `
		SELECT
			gm.id,
			gm.group_id,
			gm.user_id,
			gm.role,
			gm.joined_at,
			u.id,
			u.email,
			u.display_name,
			u.external_id,
			u.is_site_admin,
			u.invited_by_id,
			u.invitation_id,
			u.created_at,
			u.updated_at
		FROM group_memberships gm
		INNER JOIN users u ON gm.user_id = u.id
		WHERE gm.group_id = $1
		ORDER BY gm.joined_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, groupID)
	if err != nil {
		return nil, fmt.Errorf("error finding group members: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var members []*models.GroupMember
	for rows.Next() {
		var member models.GroupMember
		var user models.User
		var invitedByID, invitationID sql.NullString
		if err := rows.Scan(
			&member.ID,
			&member.GroupID,
			&member.UserID,
			&member.Role,
			&member.JoinedAt,
			&user.ID,
			&user.Email,
			&user.DisplayName,
			&user.ExternalID,
			&user.IsSiteAdmin,
			&invitedByID,
			&invitationID,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning group member: %w", err)
		}

		if invitedByID.Valid {
			user.InvitedByID = &invitedByID.String
		}
		if invitationID.Valid {
			user.InvitationID = &invitationID.String
		}

		member.User = &user
		members = append(members, &member)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating group members: %w", err)
	}

	return members, nil
}

// CountMembers returns the number of members in a group
func (r *groupRepository) CountMembers(ctx context.Context, groupID string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM group_memberships
		WHERE group_id = $1
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, groupID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting group members: %w", err)
	}

	return count, nil
}

// AddMember adds a member to a group with a specific role
func (r *groupRepository) AddMember(ctx context.Context, groupID, userID, role string) (*models.GroupMember, error) {
	query := `
		INSERT INTO group_memberships (group_id, user_id, role)
		VALUES ($1, $2, $3)
		RETURNING id, group_id, user_id, role, joined_at
	`

	var member models.GroupMember
	err := r.db.QueryRowContext(ctx, query, groupID, userID, role).Scan(
		&member.ID,
		&member.GroupID,
		&member.UserID,
		&member.Role,
		&member.JoinedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error adding group member: %w", err)
	}

	return &member, nil
}

// RemoveMember removes a member from a group
func (r *groupRepository) RemoveMember(ctx context.Context, groupID, userID string) error {
	query := `
		DELETE FROM group_memberships
		WHERE group_id = $1 AND user_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, groupID, userID)
	if err != nil {
		return fmt.Errorf("error removing group member: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("group member not found for group %s and user %s", groupID, userID)
	}

	return nil
}

// Count returns the total number of groups
func (r *groupRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM groups`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting groups: %w", err)
	}

	return count, nil
}

// CountByOrganizationIDs returns the number of groups within specific organizations
func (r *groupRepository) CountByOrganizationIDs(ctx context.Context, orgIDs []string) (int64, error) {
	if len(orgIDs) == 0 {
		return 0, nil
	}

	// Build the query with placeholders for the IN clause
	query := `SELECT COUNT(*) FROM groups WHERE organization_id = ANY($1)`

	var count int64
	err := r.db.QueryRowContext(ctx, query, pq.Array(orgIDs)).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting groups by organization IDs: %w", err)
	}

	return count, nil
}

// UpdateSettings updates the group's settings
func (r *groupRepository) UpdateSettings(ctx context.Context, groupID string, requirePostApproval, allowAnonymousPosts bool) error {
	query := `
		UPDATE groups
		SET require_post_approval = $2, allow_anonymous_posts = $3, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, groupID, requirePostApproval, allowAnonymousPosts)
	if err != nil {
		return fmt.Errorf("error updating group settings: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("group not found: %s", groupID)
	}

	return nil
}

// UpdateMemberRole updates a member's role within a group
func (r *groupRepository) UpdateMemberRole(ctx context.Context, groupID, userID, role string) error {
	query := `
		UPDATE group_memberships
		SET role = $3
		WHERE group_id = $1 AND user_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, groupID, userID, role)
	if err != nil {
		return fmt.Errorf("error updating member role: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("group member not found for group %s and user %s", groupID, userID)
	}

	return nil
}
