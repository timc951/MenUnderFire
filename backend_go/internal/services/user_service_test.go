package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"testing"
	"time"

	"menunderfire/internal/models"
)

const testSecret = "test-hmac-secret"

// ===== Mock Repository =====

type mockUserRepository struct {
	findByID                  func(ctx context.Context, userID string) (*models.User, error)
	findByExternalID          func(ctx context.Context, externalID string) (*models.User, error)
	update                    func(ctx context.Context, userID, displayName string) (*models.User, error)
	isSiteAdmin               func(ctx context.Context, userID string) (bool, error)
	findAdminOrganizationIDs  func(ctx context.Context, userID string) ([]string, error)
	findOwnedGroupIDs         func(ctx context.Context, userID string) ([]string, error)
	findMemberGroupIDs        func(ctx context.Context, userID string) ([]string, error)
	recordAgreementAcceptance func(ctx context.Context, userID, version, signature, ipAddress, userAgent string) error
}

func (m *mockUserRepository) FindByID(ctx context.Context, userID string) (*models.User, error) {
	if m.findByID != nil {
		return m.findByID(ctx, userID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserRepository) FindByExternalID(ctx context.Context, externalID string) (*models.User, error) {
	if m.findByExternalID != nil {
		return m.findByExternalID(ctx, externalID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserRepository) Update(ctx context.Context, userID, displayName string) (*models.User, error) {
	if m.update != nil {
		return m.update(ctx, userID, displayName)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserRepository) UpdateExternalID(ctx context.Context, userID, externalID string) error {
	return errors.New("not implemented")
}

func (m *mockUserRepository) UpdateInvitationInfo(ctx context.Context, userID, invitedByID, invitationID string) error {
	return errors.New("not implemented")
}

func (m *mockUserRepository) RecordAgreementAcceptance(ctx context.Context, userID, version, signature, ipAddress, userAgent string) error {
	if m.recordAgreementAcceptance != nil {
		return m.recordAgreementAcceptance(ctx, userID, version, signature, ipAddress, userAgent)
	}
	return errors.New("not implemented")
}

func (m *mockUserRepository) IsSiteAdmin(ctx context.Context, userID string) (bool, error) {
	if m.isSiteAdmin != nil {
		return m.isSiteAdmin(ctx, userID)
	}
	return false, errors.New("not implemented")
}

func (m *mockUserRepository) FindAdminOrganizationIDs(ctx context.Context, userID string) ([]string, error) {
	if m.findAdminOrganizationIDs != nil {
		return m.findAdminOrganizationIDs(ctx, userID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserRepository) FindOwnedGroupIDs(ctx context.Context, userID string) ([]string, error) {
	if m.findOwnedGroupIDs != nil {
		return m.findOwnedGroupIDs(ctx, userID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserRepository) FindMemberGroupIDs(ctx context.Context, userID string) ([]string, error) {
	if m.findMemberGroupIDs != nil {
		return m.findMemberGroupIDs(ctx, userID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserRepository) Count(ctx context.Context) (int64, error) {
	return 0, errors.New("not implemented")
}

func (m *mockUserRepository) Create(ctx context.Context, externalID, email, displayName string) (*models.User, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserRepository) CreateAsSiteAdmin(ctx context.Context, externalID, email, displayName string) (*models.User, error) {
	return nil, errors.New("not implemented")
}

// ===== Test Helpers =====

func testUser(id, email, displayName, externalID string) *models.User {
	return &models.User{
		ID:          id,
		Email:       email,
		DisplayName: displayName,
		ExternalID:  externalID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// =============================================================================
// GetByID Tests
// =============================================================================

func TestUserService_GetByID_Success(t *testing.T) {
	expectedUser := testUser("user-1", "test@example.com", "Test User", "ext-123")
	repo := &mockUserRepository{
		findByID: func(ctx context.Context, userID string) (*models.User, error) {
			if userID == "user-1" {
				return expectedUser, nil
			}
			return nil, errors.New("not found")
		},
	}
	svc := NewUserService(repo, testSecret)

	user, err := svc.GetByID(context.Background(), "user-1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if user.ID != "user-1" {
		t.Errorf("expected user ID 'user-1', got '%s'", user.ID)
	}
}

func TestUserService_GetByID_EmptyUserID(t *testing.T) {
	svc := NewUserService(&mockUserRepository{}, testSecret)

	_, err := svc.GetByID(context.Background(), "")
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestUserService_GetByID_WhitespaceOnlyUserID(t *testing.T) {
	svc := NewUserService(&mockUserRepository{}, testSecret)

	_, err := svc.GetByID(context.Background(), "   ")
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestUserService_GetByID_UserNotFound(t *testing.T) {
	repo := &mockUserRepository{
		findByID: func(ctx context.Context, userID string) (*models.User, error) {
			return nil, errors.New("not found")
		},
	}
	svc := NewUserService(repo, testSecret)

	_, err := svc.GetByID(context.Background(), "nonexistent")
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserService_GetByID_TrimsUserID(t *testing.T) {
	var receivedUserID string
	repo := &mockUserRepository{
		findByID: func(ctx context.Context, userID string) (*models.User, error) {
			receivedUserID = userID
			return testUser("user-1", "test@example.com", "Test User", "ext-123"), nil
		},
	}
	svc := NewUserService(repo, testSecret)

	_, err := svc.GetByID(context.Background(), "  user-1  ")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if receivedUserID != "user-1" {
		t.Errorf("expected trimmed userID 'user-1', got '%s'", receivedUserID)
	}
}

// =============================================================================
// GetByExternalID Tests
// =============================================================================

func TestUserService_GetByExternalID_Success(t *testing.T) {
	expectedUser := testUser("user-1", "test@example.com", "Test User", "ext-123")
	repo := &mockUserRepository{
		findByExternalID: func(ctx context.Context, externalID string) (*models.User, error) {
			if externalID == "ext-123" {
				return expectedUser, nil
			}
			return nil, errors.New("not found")
		},
	}
	svc := NewUserService(repo, testSecret)

	user, err := svc.GetByExternalID(context.Background(), "ext-123")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if user.ExternalID != "ext-123" {
		t.Errorf("expected external ID 'ext-123', got '%s'", user.ExternalID)
	}
}

func TestUserService_GetByExternalID_EmptyExternalID(t *testing.T) {
	svc := NewUserService(&mockUserRepository{}, testSecret)

	_, err := svc.GetByExternalID(context.Background(), "")
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestUserService_GetByExternalID_WhitespaceOnlyExternalID(t *testing.T) {
	svc := NewUserService(&mockUserRepository{}, testSecret)

	_, err := svc.GetByExternalID(context.Background(), "   ")
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestUserService_GetByExternalID_UserNotFound(t *testing.T) {
	repo := &mockUserRepository{
		findByExternalID: func(ctx context.Context, externalID string) (*models.User, error) {
			return nil, errors.New("not found")
		},
	}
	svc := NewUserService(repo, testSecret)

	_, err := svc.GetByExternalID(context.Background(), "nonexistent")
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserService_GetByExternalID_TrimsExternalID(t *testing.T) {
	var receivedExternalID string
	repo := &mockUserRepository{
		findByExternalID: func(ctx context.Context, externalID string) (*models.User, error) {
			receivedExternalID = externalID
			return testUser("user-1", "test@example.com", "Test User", "ext-123"), nil
		},
	}
	svc := NewUserService(repo, testSecret)

	_, err := svc.GetByExternalID(context.Background(), "  ext-123  ")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if receivedExternalID != "ext-123" {
		t.Errorf("expected trimmed externalID 'ext-123', got '%s'", receivedExternalID)
	}
}

// =============================================================================
// Update Tests
// =============================================================================

func TestUserService_Update_Success(t *testing.T) {
	existingUser := testUser("user-1", "test@example.com", "Old Name", "ext-123")
	updatedUser := testUser("user-1", "test@example.com", "New Name", "ext-123")

	repo := &mockUserRepository{
		findByID: func(ctx context.Context, userID string) (*models.User, error) {
			return existingUser, nil
		},
		update: func(ctx context.Context, userID, displayName string) (*models.User, error) {
			return updatedUser, nil
		},
	}
	svc := NewUserService(repo, testSecret)

	req := &models.UpdateUserRequest{DisplayName: "New Name"}
	user, err := svc.Update(context.Background(), "user-1", req)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if user.DisplayName != "New Name" {
		t.Errorf("expected display name 'New Name', got '%s'", user.DisplayName)
	}
}

func TestUserService_Update_EmptyUserID(t *testing.T) {
	svc := NewUserService(&mockUserRepository{}, testSecret)

	req := &models.UpdateUserRequest{DisplayName: "New Name"}
	_, err := svc.Update(context.Background(), "", req)
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestUserService_Update_WhitespaceOnlyUserID(t *testing.T) {
	svc := NewUserService(&mockUserRepository{}, testSecret)

	req := &models.UpdateUserRequest{DisplayName: "New Name"}
	_, err := svc.Update(context.Background(), "   ", req)
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestUserService_Update_NilRequest(t *testing.T) {
	svc := NewUserService(&mockUserRepository{}, testSecret)

	_, err := svc.Update(context.Background(), "user-1", nil)
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestUserService_Update_EmptyDisplayName(t *testing.T) {
	svc := NewUserService(&mockUserRepository{}, testSecret)

	req := &models.UpdateUserRequest{DisplayName: ""}
	_, err := svc.Update(context.Background(), "user-1", req)
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestUserService_Update_WhitespaceOnlyDisplayName(t *testing.T) {
	svc := NewUserService(&mockUserRepository{}, testSecret)

	req := &models.UpdateUserRequest{DisplayName: "   "}
	_, err := svc.Update(context.Background(), "user-1", req)
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestUserService_Update_UserNotFound(t *testing.T) {
	repo := &mockUserRepository{
		findByID: func(ctx context.Context, userID string) (*models.User, error) {
			return nil, errors.New("not found")
		},
	}
	svc := NewUserService(repo, testSecret)

	req := &models.UpdateUserRequest{DisplayName: "New Name"}
	_, err := svc.Update(context.Background(), "nonexistent", req)
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserService_Update_UpdateError(t *testing.T) {
	repoErr := errors.New("database error")
	repo := &mockUserRepository{
		findByID: func(ctx context.Context, userID string) (*models.User, error) {
			return testUser("user-1", "test@example.com", "Old Name", "ext-123"), nil
		},
		update: func(ctx context.Context, userID, displayName string) (*models.User, error) {
			return nil, repoErr
		},
	}
	svc := NewUserService(repo, testSecret)

	req := &models.UpdateUserRequest{DisplayName: "New Name"}
	_, err := svc.Update(context.Background(), "user-1", req)
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repository error, got %v", err)
	}
}

func TestUserService_Update_TrimsUserID(t *testing.T) {
	var receivedUserID string
	repo := &mockUserRepository{
		findByID: func(ctx context.Context, userID string) (*models.User, error) {
			receivedUserID = userID
			return testUser("user-1", "test@example.com", "Old Name", "ext-123"), nil
		},
		update: func(ctx context.Context, userID, displayName string) (*models.User, error) {
			return testUser("user-1", "test@example.com", displayName, "ext-123"), nil
		},
	}
	svc := NewUserService(repo, testSecret)

	req := &models.UpdateUserRequest{DisplayName: "New Name"}
	_, err := svc.Update(context.Background(), "  user-1  ", req)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if receivedUserID != "user-1" {
		t.Errorf("expected trimmed userID 'user-1', got '%s'", receivedUserID)
	}
}

func TestUserService_Update_TrimsDisplayName(t *testing.T) {
	var receivedDisplayName string
	repo := &mockUserRepository{
		findByID: func(ctx context.Context, userID string) (*models.User, error) {
			return testUser("user-1", "test@example.com", "Old Name", "ext-123"), nil
		},
		update: func(ctx context.Context, userID, displayName string) (*models.User, error) {
			receivedDisplayName = displayName
			return testUser("user-1", "test@example.com", displayName, "ext-123"), nil
		},
	}
	svc := NewUserService(repo, testSecret)

	req := &models.UpdateUserRequest{DisplayName: "  New Name  "}
	_, err := svc.Update(context.Background(), "user-1", req)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if receivedDisplayName != "New Name" {
		t.Errorf("expected trimmed displayName 'New Name', got '%s'", receivedDisplayName)
	}
}

// =============================================================================
// GetPermissions Tests
// =============================================================================

func TestUserService_GetPermissions_Success(t *testing.T) {
	repo := &mockUserRepository{
		findByID: func(ctx context.Context, userID string) (*models.User, error) {
			return testUser("user-1", "test@example.com", "Test User", "ext-123"), nil
		},
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			return true, nil
		},
		findAdminOrganizationIDs: func(ctx context.Context, userID string) ([]string, error) {
			return []string{"org-1", "org-2"}, nil
		},
		findOwnedGroupIDs: func(ctx context.Context, userID string) ([]string, error) {
			return []string{"group-1"}, nil
		},
		findMemberGroupIDs: func(ctx context.Context, userID string) ([]string, error) {
			return []string{"group-2", "group-3"}, nil
		},
	}
	svc := NewUserService(repo, testSecret)

	perms, err := svc.GetPermissions(context.Background(), "user-1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !perms.IsSiteAdmin {
		t.Errorf("expected IsSiteAdmin to be true")
	}
	if len(perms.AdminOfOrganizationIDs) != 2 {
		t.Errorf("expected 2 admin orgs, got %d", len(perms.AdminOfOrganizationIDs))
	}
	if len(perms.OwnedGroupIDs) != 1 {
		t.Errorf("expected 1 owned group, got %d", len(perms.OwnedGroupIDs))
	}
	if len(perms.MemberGroupIDs) != 2 {
		t.Errorf("expected 2 member groups, got %d", len(perms.MemberGroupIDs))
	}
}

func TestUserService_GetPermissions_EmptyUserID(t *testing.T) {
	svc := NewUserService(&mockUserRepository{}, testSecret)

	_, err := svc.GetPermissions(context.Background(), "")
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestUserService_GetPermissions_WhitespaceOnlyUserID(t *testing.T) {
	svc := NewUserService(&mockUserRepository{}, testSecret)

	_, err := svc.GetPermissions(context.Background(), "   ")
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestUserService_GetPermissions_UserNotFound(t *testing.T) {
	repo := &mockUserRepository{
		findByID: func(ctx context.Context, userID string) (*models.User, error) {
			return nil, errors.New("not found")
		},
	}
	svc := NewUserService(repo, testSecret)

	_, err := svc.GetPermissions(context.Background(), "nonexistent")
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserService_GetPermissions_IsSiteAdminError(t *testing.T) {
	repoErr := errors.New("database error")
	repo := &mockUserRepository{
		findByID: func(ctx context.Context, userID string) (*models.User, error) {
			return testUser("user-1", "test@example.com", "Test User", "ext-123"), nil
		},
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			return false, repoErr
		},
	}
	svc := NewUserService(repo, testSecret)

	_, err := svc.GetPermissions(context.Background(), "user-1")
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repository error, got %v", err)
	}
}

func TestUserService_GetPermissions_FindAdminOrgsError(t *testing.T) {
	repoErr := errors.New("database error")
	repo := &mockUserRepository{
		findByID: func(ctx context.Context, userID string) (*models.User, error) {
			return testUser("user-1", "test@example.com", "Test User", "ext-123"), nil
		},
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			return false, nil
		},
		findAdminOrganizationIDs: func(ctx context.Context, userID string) ([]string, error) {
			return nil, repoErr
		},
	}
	svc := NewUserService(repo, testSecret)

	_, err := svc.GetPermissions(context.Background(), "user-1")
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repository error, got %v", err)
	}
}

func TestUserService_GetPermissions_FindOwnedGroupsError(t *testing.T) {
	repoErr := errors.New("database error")
	repo := &mockUserRepository{
		findByID: func(ctx context.Context, userID string) (*models.User, error) {
			return testUser("user-1", "test@example.com", "Test User", "ext-123"), nil
		},
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			return false, nil
		},
		findAdminOrganizationIDs: func(ctx context.Context, userID string) ([]string, error) {
			return []string{}, nil
		},
		findOwnedGroupIDs: func(ctx context.Context, userID string) ([]string, error) {
			return nil, repoErr
		},
	}
	svc := NewUserService(repo, testSecret)

	_, err := svc.GetPermissions(context.Background(), "user-1")
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repository error, got %v", err)
	}
}

func TestUserService_GetPermissions_FindMemberGroupsError(t *testing.T) {
	repoErr := errors.New("database error")
	repo := &mockUserRepository{
		findByID: func(ctx context.Context, userID string) (*models.User, error) {
			return testUser("user-1", "test@example.com", "Test User", "ext-123"), nil
		},
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			return false, nil
		},
		findAdminOrganizationIDs: func(ctx context.Context, userID string) ([]string, error) {
			return []string{}, nil
		},
		findOwnedGroupIDs: func(ctx context.Context, userID string) ([]string, error) {
			return []string{}, nil
		},
		findMemberGroupIDs: func(ctx context.Context, userID string) ([]string, error) {
			return nil, repoErr
		},
	}
	svc := NewUserService(repo, testSecret)

	_, err := svc.GetPermissions(context.Background(), "user-1")
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repository error, got %v", err)
	}
}

func TestUserService_GetPermissions_NoPermissions(t *testing.T) {
	repo := &mockUserRepository{
		findByID: func(ctx context.Context, userID string) (*models.User, error) {
			return testUser("user-1", "test@example.com", "Test User", "ext-123"), nil
		},
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			return false, nil
		},
		findAdminOrganizationIDs: func(ctx context.Context, userID string) ([]string, error) {
			return []string{}, nil
		},
		findOwnedGroupIDs: func(ctx context.Context, userID string) ([]string, error) {
			return []string{}, nil
		},
		findMemberGroupIDs: func(ctx context.Context, userID string) ([]string, error) {
			return []string{}, nil
		},
	}
	svc := NewUserService(repo, testSecret)

	perms, err := svc.GetPermissions(context.Background(), "user-1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if perms.IsSiteAdmin {
		t.Errorf("expected IsSiteAdmin to be false")
	}
	if len(perms.AdminOfOrganizationIDs) != 0 {
		t.Errorf("expected 0 admin orgs, got %d", len(perms.AdminOfOrganizationIDs))
	}
	if len(perms.OwnedGroupIDs) != 0 {
		t.Errorf("expected 0 owned groups, got %d", len(perms.OwnedGroupIDs))
	}
	if len(perms.MemberGroupIDs) != 0 {
		t.Errorf("expected 0 member groups, got %d", len(perms.MemberGroupIDs))
	}
}

func TestUserService_GetPermissions_NilArraysConvertedToEmpty(t *testing.T) {
	repo := &mockUserRepository{
		findByID: func(ctx context.Context, userID string) (*models.User, error) {
			return testUser("user-1", "test@example.com", "Test User", "ext-123"), nil
		},
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			return false, nil
		},
		findAdminOrganizationIDs: func(ctx context.Context, userID string) ([]string, error) {
			return nil, nil // Return nil instead of empty slice
		},
		findOwnedGroupIDs: func(ctx context.Context, userID string) ([]string, error) {
			return nil, nil // Return nil instead of empty slice
		},
		findMemberGroupIDs: func(ctx context.Context, userID string) ([]string, error) {
			return nil, nil // Return nil instead of empty slice
		},
	}
	svc := NewUserService(repo, testSecret)

	perms, err := svc.GetPermissions(context.Background(), "user-1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	// Verify that nil slices were converted to empty slices (not nil)
	if perms.AdminOfOrganizationIDs == nil {
		t.Errorf("expected AdminOfOrganizationIDs to be empty slice, not nil")
	}
	if perms.OwnedGroupIDs == nil {
		t.Errorf("expected OwnedGroupIDs to be empty slice, not nil")
	}
	if perms.MemberGroupIDs == nil {
		t.Errorf("expected MemberGroupIDs to be empty slice, not nil")
	}
}

func TestUserService_GetPermissions_TrimsUserID(t *testing.T) {
	var receivedUserID string
	repo := &mockUserRepository{
		findByID: func(ctx context.Context, userID string) (*models.User, error) {
			receivedUserID = userID
			return testUser("user-1", "test@example.com", "Test User", "ext-123"), nil
		},
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			return false, nil
		},
		findAdminOrganizationIDs: func(ctx context.Context, userID string) ([]string, error) {
			return []string{}, nil
		},
		findOwnedGroupIDs: func(ctx context.Context, userID string) ([]string, error) {
			return []string{}, nil
		},
		findMemberGroupIDs: func(ctx context.Context, userID string) ([]string, error) {
			return []string{}, nil
		},
	}
	svc := NewUserService(repo, testSecret)

	_, err := svc.GetPermissions(context.Background(), "  user-1  ")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if receivedUserID != "user-1" {
		t.Errorf("expected trimmed userID 'user-1', got '%s'", receivedUserID)
	}
}

// =============================================================================
// AcceptAgreement Tests
// =============================================================================

func TestUserService_AcceptAgreement_Success(t *testing.T) {
	var capturedSignature string
	repo := &mockUserRepository{
		findByID: func(ctx context.Context, userID string) (*models.User, error) {
			return testUser("user-1", "test@example.com", "Test User", "ext-123"), nil
		},
		recordAgreementAcceptance: func(ctx context.Context, userID, version, signature, ipAddress, userAgent string) error {
			capturedSignature = signature
			if userID != "user-1" {
				t.Errorf("expected userID 'user-1', got '%s'", userID)
			}
			if version != "1.0" {
				t.Errorf("expected version '1.0', got '%s'", version)
			}
			if ipAddress != "127.0.0.1" {
				t.Errorf("expected IP '127.0.0.1', got '%s'", ipAddress)
			}
			return nil
		},
	}
	svc := NewUserService(repo, testSecret)

	resp, err := svc.AcceptAgreement(context.Background(), "user-1", "1.0", "127.0.0.1", "TestAgent")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Version != "1.0" {
		t.Errorf("expected version '1.0', got '%s'", resp.Version)
	}
	if resp.Signature == "" {
		t.Error("expected non-empty signature")
	}
	if resp.Signature != capturedSignature {
		t.Error("response signature should match what was stored in DB")
	}
	if resp.AcceptedAt == "" {
		t.Error("expected non-empty acceptedAt")
	}
}

func TestUserService_AcceptAgreement_EmptyUserID(t *testing.T) {
	svc := NewUserService(&mockUserRepository{}, testSecret)

	_, err := svc.AcceptAgreement(context.Background(), "", "1.0", "127.0.0.1", "TestAgent")
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestUserService_AcceptAgreement_EmptyVersion(t *testing.T) {
	svc := NewUserService(&mockUserRepository{}, testSecret)

	_, err := svc.AcceptAgreement(context.Background(), "user-1", "", "127.0.0.1", "TestAgent")
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestUserService_AcceptAgreement_UserNotFound(t *testing.T) {
	repo := &mockUserRepository{
		findByID: func(ctx context.Context, userID string) (*models.User, error) {
			return nil, errors.New("not found")
		},
	}
	svc := NewUserService(repo, testSecret)

	_, err := svc.AcceptAgreement(context.Background(), "nonexistent", "1.0", "127.0.0.1", "TestAgent")
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserService_AcceptAgreement_AlreadyAccepted(t *testing.T) {
	now := time.Now()
	version := "1.0"
	repo := &mockUserRepository{
		findByID: func(ctx context.Context, userID string) (*models.User, error) {
			user := testUser("user-1", "test@example.com", "Test User", "ext-123")
			user.AgreementAcceptedAt = &now
			user.AgreementVersion = &version
			return user, nil
		},
	}
	svc := NewUserService(repo, testSecret)

	_, err := svc.AcceptAgreement(context.Background(), "user-1", "1.0", "127.0.0.1", "TestAgent")
	if !errors.Is(err, ErrConflict) {
		t.Errorf("expected ErrConflict, got %v", err)
	}
}

func TestUserService_AcceptAgreement_SignatureIsValidHMAC(t *testing.T) {
	repo := &mockUserRepository{
		findByID: func(ctx context.Context, userID string) (*models.User, error) {
			return testUser("user-1", "test@example.com", "Test User", "ext-123"), nil
		},
		recordAgreementAcceptance: func(ctx context.Context, userID, version, signature, ipAddress, userAgent string) error {
			return nil
		},
	}
	svc := NewUserService(repo, testSecret)

	resp, err := svc.AcceptAgreement(context.Background(), "user-1", "1.0", "127.0.0.1", "TestAgent")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify the signature is a valid hex-encoded HMAC-SHA256
	sigBytes, err := hex.DecodeString(resp.Signature)
	if err != nil {
		t.Fatalf("signature is not valid hex: %v", err)
	}
	if len(sigBytes) != sha256.Size {
		t.Errorf("expected %d byte signature, got %d", sha256.Size, len(sigBytes))
	}

	// Recompute and verify
	message := resp.AcceptedAt // The timestamp used in the message
	expectedMessage := "user-1|1.0|" + message
	mac := hmac.New(sha256.New, []byte(testSecret))
	mac.Write([]byte(expectedMessage))
	expectedSig := hex.EncodeToString(mac.Sum(nil))
	if resp.Signature != expectedSig {
		t.Errorf("signature mismatch: recomputed doesn't match")
	}
}

func TestUserService_AcceptAgreement_RepoError(t *testing.T) {
	repoErr := errors.New("database error")
	repo := &mockUserRepository{
		findByID: func(ctx context.Context, userID string) (*models.User, error) {
			return testUser("user-1", "test@example.com", "Test User", "ext-123"), nil
		},
		recordAgreementAcceptance: func(ctx context.Context, userID, version, signature, ipAddress, userAgent string) error {
			return repoErr
		},
	}
	svc := NewUserService(repo, testSecret)

	_, err := svc.AcceptAgreement(context.Background(), "user-1", "1.0", "127.0.0.1", "TestAgent")
	if err == nil {
		t.Error("expected error, got nil")
	}
}
