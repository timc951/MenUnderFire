package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"menunderfire/internal/models"
)

// mockAPIKeyRepo is a mock implementation of repositories.APIKeyRepository
type mockAPIKeyRepo struct {
	createFn         func(ctx context.Context, userID, keyHash, name string, permissions []string, expiresAt *time.Time) (*models.APIKey, error)
	findByIDFn       func(ctx context.Context, keyID string) (*models.APIKey, error)
	findByUserIDFn   func(ctx context.Context, userID string) ([]*models.APIKey, error)
	findByKeyHashFn  func(ctx context.Context, keyHash string) (*models.APIKey, error)
	findAllFn        func(ctx context.Context) ([]*models.APIKey, error)
	deleteFn         func(ctx context.Context, keyID string) error
	updateLastUsedFn func(ctx context.Context, keyID string) error
}

func (m *mockAPIKeyRepo) Create(ctx context.Context, userID, keyHash, name string, permissions []string, expiresAt *time.Time) (*models.APIKey, error) {
	if m.createFn != nil {
		return m.createFn(ctx, userID, keyHash, name, permissions, expiresAt)
	}
	return nil, errors.New("createFn not set")
}

func (m *mockAPIKeyRepo) FindByID(ctx context.Context, keyID string) (*models.APIKey, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, keyID)
	}
	return nil, errors.New("findByIDFn not set")
}

func (m *mockAPIKeyRepo) FindByUserID(ctx context.Context, userID string) ([]*models.APIKey, error) {
	if m.findByUserIDFn != nil {
		return m.findByUserIDFn(ctx, userID)
	}
	return nil, errors.New("findByUserIDFn not set")
}

func (m *mockAPIKeyRepo) FindByKeyHash(ctx context.Context, keyHash string) (*models.APIKey, error) {
	if m.findByKeyHashFn != nil {
		return m.findByKeyHashFn(ctx, keyHash)
	}
	return nil, errors.New("findByKeyHashFn not set")
}

func (m *mockAPIKeyRepo) FindAll(ctx context.Context) ([]*models.APIKey, error) {
	if m.findAllFn != nil {
		return m.findAllFn(ctx)
	}
	return nil, errors.New("findAllFn not set")
}

func (m *mockAPIKeyRepo) Delete(ctx context.Context, keyID string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, keyID)
	}
	return errors.New("deleteFn not set")
}

func (m *mockAPIKeyRepo) UpdateLastUsed(ctx context.Context, keyID string) error {
	if m.updateLastUsedFn != nil {
		return m.updateLastUsedFn(ctx, keyID)
	}
	return errors.New("updateLastUsedFn not set")
}

