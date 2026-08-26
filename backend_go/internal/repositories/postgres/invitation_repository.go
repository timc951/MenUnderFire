package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"menunderfire/internal/models"
	"menunderfire/internal/repositories"
)

// invitationRepository is the PostgreSQL implementation of InvitationRepository
type invitationRepository struct {
	db *sql.DB
}

// NewInvitationRepository creates a new instance of InvitationRepository
func NewInvitationRepository(db *sql.DB) repositories.InvitationRepository {
	return &invitationRepository{db: db}
}

// invitationTypeToDb converts InvitationType to database format (lowercase with underscores)
func invitationTypeToDb(t models.InvitationType) string {
	return strings.ToLower(string(t))
}

// invitationTypeFromDb converts database format to InvitationType (uppercase with underscores)
func invitationTypeFromDb(s string) models.InvitationType {
	return models.InvitationType(strings.ToUpper(s))
}

// invitationStatusToDb converts InvitationStatus to database format (lowercase)
func invitationStatusToDb(s models.InvitationStatus) string {
	return strings.ToLower(string(s))
}

// invitationStatusFromDb converts database format to InvitationStatus (uppercase)
func invitationStatusFromDb(s string) models.InvitationStatus {
	return models.InvitationStatus(strings.ToUpper(s))
}

// FindByID retrieves an invitation by its ID
func (r *invitationRepository) FindByID(ctx context.Context, invitationID string) (*models.Invitation, error) {
	query := `
		SELECT id, token, email, invitation_type, organization_id, group_id, invited_by_id, status, expires_at, created_at
		FROM invitations
		WHERE id = $1
	`

	var invitation models.Invitation
	var invType, status string
	var orgID, groupID sql.NullString

	err := r.db.QueryRowContext(ctx, query, invitationID).Scan(
		&invitation.ID,
		&invitation.Token,
		&invitation.Email,
		&invType,
		&orgID,
		&groupID,
		&invitation.InviterID,
		&status,
		&invitation.ExpiresAt,
		&invitation.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("invitation not found: %s", invitationID)
	}
	if err != nil {
		return nil, fmt.Errorf("error finding invitation by ID: %w", err)
	}

	invitation.Type = invitationTypeFromDb(invType)
	invitation.Status = invitationStatusFromDb(status)
	if orgID.Valid {
		invitation.OrganizationID = &orgID.String
	}
	if groupID.Valid {
		invitation.GroupID = &groupID.String
	}

	return &invitation, nil
}

// FindByToken retrieves an invitation by its token
func (r *invitationRepository) FindByToken(ctx context.Context, token string) (*models.Invitation, error) {
	query := `
		SELECT id, token, email, invitation_type, organization_id, group_id, invited_by_id, status, expires_at, created_at
		FROM invitations
		WHERE token = $1
	`

	var invitation models.Invitation
	var invType, status string
	var orgID, groupID sql.NullString

	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&invitation.ID,
		&invitation.Token,
		&invitation.Email,
		&invType,
		&orgID,
		&groupID,
		&invitation.InviterID,
		&status,
		&invitation.ExpiresAt,
		&invitation.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("invitation not found with token: %s", token)
	}
	if err != nil {
		return nil, fmt.Errorf("error finding invitation by token: %w", err)
	}

	invitation.Type = invitationTypeFromDb(invType)
	invitation.Status = invitationStatusFromDb(status)
	if orgID.Valid {
		invitation.OrganizationID = &orgID.String
	}
	if groupID.Valid {
		invitation.GroupID = &groupID.String
	}

	return &invitation, nil
}

