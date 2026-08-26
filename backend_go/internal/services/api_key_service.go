package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
	"strings"
	"time"

	"menunderfire/internal/logger"
	"menunderfire/internal/models"
	"menunderfire/internal/repositories"
)

// APIKeyService defines the interface for API key business logic
type APIKeyService interface {
	// Create creates a new API key (returns the raw key only once)
	Create(ctx context.Context, userID string, req *models.CreateAPIKeyRequest) (*models.APIKey, string, error)

	// List returns all API keys for a user
	List(ctx context.Context, userID string) ([]*models.APIKey, error)

	// Delete deletes an API key
	Delete(ctx context.Context, keyID string, userID string) error

	// GetByID retrieves an API key by its ID
	GetByID(ctx context.Context, keyID string) (*models.APIKey, error)

	// ValidateKey validates an API key and returns the associated user ID
	ValidateKey(ctx context.Context, rawKey string) (string, error)
}

const (
	// APIKeyPrefix is prepended to all generated API keys for easy identification
	APIKeyPrefix = "muf_"
	// APIKeyBytes is the number of random bytes to generate (32 bytes = 256 bits)
	APIKeyBytes = 32
	// BcryptCost is the cost factor for bcrypt hashing
	BcryptCost = 12
)

// Function variables for testing - allows injecting failures in key generation/hashing
var (
	generateAPIKeyFn           = generateAPIKey
	hashAPIKeyFn               = hashAPIKey
	randReadFn       io.Reader = rand.Reader
)

// apiKeyService is the concrete implementation of APIKeyService
type apiKeyService struct {
	apiKeyRepo repositories.APIKeyRepository
}

// NewAPIKeyService creates a new instance of APIKeyService with the provided repository
func NewAPIKeyService(
	apiKeyRepo repositories.APIKeyRepository,
) APIKeyService {
	return &apiKeyService{
		apiKeyRepo: apiKeyRepo,
	}
}

// Create creates a new API key and returns the APIKey object and raw key.
// The raw key is only returned once and should be displayed to the user immediately.
//
// Edge cases handled:
//   - Empty name in request: Will fail validation at the handler level (validate:"required")
//     but this implementation would store it if passed through
//   - Name too long (>100 chars): Will fail validation at the handler level (validate:"max=100")
//     but this implementation would store it if passed through
//   - Nil permissions: Defaults to empty slice to ensure consistent JSON serialization
//   - Zero or negative expiresInDays: No expiration is set (expiresAt remains nil)
//   - Positive expiresInDays: Calculates expiration date from current time
//   - Database error when creating key: Returns error from repository, propagates to caller
//   - Random key generation failure: Returns error if crypto/rand fails (extremely rare)
//   - Bcrypt hashing failure: Returns error if bcrypt fails (e.g., key too long, though unlikely)
func (s *apiKeyService) Create(ctx context.Context, userID string, req *models.CreateAPIKeyRequest) (*models.APIKey, string, error) {
	// Generate a secure random key
	rawKey, err := generateAPIKeyFn()
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate API key: %w", err)
	}

	// Hash the key for secure storage
	keyHash, err := hashAPIKeyFn(rawKey)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash API key: %w", err)
	}

	// Handle nil permissions by defaulting to empty slice
	permissions := req.Permissions
	if permissions == nil {
		permissions = []string{}
	}

	// Calculate expiration date if expiresInDays is positive
	var expiresAt *time.Time
	if req.ExpiresInDays > 0 {
		expiry := time.Now().AddDate(0, 0, req.ExpiresInDays)
		expiresAt = &expiry
	}

	// Create the API key in the database
	apiKey, err := s.apiKeyRepo.Create(ctx, userID, keyHash, req.Name, permissions, expiresAt)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create API key in database: %w", err)
	}

	return apiKey, rawKey, nil
}

// List retrieves all API keys for a user
//
// Edge cases handled:
//   - User has no API keys: Returns empty slice, not an error
//   - Empty userID: Returns empty slice (no keys will match)
//   - Database error: Returns wrapped error
//   - User does not exist: Returns empty slice (not an error - no keys for that user)
func (s *apiKeyService) List(ctx context.Context, userID string) ([]*models.APIKey, error) {
	keys, err := s.apiKeyRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list API keys for user %s: %w", userID, err)
	}

	return keys, nil
}