func TestAPIKeyServiceImpl_Create(t *testing.T) {
	tests := []struct {
		name            string
		userID          string
		req             *models.CreateAPIKeyRequest
		mockCreateKey   func(ctx context.Context, userID, keyHash, name string, permissions []string, expiresAt *time.Time) (*models.APIKey, error)
		wantErr         bool
		wantErrContains string
		validate        func(t *testing.T, key *models.APIKey, rawKey string)
	}{
		{
			name:   "successful creation with all fields populated",
			userID: "user-123",
			req: &models.CreateAPIKeyRequest{
				Name:          "Production API Key",
				Permissions:   []string{"read:users", "write:reports"},
				ExpiresInDays: 90,
			},
			mockCreateKey: func(ctx context.Context, userID, keyHash, name string, permissions []string, expiresAt *time.Time) (*models.APIKey, error) {
				return &models.APIKey{
					ID:          "key-456",
					UserID:      userID,
					Name:        name,
					KeyHash:     keyHash,
					Permissions: permissions,
					ExpiresAt:   expiresAt,
					CreatedAt:   time.Now(),
				}, nil
			},
			wantErr: false,
			validate: func(t *testing.T, key *models.APIKey, rawKey string) {
				if key == nil {
					t.Fatal("expected key to be non-nil")
				}
				if key.ID != "key-456" {
					t.Errorf("expected ID to be 'key-456', got %s", key.ID)
				}
				if key.UserID != "user-123" {
					t.Errorf("expected UserID to be 'user-123', got %s", key.UserID)
				}
				if key.Name != "Production API Key" {
					t.Errorf("expected Name to be 'Production API Key', got %s", key.Name)
				}
				if len(key.Permissions) != 2 {
					t.Errorf("expected 2 permissions, got %d", len(key.Permissions))
				}
				if key.ExpiresAt == nil {
					t.Error("expected ExpiresAt to be non-nil")
				}
				if rawKey == "" {
					t.Error("expected rawKey to be non-empty")
				}
				if !strings.HasPrefix(rawKey, "muf_") {
					t.Errorf("expected rawKey to have 'muf_' prefix, got %s", rawKey)
				}
				// Verify bcrypt format
				err := bcrypt.CompareHashAndPassword([]byte(key.KeyHash), []byte(rawKey))
				if err != nil {
					t.Errorf("keyHash is not valid bcrypt hash of rawKey: %v", err)
				}
			},
		},
		{
			name:   "successful creation with nil permissions defaults to empty slice",
			userID: "user-123",
			req: &models.CreateAPIKeyRequest{
				Name:          "Basic Key",
				Permissions:   nil,
				ExpiresInDays: 30,
			},
			mockCreateKey: func(ctx context.Context, userID, keyHash, name string, permissions []string, expiresAt *time.Time) (*models.APIKey, error) {
				return &models.APIKey{
					ID:          "key-789",
					UserID:      userID,
					Name:        name,
					KeyHash:     keyHash,
					Permissions: permissions,
					ExpiresAt:   expiresAt,
					CreatedAt:   time.Now(),
				}, nil
			},
			wantErr: false,
			validate: func(t *testing.T, key *models.APIKey, rawKey string) {
				if key == nil {
					t.Fatal("expected key to be non-nil")
				}
				if key.Permissions == nil {
					t.Error("expected Permissions to default to empty slice, got nil")
				}
				if len(key.Permissions) != 0 {
					t.Errorf("expected 0 permissions, got %d", len(key.Permissions))
				}
			},
		},
		{
			name:   "successful creation with zero expiresInDays has no expiration",
			userID: "user-123",
			req: &models.CreateAPIKeyRequest{
				Name:          "Never Expires Key",
				Permissions:   []string{"read:all"},
				ExpiresInDays: 0,
			},
			mockCreateKey: func(ctx context.Context, userID, keyHash, name string, permissions []string, expiresAt *time.Time) (*models.APIKey, error) {
				return &models.APIKey{
					ID:          "key-no-expire",
					UserID:      userID,
					Name:        name,
					KeyHash:     keyHash,
					Permissions: permissions,
					ExpiresAt:   expiresAt,
					CreatedAt:   time.Now(),
				}, nil
			},
			wantErr: false,
			validate: func(t *testing.T, key *models.APIKey, rawKey string) {
				if key == nil {
					t.Fatal("expected key to be non-nil")
				}
				if key.ExpiresAt != nil {
					t.Errorf("expected ExpiresAt to be nil for zero expiresInDays, got %v", key.ExpiresAt)
				}
			},
		},
		{
			name:   "successful creation with positive expiresInDays has expiration",
			userID: "user-123",
			req: &models.CreateAPIKeyRequest{
				Name:          "Temporary Key",
				Permissions:   []string{"write:data"},
				ExpiresInDays: 7,
			},
			mockCreateKey: func(ctx context.Context, userID, keyHash, name string, permissions []string, expiresAt *time.Time) (*models.APIKey, error) {
				return &models.APIKey{
					ID:          "key-temp",
					UserID:      userID,
					Name:        name,
					KeyHash:     keyHash,
					Permissions: permissions,
					ExpiresAt:   expiresAt,
					CreatedAt:   time.Now(),
				}, nil
			},
			wantErr: false,
			validate: func(t *testing.T, key *models.APIKey, rawKey string) {
				if key == nil {
					t.Fatal("expected key to be non-nil")
				}
				if key.ExpiresAt == nil {
					t.Error("expected ExpiresAt to be non-nil for positive expiresInDays")
				} else {
					expectedExpiry := time.Now().Add(7 * 24 * time.Hour)
					diff := key.ExpiresAt.Sub(expectedExpiry)
					if diff < -time.Minute || diff > time.Minute {
						t.Errorf("expected ExpiresAt to be ~7 days from now, got %v (diff: %v)", key.ExpiresAt, diff)
					}
				}
			},
		},
		{
			name:   "database error when creating key returns wrapped error",
			userID: "user-123",
			req: &models.CreateAPIKeyRequest{
				Name:          "Failed Key",
				Permissions:   []string{"admin"},
				ExpiresInDays: 30,
			},
			mockCreateKey: func(ctx context.Context, userID, keyHash, name string, permissions []string, expiresAt *time.Time) (*models.APIKey, error) {
				return nil, fmt.Errorf("database connection failed")
			},
			wantErr:         true,
			wantErrContains: "failed to create API key",
			validate:        nil,
		},
		{
			name:   "raw key has muf_ prefix and sufficient length",
			userID: "user-prefix-test",
			req: &models.CreateAPIKeyRequest{
				Name:          "Prefix Test",
				Permissions:   []string{},
				ExpiresInDays: 0,
			},
			mockCreateKey: func(ctx context.Context, userID, keyHash, name string, permissions []string, expiresAt *time.Time) (*models.APIKey, error) {
				return &models.APIKey{
					ID:          "key-prefix",
					UserID:      userID,
					Name:        name,
					KeyHash:     keyHash,
					Permissions: permissions,
					ExpiresAt:   expiresAt,
					CreatedAt:   time.Now(),
				}, nil
			},
			wantErr: false,
			validate: func(t *testing.T, key *models.APIKey, rawKey string) {
				if !strings.HasPrefix(rawKey, "muf_") {
					t.Errorf("expected rawKey to start with 'muf_', got %s", rawKey)
				}
				if len(rawKey) < 40 {
					t.Errorf("expected rawKey to be at least 40 characters, got %d", len(rawKey))
				}
			},
		},
		{
			name:   "key hash is valid bcrypt format",
			userID: "user-bcrypt-test",
			req: &models.CreateAPIKeyRequest{
				Name:          "Bcrypt Test",
				Permissions:   []string{},
				ExpiresInDays: 0,
			},
			mockCreateKey: func(ctx context.Context, userID, keyHash, name string, permissions []string, expiresAt *time.Time) (*models.APIKey, error) {
				return &models.APIKey{
					ID:          "key-bcrypt",
					UserID:      userID,
					Name:        name,
					KeyHash:     keyHash,
					Permissions: permissions,
					ExpiresAt:   expiresAt,
					CreatedAt:   time.Now(),
				}, nil
			},
			wantErr: false,
			validate: func(t *testing.T, key *models.APIKey, rawKey string) {
				if !strings.HasPrefix(key.KeyHash, "$2") {
					t.Errorf("expected KeyHash to be bcrypt format (start with $2), got %s", key.KeyHash)
				}
				err := bcrypt.CompareHashAndPassword([]byte(key.KeyHash), []byte(rawKey))
				if err != nil {
					t.Errorf("KeyHash does not match rawKey: %v", err)
				}
				if len(key.KeyHash) != 60 {
					t.Errorf("expected KeyHash to be 60 characters (bcrypt), got %d", len(key.KeyHash))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockAPIKeyRepo{
				createFn: tt.mockCreateKey,
			}
			s := NewAPIKeyService(repo)

			key, rawKey, err := s.Create(context.Background(), tt.userID, tt.req)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				if tt.wantErrContains != "" && !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("expected error to contain %q, got %q", tt.wantErrContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.validate != nil {
				tt.validate(t, key, rawKey)
			}
		})
	}
}

func TestAPIKeyServiceImpl_Create_RawKeyRandomness(t *testing.T) {
	mockCreate := func(ctx context.Context, userID, keyHash, name string, permissions []string, expiresAt *time.Time) (*models.APIKey, error) {
		return &models.APIKey{
			ID:          "key-random",
			UserID:      userID,
			Name:        name,
			KeyHash:     keyHash,
			Permissions: permissions,
			ExpiresAt:   expiresAt,
			CreatedAt:   time.Now(),
		}, nil
	}

	repo := &mockAPIKeyRepo{
		createFn: mockCreate,
	}
	s := NewAPIKeyService(repo)

	req := &models.CreateAPIKeyRequest{
		Name:          "Randomness Test",
		Permissions:   []string{},
		ExpiresInDays: 0,
	}

	keys := make(map[string]bool)
	iterations := 10

	for i := 0; i < iterations; i++ {
		_, rawKey, err := s.Create(context.Background(), "user-123", req)
		if err != nil {
			t.Fatalf("unexpected error on iteration %d: %v", i, err)
		}

		if keys[rawKey] {
			t.Errorf("duplicate raw key generated: %s", rawKey)
		}
		keys[rawKey] = true

		if !strings.HasPrefix(rawKey, "muf_") {
			t.Errorf("iteration %d: raw key missing muf_ prefix: %s", i, rawKey)
		}
	}

	if len(keys) != iterations {
		t.Errorf("expected %d unique keys, got %d", iterations, len(keys))
	}
}

func TestAPIKeyServiceImpl_List(t *testing.T) {
	tests := []struct {
		name             string
		userID           string
		mockFindByUserID func(ctx context.Context, userID string) ([]*models.APIKey, error)
		wantCount        int
		wantErr          bool
		wantErrContains  string
	}{
		{
			name:   "user has multiple API keys returns all keys",
			userID: "user-123",
			mockFindByUserID: func(ctx context.Context, userID string) ([]*models.APIKey, error) {
				return []*models.APIKey{
					{ID: "key-1", UserID: "user-123", Name: "Production Key"},
					{ID: "key-2", UserID: "user-123", Name: "Development Key"},
					{ID: "key-3", UserID: "user-123", Name: "Test Key"},
				}, nil
			},
			wantCount: 3,
			wantErr:   false,
		},
		{
			name:   "user has no API keys returns empty slice",
			userID: "user-456",
			mockFindByUserID: func(ctx context.Context, userID string) ([]*models.APIKey, error) {
				return []*models.APIKey{}, nil
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:   "empty userID returns empty slice",
			userID: "",
			mockFindByUserID: func(ctx context.Context, userID string) ([]*models.APIKey, error) {
				return []*models.APIKey{}, nil
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:   "database error returns wrapped error",
			userID: "user-789",
			mockFindByUserID: func(ctx context.Context, userID string) ([]*models.APIKey, error) {
				return nil, errors.New("connection timeout")
			},
			wantCount:       0,
			wantErr:         true,
			wantErrContains: "failed to list API keys",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockAPIKeyRepo{
				findByUserIDFn: tt.mockFindByUserID,
			}
			s := NewAPIKeyService(repo)

			got, err := s.List(context.Background(), tt.userID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("List() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.wantErrContains != "" && !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("List() error = %v, wantErrContains %v", err, tt.wantErrContains)
				}
				return
			}

			if err != nil {
				t.Errorf("List() unexpected error = %v", err)
				return
			}

			if len(got) != tt.wantCount {
				t.Errorf("List() got %d keys, want %d keys", len(got), tt.wantCount)
			}
		})
	}
}

func TestAPIKeyServiceImpl_Delete(t *testing.T) {
	tests := []struct {
		name          string
		keyID         string
		userID        string
		mockFindByID  func(ctx context.Context, keyID string) (*models.APIKey, error)
		mockDeleteKey func(ctx context.Context, keyID string) error
		wantErr       error
	}{
		{
			name:   "successful delete when key exists and belongs to user",
			keyID:  "key-123",
			userID: "user-456",
			mockFindByID: func(ctx context.Context, keyID string) (*models.APIKey, error) {
				return &models.APIKey{
					ID:     "key-123",
					UserID: "user-456",
					Name:   "Test Key",
				}, nil
			},
			mockDeleteKey: func(ctx context.Context, keyID string) error {
				return nil
			},
			wantErr: nil,
		},
		{
			name:   "key not found returns ErrNotFound",
			keyID:  "nonexistent-key",
			userID: "user-456",
			mockFindByID: func(ctx context.Context, keyID string) (*models.APIKey, error) {
				return nil, ErrNotFound
			},
			mockDeleteKey: func(ctx context.Context, keyID string) error {
				t.Fatal("deleteKey should not be called when key is not found")
				return nil
			},
			wantErr: ErrNotFound,
		},
		{
			name:   "key exists but owned by different user returns ErrNotFound",
			keyID:  "key-123",
			userID: "user-999",
			mockFindByID: func(ctx context.Context, keyID string) (*models.APIKey, error) {
				return &models.APIKey{
					ID:     "key-123",
					UserID: "user-456",
					Name:   "Test Key",
				}, nil
			},
			mockDeleteKey: func(ctx context.Context, keyID string) error {
				t.Fatal("deleteKey should not be called when user doesn't own the key")
				return nil
			},
			wantErr: ErrNotFound,
		},
		{
			name:   "database error when finding key",
			keyID:  "key-123",
			userID: "user-456",
			mockFindByID: func(ctx context.Context, keyID string) (*models.APIKey, error) {
				return nil, errors.New("database connection failed")
			},
			mockDeleteKey: func(ctx context.Context, keyID string) error {
				t.Fatal("deleteKey should not be called when findByID fails")
				return nil
			},
			wantErr: errors.New("database connection failed"),
		},
		{
			name:   "database error when deleting key",
			keyID:  "key-123",
			userID: "user-456",
			mockFindByID: func(ctx context.Context, keyID string) (*models.APIKey, error) {
				return &models.APIKey{
					ID:     "key-123",
					UserID: "user-456",
					Name:   "Test Key",
				}, nil
			},
			mockDeleteKey: func(ctx context.Context, keyID string) error {
				return errors.New("failed to delete from database")
			},
			wantErr: errors.New("failed to delete from database"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockAPIKeyRepo{
				findByIDFn: tt.mockFindByID,
				deleteFn:   tt.mockDeleteKey,
			}
			s := NewAPIKeyService(repo)

			err := s.Delete(context.Background(), tt.keyID, tt.userID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("Delete() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if !errors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error() {
					t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("Delete() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestAPIKeyServiceImpl_GetByID(t *testing.T) {
	tests := []struct {
		name         string
		keyID        string
		mockFindByID func(ctx context.Context, keyID string) (*models.APIKey, error)
		wantKey      *models.APIKey
		wantErr      error
	}{
		{
			name:  "key exists returns the key",
			keyID: "key-123",
			mockFindByID: func(ctx context.Context, keyID string) (*models.APIKey, error) {
				return &models.APIKey{
					ID:     "key-123",
					UserID: "user-456",
					Name:   "Test Key",
				}, nil
			},
			wantKey: &models.APIKey{
				ID:     "key-123",
				UserID: "user-456",
				Name:   "Test Key",
			},
			wantErr: nil,
		},
		{
			name:  "key not found returns ErrNotFound",
			keyID: "nonexistent-key",
			mockFindByID: func(ctx context.Context, keyID string) (*models.APIKey, error) {
				return nil, ErrNotFound
			},
			wantKey: nil,
			wantErr: ErrNotFound,
		},
		{
			name:  "empty keyID returns ErrNotFound",
			keyID: "",
			mockFindByID: func(ctx context.Context, keyID string) (*models.APIKey, error) {
				t.Fatal("findByID should not be called for empty keyID")
				return nil, nil
			},
			wantKey: nil,
			wantErr: ErrNotFound,
		},
		{
			name:  "database error returns error",
			keyID: "key-123",
			mockFindByID: func(ctx context.Context, keyID string) (*models.APIKey, error) {
				return nil, errors.New("database connection failed")
			},
			wantKey: nil,
			wantErr: errors.New("database connection failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockAPIKeyRepo{
				findByIDFn: tt.mockFindByID,
			}
			s := NewAPIKeyService(repo)

			got, err := s.GetByID(context.Background(), tt.keyID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("GetByID() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if !errors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error() {
					t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("GetByID() unexpected error = %v", err)
				return
			}

			if got == nil && tt.wantKey != nil {
				t.Errorf("GetByID() got = nil, want %v", tt.wantKey)
				return
			}

			if got != nil && tt.wantKey != nil {
				if got.ID != tt.wantKey.ID || got.UserID != tt.wantKey.UserID || got.Name != tt.wantKey.Name {
					t.Errorf("GetByID() got = %v, want %v", got, tt.wantKey)
				}
			}
		})
	}
}

func TestAPIKeyServiceImpl_ValidateKey(t *testing.T) {
	// Helper to create bcrypt hash for testing
	hashKey := func(key string) string {
		hash, _ := bcrypt.GenerateFromPassword([]byte(key), bcrypt.MinCost)
		return string(hash)
	}

	validRawKey := "muf_test123456789"
	validKeyHash := hashKey(validRawKey)
	testUserID := "user-123"
	testKeyID := "key-456"

	now := time.Now()
	futureExpiry := now.Add(24 * time.Hour)
	pastExpiry := now.Add(-24 * time.Hour)

	tests := []struct {
		name               string
		rawKey             string
		mockFindAllKeys    func(ctx context.Context) ([]*models.APIKey, error)
		mockUpdateLastUsed func(ctx context.Context, keyID string) error
		wantUserID         string
		wantErr            error
		wantErrContains    string
	}{
		{
			name:   "valid key returns userID and updates lastUsedAt",
			rawKey: validRawKey,
			mockFindAllKeys: func(ctx context.Context) ([]*models.APIKey, error) {
				return []*models.APIKey{
					{
						ID:        testKeyID,
						UserID:    testUserID,
						KeyHash:   validKeyHash,
						ExpiresAt: &futureExpiry,
					},
				}, nil
			},
			mockUpdateLastUsed: func(ctx context.Context, keyID string) error {
				if keyID != testKeyID {
					t.Errorf("expected keyID %s, got %s", testKeyID, keyID)
				}
				return nil
			},
			wantUserID: testUserID,
			wantErr:    nil,
		},
		{
			name:   "empty rawKey returns ErrInvalidToken",
			rawKey: "",
			mockFindAllKeys: func(ctx context.Context) ([]*models.APIKey, error) {
				t.Error("findAllKeys should not be called for empty key")
				return nil, nil
			},
			mockUpdateLastUsed: func(ctx context.Context, keyID string) error {
				t.Error("updateLastUsed should not be called for empty key")
				return nil
			},
			wantUserID: "",
			wantErr:    ErrInvalidToken,
		},
		{
			name:   "key missing muf_ prefix returns ErrInvalidToken",
			rawKey: "invalid_prefix_test123",
			mockFindAllKeys: func(ctx context.Context) ([]*models.APIKey, error) {
				t.Error("findAllKeys should not be called for invalid prefix")
				return nil, nil
			},
			mockUpdateLastUsed: func(ctx context.Context, keyID string) error {
				t.Error("updateLastUsed should not be called for invalid prefix")
				return nil
			},
			wantUserID: "",
			wantErr:    ErrInvalidToken,
		},
		{
			name:   "key not found returns ErrInvalidToken",
			rawKey: validRawKey,
			mockFindAllKeys: func(ctx context.Context) ([]*models.APIKey, error) {
				differentHash := hashKey("muf_different123")
				return []*models.APIKey{
					{
						ID:      "other-key",
						UserID:  "other-user",
						KeyHash: differentHash,
					},
				}, nil
			},
			mockUpdateLastUsed: func(ctx context.Context, keyID string) error {
				t.Error("updateLastUsed should not be called when key not found")
				return nil
			},
			wantUserID: "",
			wantErr:    ErrInvalidToken,
		},
		{
			name:   "expired key returns ErrInvalidToken",
			rawKey: validRawKey,
			mockFindAllKeys: func(ctx context.Context) ([]*models.APIKey, error) {
				return []*models.APIKey{
					{
						ID:        testKeyID,
						UserID:    testUserID,
						KeyHash:   validKeyHash,
						ExpiresAt: &pastExpiry,
					},
				}, nil
			},
			mockUpdateLastUsed: func(ctx context.Context, keyID string) error {
				t.Error("updateLastUsed should not be called for expired key")
				return nil
			},
			wantUserID: "",
			wantErr:    ErrInvalidToken,
		},
		{
			name:   "key with nil expiresAt (never expires) returns userID",
			rawKey: validRawKey,
			mockFindAllKeys: func(ctx context.Context) ([]*models.APIKey, error) {
				return []*models.APIKey{
					{
						ID:        testKeyID,
						UserID:    testUserID,
						KeyHash:   validKeyHash,
						ExpiresAt: nil,
					},
				}, nil
			},
			mockUpdateLastUsed: func(ctx context.Context, keyID string) error {
				return nil
			},
			wantUserID: testUserID,
			wantErr:    nil,
		},
		{
			name:   "database error when finding keys returns wrapped error",
			rawKey: validRawKey,
			mockFindAllKeys: func(ctx context.Context) ([]*models.APIKey, error) {
				return nil, errors.New("database connection failed")
			},
			mockUpdateLastUsed: func(ctx context.Context, keyID string) error {
				t.Error("updateLastUsed should not be called when findAllKeys fails")
				return nil
			},
			wantUserID:      "",
			wantErr:         nil,
			wantErrContains: "failed to fetch API keys",
		},
		{
			name:   "database error when updating lastUsedAt logs but still returns success",
			rawKey: validRawKey,
			mockFindAllKeys: func(ctx context.Context) ([]*models.APIKey, error) {
				return []*models.APIKey{
					{
						ID:        testKeyID,
						UserID:    testUserID,
						KeyHash:   validKeyHash,
						ExpiresAt: &futureExpiry,
					},
				}, nil
			},
			mockUpdateLastUsed: func(ctx context.Context, keyID string) error {
				return errors.New("update failed")
			},
			wantUserID: testUserID,
			wantErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockAPIKeyRepo{
				findAllFn:        tt.mockFindAllKeys,
				updateLastUsedFn: tt.mockUpdateLastUsed,
			}
			s := NewAPIKeyService(repo)

			gotUserID, err := s.ValidateKey(context.Background(), tt.rawKey)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("ValidateKey() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else if tt.wantErrContains != "" {
				if err == nil {
					t.Errorf("ValidateKey() expected error containing %q, got nil", tt.wantErrContains)
				} else if !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("ValidateKey() error = %v, want error containing %q", err, tt.wantErrContains)
				}
			} else if err != nil {
				t.Errorf("ValidateKey() unexpected error = %v", err)
			}

			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateKey() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestAPIKeyServiceImpl_Create_GenerateKeyFailure(t *testing.T) {
	// Save original and restore after test
	origFn := generateAPIKeyFn
	defer func() { generateAPIKeyFn = origFn }()

	generateAPIKeyFn = func() (string, error) {
		return "", errors.New("crypto/rand unavailable")
	}

	repo := &mockAPIKeyRepo{
		createFn: func(ctx context.Context, userID, keyHash, name string, permissions []string, expiresAt *time.Time) (*models.APIKey, error) {
			t.Fatal("create should not be called when key generation fails")
			return nil, nil
		},
	}
	s := NewAPIKeyService(repo)

	key, rawKey, err := s.Create(context.Background(), "user-123", &models.CreateAPIKeyRequest{
		Name:        "Test Key",
		Permissions: []string{"read"},
	})

	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if !strings.Contains(err.Error(), "failed to generate API key") {
		t.Errorf("expected error to contain 'failed to generate API key', got %q", err.Error())
	}
	if key != nil {
		t.Errorf("expected key to be nil, got %v", key)
	}
	if rawKey != "" {
		t.Errorf("expected rawKey to be empty, got %q", rawKey)
	}
}

func TestAPIKeyServiceImpl_Create_HashKeyFailure(t *testing.T) {
	// Save original and restore after test
	origFn := hashAPIKeyFn
	defer func() { hashAPIKeyFn = origFn }()

	hashAPIKeyFn = func(rawKey string) (string, error) {
		return "", errors.New("bcrypt internal error")
	}

	repo := &mockAPIKeyRepo{
		createFn: func(ctx context.Context, userID, keyHash, name string, permissions []string, expiresAt *time.Time) (*models.APIKey, error) {
			t.Fatal("create should not be called when hashing fails")
			return nil, nil
		},
	}
	s := NewAPIKeyService(repo)

	key, rawKey, err := s.Create(context.Background(), "user-123", &models.CreateAPIKeyRequest{
		Name:        "Test Key",
		Permissions: []string{"read"},
	})

	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if !strings.Contains(err.Error(), "failed to hash API key") {
		t.Errorf("expected error to contain 'failed to hash API key', got %q", err.Error())
	}
	if key != nil {
		t.Errorf("expected key to be nil, got %v", key)
	}
	if rawKey != "" {
		t.Errorf("expected rawKey to be empty, got %q", rawKey)
	}
}

func TestAPIKeyServiceImpl_Create_NegativeExpiresInDays(t *testing.T) {
	repo := &mockAPIKeyRepo{
		createFn: func(ctx context.Context, userID, keyHash, name string, permissions []string, expiresAt *time.Time) (*models.APIKey, error) {
			if expiresAt != nil {
				t.Error("expected expiresAt to be nil for negative expiresInDays")
			}
			return &models.APIKey{
				ID:          "key-neg",
				UserID:      userID,
				Name:        name,
				KeyHash:     keyHash,
				Permissions: permissions,
				ExpiresAt:   expiresAt,
			}, nil
		},
	}
	s := NewAPIKeyService(repo)

	key, _, err := s.Create(context.Background(), "user-123", &models.CreateAPIKeyRequest{
		Name:          "Negative Expiry Key",
		Permissions:   []string{},
		ExpiresInDays: -5,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if key.ExpiresAt != nil {
		t.Errorf("expected ExpiresAt to be nil for negative expiresInDays, got %v", key.ExpiresAt)
	}
}

func TestAPIKeyServiceImpl_Create_EmptyPermissions(t *testing.T) {
	repo := &mockAPIKeyRepo{
		createFn: func(ctx context.Context, userID, keyHash, name string, permissions []string, expiresAt *time.Time) (*models.APIKey, error) {
			if permissions == nil {
				t.Error("expected permissions to be non-nil empty slice")
			}
			return &models.APIKey{
				ID:          "key-empty-perms",
				UserID:      userID,
				Name:        name,
				KeyHash:     keyHash,
				Permissions: permissions,
			}, nil
		},
	}
	s := NewAPIKeyService(repo)

	_, _, err := s.Create(context.Background(), "user-123", &models.CreateAPIKeyRequest{
		Name:        "Empty Perms Key",
		Permissions: []string{},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAPIKeyServiceImpl_ValidateKey_MultipleKeysMatchesCorrectOne(t *testing.T) {
	hashKey := func(key string) string {
		hash, _ := bcrypt.GenerateFromPassword([]byte(key), bcrypt.MinCost)
		return string(hash)
	}

	targetKey := "muf_target_key_123"
	targetHash := hashKey(targetKey)
	otherHash1 := hashKey("muf_other_key_1")
	otherHash2 := hashKey("muf_other_key_2")

	var updatedKeyID string
	repo := &mockAPIKeyRepo{
		findAllFn: func(ctx context.Context) ([]*models.APIKey, error) {
			return []*models.APIKey{
				{ID: "key-1", UserID: "user-wrong1", KeyHash: otherHash1},
				{ID: "key-2", UserID: "user-correct", KeyHash: targetHash},
				{ID: "key-3", UserID: "user-wrong2", KeyHash: otherHash2},
			}, nil
		},
		updateLastUsedFn: func(ctx context.Context, keyID string) error {
			updatedKeyID = keyID
			return nil
		},
	}
	s := NewAPIKeyService(repo)

	userID, err := s.ValidateKey(context.Background(), targetKey)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if userID != "user-correct" {
		t.Errorf("expected userID 'user-correct', got %q", userID)
	}
	if updatedKeyID != "key-2" {
		t.Errorf("expected lastUsedAt update for 'key-2', got %q", updatedKeyID)
	}
}

func TestAPIKeyServiceImpl_ValidateKey_EmptyKeyList(t *testing.T) {
	repo := &mockAPIKeyRepo{
		findAllFn: func(ctx context.Context) ([]*models.APIKey, error) {
			return []*models.APIKey{}, nil
		},
		updateLastUsedFn: func(ctx context.Context, keyID string) error {
			t.Error("updateLastUsed should not be called when no keys exist")
			return nil
		},
	}
	s := NewAPIKeyService(repo)

	userID, err := s.ValidateKey(context.Background(), "muf_some_key")
	if !errors.Is(err, ErrInvalidToken) {
		t.Errorf("expected ErrInvalidToken, got %v", err)
	}
	if userID != "" {
		t.Errorf("expected empty userID, got %q", userID)
	}
}

func TestAPIKeyServiceImpl_ValidateKey_OnlyPrefixNoPayload(t *testing.T) {
	repo := &mockAPIKeyRepo{
		findAllFn: func(ctx context.Context) ([]*models.APIKey, error) {
			return []*models.APIKey{}, nil
		},
		updateLastUsedFn: func(ctx context.Context, keyID string) error {
			return nil
		},
	}
	s := NewAPIKeyService(repo)

	// "muf_" alone with no payload — still has the prefix so passes prefix check
	userID, err := s.ValidateKey(context.Background(), "muf_")
	if !errors.Is(err, ErrInvalidToken) {
		t.Errorf("expected ErrInvalidToken for prefix-only key, got %v", err)
	}
	if userID != "" {
		t.Errorf("expected empty userID, got %q", userID)
	}
}

func TestAPIKeyServiceImpl_Delete_EmptyKeyID(t *testing.T) {
	repo := &mockAPIKeyRepo{
		findByIDFn: func(ctx context.Context, keyID string) (*models.APIKey, error) {
			return nil, ErrNotFound
		},
		deleteFn: func(ctx context.Context, keyID string) error {
			t.Fatal("delete should not be called for empty keyID")
			return nil
		},
	}
	s := NewAPIKeyService(repo)

	err := s.Delete(context.Background(), "", "user-123")
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound for empty keyID, got %v", err)
	}
}

func TestAPIKeyServiceImpl_Delete_EmptyUserID(t *testing.T) {
	repo := &mockAPIKeyRepo{
		findByIDFn: func(ctx context.Context, keyID string) (*models.APIKey, error) {
			return &models.APIKey{
				ID:     "key-123",
				UserID: "actual-owner",
			}, nil
		},
		deleteFn: func(ctx context.Context, keyID string) error {
			t.Fatal("delete should not be called when userID doesn't match")
			return nil
		},
	}
	s := NewAPIKeyService(repo)

	err := s.Delete(context.Background(), "key-123", "")
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound for empty userID, got %v", err)
	}
}

func TestAPIKeyServiceImpl_List_ReturnsNilFromRepo(t *testing.T) {
	repo := &mockAPIKeyRepo{
		findByUserIDFn: func(ctx context.Context, userID string) ([]*models.APIKey, error) {
			return nil, nil
		},
	}
	s := NewAPIKeyService(repo)

	keys, err := s.List(context.Background(), "user-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if keys != nil {
		t.Errorf("expected nil keys when repo returns nil, got %v", keys)
	}
}

func TestGenerateAPIKey(t *testing.T) {
	key1, err := generateAPIKey()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.HasPrefix(key1, APIKeyPrefix) {
		t.Errorf("expected key to have prefix %q, got %q", APIKeyPrefix, key1)
	}

	// Key should be prefix + base64(32 bytes) = "muf_" + 44 chars = 48 chars
	expectedLen := len(APIKeyPrefix) + 44 // base64 of 32 bytes = 44 chars (with padding)
	if len(key1) != expectedLen {
		t.Errorf("expected key length %d, got %d (key: %q)", expectedLen, len(key1), key1)
	}

	// Generate another and ensure uniqueness
	key2, err := generateAPIKey()
	if err != nil {
		t.Fatalf("unexpected error generating second key: %v", err)
	}
	if key1 == key2 {
		t.Error("two generated keys should not be identical")
	}
}

func TestHashAPIKey(t *testing.T) {
	rawKey := "muf_test_key_12345"

	hash, err := hashAPIKey(rawKey)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should be valid bcrypt
	if !strings.HasPrefix(hash, "$2") {
		t.Errorf("expected bcrypt hash format, got %q", hash)
	}

	// Should verify against original key
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(rawKey)); err != nil {
		t.Errorf("hash does not verify against original key: %v", err)
	}

	// Different keys should produce different hashes
	hash2, err := hashAPIKey(rawKey)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hash == hash2 {
		t.Error("bcrypt should produce different hashes for same input (different salts)")
	}
}

func TestHashAPIKey_ExtremelyLongKey(t *testing.T) {
	// bcrypt has a 72-byte input limit; Go's bcrypt rejects keys over 72 bytes
	longKey := "muf_" + strings.Repeat("a", 200)

	_, err := hashAPIKey(longKey)
	if err == nil {
		t.Fatal("expected error for key exceeding bcrypt 72-byte limit")
	}
	if !strings.Contains(err.Error(), "bcrypt") {
		t.Errorf("expected bcrypt-related error, got %q", err.Error())
	}
}

func TestAPIKeyServiceImpl_ValidateKey_ExpiresAtExactlyNow(t *testing.T) {
	hashKey := func(key string) string {
		hash, _ := bcrypt.GenerateFromPassword([]byte(key), bcrypt.MinCost)
		return string(hash)
	}

	rawKey := "muf_expiry_edge_case"
	keyHash := hashKey(rawKey)

	// Set expiration to just barely in the past (1 second ago)
	justExpired := time.Now().Add(-1 * time.Second)

	repo := &mockAPIKeyRepo{
		findAllFn: func(ctx context.Context) ([]*models.APIKey, error) {
			return []*models.APIKey{
				{
					ID:        "key-edge",
					UserID:    "user-123",
					KeyHash:   keyHash,
					ExpiresAt: &justExpired,
				},
			}, nil
		},
		updateLastUsedFn: func(ctx context.Context, keyID string) error {
			t.Error("updateLastUsed should not be called for expired key")
			return nil
		},
	}
	s := NewAPIKeyService(repo)

	userID, err := s.ValidateKey(context.Background(), rawKey)
	if !errors.Is(err, ErrInvalidToken) {
		t.Errorf("expected ErrInvalidToken for just-expired key, got %v", err)
	}
	if userID != "" {
		t.Errorf("expected empty userID, got %q", userID)
	}
}

func TestGenerateAPIKey_RandReadFailure(t *testing.T) {
	origReader := randReadFn
	defer func() { randReadFn = origReader }()

	randReadFn = &failingReader{}

	key, err := generateAPIKey()
	if err == nil {
		t.Fatal("expected error when rand reader fails")
	}
	if !strings.Contains(err.Error(), "crypto/rand read failed") {
		t.Errorf("expected error to contain 'crypto/rand read failed', got %q", err.Error())
	}
	if key != "" {
		t.Errorf("expected empty key, got %q", key)
	}
}

// failingReader is an io.Reader that always returns an error
type failingReader struct{}

func (f *failingReader) Read(p []byte) (int, error) {
	return 0, errors.New("simulated rand failure")
}

func TestNewAPIKeyService_ReturnsNonNil(t *testing.T) {
	repo := &mockAPIKeyRepo{}
	svc := NewAPIKeyService(repo)
	if svc == nil {
		t.Fatal("expected non-nil service")
	}
}