// FindByEmail retrieves an invitation by email, type, and target ID
func (r *invitationRepository) FindByEmail(ctx context.Context, email string, invType models.InvitationType, targetID string) (*models.Invitation, error) {
	// Build query based on invitation type
	var query string
	dbInvType := invitationTypeToDb(invType)

	if invType == models.InvitationTypeOrgAdmin {
		query = `
			SELECT id, token, email, invitation_type, organization_id, group_id, invited_by_id, status, expires_at, created_at
			FROM invitations
			WHERE email = $1 AND invitation_type = $2 AND organization_id = $3
		`
	} else {
		query = `
			SELECT id, token, email, invitation_type, organization_id, group_id, invited_by_id, status, expires_at, created_at
			FROM invitations
			WHERE email = $1 AND invitation_type = $2 AND group_id = $3
		`
	}

	var invitation models.Invitation
	var dbInvTypeResult, status string
	var orgID, groupID sql.NullString

	err := r.db.QueryRowContext(ctx, query, email, dbInvType, targetID).Scan(
		&invitation.ID,
		&invitation.Token,
		&invitation.Email,
		&dbInvTypeResult,
		&orgID,
		&groupID,
		&invitation.InviterID,
		&status,
		&invitation.ExpiresAt,
		&invitation.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("invitation not found for email %s, type %s, target %s", email, invType, targetID)
	}
	if err != nil {
		return nil, fmt.Errorf("error finding invitation by email: %w", err)
	}

	invitation.Type = invitationTypeFromDb(dbInvTypeResult)
	invitation.Status = invitationStatusFromDb(status)
	if orgID.Valid {
		invitation.OrganizationID = &orgID.String
	}
	if groupID.Valid {
		invitation.GroupID = &groupID.String
	}

	return &invitation, nil
}

// Create creates a new invitation
func (r *invitationRepository) Create(ctx context.Context, email string, invType models.InvitationType, orgID, groupID *string, inviterID string, token string, expiresAt time.Time) (*models.Invitation, error) {
	query := `
		INSERT INTO invitations (email, invitation_type, organization_id, group_id, invited_by_id, token, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, token, email, invitation_type, organization_id, group_id, invited_by_id, status, expires_at, created_at
	`

	var invitation models.Invitation
	var dbInvType, status string
	var nullOrgID, nullGroupID sql.NullString

	if orgID != nil {
		nullOrgID = sql.NullString{String: *orgID, Valid: true}
	}
	if groupID != nil {
		nullGroupID = sql.NullString{String: *groupID, Valid: true}
	}

	dbInvTypeParam := invitationTypeToDb(invType)

	err := r.db.QueryRowContext(ctx, query, email, dbInvTypeParam, nullOrgID, nullGroupID, inviterID, token, expiresAt).Scan(
		&invitation.ID,
		&invitation.Token,
		&invitation.Email,
		&dbInvType,
		&nullOrgID,
		&nullGroupID,
		&invitation.InviterID,
		&status,
		&invitation.ExpiresAt,
		&invitation.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error creating invitation: %w", err)
	}

	invitation.Type = invitationTypeFromDb(dbInvType)
	invitation.Status = invitationStatusFromDb(status)
	if nullOrgID.Valid {
		invitation.OrganizationID = &nullOrgID.String
	}
	if nullGroupID.Valid {
		invitation.GroupID = &nullGroupID.String
	}

	return &invitation, nil
}

