package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"menunderfire/internal/logger"
	"menunderfire/internal/models"
	"menunderfire/internal/repositories"
)

// userRepository is the PostgreSQL implementation of UserRepository
type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *sql.DB) repositories.UserRepository {
	return &userRepository{db: db}
}

// scanUserNullables handles the nullable fields from a user row scan
func scanUserNullables(user *models.User, invitedByID, invitationID, agreementVersion sql.NullString, agreementAcceptedAt sql.NullTime) {
	if invitedByID.Valid {
		user.InvitedByID = &invitedByID.String
	}
	if invitationID.Valid {
		user.InvitationID = &invitationID.String
	}
	if agreementAcceptedAt.Valid {
		user.AgreementAcceptedAt = &agreementAcceptedAt.Time
	}
	if agreementVersion.Valid {
		user.AgreementVersion = &agreementVersion.String
	}
}

// FindByID retrieves a user by their ID
func (r *userRepository) FindByID(ctx context.Context, userID string) (*models.User, error) {
	logger.Debug().Str("user_id", userID).Msg("FindByID")
	query := `
		SELECT id, email, display_name, external_id, is_site_admin, invited_by_id, invitation_id,
		       agreement_accepted_at, agreement_version, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	var invitedByID, invitationID, agreementVersion sql.NullString
	var agreementAcceptedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.DisplayName,
		&user.ExternalID,
		&user.IsSiteAdmin,
		&invitedByID,
		&invitationID,
		&agreementAcceptedAt,
		&agreementVersion,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	scanUserNullables(&user, invitedByID, invitationID, agreementVersion, agreementAcceptedAt)

	if err == sql.ErrNoRows {
		logger.Debug().Str("user_id", userID).Msg("User not found")
		return nil, fmt.Errorf("user not found: %s", userID)
	}
	if err != nil {
		logger.Error().Err(err).Str("user_id", userID).Msg("Error finding user by ID")
		return nil, fmt.Errorf("error finding user by ID: %w", err)
	}

	logger.Debug().Str("user_id", user.ID).Str("email", user.Email).Msg("Found user")
	return &user, nil
}

// FindByExternalID retrieves a user by their external authentication provider ID
func (r *userRepository) FindByExternalID(ctx context.Context, externalID string) (*models.User, error) {
	logger.Debug().Str("external_id", externalID).Msg("FindByExternalID")
	query := `
		SELECT id, email, display_name, external_id, is_site_admin, invited_by_id, invitation_id,
		       agreement_accepted_at, agreement_version, created_at, updated_at
		FROM users
		WHERE external_id = $1
	`

	var user models.User
	var invitedByID, invitationID, agreementVersion sql.NullString
	var agreementAcceptedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, externalID).Scan(
		&user.ID,
		&user.Email,
		&user.DisplayName,
		&user.ExternalID,
		&user.IsSiteAdmin,
		&invitedByID,
		&invitationID,
		&agreementAcceptedAt,
		&agreementVersion,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	scanUserNullables(&user, invitedByID, invitationID, agreementVersion, agreementAcceptedAt)

	if err == sql.ErrNoRows {
		logger.Debug().Str("external_id", externalID).Msg("No user found with external ID")
		return nil, fmt.Errorf("user not found with external ID: %s", externalID)
	}
	if err != nil {
		logger.Error().Err(err).Str("external_id", externalID).Msg("Error finding user by external ID")
		return nil, fmt.Errorf("error finding user by external ID: %w", err)
	}

	logger.Debug().Str("user_id", user.ID).Str("email", user.Email).Msg("Found user by external ID")
	return &user, nil
}

// FindByEmail retrieves a user by their email address
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, display_name, external_id, is_site_admin, invited_by_id, invitation_id,
		       agreement_accepted_at, agreement_version, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user models.User
	var invitedByID, invitationID, agreementVersion sql.NullString
	var agreementAcceptedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.DisplayName,
		&user.ExternalID,
		&user.IsSiteAdmin,
		&invitedByID,
		&invitationID,
		&agreementAcceptedAt,
		&agreementVersion,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	scanUserNullables(&user, invitedByID, invitationID, agreementVersion, agreementAcceptedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found with email: %s", email)
	}
	if err != nil {
		return nil, fmt.Errorf("error finding user by email: %w", err)
	}

	return &user, nil
}

// Count returns the total number of users
func (r *userRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM users`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting users: %w", err)
	}

	return count, nil
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, email, displayName, externalID string) (*models.User, error) {
	logger.Debug().Str("email", email).Str("display_name", displayName).Str("external_id", externalID).Msg("Creating user")
	query := `
		INSERT INTO users (email, display_name, external_id)
		VALUES ($1, $2, $3)
		RETURNING id, email, display_name, external_id, is_site_admin, invited_by_id, invitation_id,
		          agreement_accepted_at, agreement_version, created_at, updated_at
	`

	var user models.User
	var invitedByID, invitationID, agreementVersion sql.NullString
	var agreementAcceptedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, email, displayName, externalID).Scan(
		&user.ID,
		&user.Email,
		&user.DisplayName,
		&user.ExternalID,
		&user.IsSiteAdmin,
		&invitedByID,
		&invitationID,
		&agreementAcceptedAt,
		&agreementVersion,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	scanUserNullables(&user, invitedByID, invitationID, agreementVersion, agreementAcceptedAt)

	if err != nil {
		logger.Error().Err(err).Msg("Error creating user")
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	logger.Info().Str("user_id", user.ID).Str("email", user.Email).Msg("User created successfully")
	return &user, nil
}

// CreateAsSiteAdmin creates a new user with site admin privileges
func (r *userRepository) CreateAsSiteAdmin(ctx context.Context, email, displayName, externalID string) (*models.User, error) {
	logger.Debug().Str("email", email).Str("display_name", displayName).Str("external_id", externalID).Msg("Creating user as site admin")
	query := `
		INSERT INTO users (email, display_name, external_id, is_site_admin)
		VALUES ($1, $2, $3, true)
		RETURNING id, email, display_name, external_id, is_site_admin, invited_by_id, invitation_id,
		          agreement_accepted_at, agreement_version, created_at, updated_at
	`

	var user models.User
	var invitedByID, invitationID, agreementVersion sql.NullString
	var agreementAcceptedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, email, displayName, externalID).Scan(
		&user.ID,
		&user.Email,
		&user.DisplayName,
		&user.ExternalID,
		&user.IsSiteAdmin,
		&invitedByID,
		&invitationID,
		&agreementAcceptedAt,
		&agreementVersion,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	scanUserNullables(&user, invitedByID, invitationID, agreementVersion, agreementAcceptedAt)

	if err != nil {
		logger.Error().Err(err).Msg("Error creating user as site admin")
		return nil, fmt.Errorf("error creating user as site admin: %w", err)
	}

	logger.Info().Str("user_id", user.ID).Str("email", user.Email).Msg("Site admin user created successfully")
	return &user, nil
}

// Update updates a user's display name
func (r *userRepository) Update(ctx context.Context, userID, displayName string) (*models.User, error) {
	query := `
		UPDATE users
		SET display_name = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, email, display_name, external_id, is_site_admin, invited_by_id, invitation_id,
		          agreement_accepted_at, agreement_version, created_at, updated_at
	`

	var user models.User
	var invitedByID, invitationID, agreementVersion sql.NullString
	var agreementAcceptedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, displayName, userID).Scan(
		&user.ID,
		&user.Email,
		&user.DisplayName,
		&user.ExternalID,
		&user.IsSiteAdmin,
		&invitedByID,
		&invitationID,
		&agreementAcceptedAt,
		&agreementVersion,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	scanUserNullables(&user, invitedByID, invitationID, agreementVersion, agreementAcceptedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found: %s", userID)
	}
	if err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return &user, nil
}

// UpdateExternalID updates a user's external authentication provider ID
func (r *userRepository) UpdateExternalID(ctx context.Context, userID, externalID string) error {
	query := `
		UPDATE users
		SET external_id = $1, updated_at = NOW()
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, externalID, userID)
	if err != nil {
		logger.Error().Err(err).Str("user_id", userID).Msg("Error updating external ID")
		return fmt.Errorf("error updating external ID: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user not found: %s", userID)
	}

	logger.Info().Str("user_id", userID).Str("external_id", externalID).Msg("External ID updated")
	return nil
}

// UpdateInvitationInfo sets the invited_by_id and invitation_id on a user
func (r *userRepository) UpdateInvitationInfo(ctx context.Context, userID, invitedByID, invitationID string) error {
	query := `
		UPDATE users
		SET invited_by_id = $1, invitation_id = $2, updated_at = NOW()
		WHERE id = $3
	`

	result, err := r.db.ExecContext(ctx, query, invitedByID, invitationID, userID)
	if err != nil {
		logger.Error().Err(err).Str("user_id", userID).Msg("Error updating invitation info")
		return fmt.Errorf("error updating invitation info: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user not found: %s", userID)
	}

	return nil
}

// RecordAgreementAcceptance stores the user's agreement acceptance with HMAC signature
func (r *userRepository) RecordAgreementAcceptance(ctx context.Context, userID, version, signature, ipAddress, userAgent string) error {
	query := `
		UPDATE users
		SET agreement_accepted_at = NOW(),
		    agreement_version = $1,
		    agreement_signature = $2,
		    agreement_ip = $3,
		    agreement_user_agent = $4,
		    updated_at = NOW()
		WHERE id = $5
	`

	result, err := r.db.ExecContext(ctx, query, version, signature, ipAddress, userAgent, userID)
	if err != nil {
		return fmt.Errorf("error recording agreement acceptance: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user not found: %s", userID)
	}

	logger.Info().Str("user_id", userID).Str("version", version).Msg("Agreement acceptance recorded")
	return nil
}

// IsSiteAdmin checks if a user has site admin privileges
func (r *userRepository) IsSiteAdmin(ctx context.Context, userID string) (bool, error) {
	query := `
		SELECT is_site_admin
		FROM users
		WHERE id = $1
	`

	var isSiteAdmin bool
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&isSiteAdmin)

	if err == sql.ErrNoRows {
		return false, fmt.Errorf("user not found: %s", userID)
	}
	if err != nil {
		return false, fmt.Errorf("error checking site admin status: %w", err)
	}

	return isSiteAdmin, nil
}

// FindAdminOrganizationIDs returns the IDs of organizations where the user is an admin
func (r *userRepository) FindAdminOrganizationIDs(ctx context.Context, userID string) ([]string, error) {
	query := `
		SELECT organization_id
		FROM organization_admins
		WHERE user_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error finding admin organization IDs: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var orgIDs []string
	for rows.Next() {
		var orgID string
		if err := rows.Scan(&orgID); err != nil {
			return nil, fmt.Errorf("error scanning organization ID: %w", err)
		}
		orgIDs = append(orgIDs, orgID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating organization IDs: %w", err)
	}

	return orgIDs, nil
}

// FindOwnedGroupIDs returns the IDs of groups where the user is an owner
func (r *userRepository) FindOwnedGroupIDs(ctx context.Context, userID string) ([]string, error) {
	query := `
		SELECT id
		FROM groups
		WHERE owner_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error finding owned group IDs: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var groupIDs []string
	for rows.Next() {
		var groupID string
		if err := rows.Scan(&groupID); err != nil {
			return nil, fmt.Errorf("error scanning group ID: %w", err)
		}
		groupIDs = append(groupIDs, groupID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating group IDs: %w", err)
	}

	return groupIDs, nil
}

// FindMemberGroupIDs returns the IDs of groups where the user is a member
func (r *userRepository) FindMemberGroupIDs(ctx context.Context, userID string) ([]string, error) {
	query := `
		SELECT group_id
		FROM group_memberships
		WHERE user_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error finding member group IDs: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var groupIDs []string
	for rows.Next() {
		var groupID string
		if err := rows.Scan(&groupID); err != nil {
			return nil, fmt.Errorf("error scanning group ID: %w", err)
		}
		groupIDs = append(groupIDs, groupID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating group IDs: %w", err)
	}

	return groupIDs, nil
}
