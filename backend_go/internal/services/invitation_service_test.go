package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"menunderfire/internal/models"
	"menunderfire/internal/repositories"
)

// Mock repository types for invitation service tests
type mockInvitationRepo struct {
	findByIDFn     func(ctx context.Context, invitationID string) (*models.Invitation, error)
	findByTokenFn  func(ctx context.Context, token string) (*models.Invitation, error)
	findByEmailFn  func(ctx context.Context, email string, invType models.InvitationType, targetID string) (*models.Invitation, error)
	createFn       func(ctx context.Context, email string, invType models.InvitationType, orgID, groupID *string, inviterID string, token string, expiresAt time.Time) (*models.Invitation, error)
	updateStatusFn func(ctx context.Context, invitationID string, status models.InvitationStatus, acceptedByID *string) error
	deleteFn       func(ctx context.Context, invitationID string) error
}

func (m *mockInvitationRepo) FindByID(ctx context.Context, invitationID string) (*models.Invitation, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, invitationID)
	}
	return nil, errors.New("findByIDFn not set")
}

func (m *mockInvitationRepo) FindByToken(ctx context.Context, token string) (*models.Invitation, error) {
	if m.findByTokenFn != nil {
		return m.findByTokenFn(ctx, token)
	}
	return nil, errors.New("findByTokenFn not set")
}

func (m *mockInvitationRepo) FindByEmail(ctx context.Context, email string, invType models.InvitationType, targetID string) (*models.Invitation, error) {
	if m.findByEmailFn != nil {
		return m.findByEmailFn(ctx, email, invType, targetID)
	}
	return nil, errors.New("findByEmailFn not set")
}

func (m *mockInvitationRepo) Create(ctx context.Context, email string, invType models.InvitationType, orgID, groupID *string, inviterID string, token string, expiresAt time.Time) (*models.Invitation, error) {
	if m.createFn != nil {
		return m.createFn(ctx, email, invType, orgID, groupID, inviterID, token, expiresAt)
	}
	return nil, errors.New("createFn not set")
}

func (m *mockInvitationRepo) UpdateStatus(ctx context.Context, invitationID string, status models.InvitationStatus, acceptedByID *string) error {
	if m.updateStatusFn != nil {
		return m.updateStatusFn(ctx, invitationID, status, acceptedByID)
	}
	return errors.New("updateStatusFn not set")
}

func (m *mockInvitationRepo) Delete(ctx context.Context, invitationID string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, invitationID)
	}
	return errors.New("deleteFn not set")
}

func (m *mockInvitationRepo) ListAll(ctx context.Context) ([]*models.Invitation, error) {
	return nil, errors.New("listAllFn not set")
}

func (m *mockInvitationRepo) ListByOrganizationIDs(ctx context.Context, orgIDs []string) ([]*models.Invitation, error) {
	return nil, errors.New("listByOrganizationIDsFn not set")
}

func (m *mockInvitationRepo) ListByGroupIDs(ctx context.Context, groupIDs []string) ([]*models.Invitation, error) {
	return nil, errors.New("listByGroupIDsFn not set")
}

type mockUserRepoForInv struct {
	findByIDFn    func(ctx context.Context, userID string) (*models.User, error)
	findByEmailFn func(ctx context.Context, email string) (*models.User, error)
	createFn      func(ctx context.Context, email, displayName, externalID string) (*models.User, error)
	isSiteAdminFn func(ctx context.Context, userID string) (bool, error)
}

func (m *mockUserRepoForInv) FindByID(ctx context.Context, userID string) (*models.User, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, userID)
	}
	return nil, errors.New("findByIDFn not set")
}

func (m *mockUserRepoForInv) FindByExternalID(ctx context.Context, externalID string) (*models.User, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserRepoForInv) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	if m.findByEmailFn != nil {
		return m.findByEmailFn(ctx, email)
	}
	return nil, errors.New("findByEmailFn not set")
}

func (m *mockUserRepoForInv) Count(ctx context.Context) (int64, error) {
	return 0, errors.New("not implemented")
}

func (m *mockUserRepoForInv) Create(ctx context.Context, email, displayName, externalID string) (*models.User, error) {
	if m.createFn != nil {
		return m.createFn(ctx, email, displayName, externalID)
	}
	return nil, errors.New("createFn not set")
}

func (m *mockUserRepoForInv) CreateAsSiteAdmin(ctx context.Context, email, displayName, externalID string) (*models.User, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserRepoForInv) Update(ctx context.Context, userID, displayName string) (*models.User, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserRepoForInv) UpdateExternalID(ctx context.Context, userID, externalID string) error {
	return errors.New("not implemented")
}

func (m *mockUserRepoForInv) UpdateInvitationInfo(ctx context.Context, userID, invitedByID, invitationID string) error {
	return nil
}

func (m *mockUserRepoForInv) IsSiteAdmin(ctx context.Context, userID string) (bool, error) {
	if m.isSiteAdminFn != nil {
		return m.isSiteAdminFn(ctx, userID)
	}
	return false, errors.New("isSiteAdminFn not set")
}

func (m *mockUserRepoForInv) FindAdminOrganizationIDs(ctx context.Context, userID string) ([]string, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserRepoForInv) FindOwnedGroupIDs(ctx context.Context, userID string) ([]string, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserRepoForInv) FindMemberGroupIDs(ctx context.Context, userID string) ([]string, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserRepoForInv) RecordAgreementAcceptance(ctx context.Context, userID, version, signature, ipAddress, userAgent string) error {
	return errors.New("not implemented")
}

type mockOrgRepoForInv struct {
	findByIDFn func(ctx context.Context, orgID string) (*models.Organization, error)
	addAdminFn func(ctx context.Context, orgID, userID string) error
	isAdminFn  func(ctx context.Context, orgID, userID string) (bool, error)
}

func (m *mockOrgRepoForInv) FindByID(ctx context.Context, orgID string) (*models.Organization, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, orgID)
	}
	return nil, errors.New("findByIDFn not set")
}

func (m *mockOrgRepoForInv) FindByUserID(ctx context.Context, userID string) ([]*models.Organization, error) {
	return nil, errors.New("not implemented")
}

func (m *mockOrgRepoForInv) FindAll(ctx context.Context) ([]*models.Organization, error) {
	return nil, errors.New("not implemented")
}

func (m *mockOrgRepoForInv) Create(ctx context.Context, name string, description *string, createdByID string) (*models.Organization, error) {
	return nil, errors.New("not implemented")
}

func (m *mockOrgRepoForInv) Update(ctx context.Context, orgID string, name string, description *string) (*models.Organization, error) {
	return nil, errors.New("not implemented")
}

func (m *mockOrgRepoForInv) FindAdmins(ctx context.Context, orgID string) ([]*models.OrganizationAdmin, error) {
	return nil, errors.New("not implemented")
}

func (m *mockOrgRepoForInv) FindAdmin(ctx context.Context, orgID, userID string) (*models.OrganizationAdmin, error) {
	return nil, errors.New("not implemented")
}

func (m *mockOrgRepoForInv) AddAdmin(ctx context.Context, orgID, userID string) error {
	if m.addAdminFn != nil {
		return m.addAdminFn(ctx, orgID, userID)
	}
	return errors.New("addAdminFn not set")
}

func (m *mockOrgRepoForInv) IsMember(ctx context.Context, orgID, userID string) (bool, error) {
	return false, errors.New("not implemented")
}

func (m *mockOrgRepoForInv) IsAdmin(ctx context.Context, orgID, userID string) (bool, error) {
	if m.isAdminFn != nil {
		return m.isAdminFn(ctx, orgID, userID)
	}
	return false, errors.New("isAdminFn not set")
}

func (m *mockOrgRepoForInv) Count(ctx context.Context) (int64, error) {
	return 0, errors.New("not implemented")
}

type mockGroupRepoForInv struct {
	findByIDFn  func(ctx context.Context, groupID string) (*models.Group, error)
	addMemberFn func(ctx context.Context, groupID, userID, role string) (*models.GroupMember, error)
}

func (m *mockGroupRepoForInv) FindByID(ctx context.Context, groupID string) (*models.Group, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, groupID)
	}
	return nil, errors.New("findByIDFn not set")
}

func (m *mockGroupRepoForInv) FindByInviteCode(ctx context.Context, inviteCode string) (*models.Group, error) {
	return nil, errors.New("not implemented")
}

func (m *mockGroupRepoForInv) FindByUserID(ctx context.Context, userID string) ([]*models.Group, error) {
	return nil, errors.New("not implemented")
}

func (m *mockGroupRepoForInv) FindByOrganizationID(ctx context.Context, orgID string) ([]*models.Group, error) {
	return nil, errors.New("not implemented")
}

func (m *mockGroupRepoForInv) Create(ctx context.Context, name string, description *string, orgID string, inviteCode string, createdBy string) (*models.Group, error) {
	return nil, errors.New("not implemented")
}

func (m *mockGroupRepoForInv) GenerateInviteCode() string {
	return ""
}

func (m *mockGroupRepoForInv) FindMember(ctx context.Context, groupID, userID string) (*models.GroupMember, error) {
	return nil, errors.New("not implemented")
}

func (m *mockGroupRepoForInv) FindMembers(ctx context.Context, groupID string) ([]*models.GroupMember, error) {
	return nil, errors.New("not implemented")
}

func (m *mockGroupRepoForInv) CountMembers(ctx context.Context, groupID string) (int, error) {
	return 0, errors.New("not implemented")
}

