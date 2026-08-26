package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/lib/pq"

	"menunderfire/internal/models"
	"menunderfire/internal/repositories"
)

// ===== FormRepository Implementation =====

type formRepository struct {
	db *sql.DB
}

func NewFormRepository(db *sql.DB) repositories.FormRepository {
	return &formRepository{db: db}
}

func (r *formRepository) FindByID(ctx context.Context, formID string) (*models.Form, error) {
	query := `
		SELECT id, organization_id, name, description, is_active, created_at, updated_at
		FROM forms
		WHERE id = $1
	`

	var form models.Form
	var description sql.NullString
	err := r.db.QueryRowContext(ctx, query, formID).Scan(
		&form.ID,
		&form.OrganizationID,
		&form.Name,
		&description,
		&form.IsActive,
		&form.CreatedAt,
		&form.UpdatedAt,
	)

	if description.Valid {
		form.Description = &description.String
	}

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("form not found: %s", formID)
	}
	if err != nil {
		return nil, fmt.Errorf("error finding form by ID: %w", err)
	}

	return &form, nil
}

func (r *formRepository) FindByOrganizationID(ctx context.Context, orgID string) ([]*models.Form, error) {
	query := `
		SELECT id, organization_id, name, description, is_active, created_at, updated_at
		FROM forms
		WHERE organization_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, orgID)
	if err != nil {
		return nil, fmt.Errorf("error finding forms by organization ID: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var forms []*models.Form
	for rows.Next() {
		var form models.Form
		var description sql.NullString
		if err := rows.Scan(
			&form.ID,
			&form.OrganizationID,
			&form.Name,
			&description,
			&form.IsActive,
			&form.CreatedAt,
			&form.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning form: %w", err)
		}

		if description.Valid {
			form.Description = &description.String
		}

		forms = append(forms, &form)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating forms: %w", err)
	}

	return forms, nil
}

func (r *formRepository) Create(ctx context.Context, orgID, name string, description *string, createdBy string) (*models.Form, error) {
	query := `
		INSERT INTO forms (organization_id, name, description, created_by_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, organization_id, name, description, is_active, created_at, updated_at
	`

	var form models.Form
	var nullDesc sql.NullString
	if description != nil {
		nullDesc = sql.NullString{String: *description, Valid: true}
	}

	err := r.db.QueryRowContext(ctx, query, orgID, name, nullDesc, createdBy).Scan(
		&form.ID,
		&form.OrganizationID,
		&form.Name,
		&nullDesc,
		&form.IsActive,
		&form.CreatedAt,
		&form.UpdatedAt,
	)

	if nullDesc.Valid {
		form.Description = &nullDesc.String
	}

	if err != nil {
		return nil, fmt.Errorf("error creating form: %w", err)
	}

	return &form, nil
}

func (r *formRepository) Update(ctx context.Context, formID, name string, description *string, isActive bool) (*models.Form, error) {
	query := `
		UPDATE forms
		SET name = $1, description = $2, is_active = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING id, organization_id, name, description, is_active, created_at, updated_at
	`

	var form models.Form
	var nullDesc sql.NullString
	if description != nil {
		nullDesc = sql.NullString{String: *description, Valid: true}
	}

	err := r.db.QueryRowContext(ctx, query, name, nullDesc, isActive, formID).Scan(
		&form.ID,
		&form.OrganizationID,
		&form.Name,
		&nullDesc,
		&form.IsActive,
		&form.CreatedAt,
		&form.UpdatedAt,
	)

	if nullDesc.Valid {
		form.Description = &nullDesc.String
	}

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("form not found: %s", formID)
	}
	if err != nil {
		return nil, fmt.Errorf("error updating form: %w", err)
	}

	return &form, nil
}

func (r *formRepository) Delete(ctx context.Context, formID string) error {
	query := `DELETE FROM forms WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, formID)
	if err != nil {
		return fmt.Errorf("error deleting form: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("form not found: %s", formID)
	}

	return nil
}

// ===== FormFieldRepository Implementation =====

type formFieldRepository struct {
	db *sql.DB
}

func NewFormFieldRepository(db *sql.DB) repositories.FormFieldRepository {
	return &formFieldRepository{db: db}
}

func (r *formFieldRepository) GetFields(ctx context.Context, formID string) ([]models.FormField, error) {
	query := `
		SELECT id, form_id, field_type, label, description, is_required, field_order, options
		FROM form_fields
		WHERE form_id = $1
		ORDER BY field_order ASC
	`

	rows, err := r.db.QueryContext(ctx, query, formID)
	if err != nil {
		return nil, fmt.Errorf("error getting form fields: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var fields []models.FormField
	for rows.Next() {
		var field models.FormField
		var description sql.NullString
		var optionsBytes []byte

		if err := rows.Scan(
			&field.ID,
			&field.FormID,
			&field.FieldType,
			&field.Label,
			&description,
			&field.IsRequired,
			&field.FieldOrder,
			&optionsBytes,
		); err != nil {
			return nil, fmt.Errorf("error scanning form field: %w", err)
		}

		if description.Valid {
			field.Description = &description.String
		}

		if optionsBytes != nil {
			if err := json.Unmarshal(optionsBytes, &field.Options); err != nil {
				return nil, fmt.Errorf("error unmarshaling options: %w", err)
			}
		}

		fields = append(fields, field)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating form fields: %w", err)
	}

	return fields, nil
}

func (r *formFieldRepository) AddField(ctx context.Context, formID string, fieldType models.FormFieldType, label string, description *string, isRequired bool, fieldOrder int, options []string) (*models.FormField, error) {
	query := `
		INSERT INTO form_fields (form_id, field_type, label, description, is_required, field_order, options)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, form_id, field_type, label, description, is_required, field_order, options
	`

	var field models.FormField
	var nullDesc sql.NullString
	if description != nil {
		nullDesc = sql.NullString{String: *description, Valid: true}
	}

	var optionsJSON []byte
	var err error
	if options != nil {
		optionsJSON, err = json.Marshal(options)
		if err != nil {
			return nil, fmt.Errorf("error marshaling options: %w", err)
		}
	}

	var optionsBytes []byte
	err = r.db.QueryRowContext(ctx, query, formID, fieldType, label, nullDesc, isRequired, fieldOrder, optionsJSON).Scan(
		&field.ID,
		&field.FormID,
		&field.FieldType,
		&field.Label,
		&nullDesc,
		&field.IsRequired,
		&field.FieldOrder,
		&optionsBytes,
	)

	if nullDesc.Valid {
		field.Description = &nullDesc.String
	}

	if optionsBytes != nil {
		if err := json.Unmarshal(optionsBytes, &field.Options); err != nil {
			return nil, fmt.Errorf("error unmarshaling options: %w", err)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("error adding form field: %w", err)
	}

	return &field, nil
}

func (r *formFieldRepository) UpdateField(ctx context.Context, fieldID, label string, description *string, isRequired bool, options []string) (*models.FormField, error) {
	query := `
		UPDATE form_fields
		SET label = $1, description = $2, is_required = $3, options = $4
		WHERE id = $5
		RETURNING id, form_id, field_type, label, description, is_required, field_order, options
	`

	var field models.FormField
	var nullDesc sql.NullString
	if description != nil {
		nullDesc = sql.NullString{String: *description, Valid: true}
	}

	var optionsJSON []byte
	var err error
	if options != nil {
		optionsJSON, err = json.Marshal(options)
		if err != nil {
			return nil, fmt.Errorf("error marshaling options: %w", err)
		}
	}

	var optionsBytes []byte
	err = r.db.QueryRowContext(ctx, query, label, nullDesc, isRequired, optionsJSON, fieldID).Scan(
		&field.ID,
		&field.FormID,
		&field.FieldType,
		&field.Label,
		&nullDesc,
		&field.IsRequired,
		&field.FieldOrder,
		&optionsBytes,
	)

	if nullDesc.Valid {
		field.Description = &nullDesc.String
	}

	if optionsBytes != nil {
		if err := json.Unmarshal(optionsBytes, &field.Options); err != nil {
			return nil, fmt.Errorf("error unmarshaling options: %w", err)
		}
	}

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("form field not found: %s", fieldID)
	}
	if err != nil {
		return nil, fmt.Errorf("error updating form field: %w", err)
	}

	return &field, nil
}

func (r *formFieldRepository) DeleteField(ctx context.Context, fieldID string) error {
	query := `DELETE FROM form_fields WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, fieldID)
	if err != nil {
		return fmt.Errorf("error deleting form field: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("form field not found: %s", fieldID)
	}

	return nil
}

func (r *formFieldRepository) ReorderFields(ctx context.Context, formID string, fieldIDs []string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}
	// Rollback is a no-op once the tx is committed; the error is not actionable.
	defer func() { _ = tx.Rollback() }()

	query := `UPDATE form_fields SET field_order = $1 WHERE id = $2 AND form_id = $3`

	for i, fieldID := range fieldIDs {
		_, err := tx.ExecContext(ctx, query, i+1, fieldID, formID)
		if err != nil {
			return fmt.Errorf("error reordering field %s: %w", fieldID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

// ===== FormAnswerRepository Implementation =====

type formAnswerRepository struct {
	db *sql.DB
}

func NewFormAnswerRepository(db *sql.DB) repositories.FormAnswerRepository {
	return &formAnswerRepository{db: db}
}

func (r *formAnswerRepository) Submit(ctx context.Context, formID, userID string, answers json.RawMessage) (*models.FormAnswer, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", err)
	}
	// Rollback is a no-op once the tx is committed; the error is not actionable.
	defer func() { _ = tx.Rollback() }()

	// Mark all previous answers as not current
	updateQuery := `
		UPDATE form_answers
		SET is_current = false
		WHERE form_id = $1 AND user_id = $2 AND is_current = true
	`
	_, err = tx.ExecContext(ctx, updateQuery, formID, userID)
	if err != nil {
		return nil, fmt.Errorf("error updating previous answers: %w", err)
	}

	// Get the next version number
	versionQuery := `
		SELECT COALESCE(MAX(version), 0) + 1
		FROM form_answers
		WHERE form_id = $1 AND user_id = $2
	`
	var version int
	err = tx.QueryRowContext(ctx, versionQuery, formID, userID).Scan(&version)
	if err != nil {
		return nil, fmt.Errorf("error getting version: %w", err)
	}

	// Insert the new answer
	insertQuery := `
		INSERT INTO form_answers (form_id, user_id, version, is_current, answers)
		VALUES ($1, $2, $3, true, $4)
		RETURNING id, form_id, user_id, version, is_current, answers, submitted_at
	`

	var answer models.FormAnswer
	err = tx.QueryRowContext(ctx, insertQuery, formID, userID, version, answers).Scan(
		&answer.ID,
		&answer.FormID,
		&answer.UserID,
		&answer.Version,
		&answer.IsCurrent,
		&answer.Answers,
		&answer.SubmittedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error submitting form answer: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return &answer, nil
}

func (r *formAnswerRepository) FindByForm(ctx context.Context, formID string) ([]*models.FormAnswer, error) {
	query := `
		SELECT
			fa.id,
			fa.form_id,
			fa.user_id,
			fa.version,
			fa.is_current,
			fa.answers,
			fa.submitted_at,
			u.id,
			u.email,
			u.display_name,
			u.external_id,
			u.created_at,
			u.updated_at
		FROM form_answers fa
		INNER JOIN users u ON fa.user_id = u.id
		WHERE fa.form_id = $1 AND fa.is_current = true
		ORDER BY fa.submitted_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, formID)
	if err != nil {
		return nil, fmt.Errorf("error finding form answers: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var answers []*models.FormAnswer
	for rows.Next() {
		var answer models.FormAnswer
		var user models.User

		if err := rows.Scan(
			&answer.ID,
			&answer.FormID,
			&answer.UserID,
			&answer.Version,
			&answer.IsCurrent,
			&answer.Answers,
			&answer.SubmittedAt,
			&user.ID,
			&user.Email,
			&user.DisplayName,
			&user.ExternalID,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning form answer: %w", err)
		}

		answer.User = &user
		answers = append(answers, &answer)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating form answers: %w", err)
	}

	return answers, nil
}

func (r *formAnswerRepository) FindByFormAndUserIDs(ctx context.Context, formID string, userIDs []string) ([]*models.FormAnswer, error) {
	query := `
		SELECT
			fa.id,
			fa.form_id,
			fa.user_id,
			fa.version,
			fa.is_current,
			fa.answers,
			fa.submitted_at,
			u.id,
			u.email,
			u.display_name,
			u.external_id,
			u.created_at,
			u.updated_at
		FROM form_answers fa
		INNER JOIN users u ON fa.user_id = u.id
		WHERE fa.form_id = $1 AND fa.is_current = true AND fa.user_id = ANY($2)
		ORDER BY fa.submitted_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, formID, pq.Array(userIDs))
	if err != nil {
		return nil, fmt.Errorf("error finding form answers by user IDs: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var answers []*models.FormAnswer
	for rows.Next() {
		var answer models.FormAnswer
		var user models.User

		if err := rows.Scan(
			&answer.ID,
			&answer.FormID,
			&answer.UserID,
			&answer.Version,
			&answer.IsCurrent,
			&answer.Answers,
			&answer.SubmittedAt,
			&user.ID,
			&user.Email,
			&user.DisplayName,
			&user.ExternalID,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning form answer: %w", err)
		}

		answer.User = &user
		answers = append(answers, &answer)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating form answers: %w", err)
	}

	return answers, nil
}

func (r *formAnswerRepository) FindCurrent(ctx context.Context, formID, userID string) (*models.FormAnswer, error) {
	query := `
		SELECT id, form_id, user_id, version, is_current, answers, submitted_at
		FROM form_answers
		WHERE form_id = $1 AND user_id = $2 AND is_current = true
	`

	var answer models.FormAnswer
	err := r.db.QueryRowContext(ctx, query, formID, userID).Scan(
		&answer.ID,
		&answer.FormID,
		&answer.UserID,
		&answer.Version,
		&answer.IsCurrent,
		&answer.Answers,
		&answer.SubmittedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error finding current answer: %w", err)
	}

	return &answer, nil
}

func (r *formAnswerRepository) FindHistory(ctx context.Context, formID, userID string) ([]*models.FormAnswer, error) {
	query := `
		SELECT id, form_id, user_id, version, is_current, answers, submitted_at
		FROM form_answers
		WHERE form_id = $1 AND user_id = $2
		ORDER BY version DESC
	`

	rows, err := r.db.QueryContext(ctx, query, formID, userID)
	if err != nil {
		return nil, fmt.Errorf("error finding form answer history: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var answers []*models.FormAnswer
	for rows.Next() {
		var answer models.FormAnswer

		if err := rows.Scan(
			&answer.ID,
			&answer.FormID,
			&answer.UserID,
			&answer.Version,
			&answer.IsCurrent,
			&answer.Answers,
			&answer.SubmittedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning form answer: %w", err)
		}

		answers = append(answers, &answer)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating form answer history: %w", err)
	}

	return answers, nil
}
