package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"menunderfire/internal/models"
	"menunderfire/internal/repositories"
)

// organizationRepository is the PostgreSQL implementation of OrganizationRepository
type organizationRepository struct {
	db *sql.DB
}

// NewOrganizationRepository creates a new instance of OrganizationRepository
func NewOrganizationRepository(db *sql.DB) repositories.OrganizationRepository {
	return &organizationRepository{db: db}
}

// FindByID retrieves an organization by its ID
func (r *organizationRepository) FindByID(ctx context.Context, orgID string) (*models.Organization, error) {
	query := `
		SELECT id, name, description, created_by_id, created_at, updated_at
		FROM organizations
		WHERE id = $1
	`

	var org models.Organization
	var description sql.NullString
	err := r.db.QueryRowContext(ctx, query, orgID).Scan(
		&org.ID,
		&org.Name,
		&description,
		&org.CreatedByID,
		&org.CreatedAt,
		&org.UpdatedAt,
	)

	if description.Valid {
		org.Description = &description.String
	}

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("organization not found: %s", orgID)
	}
	if err != nil {
		return nil, fmt.Errorf("error finding organization by ID: %w", err)
	}

	return &org, nil
}

// FindByUserID retrieves all organizations a user belongs to as an admin
func (r *organizationRepository) FindByUserID(ctx context.Context, userID string) ([]*models.Organization, error) {
	query := `
		SELECT DISTINCT o.id, o.name, o.description, o.created_by_id, o.created_at, o.updated_at
		FROM organizations o
		INNER JOIN organization_admins oa ON o.id = oa.organization_id
		WHERE oa.user_id = $1
		ORDER BY o.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error finding organizations by user ID: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var organizations []*models.Organization
	for rows.Next() {
		var org models.Organization
		var description sql.NullString
		if err := rows.Scan(
			&org.ID,
			&org.Name,
			&description,
			&org.CreatedByID,
			&org.CreatedAt,
			&org.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning organization: %w", err)
		}

		if description.Valid {
			org.Description = &description.String
		}

		organizations = append(organizations, &org)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating organizations: %w", err)
	}

	return organizations, nil
}

// FindAll retrieves all organizations
func (r *organizationRepository) FindAll(ctx context.Context) ([]*models.Organization, error) {
	query := `
		SELECT id, name, description, created_by_id, created_at, updated_at
		FROM organizations
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error finding all organizations: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var organizations []*models.Organization
	for rows.Next() {
		var org models.Organization
		var description sql.NullString
		if err := rows.Scan(
			&org.ID,
			&org.Name,
			&description,
			&org.CreatedByID,
			&org.CreatedAt,
			&org.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning organization: %w", err)
		}

		if description.Valid {
			org.Description = &description.String
		}

		organizations = append(organizations, &org)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating organizations: %w", err)
	}

	return organizations, nil
}

// Create creates a new organization
func (r *organizationRepository) Create(ctx context.Context, name string, description *string, createdByID string) (*models.Organization, error) {
	query := `
		INSERT INTO organizations (name, description, created_by_id)
		VALUES ($1, $2, $3)
		RETURNING id, name, description, created_by_id, created_at, updated_at
	`

	var org models.Organization
	var nullDesc sql.NullString
	if description != nil {
		nullDesc = sql.NullString{String: *description, Valid: true}
	}

	err := r.db.QueryRowContext(ctx, query, name, nullDesc, createdByID).Scan(
		&org.ID,
		&org.Name,
		&nullDesc,
		&org.CreatedByID,
		&org.CreatedAt,
		&org.UpdatedAt,
	)

	if nullDesc.Valid {
		org.Description = &nullDesc.String
	}

	if err != nil {
		return nil, fmt.Errorf("error creating organization: %w", err)
	}

	return &org, nil
}

// Update updates an organization's name and description
func (r *organizationRepository) Update(ctx context.Context, orgID string, name string, description *string) (*models.Organization, error) {
	query := `
		UPDATE organizations
		SET name = $1, description = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id, name, description, created_by_id, created_at, updated_at
	`

	var org models.Organization
	var nullDesc sql.NullString
	if description != nil {
		nullDesc = sql.NullString{String: *description, Valid: true}
	}

	err := r.db.QueryRowContext(ctx, query, name, nullDesc, orgID).Scan(
		&org.ID,
		&org.Name,
		&nullDesc,
		&org.CreatedByID,
		&org.CreatedAt,
		&org.UpdatedAt,
	)

	if nullDesc.Valid {
		org.Description = &nullDesc.String
	}

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("organization not found: %s", orgID)
	}
	if err != nil {
		return nil, fmt.Errorf("error updating organization: %w", err)
	}

	return &org, nil
}