func (m *mockGroupRepoForInv) AddMember(ctx context.Context, groupID, userID, role string) (*models.GroupMember, error) {
	if m.addMemberFn != nil {
		return m.addMemberFn(ctx, groupID, userID, role)
	}
	return nil, errors.New("addMemberFn not set")
}

func (m *mockGroupRepoForInv) RemoveMember(ctx context.Context, groupID, userID string) error {
	return errors.New("not implemented")
}

func (m *mockGroupRepoForInv) Count(ctx context.Context) (int64, error) {
	return 0, errors.New("not implemented")
}

func (m *mockGroupRepoForInv) CountByOrganizationIDs(ctx context.Context, orgIDs []string) (int64, error) {
	return 0, errors.New("not implemented")
}

func (m *mockGroupRepoForInv) UpdateMemberRole(ctx context.Context, groupID, userID, role string) error {
	return nil
}

func (m *mockGroupRepoForInv) UpdateSettings(ctx context.Context, groupID string, requirePostApproval, allowAnonymousPosts bool) error {
	return nil
}

type mockGroupServiceForInv struct {
	getUserRoleFn func(ctx context.Context, groupID, userID string) (string, error)
}

func (m *mockGroupServiceForInv) Create(ctx context.Context, userID string, req *models.CreateGroupRequest) (*models.Group, error) {
	return nil, errors.New("not implemented")
}

func (m *mockGroupServiceForInv) List(ctx context.Context, userID string) ([]*models.Group, error) {
	return nil, errors.New("not implemented")
}

func (m *mockGroupServiceForInv) GetByID(ctx context.Context, groupID string) (*models.Group, error) {
	return nil, errors.New("not implemented")
}

func (m *mockGroupServiceForInv) GetDetailByID(ctx context.Context, groupID string, userID string) (*models.GroupDetailResponse, error) {
	return nil, errors.New("not implemented")
}

func (m *mockGroupServiceForInv) Join(ctx context.Context, groupID string, userID string, inviteCode string) error {
	return errors.New("not implemented")
}

func (m *mockGroupServiceForInv) JoinByInviteCode(ctx context.Context, userID string, inviteCode string) error {
	return errors.New("not implemented")
}

func (m *mockGroupServiceForInv) RemoveMember(ctx context.Context, groupID string, targetUserID string, requestingUserID string) error {
	return errors.New("not implemented")
}

func (m *mockGroupServiceForInv) GetUserRole(ctx context.Context, groupID string, userID string) (string, error) {
	if m.getUserRoleFn != nil {
		return m.getUserRoleFn(ctx, groupID, userID)
	}
	return "", errors.New("getUserRoleFn not set")
}

func (m *mockGroupServiceForInv) GetMembers(ctx context.Context, groupID string) ([]*models.GroupMember, error) {
	return nil, errors.New("not implemented")
}

func (m *mockGroupServiceForInv) GetMemberCount(ctx context.Context, groupID string) (int, error) {
	return 0, errors.New("not implemented")
}

func (m *mockGroupServiceForInv) UpdateMemberRole(ctx context.Context, groupID string, targetUserID string, newRole string, requestingUserID string) error {
	return nil
}

func (m *mockGroupServiceForInv) UpdateSettings(ctx context.Context, groupID string, userID string, req *models.UpdateGroupSettingsRequest) error {
	return nil
}

// Helper functions

func noopUserRepoForInv() repositories.UserRepository {
	return &mockUserRepoForInv{}
}

func noopOrgRepoForInv() repositories.OrganizationRepository {
	return &mockOrgRepoForInv{}
}

func noopGroupRepoForInv() repositories.GroupRepository {
	return &mockGroupRepoForInv{}
}

func noopGroupServiceForInv() GroupService {
	return &mockGroupServiceForInv{}
}

