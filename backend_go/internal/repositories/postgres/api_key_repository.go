package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"menunderfire/internal/models"
	"menunderfire/internal/repositories"
)

// apiKeyRepository is the PostgreSQL implementation of APIKeyRepository
type apiKeyRepository struct {
	db *sql.DB
}

// NewAPIKeyRepository creates a new instance of APIKeyRepository
func NewAPIKeyRepository(db *sql.DB) repositories.APIKeyRepository {
	return &apiKeyRepository{db: db}
}

// Create creates a new API key
func (r *apiKeyRepository) Create(ctx context.Context, userID, keyHash, name string, permissions []string, expiresAt *time.Time) (*models.APIKey, error) {
	query := `
		INSERT INTO api_keys (user_id, key_hash, name, permissions, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, key_hash, name, permissions, expires_at, created_at
	`

	// Convert permissions to JSON
	permissionsJSON, err := json.Marshal(permissions)
	if err != nil {
		return nil, fmt.Errorf("error marshaling permissions: %w", err)
	}

	var apiKey models.APIKey
	var permissionsBytes []byte
	var nullExpiresAt sql.NullTime
	if expiresAt != nil {
		nullExpiresAt = sql.NullTime{Time: *expiresAt, Valid: true}
	}

	err = r.db.QueryRowContext(ctx, query, userID, keyHash, name, permissionsJSON, nullExpiresAt).Scan(
		&apiKey.ID,
		&apiKey.UserID,
		&apiKey.KeyHash,
		&apiKey.Name,
		&permissionsBytes,
		&nullExpiresAt,
		&apiKey.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error creating API key: %w", err)
	}

	// Parse permissions from JSON
	if err := json.Unmarshal(permissionsBytes, &apiKey.Permissions); err != nil {
		return nil, fmt.Errorf("error unmarshaling permissions: %w", err)
	}

	if nullExpiresAt.Valid {
		apiKey.ExpiresAt = &nullExpiresAt.Time
	}

	return &apiKey, nil
}

// FindByID retrieves an API key by its ID
func (r *apiKeyRepository) FindByID(ctx context.Context, keyID string) (*models.APIKey, error) {
	query := `
		SELECT id, user_id, key_hash, name, permissions, expires_at, created_at
		FROM api_keys
		WHERE id = $1
	`

	var apiKey models.APIKey
	var permissionsBytes []byte
	var nullExpiresAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, keyID).Scan(
		&apiKey.ID,
		&apiKey.UserID,
		&apiKey.KeyHash,
		&apiKey.Name,
		&permissionsBytes,
		&nullExpiresAt,
		&apiKey.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("API key not found: %s", keyID)
	}
	if err != nil {
		return nil, fmt.Errorf("error finding API key by ID: %w", err)
	}

	// Parse permissions from JSON
	if err := json.Unmarshal(permissionsBytes, &apiKey.Permissions); err != nil {
		return nil, fmt.Errorf("error unmarshaling permissions: %w", err)
	}

	if nullExpiresAt.Valid {
		apiKey.ExpiresAt = &nullExpiresAt.Time
	}

	return &apiKey, nil
}

// FindByUserID retrieves all API keys for a user
func (r *apiKeyRepository) FindByUserID(ctx context.Context, userID string) ([]*models.APIKey, error) {
	query := `
		SELECT id, user_id, key_hash, name, permissions, expires_at, created_at
		FROM api_keys
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error finding API keys by user ID: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var apiKeys []*models.APIKey
	for rows.Next() {
		var apiKey models.APIKey
		var permissionsBytes []byte
		var nullExpiresAt sql.NullTime

		if err := rows.Scan(
			&apiKey.ID,
			&apiKey.UserID,
			&apiKey.KeyHash,
			&apiKey.Name,
			&permissionsBytes,
			&nullExpiresAt,
			&apiKey.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning API key: %w", err)
		}

		// Parse permissions from JSON
		if err := json.Unmarshal(permissionsBytes, &apiKey.Permissions); err != nil {
			return nil, fmt.Errorf("error unmarshaling permissions: %w", err)
		}

		if nullExpiresAt.Valid {
			apiKey.ExpiresAt = &nullExpiresAt.Time
		}

		apiKeys = append(apiKeys, &apiKey)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating API keys: %w", err)
	}

	return apiKeys, nil
}

// FindByKeyHash retrieves an API key by its hash
func (r *apiKeyRepository) FindByKeyHash(ctx context.Context, keyHash string) (*models.APIKey, error) {
	query := `
		SELECT id, user_id, key_hash, name, permissions, expires_at, created_at
		FROM api_keys
		WHERE key_hash = $1
	`

	var apiKey models.APIKey
	var permissionsBytes []byte
	var nullExpiresAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, keyHash).Scan(
		&apiKey.ID,
		&apiKey.UserID,
		&apiKey.KeyHash,
		&apiKey.Name,
		&permissionsBytes,
		&nullExpiresAt,
		&apiKey.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("API key not found with hash: %s", keyHash)
	}
	if err != nil {
		return nil, fmt.Errorf("error finding API key by hash: %w", err)
	}

	// Parse permissions from JSON
	if err := json.Unmarshal(permissionsBytes, &apiKey.Permissions); err != nil {
		return nil, fmt.Errorf("error unmarshaling permissions: %w", err)
	}

	if nullExpiresAt.Valid {
		apiKey.ExpiresAt = &nullExpiresAt.Time
	}

	return &apiKey, nil
}

// FindAll retrieves all API keys
func (r *apiKeyRepository) FindAll(ctx context.Context) ([]*models.APIKey, error) {
	query := `
		SELECT id, user_id, key_hash, name, permissions, expires_at, created_at
		FROM api_keys
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error finding all API keys: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var apiKeys []*models.APIKey
	for rows.Next() {
		var apiKey models.APIKey
		var permissionsBytes []byte
		var nullExpiresAt sql.NullTime

		if err := rows.Scan(
			&apiKey.ID,
			&apiKey.UserID,
			&apiKey.KeyHash,
			&apiKey.Name,
			&permissionsBytes,
			&nullExpiresAt,
			&apiKey.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning API key: %w", err)
		}

		// Parse permissions from JSON
		if err := json.Unmarshal(permissionsBytes, &apiKey.Permissions); err != nil {
			return nil, fmt.Errorf("error unmarshaling permissions: %w", err)
		}

		if nullExpiresAt.Valid {
			apiKey.ExpiresAt = &nullExpiresAt.Time
		}

		apiKeys = append(apiKeys, &apiKey)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating API keys: %w", err)
	}

	return apiKeys, nil
}

// Delete deletes an API key
func (r *apiKeyRepository) Delete(ctx context.Context, keyID string) error {
	query := `
		DELETE FROM api_keys
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, keyID)
	if err != nil {
		return fmt.Errorf("error deleting API key: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("API key not found: %s", keyID)
	}

	return nil
}

// UpdateLastUsed updates the last used timestamp for an API key
func (r *apiKeyRepository) UpdateLastUsed(ctx context.Context, keyID string) error {
	query := `
		UPDATE api_keys
		SET last_used_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, keyID)
	if err != nil {
		return fmt.Errorf("error updating API key last used: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("API key not found: %s", keyID)
	}

	return nil
}