// UpdateStatus updates the status of an invitation.
// When acceptedByID is non-nil, also sets accepted_at and accepted_by_id.
func (r *invitationRepository) UpdateStatus(ctx context.Context, invitationID string, status models.InvitationStatus, acceptedByID *string) error {
	dbStatus := invitationStatusToDb(status)

	var result sql.Result
	var err error

	if acceptedByID != nil {
		query := `
			UPDATE invitations
			SET status = $1, accepted_at = NOW(), accepted_by_id = $2
			WHERE id = $3
		`
		result, err = r.db.ExecContext(ctx, query, dbStatus, *acceptedByID, invitationID)
	} else {
		query := `
			UPDATE invitations
			SET status = $1
			WHERE id = $2
		`
		result, err = r.db.ExecContext(ctx, query, dbStatus, invitationID)
	}

	if err != nil {
		return fmt.Errorf("error updating invitation status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("invitation not found: %s", invitationID)
	}

	return nil
}

// scanInvitation scans a single invitation row
func (r *invitationRepository) scanInvitation(scanner interface{ Scan(...interface{}) error }) (*models.Invitation, error) {
	var invitation models.Invitation
	var invType, status string
	var orgID, groupID sql.NullString

	err := scanner.Scan(
		&invitation.ID,
		&invitation.Token,
		&invitation.Email,
		&invType,
		&orgID,
		&groupID,
		&invitation.InviterID,
		&status,
		&invitation.ExpiresAt,
		&invitation.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	invitation.Type = invitationTypeFromDb(invType)
	invitation.Status = invitationStatusFromDb(status)
	if orgID.Valid {
		invitation.OrganizationID = &orgID.String
	}
	if groupID.Valid {
		invitation.GroupID = &groupID.String
	}
	return &invitation, nil
}

// scanInvitations scans multiple invitation rows
func (r *invitationRepository) scanInvitations(rows *sql.Rows) ([]*models.Invitation, error) {
	var invitations []*models.Invitation
	for rows.Next() {
		inv, err := r.scanInvitation(rows)
		if err != nil {
			return nil, fmt.Errorf("error scanning invitation row: %w", err)
		}
		invitations = append(invitations, inv)
	}
	if invitations == nil {
		invitations = []*models.Invitation{}
	}
	return invitations, nil
}

// ListAll retrieves all invitations (for site admin)
func (r *invitationRepository) ListAll(ctx context.Context) ([]*models.Invitation, error) {
	query := `
		SELECT id, token, email, invitation_type, organization_id, group_id, invited_by_id, status, expires_at, created_at
		FROM invitations
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error listing all invitations: %w", err)
	}
	defer func() { _ = rows.Close() }()
	return r.scanInvitations(rows)
}

// ListByOrganizationIDs retrieves invitations for given organization IDs
func (r *invitationRepository) ListByOrganizationIDs(ctx context.Context, orgIDs []string) ([]*models.Invitation, error) {
	if len(orgIDs) == 0 {
		return []*models.Invitation{}, nil
	}

	placeholders := make([]string, len(orgIDs))
	args := make([]interface{}, len(orgIDs))
	for i, id := range orgIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	// #nosec G201 -- only generated $N placeholders are interpolated; values are bound
	query := fmt.Sprintf(`
		SELECT id, token, email, invitation_type, organization_id, group_id, invited_by_id, status, expires_at, created_at
		FROM invitations
		WHERE organization_id IN (%s)
		   OR group_id IN (SELECT g.id FROM groups g WHERE g.organization_id IN (%s))
		ORDER BY created_at DESC
	`, strings.Join(placeholders, ","), strings.Join(placeholders, ","))

	// Both IN clauses reuse the same $1..$N placeholders, so the statement
	// takes N parameters, not 2N.
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error listing invitations by org IDs: %w", err)
	}
	defer func() { _ = rows.Close() }()
	return r.scanInvitations(rows)
}

// ListByGroupIDs retrieves invitations for given group IDs
func (r *invitationRepository) ListByGroupIDs(ctx context.Context, groupIDs []string) ([]*models.Invitation, error) {
	if len(groupIDs) == 0 {
		return []*models.Invitation{}, nil
	}

	placeholders := make([]string, len(groupIDs))
	args := make([]interface{}, len(groupIDs))
	for i, id := range groupIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	// #nosec G201 -- only generated $N placeholders are interpolated; values are bound
	query := fmt.Sprintf(`
		SELECT id, token, email, invitation_type, organization_id, group_id, invited_by_id, status, expires_at, created_at
		FROM invitations
		WHERE group_id IN (%s)
		ORDER BY created_at DESC
	`, strings.Join(placeholders, ","))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error listing invitations by group IDs: %w", err)
	}
	defer func() { _ = rows.Close() }()
	return r.scanInvitations(rows)
}

// Delete deletes an invitation
func (r *invitationRepository) Delete(ctx context.Context, invitationID string) error {
	query := `
		DELETE FROM invitations
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, invitationID)
	if err != nil {
		return fmt.Errorf("error deleting invitation: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("invitation not found: %s", invitationID)
	}

	return nil
}