func TestCreateOrgAdmin(t *testing.T) {
	tests := []struct {
		name                   string
		inviterID              string
		req                    *models.CreateOrgAdminInvitationRequest
		isSiteAdminResult      bool
		isSiteAdminErr         error
		findOrgResult          *models.Organization
		findOrgErr             error
		findInvitationResult   *models.Invitation
		findInvitationErr      error
		createInvitationResult *models.Invitation
		createInvitationErr    error
		wantErr                bool
		wantErrType            error
	}{
		{
			name:        "nil request returns validation error",
			inviterID:   "inviter-1",
			req:         nil,
			wantErr:     true,
			wantErrType: ErrValidation,
		},
		{
			name:        "empty email returns validation error",
			inviterID:   "inviter-1",
			req:         &models.CreateOrgAdminInvitationRequest{Email: "", OrganizationID: "org-1"},
			wantErr:     true,
			wantErrType: ErrValidation,
		},
		{
			name:        "whitespace email returns validation error",
			inviterID:   "inviter-1",
			req:         &models.CreateOrgAdminInvitationRequest{Email: "   ", OrganizationID: "org-1"},
			wantErr:     true,
			wantErrType: ErrValidation,
		},
		{
			name:        "empty organizationID returns validation error",
			inviterID:   "inviter-1",
			req:         &models.CreateOrgAdminInvitationRequest{Email: "test@example.com", OrganizationID: ""},
			wantErr:     true,
			wantErrType: ErrValidation,
		},
		{
			name:        "whitespace organizationID returns validation error",
			inviterID:   "inviter-1",
			req:         &models.CreateOrgAdminInvitationRequest{Email: "test@example.com", OrganizationID: "   "},
			wantErr:     true,
			wantErrType: ErrValidation,
		},
		{
			name:           "isSiteAdmin error returns wrapped error",
			inviterID:      "inviter-1",
			req:            &models.CreateOrgAdminInvitationRequest{Email: "test@example.com", OrganizationID: "org-1"},
			isSiteAdminErr: errors.New("db error"),
			wantErr:        true,
		},
		{
			name:              "non-site admin returns forbidden",
			inviterID:         "inviter-1",
			req:               &models.CreateOrgAdminInvitationRequest{Email: "test@example.com", OrganizationID: "org-1"},
			isSiteAdminResult: false,
			wantErr:           true,
			wantErrType:       ErrForbidden,
		},
		{
			name:              "organization not found returns error",
			inviterID:         "inviter-1",
			req:               &models.CreateOrgAdminInvitationRequest{Email: "test@example.com", OrganizationID: "org-1"},
			isSiteAdminResult: true,
			findOrgErr:        errors.New("not found"),
			wantErr:           true,
			wantErrType:       ErrOrganizationNotFound,
		},
		{
			name:              "existing pending invitation returns conflict",
			inviterID:         "inviter-1",
			req:               &models.CreateOrgAdminInvitationRequest{Email: "test@example.com", OrganizationID: "org-1"},
			isSiteAdminResult: true,
			findOrgResult:     &models.Organization{ID: "org-1", Name: "Test Org"},
			findInvitationResult: &models.Invitation{
				ID:     "inv-1",
				Email:  "test@example.com",
				Status: models.InvitationStatusPending,
			},
			wantErr:     true,
			wantErrType: ErrConflict,
		},
		{
			name:              "existing accepted invitation allows new invitation",
			inviterID:         "inviter-1",
			req:               &models.CreateOrgAdminInvitationRequest{Email: "test@example.com", OrganizationID: "org-1"},
			isSiteAdminResult: true,
			findOrgResult:     &models.Organization{ID: "org-1", Name: "Test Org"},
			findInvitationResult: &models.Invitation{
				ID:     "inv-1",
				Email:  "test@example.com",
				Status: models.InvitationStatusAccepted,
			},
			createInvitationResult: &models.Invitation{
				ID:    "inv-2",
				Email: "test@example.com",
				Type:  models.InvitationTypeOrgAdmin,
			},
			wantErr: false,
		},
		{
			name:              "no existing invitation allows creation",
			inviterID:         "inviter-1",
			req:               &models.CreateOrgAdminInvitationRequest{Email: "test@example.com", OrganizationID: "org-1"},
			isSiteAdminResult: true,
			findOrgResult:     &models.Organization{ID: "org-1", Name: "Test Org"},
			findInvitationErr: errors.New("not found"),
			createInvitationResult: &models.Invitation{
				ID:    "inv-1",
				Email: "test@example.com",
				Type:  models.InvitationTypeOrgAdmin,
			},
			wantErr: false,
		},
		{
			name:                "createInvitation error returns wrapped error",
			inviterID:           "inviter-1",
			req:                 &models.CreateOrgAdminInvitationRequest{Email: "test@example.com", OrganizationID: "org-1"},
			isSiteAdminResult:   true,
			findOrgResult:       &models.Organization{ID: "org-1", Name: "Test Org"},
			findInvitationErr:   errors.New("not found"),
			createInvitationErr: errors.New("db error"),
			wantErr:             true,
		},
		{
			name:              "successful creation returns invitation",
			inviterID:         "inviter-1",
			req:               &models.CreateOrgAdminInvitationRequest{Email: "test@example.com", OrganizationID: "org-1"},
			isSiteAdminResult: true,
			findOrgResult:     &models.Organization{ID: "org-1", Name: "Test Org"},
			findInvitationErr: errors.New("not found"),
			createInvitationResult: &models.Invitation{
				ID:             "inv-1",
				Email:          "test@example.com",
				Type:           models.InvitationTypeOrgAdmin,
				OrganizationID: strPtr("org-1"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			invitationRepo := &mockInvitationRepo{
				findByEmailFn: func(ctx context.Context, email string, invType models.InvitationType, targetID string) (*models.Invitation, error) {
					return tt.findInvitationResult, tt.findInvitationErr
				},
				createFn: func(ctx context.Context, email string, invType models.InvitationType, orgID, groupID *string, inviterID string, token string, expiresAt time.Time) (*models.Invitation, error) {
					return tt.createInvitationResult, tt.createInvitationErr
				},
			}
			userRepo := &mockUserRepoForInv{
				isSiteAdminFn: func(ctx context.Context, userID string) (bool, error) {
					return tt.isSiteAdminResult, tt.isSiteAdminErr
				},
			}
			orgRepo := &mockOrgRepoForInv{
				findByIDFn: func(ctx context.Context, orgID string) (*models.Organization, error) {
					return tt.findOrgResult, tt.findOrgErr
				},
			}

			svc := NewInvitationService(
				invitationRepo,
				userRepo,
				orgRepo,
				noopGroupRepoForInv(),
				noopGroupServiceForInv(),
			)

			result, err := svc.CreateOrgAdmin(context.Background(), tt.inviterID, tt.req)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				if tt.wantErrType != nil && !errors.Is(err, tt.wantErrType) {
					t.Errorf("expected error type %v, got %v", tt.wantErrType, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result == nil {
					t.Error("expected result, got nil")
				}
			}
		})
	}
}

func TestCreateOrgAdmin_ContextPropagation(t *testing.T) {
	type ctxKey string
	key := ctxKey("test")
	ctx := context.WithValue(context.Background(), key, "value")

	var capturedCtx context.Context

	invitationRepo := &mockInvitationRepo{
		findByEmailFn: func(ctx context.Context, email string, invType models.InvitationType, targetID string) (*models.Invitation, error) {
			return nil, errors.New("not found")
		},
		createFn: func(ctx context.Context, email string, invType models.InvitationType, orgID, groupID *string, inviterID string, token string, expiresAt time.Time) (*models.Invitation, error) {
			capturedCtx = ctx
			return &models.Invitation{ID: "inv-1"}, nil
		},
	}
	userRepo := &mockUserRepoForInv{
		isSiteAdminFn: func(ctx context.Context, userID string) (bool, error) {
			return true, nil
		},
	}
	orgRepo := &mockOrgRepoForInv{
		findByIDFn: func(ctx context.Context, orgID string) (*models.Organization, error) {
			return &models.Organization{ID: "org-1"}, nil
		},
	}

	svc := NewInvitationService(
		invitationRepo,
		userRepo,
		orgRepo,
		noopGroupRepoForInv(),
		noopGroupServiceForInv(),
	)

	req := &models.CreateOrgAdminInvitationRequest{
		Email:          "test@example.com",
		OrganizationID: "org-1",
	}

	_, _ = svc.CreateOrgAdmin(ctx, "inviter-1", req)

	if capturedCtx.Value(key) != "value" {
		t.Error("context was not properly propagated")
	}
}

func TestCreateGroupOwner(t *testing.T) {
	tests := []struct {
		name                   string
		inviterID              string
		req                    *models.CreateGroupOwnerInvitationRequest
		findGroupResult        *models.Group
		findGroupErr           error
		isSiteAdminResult      bool
		isSiteAdminErr         error
		isOrgAdminResult       bool
		isOrgAdminErr          error
		findInvitationResult   *models.Invitation
		findInvitationErr      error
		createInvitationResult *models.Invitation
		createInvitationErr    error
		wantErr                bool
		wantErrType            error
	}{
		{
			name:        "nil request returns validation error",
			inviterID:   "inviter-1",
			req:         nil,
			wantErr:     true,
			wantErrType: ErrValidation,
		},
		{
			name:        "empty email returns validation error",
			inviterID:   "inviter-1",
			req:         &models.CreateGroupOwnerInvitationRequest{Email: "", GroupID: "group-1"},
			wantErr:     true,
			wantErrType: ErrValidation,
		},
		{
			name:        "whitespace email returns validation error",
			inviterID:   "inviter-1",
			req:         &models.CreateGroupOwnerInvitationRequest{Email: "   ", GroupID: "group-1"},
			wantErr:     true,
			wantErrType: ErrValidation,
		},
		{
			name:        "empty groupID returns validation error",
			inviterID:   "inviter-1",
			req:         &models.CreateGroupOwnerInvitationRequest{Email: "test@example.com", GroupID: ""},
			wantErr:     true,
			wantErrType: ErrValidation,
		},
		{
			name:        "whitespace groupID returns validation error",
			inviterID:   "inviter-1",
			req:         &models.CreateGroupOwnerInvitationRequest{Email: "test@example.com", GroupID: "   "},
			wantErr:     true,
			wantErrType: ErrValidation,
		},
		{
			name:         "group not found returns error",
			inviterID:    "inviter-1",
			req:          &models.CreateGroupOwnerInvitationRequest{Email: "test@example.com", GroupID: "group-1"},
			findGroupErr: errors.New("not found"),
			wantErr:      true,
			wantErrType:  ErrGroupNotFound,
		},
		{
			name:            "isSiteAdmin error returns wrapped error",
			inviterID:       "inviter-1",
			req:             &models.CreateGroupOwnerInvitationRequest{Email: "test@example.com", GroupID: "group-1"},
			findGroupResult: &models.Group{ID: "group-1", OrganizationID: "org-1"},
			isSiteAdminErr:  errors.New("db error"),
			wantErr:         true,
		},
		{
			name:              "non-site admin, isOrgAdmin error returns wrapped error",
			inviterID:         "inviter-1",
			req:               &models.CreateGroupOwnerInvitationRequest{Email: "test@example.com", GroupID: "group-1"},
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			isSiteAdminResult: false,
			isOrgAdminErr:     errors.New("db error"),
			wantErr:           true,
		},
		{
			name:              "non-site admin and non-org admin returns forbidden",
			inviterID:         "inviter-1",
			req:               &models.CreateGroupOwnerInvitationRequest{Email: "test@example.com", GroupID: "group-1"},
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			isSiteAdminResult: false,
			isOrgAdminResult:  false,
			wantErr:           true,
			wantErrType:       ErrForbidden,
		},
		{
			name:              "site admin can create invitation",
			inviterID:         "inviter-1",
			req:               &models.CreateGroupOwnerInvitationRequest{Email: "test@example.com", GroupID: "group-1"},
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			isSiteAdminResult: true,
			findInvitationErr: errors.New("not found"),
			createInvitationResult: &models.Invitation{
				ID:      "inv-1",
				Email:   "test@example.com",
				Type:    models.InvitationTypeGroupOwner,
				GroupID: strPtr("group-1"),
			},
			wantErr: false,
		},
		{
			name:              "org admin can create invitation",
			inviterID:         "inviter-1",
			req:               &models.CreateGroupOwnerInvitationRequest{Email: "test@example.com", GroupID: "group-1"},
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			isSiteAdminResult: false,
			isOrgAdminResult:  true,
			findInvitationErr: errors.New("not found"),
			createInvitationResult: &models.Invitation{
				ID:      "inv-1",
				Email:   "test@example.com",
				Type:    models.InvitationTypeGroupOwner,
				GroupID: strPtr("group-1"),
			},
			wantErr: false,
		},
		{
			name:              "existing pending invitation returns conflict",
			inviterID:         "inviter-1",
			req:               &models.CreateGroupOwnerInvitationRequest{Email: "test@example.com", GroupID: "group-1"},
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			isSiteAdminResult: true,
			findInvitationResult: &models.Invitation{
				ID:     "inv-1",
				Email:  "test@example.com",
				Status: models.InvitationStatusPending,
			},
			wantErr:     true,
			wantErrType: ErrConflict,
		},
		{
			name:                "createInvitation error returns wrapped error",
			inviterID:           "inviter-1",
			req:                 &models.CreateGroupOwnerInvitationRequest{Email: "test@example.com", GroupID: "group-1"},
			findGroupResult:     &models.Group{ID: "group-1", OrganizationID: "org-1"},
			isSiteAdminResult:   true,
			findInvitationErr:   errors.New("not found"),
			createInvitationErr: errors.New("db error"),
			wantErr:             true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			invitationRepo := &mockInvitationRepo{
				findByEmailFn: func(ctx context.Context, email string, invType models.InvitationType, targetID string) (*models.Invitation, error) {
					return tt.findInvitationResult, tt.findInvitationErr
				},
				createFn: func(ctx context.Context, email string, invType models.InvitationType, orgID, groupID *string, inviterID string, token string, expiresAt time.Time) (*models.Invitation, error) {
					return tt.createInvitationResult, tt.createInvitationErr
				},
			}
			userRepo := &mockUserRepoForInv{
				isSiteAdminFn: func(ctx context.Context, userID string) (bool, error) {
					return tt.isSiteAdminResult, tt.isSiteAdminErr
				},
			}
			orgRepo := &mockOrgRepoForInv{
				isAdminFn: func(ctx context.Context, orgID, userID string) (bool, error) {
					return tt.isOrgAdminResult, tt.isOrgAdminErr
				},
			}
			groupRepo := &mockGroupRepoForInv{
				findByIDFn: func(ctx context.Context, groupID string) (*models.Group, error) {
					return tt.findGroupResult, tt.findGroupErr
				},
			}

			svc := NewInvitationService(
				invitationRepo,
				userRepo,
				orgRepo,
				groupRepo,
				noopGroupServiceForInv(),
			)

			result, err := svc.CreateGroupOwner(context.Background(), tt.inviterID, tt.req)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				if tt.wantErrType != nil && !errors.Is(err, tt.wantErrType) {
					t.Errorf("expected error type %v, got %v", tt.wantErrType, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result == nil {
					t.Error("expected result, got nil")
				}
			}
		})
	}
}

func TestCreateGroupMember(t *testing.T) {
	tests := []struct {
		name                   string
		inviterID              string
		req                    *models.CreateGroupMemberInvitationRequest
		findGroupResult        *models.Group
		findGroupErr           error
		isSiteAdminResult      bool
		isSiteAdminErr         error
		isOrgAdminResult       bool
		isOrgAdminErr          error
		getUserRoleResult      string
		getUserRoleErr         error
		findInvitationResult   *models.Invitation
		findInvitationErr      error
		createInvitationResult *models.Invitation
		createInvitationErr    error
		wantErr                bool
		wantErrType            error
	}{
		{
			name:        "nil request returns validation error",
			inviterID:   "inviter-1",
			req:         nil,
			wantErr:     true,
			wantErrType: ErrValidation,
		},
		{
			name:        "empty email returns validation error",
			inviterID:   "inviter-1",
			req:         &models.CreateGroupMemberInvitationRequest{Email: "", GroupID: "group-1"},
			wantErr:     true,
			wantErrType: ErrValidation,
		},
		{
			name:        "whitespace email returns validation error",
			inviterID:   "inviter-1",
			req:         &models.CreateGroupMemberInvitationRequest{Email: "   ", GroupID: "group-1"},
			wantErr:     true,
			wantErrType: ErrValidation,
		},
		{
			name:        "empty groupID returns validation error",
			inviterID:   "inviter-1",
			req:         &models.CreateGroupMemberInvitationRequest{Email: "test@example.com", GroupID: ""},
			wantErr:     true,
			wantErrType: ErrValidation,
		},
		{
			name:        "whitespace groupID returns validation error",
			inviterID:   "inviter-1",
			req:         &models.CreateGroupMemberInvitationRequest{Email: "test@example.com", GroupID: "   "},
			wantErr:     true,
			wantErrType: ErrValidation,
		},
		{
			name:         "group not found returns error",
			inviterID:    "inviter-1",
			req:          &models.CreateGroupMemberInvitationRequest{Email: "test@example.com", GroupID: "group-1"},
			findGroupErr: errors.New("not found"),
			wantErr:      true,
			wantErrType:  ErrGroupNotFound,
		},
		{
			name:            "isSiteAdmin error returns wrapped error",
			inviterID:       "inviter-1",
			req:             &models.CreateGroupMemberInvitationRequest{Email: "test@example.com", GroupID: "group-1"},
			findGroupResult: &models.Group{ID: "group-1", OrganizationID: "org-1"},
			isSiteAdminErr:  errors.New("db error"),
			wantErr:         true,
		},
		{
			name:              "non-site admin, isOrgAdmin error returns wrapped error",
			inviterID:         "inviter-1",
			req:               &models.CreateGroupMemberInvitationRequest{Email: "test@example.com", GroupID: "group-1"},
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			isSiteAdminResult: false,
			isOrgAdminErr:     errors.New("db error"),
			wantErr:           true,
		},
		{
			name:              "non-site admin, non-org admin, getUserRole error returns wrapped error",
			inviterID:         "inviter-1",
			req:               &models.CreateGroupMemberInvitationRequest{Email: "test@example.com", GroupID: "group-1"},
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			isSiteAdminResult: false,
			isOrgAdminResult:  false,
			getUserRoleErr:    errors.New("db error"),
			wantErr:           true,
		},
		{
			name:              "non-site admin, non-org admin, non-owner returns forbidden",
			inviterID:         "inviter-1",
			req:               &models.CreateGroupMemberInvitationRequest{Email: "test@example.com", GroupID: "group-1"},
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			isSiteAdminResult: false,
			isOrgAdminResult:  false,
			getUserRoleResult: "MEMBER",
			wantErr:           true,
			wantErrType:       ErrForbidden,
		},
		{
			name:              "leader cannot create member invitation",
			inviterID:         "inviter-1",
			req:               &models.CreateGroupMemberInvitationRequest{Email: "test@example.com", GroupID: "group-1"},
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			isSiteAdminResult: false,
			isOrgAdminResult:  false,
			getUserRoleResult: "LEADER",
			wantErr:           true,
			wantErrType:       ErrForbidden,
		},
		{
			name:              "site admin can create invitation",
			inviterID:         "inviter-1",
			req:               &models.CreateGroupMemberInvitationRequest{Email: "test@example.com", GroupID: "group-1"},
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			isSiteAdminResult: true,
			findInvitationErr: errors.New("not found"),
			createInvitationResult: &models.Invitation{
				ID:      "inv-1",
				Email:   "test@example.com",
				Type:    models.InvitationTypeGroupMember,
				GroupID: strPtr("group-1"),
			},
			wantErr: false,
		},
		{
			name:              "org admin can create invitation",
			inviterID:         "inviter-1",
			req:               &models.CreateGroupMemberInvitationRequest{Email: "test@example.com", GroupID: "group-1"},
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			isSiteAdminResult: false,
			isOrgAdminResult:  true,
			findInvitationErr: errors.New("not found"),
			createInvitationResult: &models.Invitation{
				ID:      "inv-1",
				Email:   "test@example.com",
				Type:    models.InvitationTypeGroupMember,
				GroupID: strPtr("group-1"),
			},
			wantErr: false,
		},
		{
			name:              "group owner can create invitation",
			inviterID:         "inviter-1",
			req:               &models.CreateGroupMemberInvitationRequest{Email: "test@example.com", GroupID: "group-1"},
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			isSiteAdminResult: false,
			isOrgAdminResult:  false,
			getUserRoleResult: "OWNER",
			findInvitationErr: errors.New("not found"),
			createInvitationResult: &models.Invitation{
				ID:      "inv-1",
				Email:   "test@example.com",
				Type:    models.InvitationTypeGroupMember,
				GroupID: strPtr("group-1"),
			},
			wantErr: false,
		},
		{
			name:              "existing pending invitation returns conflict",
			inviterID:         "inviter-1",
			req:               &models.CreateGroupMemberInvitationRequest{Email: "test@example.com", GroupID: "group-1"},
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			isSiteAdminResult: true,
			findInvitationResult: &models.Invitation{
				ID:     "inv-1",
				Email:  "test@example.com",
				Status: models.InvitationStatusPending,
			},
			wantErr:     true,
			wantErrType: ErrConflict,
		},
		{
			name:                "createInvitation error returns wrapped error",
			inviterID:           "inviter-1",
			req:                 &models.CreateGroupMemberInvitationRequest{Email: "test@example.com", GroupID: "group-1"},
			findGroupResult:     &models.Group{ID: "group-1", OrganizationID: "org-1"},
			isSiteAdminResult:   true,
			findInvitationErr:   errors.New("not found"),
			createInvitationErr: errors.New("db error"),
			wantErr:             true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			invitationRepo := &mockInvitationRepo{
				findByEmailFn: func(ctx context.Context, email string, invType models.InvitationType, targetID string) (*models.Invitation, error) {
					return tt.findInvitationResult, tt.findInvitationErr
				},
				createFn: func(ctx context.Context, email string, invType models.InvitationType, orgID, groupID *string, inviterID string, token string, expiresAt time.Time) (*models.Invitation, error) {
					return tt.createInvitationResult, tt.createInvitationErr
				},
			}
			userRepo := &mockUserRepoForInv{
				isSiteAdminFn: func(ctx context.Context, userID string) (bool, error) {
					return tt.isSiteAdminResult, tt.isSiteAdminErr
				},
			}
			orgRepo := &mockOrgRepoForInv{
				isAdminFn: func(ctx context.Context, orgID, userID string) (bool, error) {
					return tt.isOrgAdminResult, tt.isOrgAdminErr
				},
			}
			groupRepo := &mockGroupRepoForInv{
				findByIDFn: func(ctx context.Context, groupID string) (*models.Group, error) {
					return tt.findGroupResult, tt.findGroupErr
				},
			}
			groupService := &mockGroupServiceForInv{
				getUserRoleFn: func(ctx context.Context, groupID, userID string) (string, error) {
					return tt.getUserRoleResult, tt.getUserRoleErr
				},
			}

			svc := NewInvitationService(
				invitationRepo,
				userRepo,
				orgRepo,
				groupRepo,
				groupService,
			)

			result, err := svc.CreateGroupMember(context.Background(), tt.inviterID, tt.req)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				if tt.wantErrType != nil && !errors.Is(err, tt.wantErrType) {
					t.Errorf("expected error type %v, got %v", tt.wantErrType, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result == nil {
					t.Error("expected result, got nil")
				}
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name                 string
		invitationID         string
		userID               string
		findInvitationResult *models.Invitation
		findInvitationErr    error
		isSiteAdminResult    bool
		isSiteAdminErr       error
		findGroupResult      *models.Group
		findGroupErr         error
		isOrgAdminResult     bool
		isOrgAdminErr        error
		getUserRoleResult    string
		getUserRoleErr       error
		deleteInvitationErr  error
		wantErr              bool
		wantErrType          error
	}{
		{
			name:         "empty invitationID returns validation error",
			invitationID: "",
			userID:       "user-1",
			wantErr:      true,
			wantErrType:  ErrValidation,
		},
		{
			name:         "whitespace invitationID returns validation error",
			invitationID: "   ",
			userID:       "user-1",
			wantErr:      true,
			wantErrType:  ErrValidation,
		},
		{
			name:              "invitation not found returns error",
			invitationID:      "inv-1",
			userID:            "user-1",
			findInvitationErr: errors.New("not found"),
			wantErr:           true,
			wantErrType:       ErrInvitationNotFound,
		},
		{
			name:         "accepted invitation returns conflict",
			invitationID: "inv-1",
			userID:       "user-1",
			findInvitationResult: &models.Invitation{
				ID:     "inv-1",
				Status: models.InvitationStatusAccepted,
				Type:   models.InvitationTypeOrgAdmin,
			},
			wantErr:     true,
			wantErrType: ErrConflict,
		},
		{
			name:         "isSiteAdmin error returns wrapped error",
			invitationID: "inv-1",
			userID:       "user-1",
			findInvitationResult: &models.Invitation{
				ID:     "inv-1",
				Status: models.InvitationStatusPending,
				Type:   models.InvitationTypeOrgAdmin,
			},
			isSiteAdminErr: errors.New("db error"),
			wantErr:        true,
		},
		{
			name:         "site admin can delete org admin invitation",
			invitationID: "inv-1",
			userID:       "site-admin",
			findInvitationResult: &models.Invitation{
				ID:             "inv-1",
				Status:         models.InvitationStatusPending,
				Type:           models.InvitationTypeOrgAdmin,
				OrganizationID: strPtr("org-1"),
			},
			isSiteAdminResult: true,
			wantErr:           false,
		},
		{
			name:         "non-site admin cannot delete org admin invitation",
			invitationID: "inv-1",
			userID:       "user-1",
			findInvitationResult: &models.Invitation{
				ID:             "inv-1",
				Status:         models.InvitationStatusPending,
				Type:           models.InvitationTypeOrgAdmin,
				OrganizationID: strPtr("org-1"),
			},
			isSiteAdminResult: false,
			wantErr:           true,
			wantErrType:       ErrForbidden,
		},
		{
			name:         "site admin can delete group owner invitation",
			invitationID: "inv-1",
			userID:       "site-admin",
			findInvitationResult: &models.Invitation{
				ID:      "inv-1",
				Status:  models.InvitationStatusPending,
				Type:    models.InvitationTypeGroupOwner,
				GroupID: strPtr("group-1"),
			},
			isSiteAdminResult: true,
			wantErr:           false,
		},
		{
			name:         "group owner invitation with nil groupID returns forbidden",
			invitationID: "inv-1",
			userID:       "user-1",
			findInvitationResult: &models.Invitation{
				ID:      "inv-1",
				Status:  models.InvitationStatusPending,
				Type:    models.InvitationTypeGroupOwner,
				GroupID: nil,
			},
			isSiteAdminResult: false,
			wantErr:           true,
			wantErrType:       ErrForbidden,
		},
		{
			name:         "group owner invitation findGroup error returns wrapped error",
			invitationID: "inv-1",
			userID:       "user-1",
			findInvitationResult: &models.Invitation{
				ID:      "inv-1",
				Status:  models.InvitationStatusPending,
				Type:    models.InvitationTypeGroupOwner,
				GroupID: strPtr("group-1"),
			},
			isSiteAdminResult: false,
			findGroupErr:      errors.New("db error"),
			wantErr:           true,
		},
		{
			name:         "group owner invitation isOrgAdmin error returns wrapped error",
			invitationID: "inv-1",
			userID:       "user-1",
			findInvitationResult: &models.Invitation{
				ID:      "inv-1",
				Status:  models.InvitationStatusPending,
				Type:    models.InvitationTypeGroupOwner,
				GroupID: strPtr("group-1"),
			},
			isSiteAdminResult: false,
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			isOrgAdminErr:     errors.New("db error"),
			wantErr:           true,
		},
		{
			name:         "org admin can delete group owner invitation",
			invitationID: "inv-1",
			userID:       "org-admin",
			findInvitationResult: &models.Invitation{
				ID:      "inv-1",
				Status:  models.InvitationStatusPending,
				Type:    models.InvitationTypeGroupOwner,
				GroupID: strPtr("group-1"),
			},
			isSiteAdminResult: false,
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			isOrgAdminResult:  true,
			wantErr:           false,
		},
		{
			name:         "non-org admin cannot delete group owner invitation",
			invitationID: "inv-1",
			userID:       "user-1",
			findInvitationResult: &models.Invitation{
				ID:      "inv-1",
				Status:  models.InvitationStatusPending,
				Type:    models.InvitationTypeGroupOwner,
				GroupID: strPtr("group-1"),
			},
			isSiteAdminResult: false,
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			isOrgAdminResult:  false,
			wantErr:           true,
			wantErrType:       ErrForbidden,
		},
		{
			name:         "site admin can delete group member invitation",
			invitationID: "inv-1",
			userID:       "site-admin",
			findInvitationResult: &models.Invitation{
				ID:      "inv-1",
				Status:  models.InvitationStatusPending,
				Type:    models.InvitationTypeGroupMember,
				GroupID: strPtr("group-1"),
			},
			isSiteAdminResult: true,
			wantErr:           false,
		},
		{
			name:         "group member invitation with nil groupID returns forbidden",
			invitationID: "inv-1",
			userID:       "user-1",
			findInvitationResult: &models.Invitation{
				ID:      "inv-1",
				Status:  models.InvitationStatusPending,
				Type:    models.InvitationTypeGroupMember,
				GroupID: nil,
			},
			isSiteAdminResult: false,
			wantErr:           true,
			wantErrType:       ErrForbidden,
		},
		{
			name:         "group member invitation findGroup error returns wrapped error",
			invitationID: "inv-1",
			userID:       "user-1",
			findInvitationResult: &models.Invitation{
				ID:      "inv-1",
				Status:  models.InvitationStatusPending,
				Type:    models.InvitationTypeGroupMember,
				GroupID: strPtr("group-1"),
			},
			isSiteAdminResult: false,
			findGroupErr:      errors.New("db error"),
			wantErr:           true,
		},
		{
			name:         "group member invitation isOrgAdmin error returns wrapped error",
			invitationID: "inv-1",
			userID:       "user-1",
			findInvitationResult: &models.Invitation{
				ID:      "inv-1",
				Status:  models.InvitationStatusPending,
				Type:    models.InvitationTypeGroupMember,
				GroupID: strPtr("group-1"),
			},
			isSiteAdminResult: false,
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			isOrgAdminErr:     errors.New("db error"),
			wantErr:           true,
		},
		{
			name:         "org admin can delete group member invitation",
			invitationID: "inv-1",
			userID:       "org-admin",
			findInvitationResult: &models.Invitation{
				ID:      "inv-1",
				Status:  models.InvitationStatusPending,
				Type:    models.InvitationTypeGroupMember,
				GroupID: strPtr("group-1"),
			},
			isSiteAdminResult: false,
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			isOrgAdminResult:  true,
			wantErr:           false,
		},
		{
			name:         "group member invitation getUserRole error returns wrapped error",
			invitationID: "inv-1",
			userID:       "user-1",
			findInvitationResult: &models.Invitation{
				ID:      "inv-1",
				Status:  models.InvitationStatusPending,
				Type:    models.InvitationTypeGroupMember,
				GroupID: strPtr("group-1"),
			},
			isSiteAdminResult: false,
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			isOrgAdminResult:  false,
			getUserRoleErr:    errors.New("db error"),
			wantErr:           true,
		},
		{
			name:         "group owner can delete group member invitation",
			invitationID: "inv-1",
			userID:       "group-owner",
			findInvitationResult: &models.Invitation{
				ID:      "inv-1",
				Status:  models.InvitationStatusPending,
				Type:    models.InvitationTypeGroupMember,
				GroupID: strPtr("group-1"),
			},
			isSiteAdminResult: false,
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			isOrgAdminResult:  false,
			getUserRoleResult: "OWNER",
			wantErr:           false,
		},
		{
			name:         "group member cannot delete group member invitation",
			invitationID: "inv-1",
			userID:       "user-1",
			findInvitationResult: &models.Invitation{
				ID:      "inv-1",
				Status:  models.InvitationStatusPending,
				Type:    models.InvitationTypeGroupMember,
				GroupID: strPtr("group-1"),
			},
			isSiteAdminResult: false,
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			isOrgAdminResult:  false,
			getUserRoleResult: "MEMBER",
			wantErr:           true,
			wantErrType:       ErrForbidden,
		},
		{
			name:         "deleteInvitation error returns wrapped error",
			invitationID: "inv-1",
			userID:       "site-admin",
			findInvitationResult: &models.Invitation{
				ID:     "inv-1",
				Status: models.InvitationStatusPending,
				Type:   models.InvitationTypeOrgAdmin,
			},
			isSiteAdminResult:   true,
			deleteInvitationErr: errors.New("db error"),
			wantErr:             true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			invitationRepo := &mockInvitationRepo{
				findByIDFn: func(ctx context.Context, invitationID string) (*models.Invitation, error) {
					return tt.findInvitationResult, tt.findInvitationErr
				},
				deleteFn: func(ctx context.Context, invitationID string) error {
					return tt.deleteInvitationErr
				},
			}
			userRepo := &mockUserRepoForInv{
				isSiteAdminFn: func(ctx context.Context, userID string) (bool, error) {
					return tt.isSiteAdminResult, tt.isSiteAdminErr
				},
			}
			orgRepo := &mockOrgRepoForInv{
				isAdminFn: func(ctx context.Context, orgID, userID string) (bool, error) {
					return tt.isOrgAdminResult, tt.isOrgAdminErr
				},
			}
			groupRepo := &mockGroupRepoForInv{
				findByIDFn: func(ctx context.Context, groupID string) (*models.Group, error) {
					return tt.findGroupResult, tt.findGroupErr
				},
			}
			groupService := &mockGroupServiceForInv{
				getUserRoleFn: func(ctx context.Context, groupID, userID string) (string, error) {
					return tt.getUserRoleResult, tt.getUserRoleErr
				},
			}

			svc := NewInvitationService(
				invitationRepo,
				userRepo,
				orgRepo,
				groupRepo,
				groupService,
			)

			err := svc.Delete(context.Background(), tt.invitationID, tt.userID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				if tt.wantErrType != nil && !errors.Is(err, tt.wantErrType) {
					t.Errorf("expected error type %v, got %v", tt.wantErrType, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestAccept(t *testing.T) {
	validExpiresAt := time.Now().Add(24 * time.Hour)
	expiredAt := time.Now().Add(-24 * time.Hour)

	tests := []struct {
		name                      string
		req                       *models.AcceptInvitationRequest
		findInvitationResult      *models.Invitation
		findInvitationErr         error
		findUserByEmailResult     *models.User
		findUserByEmailErr        error
		createUserResult          *models.User
		createUserErr             error
		addOrgAdminErr            error
		addGroupMemberResult      *models.GroupMember
		addGroupMemberErr         error
		updateInvitationStatusErr error
		wantErr                   bool
		wantErrType               error
	}{
		{
			name:        "nil request returns validation error",
			req:         nil,
			wantErr:     true,
			wantErrType: ErrValidation,
		},
		{
			name:        "empty token returns validation error",
			req:         &models.AcceptInvitationRequest{Token: "", ExternalID: "ext-1", DisplayName: "Test"},
			wantErr:     true,
			wantErrType: ErrValidation,
		},
		{
			name:        "whitespace token returns validation error",
			req:         &models.AcceptInvitationRequest{Token: "   ", ExternalID: "ext-1", DisplayName: "Test"},
			wantErr:     true,
			wantErrType: ErrValidation,
		},
		{
			name:        "empty externalID returns validation error",
			req:         &models.AcceptInvitationRequest{Token: "token-1", ExternalID: "", DisplayName: "Test"},
			wantErr:     true,
			wantErrType: ErrValidation,
		},
		{
			name:        "whitespace externalID returns validation error",
			req:         &models.AcceptInvitationRequest{Token: "token-1", ExternalID: "   ", DisplayName: "Test"},
			wantErr:     true,
			wantErrType: ErrValidation,
		},
		{
			name:        "empty displayName returns validation error",
			req:         &models.AcceptInvitationRequest{Token: "token-1", ExternalID: "ext-1", DisplayName: ""},
			wantErr:     true,
			wantErrType: ErrValidation,
		},
		{
			name:        "whitespace displayName returns validation error",
			req:         &models.AcceptInvitationRequest{Token: "token-1", ExternalID: "ext-1", DisplayName: "   "},
			wantErr:     true,
			wantErrType: ErrValidation,
		},
		{
			name:              "invitation not found returns error",
			req:               &models.AcceptInvitationRequest{Token: "token-1", ExternalID: "ext-1", DisplayName: "Test"},
			findInvitationErr: errors.New("not found"),
			wantErr:           true,
			wantErrType:       ErrInvitationNotFound,
		},
		{
			name: "already accepted invitation returns conflict",
			req:  &models.AcceptInvitationRequest{Token: "token-1", ExternalID: "ext-1", DisplayName: "Test"},
			findInvitationResult: &models.Invitation{
				ID:        "inv-1",
				Status:    models.InvitationStatusAccepted,
				ExpiresAt: validExpiresAt,
			},
			wantErr:     true,
			wantErrType: ErrConflict,
		},
		{
			name: "expired status invitation returns invalid token",
			req:  &models.AcceptInvitationRequest{Token: "token-1", ExternalID: "ext-1", DisplayName: "Test"},
			findInvitationResult: &models.Invitation{
				ID:        "inv-1",
				Status:    models.InvitationStatusExpired,
				ExpiresAt: validExpiresAt,
			},
			wantErr:     true,
			wantErrType: ErrInvalidToken,
		},
		{
			name: "revoked invitation returns invalid token",
			req:  &models.AcceptInvitationRequest{Token: "token-1", ExternalID: "ext-1", DisplayName: "Test"},
			findInvitationResult: &models.Invitation{
				ID:        "inv-1",
				Status:    models.InvitationStatusRevoked,
				ExpiresAt: validExpiresAt,
			},
			wantErr:     true,
			wantErrType: ErrInvalidToken,
		},
		{
			name: "expired by time invitation returns invalid token",
			req:  &models.AcceptInvitationRequest{Token: "token-1", ExternalID: "ext-1", DisplayName: "Test"},
			findInvitationResult: &models.Invitation{
				ID:        "inv-1",
				Status:    models.InvitationStatusPending,
				ExpiresAt: expiredAt,
			},
			wantErr:     true,
			wantErrType: ErrInvalidToken,
		},
		{
			name: "createUser error returns wrapped error",
			req:  &models.AcceptInvitationRequest{Token: "token-1", ExternalID: "ext-1", DisplayName: "Test"},
			findInvitationResult: &models.Invitation{
				ID:             "inv-1",
				Email:          "test@example.com",
				Status:         models.InvitationStatusPending,
				Type:           models.InvitationTypeOrgAdmin,
				OrganizationID: strPtr("org-1"),
				ExpiresAt:      validExpiresAt,
			},
			findUserByEmailErr: errors.New("not found"),
			createUserErr:      errors.New("db error"),
			wantErr:            true,
		},
		{
			name: "org admin with nil organizationID returns validation error",
			req:  &models.AcceptInvitationRequest{Token: "token-1", ExternalID: "ext-1", DisplayName: "Test"},
			findInvitationResult: &models.Invitation{
				ID:             "inv-1",
				Email:          "test@example.com",
				Status:         models.InvitationStatusPending,
				Type:           models.InvitationTypeOrgAdmin,
				OrganizationID: nil,
				ExpiresAt:      validExpiresAt,
			},
			findUserByEmailResult: &models.User{ID: "user-1", Email: "test@example.com"},
			wantErr:               true,
			wantErrType:           ErrValidation,
		},
		{
			name: "addOrgAdmin error returns wrapped error",
			req:  &models.AcceptInvitationRequest{Token: "token-1", ExternalID: "ext-1", DisplayName: "Test"},
			findInvitationResult: &models.Invitation{
				ID:             "inv-1",
				Email:          "test@example.com",
				Status:         models.InvitationStatusPending,
				Type:           models.InvitationTypeOrgAdmin,
				OrganizationID: strPtr("org-1"),
				ExpiresAt:      validExpiresAt,
			},
			findUserByEmailResult: &models.User{ID: "user-1", Email: "test@example.com"},
			addOrgAdminErr:        errors.New("db error"),
			wantErr:               true,
		},
		{
			name: "group owner with nil groupID returns validation error",
			req:  &models.AcceptInvitationRequest{Token: "token-1", ExternalID: "ext-1", DisplayName: "Test"},
			findInvitationResult: &models.Invitation{
				ID:        "inv-1",
				Email:     "test@example.com",
				Status:    models.InvitationStatusPending,
				Type:      models.InvitationTypeGroupOwner,
				GroupID:   nil,
				ExpiresAt: validExpiresAt,
			},
			findUserByEmailResult: &models.User{ID: "user-1", Email: "test@example.com"},
			wantErr:               true,
			wantErrType:           ErrValidation,
		},
		{
			name: "group owner addGroupMember error returns wrapped error",
			req:  &models.AcceptInvitationRequest{Token: "token-1", ExternalID: "ext-1", DisplayName: "Test"},
			findInvitationResult: &models.Invitation{
				ID:        "inv-1",
				Email:     "test@example.com",
				Status:    models.InvitationStatusPending,
				Type:      models.InvitationTypeGroupOwner,
				GroupID:   strPtr("group-1"),
				ExpiresAt: validExpiresAt,
			},
			findUserByEmailResult: &models.User{ID: "user-1", Email: "test@example.com"},
			addGroupMemberErr:     errors.New("db error"),
			wantErr:               true,
		},
		{
			name: "group member with nil groupID returns validation error",
			req:  &models.AcceptInvitationRequest{Token: "token-1", ExternalID: "ext-1", DisplayName: "Test"},
			findInvitationResult: &models.Invitation{
				ID:        "inv-1",
				Email:     "test@example.com",
				Status:    models.InvitationStatusPending,
				Type:      models.InvitationTypeGroupMember,
				GroupID:   nil,
				ExpiresAt: validExpiresAt,
			},
			findUserByEmailResult: &models.User{ID: "user-1", Email: "test@example.com"},
			wantErr:               true,
			wantErrType:           ErrValidation,
		},
		{
			name: "group member addGroupMember error returns wrapped error",
			req:  &models.AcceptInvitationRequest{Token: "token-1", ExternalID: "ext-1", DisplayName: "Test"},
			findInvitationResult: &models.Invitation{
				ID:        "inv-1",
				Email:     "test@example.com",
				Status:    models.InvitationStatusPending,
				Type:      models.InvitationTypeGroupMember,
				GroupID:   strPtr("group-1"),
				ExpiresAt: validExpiresAt,
			},
			findUserByEmailResult: &models.User{ID: "user-1", Email: "test@example.com"},
			addGroupMemberErr:     errors.New("db error"),
			wantErr:               true,
		},
		{
			name: "updateInvitationStatus error returns wrapped error",
			req:  &models.AcceptInvitationRequest{Token: "token-1", ExternalID: "ext-1", DisplayName: "Test"},
			findInvitationResult: &models.Invitation{
				ID:             "inv-1",
				Email:          "test@example.com",
				Status:         models.InvitationStatusPending,
				Type:           models.InvitationTypeOrgAdmin,
				OrganizationID: strPtr("org-1"),
				ExpiresAt:      validExpiresAt,
			},
			findUserByEmailResult:     &models.User{ID: "user-1", Email: "test@example.com"},
			updateInvitationStatusErr: errors.New("db error"),
			wantErr:                   true,
		},
		{
			name: "successful org admin accept with existing user",
			req:  &models.AcceptInvitationRequest{Token: "token-1", ExternalID: "ext-1", DisplayName: "Test"},
			findInvitationResult: &models.Invitation{
				ID:             "inv-1",
				Email:          "test@example.com",
				Status:         models.InvitationStatusPending,
				Type:           models.InvitationTypeOrgAdmin,
				OrganizationID: strPtr("org-1"),
				ExpiresAt:      validExpiresAt,
			},
			findUserByEmailResult: &models.User{ID: "user-1", Email: "test@example.com"},
			wantErr:               false,
		},
		{
			name: "successful org admin accept with new user",
			req:  &models.AcceptInvitationRequest{Token: "token-1", ExternalID: "ext-1", DisplayName: "Test"},
			findInvitationResult: &models.Invitation{
				ID:             "inv-1",
				Email:          "test@example.com",
				Status:         models.InvitationStatusPending,
				Type:           models.InvitationTypeOrgAdmin,
				OrganizationID: strPtr("org-1"),
				ExpiresAt:      validExpiresAt,
			},
			findUserByEmailErr: errors.New("not found"),
			createUserResult:   &models.User{ID: "user-1", Email: "test@example.com"},
			wantErr:            false,
		},
		{
			name: "successful group owner accept",
			req:  &models.AcceptInvitationRequest{Token: "token-1", ExternalID: "ext-1", DisplayName: "Test"},
			findInvitationResult: &models.Invitation{
				ID:        "inv-1",
				Email:     "test@example.com",
				Status:    models.InvitationStatusPending,
				Type:      models.InvitationTypeGroupOwner,
				GroupID:   strPtr("group-1"),
				ExpiresAt: validExpiresAt,
			},
			findUserByEmailResult: &models.User{ID: "user-1", Email: "test@example.com"},
			addGroupMemberResult:  &models.GroupMember{ID: "member-1", GroupID: "group-1", UserID: "user-1", Role: "OWNER"},
			wantErr:               false,
		},
		{
			name: "successful group member accept",
			req:  &models.AcceptInvitationRequest{Token: "token-1", ExternalID: "ext-1", DisplayName: "Test"},
			findInvitationResult: &models.Invitation{
				ID:        "inv-1",
				Email:     "test@example.com",
				Status:    models.InvitationStatusPending,
				Type:      models.InvitationTypeGroupMember,
				GroupID:   strPtr("group-1"),
				ExpiresAt: validExpiresAt,
			},
			findUserByEmailResult: &models.User{ID: "user-1", Email: "test@example.com"},
			addGroupMemberResult:  &models.GroupMember{ID: "member-1", GroupID: "group-1", UserID: "user-1", Role: "MEMBER"},
			wantErr:               false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			invitationRepo := &mockInvitationRepo{
				findByTokenFn: func(ctx context.Context, token string) (*models.Invitation, error) {
					return tt.findInvitationResult, tt.findInvitationErr
				},
				updateStatusFn: func(ctx context.Context, invitationID string, status models.InvitationStatus, acceptedByID *string) error {
					return tt.updateInvitationStatusErr
				},
			}
			userRepo := &mockUserRepoForInv{
				findByEmailFn: func(ctx context.Context, email string) (*models.User, error) {
					return tt.findUserByEmailResult, tt.findUserByEmailErr
				},
				createFn: func(ctx context.Context, email, displayName, externalID string) (*models.User, error) {
					return tt.createUserResult, tt.createUserErr
				},
			}
			orgRepo := &mockOrgRepoForInv{
				addAdminFn: func(ctx context.Context, orgID, userID string) error {
					return tt.addOrgAdminErr
				},
			}
			groupRepo := &mockGroupRepoForInv{
				addMemberFn: func(ctx context.Context, groupID, userID, role string) (*models.GroupMember, error) {
					return tt.addGroupMemberResult, tt.addGroupMemberErr
				},
			}

			svc := NewInvitationService(
				invitationRepo,
				userRepo,
				orgRepo,
				groupRepo,
				noopGroupServiceForInv(),
			)

			invitation, user, err := svc.Accept(context.Background(), tt.req)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				if tt.wantErrType != nil && !errors.Is(err, tt.wantErrType) {
					t.Errorf("expected error type %v, got %v", tt.wantErrType, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if invitation == nil {
					t.Error("expected invitation, got nil")
				}
				if user == nil {
					t.Error("expected user, got nil")
				}
				if invitation != nil && invitation.Status != models.InvitationStatusAccepted {
					t.Errorf("expected status ACCEPTED, got %v", invitation.Status)
				}
			}
		})
	}
}

func TestValidate(t *testing.T) {
	validExpiresAt := time.Now().Add(24 * time.Hour)
	expiredAt := time.Now().Add(-24 * time.Hour)

	tests := []struct {
		name                 string
		token                string
		findInvitationResult *models.Invitation
		findInvitationErr    error
		findOrgResult        *models.Organization
		findOrgErr           error
		findGroupResult      *models.Group
		findGroupErr         error
		findUserResult       *models.User
		findUserErr          error
		wantValid            bool
		wantOrgName          bool
		wantGroupName        bool
		wantInviterName      bool
	}{
		{
			name:      "empty token returns invalid",
			token:     "",
			wantValid: false,
		},
		{
			name:      "whitespace token returns invalid",
			token:     "   ",
			wantValid: false,
		},
		{
			name:              "invitation not found returns invalid",
			token:             "token-1",
			findInvitationErr: errors.New("not found"),
			wantValid:         false,
		},
		{
			name:  "expired invitation returns invalid",
			token: "token-1",
			findInvitationResult: &models.Invitation{
				ID:        "inv-1",
				Status:    models.InvitationStatusPending,
				ExpiresAt: expiredAt,
			},
			wantValid: false,
		},
		{
			name:  "accepted invitation returns invalid",
			token: "token-1",
			findInvitationResult: &models.Invitation{
				ID:        "inv-1",
				Status:    models.InvitationStatusAccepted,
				ExpiresAt: validExpiresAt,
			},
			wantValid: false,
		},
		{
			name:  "expired status invitation returns invalid",
			token: "token-1",
			findInvitationResult: &models.Invitation{
				ID:        "inv-1",
				Status:    models.InvitationStatusExpired,
				ExpiresAt: validExpiresAt,
			},
			wantValid: false,
		},
		{
			name:  "revoked invitation returns invalid",
			token: "token-1",
			findInvitationResult: &models.Invitation{
				ID:        "inv-1",
				Status:    models.InvitationStatusRevoked,
				ExpiresAt: validExpiresAt,
			},
			wantValid: false,
		},
		{
			name:  "valid org admin invitation returns org name",
			token: "token-1",
			findInvitationResult: &models.Invitation{
				ID:             "inv-1",
				Email:          "test@example.com",
				Status:         models.InvitationStatusPending,
				Type:           models.InvitationTypeOrgAdmin,
				OrganizationID: strPtr("org-1"),
				InviterID:      "inviter-1",
				ExpiresAt:      validExpiresAt,
			},
			findOrgResult:   &models.Organization{ID: "org-1", Name: "Test Org"},
			findUserResult:  &models.User{ID: "inviter-1", DisplayName: "Inviter"},
			wantValid:       true,
			wantOrgName:     true,
			wantInviterName: true,
		},
		{
			name:  "valid org admin invitation with org lookup failure still returns valid",
			token: "token-1",
			findInvitationResult: &models.Invitation{
				ID:             "inv-1",
				Email:          "test@example.com",
				Status:         models.InvitationStatusPending,
				Type:           models.InvitationTypeOrgAdmin,
				OrganizationID: strPtr("org-1"),
				InviterID:      "inviter-1",
				ExpiresAt:      validExpiresAt,
			},
			findOrgErr:      errors.New("db error"),
			findUserResult:  &models.User{ID: "inviter-1", DisplayName: "Inviter"},
			wantValid:       true,
			wantOrgName:     false,
			wantInviterName: true,
		},
		{
			name:  "valid group owner invitation returns group name",
			token: "token-1",
			findInvitationResult: &models.Invitation{
				ID:        "inv-1",
				Email:     "test@example.com",
				Status:    models.InvitationStatusPending,
				Type:      models.InvitationTypeGroupOwner,
				GroupID:   strPtr("group-1"),
				InviterID: "inviter-1",
				ExpiresAt: validExpiresAt,
			},
			findGroupResult: &models.Group{ID: "group-1", Name: "Test Group"},
			findUserResult:  &models.User{ID: "inviter-1", DisplayName: "Inviter"},
			wantValid:       true,
			wantGroupName:   true,
			wantInviterName: true,
		},
		{
			name:  "valid group member invitation returns group name",
			token: "token-1",
			findInvitationResult: &models.Invitation{
				ID:        "inv-1",
				Email:     "test@example.com",
				Status:    models.InvitationStatusPending,
				Type:      models.InvitationTypeGroupMember,
				GroupID:   strPtr("group-1"),
				InviterID: "inviter-1",
				ExpiresAt: validExpiresAt,
			},
			findGroupResult: &models.Group{ID: "group-1", Name: "Test Group"},
			findUserResult:  &models.User{ID: "inviter-1", DisplayName: "Inviter"},
			wantValid:       true,
			wantGroupName:   true,
			wantInviterName: true,
		},
		{
			name:  "valid invitation with group lookup failure still returns valid",
			token: "token-1",
			findInvitationResult: &models.Invitation{
				ID:        "inv-1",
				Email:     "test@example.com",
				Status:    models.InvitationStatusPending,
				Type:      models.InvitationTypeGroupOwner,
				GroupID:   strPtr("group-1"),
				InviterID: "inviter-1",
				ExpiresAt: validExpiresAt,
			},
			findGroupErr:    errors.New("db error"),
			findUserResult:  &models.User{ID: "inviter-1", DisplayName: "Inviter"},
			wantValid:       true,
			wantGroupName:   false,
			wantInviterName: true,
		},
		{
			name:  "valid invitation with inviter lookup failure still returns valid",
			token: "token-1",
			findInvitationResult: &models.Invitation{
				ID:             "inv-1",
				Email:          "test@example.com",
				Status:         models.InvitationStatusPending,
				Type:           models.InvitationTypeOrgAdmin,
				OrganizationID: strPtr("org-1"),
				InviterID:      "inviter-1",
				ExpiresAt:      validExpiresAt,
			},
			findOrgResult:   &models.Organization{ID: "org-1", Name: "Test Org"},
			findUserErr:     errors.New("db error"),
			wantValid:       true,
			wantOrgName:     true,
			wantInviterName: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			invitationRepo := &mockInvitationRepo{
				findByTokenFn: func(ctx context.Context, token string) (*models.Invitation, error) {
					return tt.findInvitationResult, tt.findInvitationErr
				},
			}
			userRepo := &mockUserRepoForInv{
				findByIDFn: func(ctx context.Context, userID string) (*models.User, error) {
					return tt.findUserResult, tt.findUserErr
				},
			}
			orgRepo := &mockOrgRepoForInv{
				findByIDFn: func(ctx context.Context, orgID string) (*models.Organization, error) {
					return tt.findOrgResult, tt.findOrgErr
				},
			}
			groupRepo := &mockGroupRepoForInv{
				findByIDFn: func(ctx context.Context, groupID string) (*models.Group, error) {
					return tt.findGroupResult, tt.findGroupErr
				},
			}

			svc := NewInvitationService(
				invitationRepo,
				userRepo,
				orgRepo,
				groupRepo,
				noopGroupServiceForInv(),
			)

			result, err := svc.Validate(context.Background(), tt.token)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if result.Valid != tt.wantValid {
				t.Errorf("expected Valid=%v, got %v", tt.wantValid, result.Valid)
			}

			if tt.wantValid {
				if tt.wantOrgName && result.OrganizationName == nil {
					t.Error("expected organization name, got nil")
				}
				if !tt.wantOrgName && result.OrganizationName != nil {
					t.Errorf("expected no organization name, got %v", *result.OrganizationName)
				}
				if tt.wantGroupName && result.GroupName == nil {
					t.Error("expected group name, got nil")
				}
				if !tt.wantGroupName && result.GroupName != nil {
					t.Errorf("expected no group name, got %v", *result.GroupName)
				}
				if tt.wantInviterName && result.InviterName == nil {
					t.Error("expected inviter name, got nil")
				}
				if !tt.wantInviterName && result.InviterName != nil {
					t.Errorf("expected no inviter name, got %v", *result.InviterName)
				}
			}
		})
	}
}

func TestGetByID(t *testing.T) {
	tests := []struct {
		name                 string
		invitationID         string
		findInvitationResult *models.Invitation
		findInvitationErr    error
		wantErr              bool
		wantErrType          error
	}{
		{
			name:         "empty invitationID returns validation error",
			invitationID: "",
			wantErr:      true,
			wantErrType:  ErrValidation,
		},
		{
			name:         "whitespace invitationID returns validation error",
			invitationID: "   ",
			wantErr:      true,
			wantErrType:  ErrValidation,
		},
		{
			name:              "invitation not found returns error",
			invitationID:      "inv-1",
			findInvitationErr: errors.New("not found"),
			wantErr:           true,
			wantErrType:       ErrInvitationNotFound,
		},
		{
			name:         "successful retrieval returns invitation",
			invitationID: "inv-1",
			findInvitationResult: &models.Invitation{
				ID:    "inv-1",
				Email: "test@example.com",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			invitationRepo := &mockInvitationRepo{
				findByIDFn: func(ctx context.Context, invitationID string) (*models.Invitation, error) {
					return tt.findInvitationResult, tt.findInvitationErr
				},
			}

			svc := NewInvitationService(
				invitationRepo,
				noopUserRepoForInv(),
				noopOrgRepoForInv(),
				noopGroupRepoForInv(),
				noopGroupServiceForInv(),
			)

			result, err := svc.GetByID(context.Background(), tt.invitationID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				if tt.wantErrType != nil && !errors.Is(err, tt.wantErrType) {
					t.Errorf("expected error type %v, got %v", tt.wantErrType, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result == nil {
					t.Error("expected result, got nil")
				}
			}
		})
	}
}

func TestGetByID_ContextPropagation(t *testing.T) {
	type ctxKey string
	key := ctxKey("test")
	ctx := context.WithValue(context.Background(), key, "value")

	var capturedCtx context.Context

	invitationRepo := &mockInvitationRepo{
		findByIDFn: func(ctx context.Context, invitationID string) (*models.Invitation, error) {
			capturedCtx = ctx
			return &models.Invitation{ID: "inv-1"}, nil
		},
	}

	svc := NewInvitationService(
		invitationRepo,
		noopUserRepoForInv(),
		noopOrgRepoForInv(),
		noopGroupRepoForInv(),
		noopGroupServiceForInv(),
	)

	_, _ = svc.GetByID(ctx, "inv-1")

	if capturedCtx.Value(key) != "value" {
		t.Error("context was not properly propagated")
	}
}

func TestGetByToken(t *testing.T) {
	tests := []struct {
		name                 string
		token                string
		findInvitationResult *models.Invitation
		findInvitationErr    error
		wantErr              bool
		wantErrType          error
	}{
		{
			name:        "empty token returns validation error",
			token:       "",
			wantErr:     true,
			wantErrType: ErrValidation,
		},
		{
			name:        "whitespace token returns validation error",
			token:       "   ",
			wantErr:     true,
			wantErrType: ErrValidation,
		},
		{
			name:              "invitation not found returns error",
			token:             "token-1",
			findInvitationErr: errors.New("not found"),
			wantErr:           true,
			wantErrType:       ErrInvitationNotFound,
		},
		{
			name:  "successful retrieval returns invitation",
			token: "token-1",
			findInvitationResult: &models.Invitation{
				ID:    "inv-1",
				Token: "token-1",
				Email: "test@example.com",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			invitationRepo := &mockInvitationRepo{
				findByTokenFn: func(ctx context.Context, token string) (*models.Invitation, error) {
					return tt.findInvitationResult, tt.findInvitationErr
				},
			}

			svc := NewInvitationService(
				invitationRepo,
				noopUserRepoForInv(),
				noopOrgRepoForInv(),
				noopGroupRepoForInv(),
				noopGroupServiceForInv(),
			)

			result, err := svc.GetByToken(context.Background(), tt.token)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				if tt.wantErrType != nil && !errors.Is(err, tt.wantErrType) {
					t.Errorf("expected error type %v, got %v", tt.wantErrType, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result == nil {
					t.Error("expected result, got nil")
				}
			}
		})
	}
}

func TestGetByToken_ContextPropagation(t *testing.T) {
	type ctxKey string
	key := ctxKey("test")
	ctx := context.WithValue(context.Background(), key, "value")

	var capturedCtx context.Context

	invitationRepo := &mockInvitationRepo{
		findByTokenFn: func(ctx context.Context, token string) (*models.Invitation, error) {
			capturedCtx = ctx
			return &models.Invitation{ID: "inv-1", Token: "token-1"}, nil
		},
	}

	svc := NewInvitationService(
		invitationRepo,
		noopUserRepoForInv(),
		noopOrgRepoForInv(),
		noopGroupRepoForInv(),
		noopGroupServiceForInv(),
	)

	_, _ = svc.GetByToken(ctx, "token-1")

	if capturedCtx.Value(key) != "value" {
		t.Error("context was not properly propagated")
	}
}