// Delete deletes an API key owned by the user
//
// Edge cases handled:
//   - Key not found: Returns ErrNotFound
//   - Key exists but owned by different user: Returns ErrNotFound (don't reveal existence)
//   - Empty keyID: Returns ErrNotFound
//   - Empty userID: Returns ErrNotFound (no key will match ownership check)
//   - Database error when finding key: Returns wrapped error
//   - Database error when deleting key: Returns wrapped error
//   - Key already deleted (concurrent delete): Returns ErrNotFound
func (s *apiKeyService) Delete(ctx context.Context, keyID string, userID string) error {
	// Find the key first
	key, err := s.apiKeyRepo.FindByID(ctx, keyID)
	if err != nil {
		return err
	}

	// Security: don't reveal key existence if owned by different user
	if key.UserID != userID {
		return ErrNotFound
	}

	// Delete the key
	if err := s.apiKeyRepo.Delete(ctx, keyID); err != nil {
		return err
	}

	return nil
}

// GetByID retrieves an API key by its ID
//
// Edge cases handled:
//   - Key not found: Returns ErrNotFound
//   - Empty keyID: Returns ErrNotFound
//   - Database error: Returns wrapped error
//   - Key exists: Returns the key
func (s *apiKeyService) GetByID(ctx context.Context, keyID string) (*models.APIKey, error) {
	if keyID == "" {
		return nil, ErrNotFound
	}

	key, err := s.apiKeyRepo.FindByID(ctx, keyID)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// ValidateKey validates an API key and returns the associated user ID
//
// Edge cases handled:
//   - Empty rawKey: Returns ErrInvalidToken
//   - Key not found in database: Returns ErrInvalidToken
//   - Key found but expired: Returns ErrInvalidToken
//   - Key found but hash doesn't match: Returns ErrInvalidToken
//   - Database error when finding keys: Returns wrapped error
//   - Database error when updating lastUsedAt: Logs error but still returns success (non-critical)
//   - Valid key: Updates lastUsedAt and returns userID
//   - Malformed key (missing "muf_" prefix): Returns ErrInvalidToken
func (s *apiKeyService) ValidateKey(ctx context.Context, rawKey string) (string, error) {
	// Edge case: empty key
	if rawKey == "" {
		return "", ErrInvalidToken
	}

	// Edge case: malformed key (must have "muf_" prefix)
	if !strings.HasPrefix(rawKey, APIKeyPrefix) {
		return "", ErrInvalidToken
	}

	// Since bcrypt generates different hashes each time (due to salting),
	// we need to iterate through all API keys and use bcrypt.CompareHashAndPassword
	// to find a match. This is the correct bcrypt usage pattern.
	allKeys, err := s.apiKeyRepo.FindAll(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to fetch API keys: %w", err)
	}

	// Find matching key by bcrypt comparison
	var matchedKey *models.APIKey
	for _, key := range allKeys {
		if err := bcrypt.CompareHashAndPassword([]byte(key.KeyHash), []byte(rawKey)); err == nil {
			matchedKey = key
			break
		}
	}

	// Edge case: key not found or hash doesn't match
	if matchedKey == nil {
		return "", ErrInvalidToken
	}

	// Edge case: key is expired
	if matchedKey.ExpiresAt != nil && matchedKey.ExpiresAt.Before(time.Now()) {
		return "", ErrInvalidToken
	}

	// Update last used timestamp (non-critical, log error but don't fail)
	if err := s.apiKeyRepo.UpdateLastUsed(ctx, matchedKey.ID); err != nil {
		logger.Warn().Err(err).Str("key_id", matchedKey.ID).Msg("Failed to update lastUsedAt for key")
	}

	return matchedKey.UserID, nil
}

// generateAPIKey generates a cryptographically secure random API key
// with the "muf_" prefix for easy identification
func generateAPIKey() (string, error) {
	// Generate 32 random bytes
	bytes := make([]byte, APIKeyBytes)
	if _, err := io.ReadFull(randReadFn, bytes); err != nil {
		return "", fmt.Errorf("crypto/rand read failed: %w", err)
	}

	// Encode to base64 and prepend prefix
	encoded := base64.URLEncoding.EncodeToString(bytes)
	return APIKeyPrefix + encoded, nil
}

// hashAPIKey hashes an API key using bcrypt for secure storage
func hashAPIKey(rawKey string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(rawKey), BcryptCost)
	if err != nil {
		return "", fmt.Errorf("bcrypt hashing failed: %w", err)
	}
	return string(hash), nil
}