// FindAdmins retrieves all admins for an organization with user details
func (r *organizationRepository) FindAdmins(ctx context.Context, orgID string) ([]*models.OrganizationAdmin, error) {
	query := `
		SELECT
			oa.id,
			oa.user_id,
			oa.organization_id,
			oa.invited_by_id,
			oa.joined_at,
			u.id,
			u.email,
			u.display_name,
			u.external_id,
			u.is_site_admin,
			u.invited_by_id,
			u.invitation_id,
			u.created_at,
			u.updated_at
		FROM organization_admins oa
		INNER JOIN users u ON oa.user_id = u.id
		WHERE oa.organization_id = $1
		ORDER BY oa.joined_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, orgID)
	if err != nil {
		return nil, fmt.Errorf("error finding organization admins: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var admins []*models.OrganizationAdmin
	for rows.Next() {
		var admin models.OrganizationAdmin
		var user models.User
		var userInvitedByID, userInvitationID sql.NullString
		if err := rows.Scan(
			&admin.ID,
			&admin.UserID,
			&admin.OrganizationID,
			&admin.InvitedByID,
			&admin.JoinedAt,
			&user.ID,
			&user.Email,
			&user.DisplayName,
			&user.ExternalID,
			&user.IsSiteAdmin,
			&userInvitedByID,
			&userInvitationID,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning organization admin: %w", err)
		}

		if userInvitedByID.Valid {
			user.InvitedByID = &userInvitedByID.String
		}
		if userInvitationID.Valid {
			user.InvitationID = &userInvitationID.String
		}

		admin.User = &user
		admins = append(admins, &admin)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating organization admins: %w", err)
	}

	return admins, nil
}

// FindAdmin retrieves a specific admin by organization and user ID
func (r *organizationRepository) FindAdmin(ctx context.Context, orgID, userID string) (*models.OrganizationAdmin, error) {
	query := `
		SELECT id, user_id, organization_id, invited_by_id, joined_at
		FROM organization_admins
		WHERE organization_id = $1 AND user_id = $2
	`

	var admin models.OrganizationAdmin
	err := r.db.QueryRowContext(ctx, query, orgID, userID).Scan(
		&admin.ID,
		&admin.UserID,
		&admin.OrganizationID,
		&admin.InvitedByID,
		&admin.JoinedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("organization admin not found for org %s and user %s", orgID, userID)
	}
	if err != nil {
		return nil, fmt.Errorf("error finding organization admin: %w", err)
	}

	return &admin, nil
}

// AddAdmin adds a user as an admin of an organization
func (r *organizationRepository) AddAdmin(ctx context.Context, orgID, userID string) error {
	query := `
		INSERT INTO organization_admins (organization_id, user_id, invited_by_id)
		VALUES ($1, $2, $2)
		ON CONFLICT (user_id, organization_id) DO NOTHING
	`

	_, err := r.db.ExecContext(ctx, query, orgID, userID)
	if err != nil {
		return fmt.Errorf("error adding organization admin: %w", err)
	}

	return nil
}

// IsMember checks if a user is a member of an organization (via any group)
func (r *organizationRepository) IsMember(ctx context.Context, orgID, userID string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM group_memberships gm
			INNER JOIN groups g ON gm.group_id = g.id
			WHERE g.organization_id = $1 AND gm.user_id = $2
		)
	`

	var isMember bool
	err := r.db.QueryRowContext(ctx, query, orgID, userID).Scan(&isMember)
	if err != nil {
		return false, fmt.Errorf("error checking organization membership: %w", err)
	}

	return isMember, nil
}

// IsAdmin checks if a user is an admin of an organization
func (r *organizationRepository) IsAdmin(ctx context.Context, orgID, userID string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM organization_admins
			WHERE organization_id = $1 AND user_id = $2
		)
	`

	var isAdmin bool
	err := r.db.QueryRowContext(ctx, query, orgID, userID).Scan(&isAdmin)
	if err != nil {
		return false, fmt.Errorf("error checking organization admin status: %w", err)
	}

	return isAdmin, nil
}

// Count returns the total number of organizations
func (r *organizationRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM organizations`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting organizations: %w", err)
	}

	return count, nil
}
