package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"menunderfire/internal/models"
	"menunderfire/internal/repositories"
)

// groupMessageRepository is the PostgreSQL implementation of GroupMessageRepository
type groupMessageRepository struct {
	db *sql.DB
}

// NewGroupMessageRepository creates a new instance of GroupMessageRepository
func NewGroupMessageRepository(db *sql.DB) repositories.GroupMessageRepository {
	return &groupMessageRepository{db: db}
}

// FindByID retrieves a message by its ID
func (r *groupMessageRepository) FindByID(ctx context.Context, messageID string) (*models.GroupMessage, error) {
	query := `
		SELECT id, group_id, sender_id, content, notify_members, form_id, created_at
		FROM group_messages
		WHERE id = $1
	`

	var message models.GroupMessage
	var formID sql.NullString
	err := r.db.QueryRowContext(ctx, query, messageID).Scan(
		&message.ID,
		&message.GroupID,
		&message.SenderID,
		&message.Content,
		&message.NotifyMembers,
		&formID,
		&message.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("group message not found: %s", messageID)
	}
	if err != nil {
		return nil, fmt.Errorf("error finding group message by ID: %w", err)
	}

	if formID.Valid {
		message.FormID = &formID.String
	}

	return &message, nil
}

// FindByGroupID retrieves all messages for a group with sender details
func (r *groupMessageRepository) FindByGroupID(ctx context.Context, groupID string) ([]*models.GroupMessage, error) {
	query := `
		SELECT
			gm.id,
			gm.group_id,
			gm.sender_id,
			gm.content,
			gm.notify_members,
			gm.form_id,
			gm.created_at,
			u.id,
			u.email,
			u.display_name,
			u.external_id,
			u.created_at,
			u.updated_at,
			f.name
		FROM group_messages gm
		INNER JOIN users u ON gm.sender_id = u.id
		LEFT JOIN forms f ON gm.form_id = f.id
		WHERE gm.group_id = $1
		ORDER BY gm.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, groupID)
	if err != nil {
		return nil, fmt.Errorf("error finding group messages by group ID: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var messages []*models.GroupMessage
	for rows.Next() {
		var message models.GroupMessage
		var sender models.User
		var formID sql.NullString
		var formName sql.NullString

		if err := rows.Scan(
			&message.ID,
			&message.GroupID,
			&message.SenderID,
			&message.Content,
			&message.NotifyMembers,
			&formID,
			&message.CreatedAt,
			&sender.ID,
			&sender.Email,
			&sender.DisplayName,
			&sender.ExternalID,
			&sender.CreatedAt,
			&sender.UpdatedAt,
			&formName,
		); err != nil {
			return nil, fmt.Errorf("error scanning group message: %w", err)
		}

		message.Sender = &sender
		if formID.Valid {
			message.FormID = &formID.String
		}
		if formName.Valid {
			message.FormName = &formName.String
		}
		messages = append(messages, &message)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating group messages: %w", err)
	}

	return messages, nil
}

// FindFormMessagesByGroupID retrieves all messages with a form_id for a group
func (r *groupMessageRepository) FindFormMessagesByGroupID(ctx context.Context, groupID string) ([]*models.GroupMessage, error) {
	query := `
		SELECT
			gm.id,
			gm.group_id,
			gm.sender_id,
			gm.content,
			gm.notify_members,
			gm.form_id,
			gm.created_at,
			u.id,
			u.email,
			u.display_name,
			u.external_id,
			u.created_at,
			u.updated_at,
			f.name
		FROM group_messages gm
		INNER JOIN users u ON gm.sender_id = u.id
		INNER JOIN forms f ON gm.form_id = f.id
		WHERE gm.group_id = $1 AND gm.form_id IS NOT NULL
		ORDER BY gm.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, groupID)
	if err != nil {
		return nil, fmt.Errorf("error finding form messages by group ID: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var messages []*models.GroupMessage
	for rows.Next() {
		var message models.GroupMessage
		var sender models.User
		var formID sql.NullString
		var formName sql.NullString

		if err := rows.Scan(
			&message.ID,
			&message.GroupID,
			&message.SenderID,
			&message.Content,
			&message.NotifyMembers,
			&formID,
			&message.CreatedAt,
			&sender.ID,
			&sender.Email,
			&sender.DisplayName,
			&sender.ExternalID,
			&sender.CreatedAt,
			&sender.UpdatedAt,
			&formName,
		); err != nil {
			return nil, fmt.Errorf("error scanning form message: %w", err)
		}

		message.Sender = &sender
		if formID.Valid {
			message.FormID = &formID.String
		}
		if formName.Valid {
			message.FormName = &formName.String
		}
		messages = append(messages, &message)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating form messages: %w", err)
	}

	return messages, nil
}

// Create creates a new group message
func (r *groupMessageRepository) Create(ctx context.Context, groupID, senderID, content string, notifyMembers bool, formID *string) (*models.GroupMessage, error) {
	query := `
		INSERT INTO group_messages (group_id, sender_id, content, notify_members, form_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, group_id, sender_id, content, notify_members, form_id, created_at
	`

	var message models.GroupMessage
	var returnedFormID sql.NullString
	err := r.db.QueryRowContext(ctx, query, groupID, senderID, content, notifyMembers, formID).Scan(
		&message.ID,
		&message.GroupID,
		&message.SenderID,
		&message.Content,
		&message.NotifyMembers,
		&returnedFormID,
		&message.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error creating group message: %w", err)
	}

	if returnedFormID.Valid {
		message.FormID = &returnedFormID.String
	}

	return &message, nil
}

// Delete deletes a group message
func (r *groupMessageRepository) Delete(ctx context.Context, messageID string) error {
	query := `
		DELETE FROM group_messages
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, messageID)
	if err != nil {
		return fmt.Errorf("error deleting group message: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("group message not found: %s", messageID)
	}

	return nil
}
