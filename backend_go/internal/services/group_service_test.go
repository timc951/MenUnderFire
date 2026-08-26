package services

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"menunderfire/internal/models"
)

// Helper to create string pointer
func strPtrGrp(s string) *string {
	return &s
}

// mockGroupRepoForGroup is a mock implementation of repositories.GroupRepository for group service tests
type mockGroupRepoForGroup struct {
	findByIDFn               func(ctx context.Context, groupID string) (*models.Group, error)
	findByUserIDFn           func(ctx context.Context, userID string) ([]*models.Group, error)
	findByOrganizationIDFn   func(ctx context.Context, orgID string) ([]*models.Group, error)
	createFn                 func(ctx context.Context, name string, description *string, orgID string, inviteCode string, createdBy string) (*models.Group, error)
	generateInviteCodeFn     func() string
	findMemberFn             func(ctx context.Context, groupID, userID string) (*models.GroupMember, error)
	findMembersFn            func(ctx context.Context, groupID string) ([]*models.GroupMember, error)
	countMembersFn           func(ctx context.Context, groupID string) (int, error)
	addMemberFn              func(ctx context.Context, groupID, userID, role string) (*models.GroupMember, error)
	removeMemberFn           func(ctx context.Context, groupID, userID string) error
	countFn                  func(ctx context.Context) (int64, error)
	countByOrganizationIDsFn func(ctx context.Context, orgIDs []string) (int64, error)
}

func (m *mockGroupRepoForGroup) FindByID(ctx context.Context, groupID string) (*models.Group, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, groupID)
	}
	return nil, errors.New("findByIDFn not set")
}

func (m *mockGroupRepoForGroup) FindByInviteCode(ctx context.Context, inviteCode string) (*models.Group, error) {
	return nil, errors.New("findByInviteCodeFn not set")
}

func (m *mockGroupRepoForGroup) FindByUserID(ctx context.Context, userID string) ([]*models.Group, error) {
	if m.findByUserIDFn != nil {
		return m.findByUserIDFn(ctx, userID)
	}
	return nil, errors.New("findByUserIDFn not set")
}

func (m *mockGroupRepoForGroup) FindByOrganizationID(ctx context.Context, orgID string) ([]*models.Group, error) {
	if m.findByOrganizationIDFn != nil {
		return m.findByOrganizationIDFn(ctx, orgID)
	}
	return nil, errors.New("findByOrganizationIDFn not set")
}

func (m *mockGroupRepoForGroup) Create(ctx context.Context, name string, description *string, orgID string, inviteCode string, createdBy string) (*models.Group, error) {
	if m.createFn != nil {
		return m.createFn(ctx, name, description, orgID, inviteCode, createdBy)
	}
	return nil, errors.New("createFn not set")
}

func (m *mockGroupRepoForGroup) GenerateInviteCode() string {
	if m.generateInviteCodeFn != nil {
		return m.generateInviteCodeFn()
	}
	return "INVITE123"
}

func (m *mockGroupRepoForGroup) FindMember(ctx context.Context, groupID, userID string) (*models.GroupMember, error) {
	if m.findMemberFn != nil {
		return m.findMemberFn(ctx, groupID, userID)
	}
	return nil, errors.New("findMemberFn not set")
}

func (m *mockGroupRepoForGroup) FindMembers(ctx context.Context, groupID string) ([]*models.GroupMember, error) {
	if m.findMembersFn != nil {
		return m.findMembersFn(ctx, groupID)
	}
	return nil, errors.New("findMembersFn not set")
}

func (m *mockGroupRepoForGroup) CountMembers(ctx context.Context, groupID string) (int, error) {
	if m.countMembersFn != nil {
		return m.countMembersFn(ctx, groupID)
	}
	return 0, errors.New("countMembersFn not set")
}

func (m *mockGroupRepoForGroup) AddMember(ctx context.Context, groupID, userID, role string) (*models.GroupMember, error) {
	if m.addMemberFn != nil {
		return m.addMemberFn(ctx, groupID, userID, role)
	}
	return nil, errors.New("addMemberFn not set")
}

func (m *mockGroupRepoForGroup) RemoveMember(ctx context.Context, groupID, userID string) error {
	if m.removeMemberFn != nil {
		return m.removeMemberFn(ctx, groupID, userID)
	}
	return errors.New("removeMemberFn not set")
}

func (m *mockGroupRepoForGroup) Count(ctx context.Context) (int64, error) {
	if m.countFn != nil {
		return m.countFn(ctx)
	}
	return 0, errors.New("countFn not set")
}

func (m *mockGroupRepoForGroup) CountByOrganizationIDs(ctx context.Context, orgIDs []string) (int64, error) {
	if m.countByOrganizationIDsFn != nil {
		return m.countByOrganizationIDsFn(ctx, orgIDs)
	}
	return 0, errors.New("countByOrganizationIDsFn not set")
}

func (m *mockGroupRepoForGroup) UpdateMemberRole(ctx context.Context, groupID, userID, role string) error {
	return nil
}

func (m *mockGroupRepoForGroup) UpdateSettings(ctx context.Context, groupID string, requirePostApproval, allowAnonymousPosts bool) error {
	return nil
}

// mockOrgRepoForGroup is a mock implementation of repositories.OrganizationRepository for group service tests
type mockOrgRepoForGroup struct {
	findByIDFn     func(ctx context.Context, orgID string) (*models.Organization, error)
	findByUserIDFn func(ctx context.Context, userID string) ([]*models.Organization, error)
	findAllFn      func(ctx context.Context) ([]*models.Organization, error)
	createFn       func(ctx context.Context, name string, description *string, createdByID string) (*models.Organization, error)
	updateFn       func(ctx context.Context, orgID string, name string, description *string) (*models.Organization, error)
	findAdminsFn   func(ctx context.Context, orgID string) ([]*models.OrganizationAdmin, error)
	findAdminFn    func(ctx context.Context, orgID, userID string) (*models.OrganizationAdmin, error)
	addAdminFn     func(ctx context.Context, orgID, userID string) error
	isMemberFn     func(ctx context.Context, orgID, userID string) (bool, error)
	isAdminFn      func(ctx context.Context, orgID, userID string) (bool, error)
	countFn        func(ctx context.Context) (int64, error)
}

func (m *mockOrgRepoForGroup) FindByID(ctx context.Context, orgID string) (*models.Organization, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, orgID)
	}
	return nil, errors.New("findByIDFn not set")
}

func (m *mockOrgRepoForGroup) FindByUserID(ctx context.Context, userID string) ([]*models.Organization, error) {
	if m.findByUserIDFn != nil {
		return m.findByUserIDFn(ctx, userID)
	}
	return nil, errors.New("findByUserIDFn not set")
}

func (m *mockOrgRepoForGroup) FindAll(ctx context.Context) ([]*models.Organization, error) {
	if m.findAllFn != nil {
		return m.findAllFn(ctx)
	}
	return nil, errors.New("findAllFn not set")
}

func (m *mockOrgRepoForGroup) Create(ctx context.Context, name string, description *string, createdByID string) (*models.Organization, error) {
	if m.createFn != nil {
		return m.createFn(ctx, name, description, createdByID)
	}
	return nil, errors.New("createFn not set")
}

func (m *mockOrgRepoForGroup) Update(ctx context.Context, orgID string, name string, description *string) (*models.Organization, error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, orgID, name, description)
	}
	return nil, errors.New("updateFn not set")
}

func (m *mockOrgRepoForGroup) FindAdmins(ctx context.Context, orgID string) ([]*models.OrganizationAdmin, error) {
	if m.findAdminsFn != nil {
		return m.findAdminsFn(ctx, orgID)
	}
	return nil, errors.New("findAdminsFn not set")
}

func (m *mockOrgRepoForGroup) FindAdmin(ctx context.Context, orgID, userID string) (*models.OrganizationAdmin, error) {
	if m.findAdminFn != nil {
		return m.findAdminFn(ctx, orgID, userID)
	}
	return nil, errors.New("findAdminFn not set")
}

func (m *mockOrgRepoForGroup) AddAdmin(ctx context.Context, orgID, userID string) error {
	if m.addAdminFn != nil {
		return m.addAdminFn(ctx, orgID, userID)
	}
	return errors.New("addAdminFn not set")
}

func (m *mockOrgRepoForGroup) IsMember(ctx context.Context, orgID, userID string) (bool, error) {
	if m.isMemberFn != nil {
		return m.isMemberFn(ctx, orgID, userID)
	}
	return false, errors.New("isMemberFn not set")
}

func (m *mockOrgRepoForGroup) IsAdmin(ctx context.Context, orgID, userID string) (bool, error) {
	if m.isAdminFn != nil {
		return m.isAdminFn(ctx, orgID, userID)
	}
	return false, errors.New("isAdminFn not set")
}

func (m *mockOrgRepoForGroup) Count(ctx context.Context) (int64, error) {
	if m.countFn != nil {
		return m.countFn(ctx)
	}
	return 0, errors.New("countFn not set")
}

// mockUserRepoForGroup is a mock implementation of repositories.UserRepository for group service tests
type mockUserRepoForGroup struct {
	findByIDFn                 func(ctx context.Context, userID string) (*models.User, error)
	findByExternalIDFn         func(ctx context.Context, externalID string) (*models.User, error)
	findByEmailFn              func(ctx context.Context, email string) (*models.User, error)
	countFn                    func(ctx context.Context) (int64, error)
	createFn                   func(ctx context.Context, email, displayName, externalID string) (*models.User, error)
	createAsSiteAdminFn        func(ctx context.Context, email, displayName, externalID string) (*models.User, error)
	updateFn                   func(ctx context.Context, userID, displayName string) (*models.User, error)
	isSiteAdminFn              func(ctx context.Context, userID string) (bool, error)
	findAdminOrganizationIDsFn func(ctx context.Context, userID string) ([]string, error)
	findOwnedGroupIDsFn        func(ctx context.Context, userID string) ([]string, error)
	findMemberGroupIDsFn       func(ctx context.Context, userID string) ([]string, error)
}

func (m *mockUserRepoForGroup) FindByID(ctx context.Context, userID string) (*models.User, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, userID)
	}
	return nil, errors.New("findByIDFn not set")
}

func (m *mockUserRepoForGroup) FindByExternalID(ctx context.Context, externalID string) (*models.User, error) {
	if m.findByExternalIDFn != nil {
		return m.findByExternalIDFn(ctx, externalID)
	}
	return nil, errors.New("findByExternalIDFn not set")
}

func (m *mockUserRepoForGroup) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	if m.findByEmailFn != nil {
		return m.findByEmailFn(ctx, email)
	}
	return nil, errors.New("findByEmailFn not set")
}

func (m *mockUserRepoForGroup) Count(ctx context.Context) (int64, error) {
	if m.countFn != nil {
		return m.countFn(ctx)
	}
	return 0, errors.New("countFn not set")
}

func (m *mockUserRepoForGroup) Create(ctx context.Context, email, displayName, externalID string) (*models.User, error) {
	if m.createFn != nil {
		return m.createFn(ctx, email, displayName, externalID)
	}
	return nil, errors.New("createFn not set")
}

func (m *mockUserRepoForGroup) CreateAsSiteAdmin(ctx context.Context, email, displayName, externalID string) (*models.User, error) {
	if m.createAsSiteAdminFn != nil {
		return m.createAsSiteAdminFn(ctx, email, displayName, externalID)
	}
	return nil, errors.New("createAsSiteAdminFn not set")
}

func (m *mockUserRepoForGroup) Update(ctx context.Context, userID, displayName string) (*models.User, error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, userID, displayName)
	}
	return nil, errors.New("updateFn not set")
}

func (m *mockUserRepoForGroup) UpdateExternalID(ctx context.Context, userID, externalID string) error {
	return errors.New("updateExternalIDFn not set")
}

func (m *mockUserRepoForGroup) UpdateInvitationInfo(ctx context.Context, userID, invitedByID, invitationID string) error {
	return errors.New("updateInvitationInfoFn not set")
}

func (m *mockUserRepoForGroup) IsSiteAdmin(ctx context.Context, userID string) (bool, error) {
	if m.isSiteAdminFn != nil {
		return m.isSiteAdminFn(ctx, userID)
	}
	return false, errors.New("isSiteAdminFn not set")
}

func (m *mockUserRepoForGroup) FindAdminOrganizationIDs(ctx context.Context, userID string) ([]string, error) {
	if m.findAdminOrganizationIDsFn != nil {
		return m.findAdminOrganizationIDsFn(ctx, userID)
	}
	return nil, errors.New("findAdminOrganizationIDsFn not set")
}

func (m *mockUserRepoForGroup) FindOwnedGroupIDs(ctx context.Context, userID string) ([]string, error) {
	if m.findOwnedGroupIDsFn != nil {
		return m.findOwnedGroupIDsFn(ctx, userID)
	}
	return nil, errors.New("findOwnedGroupIDsFn not set")
}

func (m *mockUserRepoForGroup) FindMemberGroupIDs(ctx context.Context, userID string) ([]string, error) {
	if m.findMemberGroupIDsFn != nil {
		return m.findMemberGroupIDsFn(ctx, userID)
	}
	return nil, errors.New("findMemberGroupIDsFn not set")
}

func (m *mockUserRepoForGroup) RecordAgreementAcceptance(ctx context.Context, userID, version, signature, ipAddress, userAgent string) error {
	return errors.New("not implemented")
}

// ========== Create Tests ==========

func TestGroupServiceImpl_Create(t *testing.T) {
	tests := []struct {
		name              string
		userID            string
		req               *models.CreateGroupRequest
		findOrgResult     *models.Organization
		findOrgErr        error
		isOrgAdminResult  bool
		isOrgAdminErr     error
		createGroupResult *models.Group
		createGroupErr    error
		addMemberResult   *models.GroupMember
		addMemberErr      error
		want              *models.Group
		wantErr           error
		wantErrContains   string
	}{
		{
			name:    "empty name returns ErrValidation",
			userID:  "user-1",
			req:     &models.CreateGroupRequest{Name: "", OrganizationID: "org-1"},
			wantErr: ErrValidation,
		},
		{
			name:    "whitespace only name returns ErrValidation",
			userID:  "user-1",
			req:     &models.CreateGroupRequest{Name: "   ", OrganizationID: "org-1"},
			wantErr: ErrValidation,
		},
		{
			name:    "empty organizationID returns ErrValidation",
			userID:  "user-1",
			req:     &models.CreateGroupRequest{Name: "Test Group", OrganizationID: ""},
			wantErr: ErrValidation,
		},
		{
			name:            "find org error returns wrapped error",
			userID:          "user-1",
			req:             &models.CreateGroupRequest{Name: "Test Group", OrganizationID: "org-1"},
			findOrgErr:      errors.New("db error"),
			wantErrContains: "failed to find organization",
		},
		{
			name:          "organization not found returns ErrOrganizationNotFound",
			userID:        "user-1",
			req:           &models.CreateGroupRequest{Name: "Test Group", OrganizationID: "org-1"},
			findOrgResult: nil,
			wantErr:       ErrOrganizationNotFound,
		},
		{
			name:            "org admin check error returns wrapped error",
			userID:          "user-1",
			req:             &models.CreateGroupRequest{Name: "Test Group", OrganizationID: "org-1"},
			findOrgResult:   &models.Organization{ID: "org-1"},
			isOrgAdminErr:   errors.New("db error"),
			wantErrContains: "failed to check org admin status",
		},
		{
			name:             "non-org-admin returns ErrForbidden",
			userID:           "user-1",
			req:              &models.CreateGroupRequest{Name: "Test Group", OrganizationID: "org-1"},
			findOrgResult:    &models.Organization{ID: "org-1"},
			isOrgAdminResult: false,
			wantErr:          ErrForbidden,
		},
		{
			name:             "create group error returns wrapped error",
			userID:           "user-1",
			req:              &models.CreateGroupRequest{Name: "Test Group", OrganizationID: "org-1"},
			findOrgResult:    &models.Organization{ID: "org-1"},
			isOrgAdminResult: true,
			createGroupErr:   errors.New("db error"),
			wantErrContains:  "failed to create group",
		},
		{
			name:              "add member error returns wrapped error",
			userID:            "user-1",
			req:               &models.CreateGroupRequest{Name: "Test Group", OrganizationID: "org-1"},
			findOrgResult:     &models.Organization{ID: "org-1"},
			isOrgAdminResult:  true,
			createGroupResult: &models.Group{ID: "group-1", Name: "Test Group"},
			addMemberErr:      errors.New("db error"),
			wantErrContains:   "failed to add owner to group",
		},
		{
			name:              "successful creation returns group",
			userID:            "user-1",
			req:               &models.CreateGroupRequest{Name: "  Test Group  ", OrganizationID: "org-1", Description: strPtrGrp("  Desc  ")},
			findOrgResult:     &models.Organization{ID: "org-1"},
			isOrgAdminResult:  true,
			createGroupResult: &models.Group{ID: "group-1", Name: "Test Group", OrganizationID: "org-1"},
			addMemberResult:   &models.GroupMember{ID: "member-1", GroupID: "group-1", UserID: "user-1", Role: "OWNER"},
			want:              &models.Group{ID: "group-1", Name: "Test Group", OrganizationID: "org-1"},
		},
		{
			name:              "nil description is passed through",
			userID:            "user-1",
			req:               &models.CreateGroupRequest{Name: "Test Group", OrganizationID: "org-1", Description: nil},
			findOrgResult:     &models.Organization{ID: "org-1"},
			isOrgAdminResult:  true,
			createGroupResult: &models.Group{ID: "group-1", Name: "Test Group"},
			addMemberResult:   &models.GroupMember{ID: "member-1"},
			want:              &models.Group{ID: "group-1", Name: "Test Group"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			groupRepo := &mockGroupRepoForGroup{
				createFn: func(ctx context.Context, name string, description *string, orgID string, inviteCode string, createdBy string) (*models.Group, error) {
					if tt.createGroupErr != nil {
						return nil, tt.createGroupErr
					}
					return tt.createGroupResult, nil
				},
				generateInviteCodeFn: func() string {
					return "INVITE123"
				},
				addMemberFn: func(ctx context.Context, groupID, userID, role string) (*models.GroupMember, error) {
					if tt.addMemberErr != nil {
						return nil, tt.addMemberErr
					}
					return tt.addMemberResult, nil
				},
			}
			orgRepo := &mockOrgRepoForGroup{
				findByIDFn: func(ctx context.Context, orgID string) (*models.Organization, error) {
					if tt.findOrgErr != nil {
						return nil, tt.findOrgErr
					}
					return tt.findOrgResult, nil
				},
				isAdminFn: func(ctx context.Context, orgID, userID string) (bool, error) {
					if tt.isOrgAdminErr != nil {
						return false, tt.isOrgAdminErr
					}
					return tt.isOrgAdminResult, nil
				},
			}
			userRepo := &mockUserRepoForGroup{}

			svc := NewGroupService(groupRepo, orgRepo, userRepo)

			got, err := svc.Create(context.Background(), tt.userID, tt.req)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if tt.wantErrContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("Create() error = %v, wantErrContains %v", err, tt.wantErrContains)
				}
				return
			}
			if err != nil {
				t.Errorf("Create() unexpected error = %v", err)
				return
			}
			if got.ID != tt.want.ID {
				t.Errorf("Create() got ID = %v, want %v", got.ID, tt.want.ID)
			}
		})
	}
}

// Test that invite code is generated and passed to createGroup
func TestGroupServiceImpl_Create_InviteCodeGeneration(t *testing.T) {
	var capturedInviteCode string

	groupRepo := &mockGroupRepoForGroup{
		createFn: func(ctx context.Context, name string, description *string, orgID string, inviteCode string, createdBy string) (*models.Group, error) {
			capturedInviteCode = inviteCode
			return &models.Group{ID: "group-1", InviteCode: inviteCode}, nil
		},
		generateInviteCodeFn: func() string {
			return "CUSTOM-CODE-123"
		},
		addMemberFn: func(ctx context.Context, groupID, userID, role string) (*models.GroupMember, error) {
			return &models.GroupMember{}, nil
		},
	}
	orgRepo := &mockOrgRepoForGroup{
		findByIDFn: func(ctx context.Context, orgID string) (*models.Organization, error) {
			return &models.Organization{ID: "org-1"}, nil
		},
		isAdminFn: func(ctx context.Context, orgID, userID string) (bool, error) {
			return true, nil
		},
	}
	userRepo := &mockUserRepoForGroup{}

	svc := NewGroupService(groupRepo, orgRepo, userRepo)

	_, err := svc.Create(context.Background(), "user-1", &models.CreateGroupRequest{
		Name:           "Test Group",
		OrganizationID: "org-1",
	})
	if err != nil {
		t.Fatalf("Create() unexpected error = %v", err)
	}

	if capturedInviteCode != "CUSTOM-CODE-123" {
		t.Errorf("Create() invite code = %q, want %q", capturedInviteCode, "CUSTOM-CODE-123")
	}
}

// ========== List Tests ==========

func TestGroupServiceImpl_List(t *testing.T) {
	tests := []struct {
		name             string
		userID           string
		findGroupsResult []*models.Group
		findGroupsErr    error
		want             []*models.Group
		wantErr          error
		wantErrContains  string
	}{
		{
			name:            "find groups error returns wrapped error",
			userID:          "user-1",
			findGroupsErr:   errors.New("db error"),
			wantErrContains: "failed to find groups for user",
		},
		{
			name:             "nil result returns empty slice",
			userID:           "user-1",
			findGroupsResult: nil,
			want:             []*models.Group{},
		},
		{
			name:             "empty slice returns empty slice",
			userID:           "user-1",
			findGroupsResult: []*models.Group{},
			want:             []*models.Group{},
		},
		{
			name:             "user belongs to multiple groups returns all",
			userID:           "user-1",
			findGroupsResult: []*models.Group{{ID: "group-1"}, {ID: "group-2"}},
			want:             []*models.Group{{ID: "group-1"}, {ID: "group-2"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			groupRepo := &mockGroupRepoForGroup{
				findByUserIDFn: func(ctx context.Context, userID string) ([]*models.Group, error) {
					if tt.findGroupsErr != nil {
						return nil, tt.findGroupsErr
					}
					return tt.findGroupsResult, nil
				},
			}
			orgRepo := &mockOrgRepoForGroup{}
			userRepo := &mockUserRepoForGroup{
				// user is not an org admin, so List returns only member groups
				findAdminOrganizationIDsFn: func(ctx context.Context, userID string) ([]string, error) {
					return nil, nil
				},
			}

			svc := NewGroupService(groupRepo, orgRepo, userRepo)

			got, err := svc.List(context.Background(), tt.userID)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if tt.wantErrContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("List() error = %v, wantErrContains %v", err, tt.wantErrContains)
				}
				return
			}
			if err != nil {
				t.Errorf("List() unexpected error = %v", err)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("List() got %d groups, want %d", len(got), len(tt.want))
			}
		})
	}
}

// ========== GetByID Tests ==========

func TestGroupServiceImpl_GetByID(t *testing.T) {
	tests := []struct {
		name            string
		groupID         string
		findGroupResult *models.Group
		findGroupErr    error
		want            *models.Group
		wantErr         error
		wantErrContains string
	}{
		{
			name:    "empty groupID returns ErrGroupNotFound",
			groupID: "",
			wantErr: ErrGroupNotFound,
		},
		{
			name:            "find group error returns wrapped error",
			groupID:         "group-1",
			findGroupErr:    errors.New("db error"),
			wantErrContains: "failed to get group by ID",
		},
		{
			name:            "group not found returns ErrGroupNotFound",
			groupID:         "group-1",
			findGroupResult: nil,
			wantErr:         ErrGroupNotFound,
		},
		{
			name:            "group found returns the group",
			groupID:         "group-1",
			findGroupResult: &models.Group{ID: "group-1", Name: "Test Group"},
			want:            &models.Group{ID: "group-1", Name: "Test Group"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			groupRepo := &mockGroupRepoForGroup{
				findByIDFn: func(ctx context.Context, groupID string) (*models.Group, error) {
					if tt.findGroupErr != nil {
						return nil, tt.findGroupErr
					}
					return tt.findGroupResult, nil
				},
			}
			orgRepo := &mockOrgRepoForGroup{}
			userRepo := &mockUserRepoForGroup{}

			svc := NewGroupService(groupRepo, orgRepo, userRepo)

			got, err := svc.GetByID(context.Background(), tt.groupID)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if tt.wantErrContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("GetByID() error = %v, wantErrContains %v", err, tt.wantErrContains)
				}
				return
			}
			if err != nil {
				t.Errorf("GetByID() unexpected error = %v", err)
				return
			}
			if got.ID != tt.want.ID {
				t.Errorf("GetByID() got ID = %v, want %v", got.ID, tt.want.ID)
			}
		})
	}
}

// ========== GetDetailByID Tests ==========

func TestGroupServiceImpl_GetDetailByID(t *testing.T) {
	tests := []struct {
		name              string
		groupID           string
		userID            string
		findGroupResult   *models.Group
		findGroupErr      error
		findMemberResult  *models.GroupMember
		findMemberErr     error
		isSiteAdminResult bool
		isSiteAdminErr    error
		isOrgAdminResult  bool
		isOrgAdminErr     error
		findMembersResult []*models.GroupMember
		findMembersErr    error
		wantRole          string
		wantInviteCode    bool
		wantErr           error
		wantErrContains   string
	}{
		{
			name:            "find group error returns wrapped error",
			groupID:         "group-1",
			userID:          "user-1",
			findGroupErr:    errors.New("db error"),
			wantErrContains: "failed to find group",
		},
		{
			name:            "group not found returns ErrGroupNotFound",
			groupID:         "group-1",
			userID:          "user-1",
			findGroupResult: nil,
			wantErr:         ErrGroupNotFound,
		},
		{
			name:            "find member error returns wrapped error",
			groupID:         "group-1",
			userID:          "user-1",
			findGroupResult: &models.Group{ID: "group-1", OrganizationID: "org-1", InviteCode: "CODE123"},
			findMemberErr:   errors.New("db error"),
			wantErrContains: "failed to check membership",
		},
		{
			name:              "member with OWNER role sees invite code",
			groupID:           "group-1",
			userID:            "owner-1",
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1", InviteCode: "CODE123"},
			findMemberResult:  &models.GroupMember{UserID: "owner-1", Role: "OWNER"},
			findMembersResult: []*models.GroupMember{},
			wantRole:          "OWNER",
			wantInviteCode:    true,
		},
		{
			name:              "member with LEADER role sees invite code",
			groupID:           "group-1",
			userID:            "leader-1",
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1", InviteCode: "CODE123"},
			findMemberResult:  &models.GroupMember{UserID: "leader-1", Role: "LEADER"},
			findMembersResult: []*models.GroupMember{},
			wantRole:          "LEADER",
			wantInviteCode:    true,
		},
		{
			name:              "member with MEMBER role does not see invite code",
			groupID:           "group-1",
			userID:            "member-1",
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1", InviteCode: "CODE123"},
			findMemberResult:  &models.GroupMember{UserID: "member-1", Role: "MEMBER"},
			findMembersResult: []*models.GroupMember{},
			wantRole:          "MEMBER",
			wantInviteCode:    false,
		},
		{
			name:             "non-member - site admin check error returns wrapped error",
			groupID:          "group-1",
			userID:           "user-1",
			findGroupResult:  &models.Group{ID: "group-1", OrganizationID: "org-1"},
			findMemberResult: nil,
			isSiteAdminErr:   errors.New("db error"),
			wantErrContains:  "failed to check site admin status",
		},
		{
			name:              "site admin (non-member) gets ADMIN role and sees invite code",
			groupID:           "group-1",
			userID:            "site-admin",
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1", InviteCode: "CODE123"},
			findMemberResult:  nil,
			isSiteAdminResult: true,
			findMembersResult: []*models.GroupMember{},
			wantRole:          "ADMIN",
			wantInviteCode:    true,
		},
		{
			name:              "non-member non-site-admin - org admin check error returns wrapped error",
			groupID:           "group-1",
			userID:            "user-1",
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			findMemberResult:  nil,
			isSiteAdminResult: false,
			isOrgAdminErr:     errors.New("db error"),
			wantErrContains:   "failed to check org admin status",
		},
		{
			name:              "org admin (non-member) gets ADMIN role and sees invite code",
			groupID:           "group-1",
			userID:            "org-admin",
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1", InviteCode: "CODE123"},
			findMemberResult:  nil,
			isSiteAdminResult: false,
			isOrgAdminResult:  true,
			findMembersResult: []*models.GroupMember{},
			wantRole:          "ADMIN",
			wantInviteCode:    true,
		},
		{
			name:              "non-member non-admin returns detail with empty role",
			groupID:           "group-1",
			userID:            "random-user",
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			findMemberResult:  nil,
			isSiteAdminResult: false,
			isOrgAdminResult:  false,
			wantRole:          "",
			wantInviteCode:    false,
		},
		{
			name:             "find members error returns wrapped error",
			groupID:          "group-1",
			userID:           "owner-1",
			findGroupResult:  &models.Group{ID: "group-1", OrganizationID: "org-1"},
			findMemberResult: &models.GroupMember{UserID: "owner-1", Role: "OWNER"},
			findMembersErr:   errors.New("db error"),
			wantErrContains:  "failed to find group members",
		},
		{
			name:              "empty members list returns detail with empty members",
			groupID:           "group-1",
			userID:            "owner-1",
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1", InviteCode: "CODE123"},
			findMemberResult:  &models.GroupMember{UserID: "owner-1", Role: "OWNER"},
			findMembersResult: []*models.GroupMember{},
			wantRole:          "OWNER",
			wantInviteCode:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			groupRepo := &mockGroupRepoForGroup{
				findByIDFn: func(ctx context.Context, groupID string) (*models.Group, error) {
					if tt.findGroupErr != nil {
						return nil, tt.findGroupErr
					}
					return tt.findGroupResult, nil
				},
				findMemberFn: func(ctx context.Context, groupID, userID string) (*models.GroupMember, error) {
					if tt.findMemberErr != nil {
						return nil, tt.findMemberErr
					}
					return tt.findMemberResult, nil
				},
				findMembersFn: func(ctx context.Context, groupID string) ([]*models.GroupMember, error) {
					if tt.findMembersErr != nil {
						return nil, tt.findMembersErr
					}
					return tt.findMembersResult, nil
				},
			}
			orgRepo := &mockOrgRepoForGroup{
				isAdminFn: func(ctx context.Context, orgID, userID string) (bool, error) {
					if tt.isOrgAdminErr != nil {
						return false, tt.isOrgAdminErr
					}
					return tt.isOrgAdminResult, nil
				},
			}
			userRepo := &mockUserRepoForGroup{
				isSiteAdminFn: func(ctx context.Context, userID string) (bool, error) {
					if tt.isSiteAdminErr != nil {
						return false, tt.isSiteAdminErr
					}
					return tt.isSiteAdminResult, nil
				},
			}

			svc := NewGroupService(groupRepo, orgRepo, userRepo)

			got, err := svc.GetDetailByID(context.Background(), tt.groupID, tt.userID)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("GetDetailByID() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if tt.wantErrContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("GetDetailByID() error = %v, wantErrContains %v", err, tt.wantErrContains)
				}
				return
			}
			if err != nil {
				t.Errorf("GetDetailByID() unexpected error = %v", err)
				return
			}
			if got.Role != tt.wantRole {
				t.Errorf("GetDetailByID() role = %v, want %v", got.Role, tt.wantRole)
			}
			hasInviteCode := got.InviteCode != nil
			if hasInviteCode != tt.wantInviteCode {
				t.Errorf("GetDetailByID() hasInviteCode = %v, want %v", hasInviteCode, tt.wantInviteCode)
			}
		})
	}
}

// ========== Join Tests ==========

func TestGroupServiceImpl_Join(t *testing.T) {
	tests := []struct {
		name             string
		groupID          string
		userID           string
		inviteCode       string
		findGroupResult  *models.Group
		findGroupErr     error
		findMemberResult *models.GroupMember
		findMemberErr    error
		addMemberErr     error
		wantErr          error
		wantErrContains  string
	}{
		{
			name:            "find group error returns wrapped error",
			groupID:         "group-1",
			userID:          "user-1",
			inviteCode:      "CODE123",
			findGroupErr:    errors.New("db error"),
			wantErrContains: "failed to find group",
		},
		{
			name:            "group not found returns ErrGroupNotFound",
			groupID:         "group-1",
			userID:          "user-1",
			inviteCode:      "CODE123",
			findGroupResult: nil,
			wantErr:         ErrGroupNotFound,
		},
		{
			name:            "empty invite code returns ErrInvalidInviteCode",
			groupID:         "group-1",
			userID:          "user-1",
			inviteCode:      "",
			findGroupResult: &models.Group{ID: "group-1", InviteCode: "CODE123"},
			wantErr:         ErrInvalidInviteCode,
		},
		{
			name:            "invalid invite code returns ErrInvalidInviteCode",
			groupID:         "group-1",
			userID:          "user-1",
			inviteCode:      "WRONG-CODE",
			findGroupResult: &models.Group{ID: "group-1", InviteCode: "CODE123"},
			wantErr:         ErrInvalidInviteCode,
		},
		{
			name:            "find member error returns wrapped error",
			groupID:         "group-1",
			userID:          "user-1",
			inviteCode:      "CODE123",
			findGroupResult: &models.Group{ID: "group-1", InviteCode: "CODE123"},
			findMemberErr:   errors.New("db error"),
			wantErrContains: "failed to check existing membership",
		},
		{
			name:             "already a member returns ErrConflict",
			groupID:          "group-1",
			userID:           "user-1",
			inviteCode:       "CODE123",
			findGroupResult:  &models.Group{ID: "group-1", InviteCode: "CODE123"},
			findMemberResult: &models.GroupMember{UserID: "user-1"},
			wantErr:          ErrConflict,
		},
		{
			name:             "add member error returns wrapped error",
			groupID:          "group-1",
			userID:           "user-1",
			inviteCode:       "CODE123",
			findGroupResult:  &models.Group{ID: "group-1", InviteCode: "CODE123"},
			findMemberResult: nil,
			addMemberErr:     errors.New("db error"),
			wantErrContains:  "failed to add member",
		},
		{
			name:             "successful join returns nil",
			groupID:          "group-1",
			userID:           "user-1",
			inviteCode:       "CODE123",
			findGroupResult:  &models.Group{ID: "group-1", InviteCode: "CODE123"},
			findMemberResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			groupRepo := &mockGroupRepoForGroup{
				findByIDFn: func(ctx context.Context, groupID string) (*models.Group, error) {
					if tt.findGroupErr != nil {
						return nil, tt.findGroupErr
					}
					return tt.findGroupResult, nil
				},
				findMemberFn: func(ctx context.Context, groupID, userID string) (*models.GroupMember, error) {
					if tt.findMemberErr != nil {
						return nil, tt.findMemberErr
					}
					return tt.findMemberResult, nil
				},
				addMemberFn: func(ctx context.Context, groupID, userID, role string) (*models.GroupMember, error) {
					if tt.addMemberErr != nil {
						return nil, tt.addMemberErr
					}
					return &models.GroupMember{}, nil
				},
			}
			orgRepo := &mockOrgRepoForGroup{}
			userRepo := &mockUserRepoForGroup{}

			svc := NewGroupService(groupRepo, orgRepo, userRepo)

			err := svc.Join(context.Background(), tt.groupID, tt.userID, tt.inviteCode)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Join() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if tt.wantErrContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("Join() error = %v, wantErrContains %v", err, tt.wantErrContains)
				}
				return
			}
			if err != nil {
				t.Errorf("Join() unexpected error = %v", err)
			}
		})
	}
}

// Test that join adds user with MEMBER role
func TestGroupServiceImpl_Join_AddsAsMember(t *testing.T) {
	var capturedRole string

	groupRepo := &mockGroupRepoForGroup{
		findByIDFn: func(ctx context.Context, groupID string) (*models.Group, error) {
			return &models.Group{ID: "group-1", InviteCode: "CODE123"}, nil
		},
		findMemberFn: func(ctx context.Context, groupID, userID string) (*models.GroupMember, error) {
			return nil, nil
		},
		addMemberFn: func(ctx context.Context, groupID, userID, role string) (*models.GroupMember, error) {
			capturedRole = role
			return &models.GroupMember{}, nil
		},
	}
	orgRepo := &mockOrgRepoForGroup{}
	userRepo := &mockUserRepoForGroup{}

	svc := NewGroupService(groupRepo, orgRepo, userRepo)

	err := svc.Join(context.Background(), "group-1", "user-1", "CODE123")
	if err != nil {
		t.Fatalf("Join() unexpected error = %v", err)
	}

	if capturedRole != "MEMBER" {
		t.Errorf("Join() role = %q, want %q", capturedRole, "MEMBER")
	}
}

// ========== RemoveMember Tests ==========

func TestGroupServiceImpl_RemoveMember(t *testing.T) {
	tests := []struct {
		name                  string
		groupID               string
		targetUserID          string
		requestingUserID      string
		findGroupResult       *models.Group
		findGroupErr          error
		targetMemberResult    *models.GroupMember
		targetMemberErr       error
		findMembersResult     []*models.GroupMember
		findMembersErr        error
		requesterMemberResult *models.GroupMember
		requesterMemberErr    error
		isSiteAdminResult     bool
		isSiteAdminErr        error
		isOrgAdminResult      bool
		isOrgAdminErr         error
		removeMemberErr       error
		wantErr               error
		wantErrContains       string
	}{
		{
			name:             "find group error returns wrapped error",
			groupID:          "group-1",
			targetUserID:     "target-1",
			requestingUserID: "requester-1",
			findGroupErr:     errors.New("db error"),
			wantErrContains:  "failed to find group",
		},
		{
			name:             "group not found returns ErrGroupNotFound",
			groupID:          "group-1",
			targetUserID:     "target-1",
			requestingUserID: "requester-1",
			findGroupResult:  nil,
			wantErr:          ErrGroupNotFound,
		},
		{
			name:             "target member check error returns wrapped error",
			groupID:          "group-1",
			targetUserID:     "target-1",
			requestingUserID: "requester-1",
			findGroupResult:  &models.Group{ID: "group-1", OrganizationID: "org-1"},
			targetMemberErr:  errors.New("db error"),
			wantErrContains:  "failed to check target membership",
		},
		{
			name:               "target not a member returns ErrNotFound",
			groupID:            "group-1",
			targetUserID:       "target-1",
			requestingUserID:   "requester-1",
			findGroupResult:    &models.Group{ID: "group-1", OrganizationID: "org-1"},
			targetMemberResult: nil,
			wantErr:            ErrNotFound,
		},
		{
			name:               "removing last owner returns ErrValidation",
			groupID:            "group-1",
			targetUserID:       "owner-1",
			requestingUserID:   "owner-1",
			findGroupResult:    &models.Group{ID: "group-1", OrganizationID: "org-1"},
			targetMemberResult: &models.GroupMember{UserID: "owner-1", Role: "OWNER"},
			findMembersResult:  []*models.GroupMember{{UserID: "owner-1", Role: "OWNER"}},
			wantErr:            ErrValidation,
		},
		{
			name:               "removing owner when multiple owners exist is allowed",
			groupID:            "group-1",
			targetUserID:       "owner-1",
			requestingUserID:   "owner-1",
			findGroupResult:    &models.Group{ID: "group-1", OrganizationID: "org-1"},
			targetMemberResult: &models.GroupMember{UserID: "owner-1", Role: "OWNER"},
			findMembersResult:  []*models.GroupMember{{UserID: "owner-1", Role: "OWNER"}, {UserID: "owner-2", Role: "OWNER"}},
		},
		{
			name:               "find members error when removing owner returns wrapped error",
			groupID:            "group-1",
			targetUserID:       "owner-1",
			requestingUserID:   "requester-1",
			findGroupResult:    &models.Group{ID: "group-1", OrganizationID: "org-1"},
			targetMemberResult: &models.GroupMember{UserID: "owner-1", Role: "OWNER"},
			findMembersErr:     errors.New("db error"),
			wantErrContains:    "failed to find group members",
		},
		{
			name:               "user can remove themselves",
			groupID:            "group-1",
			targetUserID:       "user-1",
			requestingUserID:   "user-1",
			findGroupResult:    &models.Group{ID: "group-1", OrganizationID: "org-1"},
			targetMemberResult: &models.GroupMember{UserID: "user-1", Role: "MEMBER"},
		},
		{
			name:               "self-removal error returns wrapped error",
			groupID:            "group-1",
			targetUserID:       "user-1",
			requestingUserID:   "user-1",
			findGroupResult:    &models.Group{ID: "group-1", OrganizationID: "org-1"},
			targetMemberResult: &models.GroupMember{UserID: "user-1", Role: "MEMBER"},
			removeMemberErr:    errors.New("db error"),
			wantErrContains:    "failed to remove member",
		},
		{
			name:               "requester membership check error returns wrapped error",
			groupID:            "group-1",
			targetUserID:       "target-1",
			requestingUserID:   "requester-1",
			findGroupResult:    &models.Group{ID: "group-1", OrganizationID: "org-1"},
			targetMemberResult: &models.GroupMember{UserID: "target-1", Role: "MEMBER"},
			requesterMemberErr: errors.New("db error"),
			wantErrContains:    "failed to check requester membership",
		},
		{
			name:                  "owner can remove any member",
			groupID:               "group-1",
			targetUserID:          "target-1",
			requestingUserID:      "owner-1",
			findGroupResult:       &models.Group{ID: "group-1", OrganizationID: "org-1"},
			targetMemberResult:    &models.GroupMember{UserID: "target-1", Role: "MEMBER"},
			requesterMemberResult: &models.GroupMember{UserID: "owner-1", Role: "OWNER"},
		},
		{
			name:                  "owner removal error returns wrapped error",
			groupID:               "group-1",
			targetUserID:          "target-1",
			requestingUserID:      "owner-1",
			findGroupResult:       &models.Group{ID: "group-1", OrganizationID: "org-1"},
			targetMemberResult:    &models.GroupMember{UserID: "target-1", Role: "MEMBER"},
			requesterMemberResult: &models.GroupMember{UserID: "owner-1", Role: "OWNER"},
			removeMemberErr:       errors.New("db error"),
			wantErrContains:       "failed to remove member",
		},
		{
			name:                  "leader can remove MEMBER",
			groupID:               "group-1",
			targetUserID:          "target-1",
			requestingUserID:      "leader-1",
			findGroupResult:       &models.Group{ID: "group-1", OrganizationID: "org-1"},
			targetMemberResult:    &models.GroupMember{UserID: "target-1", Role: "MEMBER"},
			requesterMemberResult: &models.GroupMember{UserID: "leader-1", Role: "LEADER"},
		},
		{
			name:                  "leader removal error returns wrapped error",
			groupID:               "group-1",
			targetUserID:          "target-1",
			requestingUserID:      "leader-1",
			findGroupResult:       &models.Group{ID: "group-1", OrganizationID: "org-1"},
			targetMemberResult:    &models.GroupMember{UserID: "target-1", Role: "MEMBER"},
			requesterMemberResult: &models.GroupMember{UserID: "leader-1", Role: "LEADER"},
			removeMemberErr:       errors.New("db error"),
			wantErrContains:       "failed to remove member",
		},
		{
			name:                  "leader cannot remove another LEADER - check site admin error",
			groupID:               "group-1",
			targetUserID:          "leader-2",
			requestingUserID:      "leader-1",
			findGroupResult:       &models.Group{ID: "group-1", OrganizationID: "org-1"},
			targetMemberResult:    &models.GroupMember{UserID: "leader-2", Role: "LEADER"},
			requesterMemberResult: &models.GroupMember{UserID: "leader-1", Role: "LEADER"},
			isSiteAdminErr:        errors.New("db error"),
			wantErrContains:       "failed to check site admin status",
		},
		{
			name:                  "site admin can remove any member",
			groupID:               "group-1",
			targetUserID:          "target-1",
			requestingUserID:      "site-admin",
			findGroupResult:       &models.Group{ID: "group-1", OrganizationID: "org-1"},
			targetMemberResult:    &models.GroupMember{UserID: "target-1", Role: "MEMBER"},
			requesterMemberResult: nil,
			isSiteAdminResult:     true,
		},
		{
			name:                  "site admin removal error returns wrapped error",
			groupID:               "group-1",
			targetUserID:          "target-1",
			requestingUserID:      "site-admin",
			findGroupResult:       &models.Group{ID: "group-1", OrganizationID: "org-1"},
			targetMemberResult:    &models.GroupMember{UserID: "target-1", Role: "MEMBER"},
			requesterMemberResult: nil,
			isSiteAdminResult:     true,
			removeMemberErr:       errors.New("db error"),
			wantErrContains:       "failed to remove member",
		},
		{
			name:                  "non-site-admin - org admin check error returns wrapped error",
			groupID:               "group-1",
			targetUserID:          "target-1",
			requestingUserID:      "user-1",
			findGroupResult:       &models.Group{ID: "group-1", OrganizationID: "org-1"},
			targetMemberResult:    &models.GroupMember{UserID: "target-1", Role: "MEMBER"},
			requesterMemberResult: &models.GroupMember{UserID: "user-1", Role: "MEMBER"},
			isSiteAdminResult:     false,
			isOrgAdminErr:         errors.New("db error"),
			wantErrContains:       "failed to check org admin status",
		},
		{
			name:                  "org admin can remove any member",
			groupID:               "group-1",
			targetUserID:          "target-1",
			requestingUserID:      "org-admin",
			findGroupResult:       &models.Group{ID: "group-1", OrganizationID: "org-1"},
			targetMemberResult:    &models.GroupMember{UserID: "target-1", Role: "LEADER"},
			requesterMemberResult: nil,
			isSiteAdminResult:     false,
			isOrgAdminResult:      true,
		},
		{
			name:                  "org admin removal error returns wrapped error",
			groupID:               "group-1",
			targetUserID:          "target-1",
			requestingUserID:      "org-admin",
			findGroupResult:       &models.Group{ID: "group-1", OrganizationID: "org-1"},
			targetMemberResult:    &models.GroupMember{UserID: "target-1", Role: "MEMBER"},
			requesterMemberResult: nil,
			isSiteAdminResult:     false,
			isOrgAdminResult:      true,
			removeMemberErr:       errors.New("db error"),
			wantErrContains:       "failed to remove member",
		},
		{
			name:                  "regular member cannot remove another member returns ErrForbidden",
			groupID:               "group-1",
			targetUserID:          "target-1",
			requestingUserID:      "member-1",
			findGroupResult:       &models.Group{ID: "group-1", OrganizationID: "org-1"},
			targetMemberResult:    &models.GroupMember{UserID: "target-1", Role: "MEMBER"},
			requesterMemberResult: &models.GroupMember{UserID: "member-1", Role: "MEMBER"},
			isSiteAdminResult:     false,
			isOrgAdminResult:      false,
			wantErr:               ErrForbidden,
		},
		{
			name:                  "leader cannot remove OWNER returns ErrForbidden",
			groupID:               "group-1",
			targetUserID:          "owner-1",
			requestingUserID:      "leader-1",
			findGroupResult:       &models.Group{ID: "group-1", OrganizationID: "org-1"},
			targetMemberResult:    &models.GroupMember{UserID: "owner-1", Role: "OWNER"},
			findMembersResult:     []*models.GroupMember{{UserID: "owner-1", Role: "OWNER"}, {UserID: "owner-2", Role: "OWNER"}},
			requesterMemberResult: &models.GroupMember{UserID: "leader-1", Role: "LEADER"},
			isSiteAdminResult:     false,
			isOrgAdminResult:      false,
			wantErr:               ErrForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findMemberCallCount := 0
			groupRepo := &mockGroupRepoForGroup{
				findByIDFn: func(ctx context.Context, groupID string) (*models.Group, error) {
					if tt.findGroupErr != nil {
						return nil, tt.findGroupErr
					}
					return tt.findGroupResult, nil
				},
				findMemberFn: func(ctx context.Context, groupID, userID string) (*models.GroupMember, error) {
					findMemberCallCount++
					if findMemberCallCount == 1 {
						if tt.targetMemberErr != nil {
							return nil, tt.targetMemberErr
						}
						return tt.targetMemberResult, nil
					}
					if tt.requesterMemberErr != nil {
						return nil, tt.requesterMemberErr
					}
					return tt.requesterMemberResult, nil
				},
				findMembersFn: func(ctx context.Context, groupID string) ([]*models.GroupMember, error) {
					if tt.findMembersErr != nil {
						return nil, tt.findMembersErr
					}
					return tt.findMembersResult, nil
				},
				removeMemberFn: func(ctx context.Context, groupID, userID string) error {
					return tt.removeMemberErr
				},
			}
			orgRepo := &mockOrgRepoForGroup{
				isAdminFn: func(ctx context.Context, orgID, userID string) (bool, error) {
					if tt.isOrgAdminErr != nil {
						return false, tt.isOrgAdminErr
					}
					return tt.isOrgAdminResult, nil
				},
			}
			userRepo := &mockUserRepoForGroup{
				isSiteAdminFn: func(ctx context.Context, userID string) (bool, error) {
					if tt.isSiteAdminErr != nil {
						return false, tt.isSiteAdminErr
					}
					return tt.isSiteAdminResult, nil
				},
			}

			svc := NewGroupService(groupRepo, orgRepo, userRepo)

			err := svc.RemoveMember(context.Background(), tt.groupID, tt.targetUserID, tt.requestingUserID)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("RemoveMember() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if tt.wantErrContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("RemoveMember() error = %v, wantErrContains %v", err, tt.wantErrContains)
				}
				return
			}
			if err != nil {
				t.Errorf("RemoveMember() unexpected error = %v", err)
			}
		})
	}
}

// ========== GetUserRole Tests ==========

func TestGroupServiceImpl_GetUserRole(t *testing.T) {
	tests := []struct {
		name             string
		groupID          string
		userID           string
		findMemberResult *models.GroupMember
		findMemberErr    error
		wantRole         string
		wantErr          error
		wantErrContains  string
	}{
		{
			name:            "find member error returns wrapped error",
			groupID:         "group-1",
			userID:          "user-1",
			findMemberErr:   errors.New("db error"),
			wantErrContains: "failed to get user role",
		},
		{
			name:             "not a member returns empty string",
			groupID:          "group-1",
			userID:           "user-1",
			findMemberResult: nil,
			wantRole:         "",
		},
		{
			name:             "OWNER role returned",
			groupID:          "group-1",
			userID:           "owner-1",
			findMemberResult: &models.GroupMember{UserID: "owner-1", Role: "OWNER"},
			wantRole:         "OWNER",
		},
		{
			name:             "LEADER role returned",
			groupID:          "group-1",
			userID:           "leader-1",
			findMemberResult: &models.GroupMember{UserID: "leader-1", Role: "LEADER"},
			wantRole:         "LEADER",
		},
		{
			name:             "MEMBER role returned",
			groupID:          "group-1",
			userID:           "member-1",
			findMemberResult: &models.GroupMember{UserID: "member-1", Role: "MEMBER"},
			wantRole:         "MEMBER",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			groupRepo := &mockGroupRepoForGroup{
				findMemberFn: func(ctx context.Context, groupID, userID string) (*models.GroupMember, error) {
					if tt.findMemberErr != nil {
						return nil, tt.findMemberErr
					}
					return tt.findMemberResult, nil
				},
			}
			orgRepo := &mockOrgRepoForGroup{}
			userRepo := &mockUserRepoForGroup{}

			svc := NewGroupService(groupRepo, orgRepo, userRepo)

			got, err := svc.GetUserRole(context.Background(), tt.groupID, tt.userID)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("GetUserRole() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if tt.wantErrContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("GetUserRole() error = %v, wantErrContains %v", err, tt.wantErrContains)
				}
				return
			}
			if err != nil {
				t.Errorf("GetUserRole() unexpected error = %v", err)
				return
			}
			if got != tt.wantRole {
				t.Errorf("GetUserRole() = %v, want %v", got, tt.wantRole)
			}
		})
	}
}

// ========== GetMembers Tests ==========

func TestGroupServiceImpl_GetMembers(t *testing.T) {
	tests := []struct {
		name              string
		groupID           string
		findMembersResult []*models.GroupMember
		findMembersErr    error
		want              []*models.GroupMember
		wantErr           error
		wantErrContains   string
	}{
		{
			name:            "find members error returns wrapped error",
			groupID:         "group-1",
			findMembersErr:  errors.New("db error"),
			wantErrContains: "failed to get group members",
		},
		{
			name:              "nil result returns empty slice",
			groupID:           "group-1",
			findMembersResult: nil,
			want:              []*models.GroupMember{},
		},
		{
			name:              "empty slice returns empty slice",
			groupID:           "group-1",
			findMembersResult: []*models.GroupMember{},
			want:              []*models.GroupMember{},
		},
		{
			name:    "members found returns all members",
			groupID: "group-1",
			findMembersResult: []*models.GroupMember{
				{UserID: "user-1", Role: "OWNER"},
				{UserID: "user-2", Role: "MEMBER"},
			},
			want: []*models.GroupMember{
				{UserID: "user-1", Role: "OWNER"},
				{UserID: "user-2", Role: "MEMBER"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			groupRepo := &mockGroupRepoForGroup{
				findMembersFn: func(ctx context.Context, groupID string) ([]*models.GroupMember, error) {
					if tt.findMembersErr != nil {
						return nil, tt.findMembersErr
					}
					return tt.findMembersResult, nil
				},
			}
			orgRepo := &mockOrgRepoForGroup{}
			userRepo := &mockUserRepoForGroup{}

			svc := NewGroupService(groupRepo, orgRepo, userRepo)

			got, err := svc.GetMembers(context.Background(), tt.groupID)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("GetMembers() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if tt.wantErrContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("GetMembers() error = %v, wantErrContains %v", err, tt.wantErrContains)
				}
				return
			}
			if err != nil {
				t.Errorf("GetMembers() unexpected error = %v", err)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("GetMembers() got %d members, want %d", len(got), len(tt.want))
			}
		})
	}
}

// ========== GetMemberCount Tests ==========

func TestGroupServiceImpl_GetMemberCount(t *testing.T) {
	tests := []struct {
		name            string
		groupID         string
		countResult     int
		countErr        error
		want            int
		wantErr         error
		wantErrContains string
	}{
		{
			name:            "count error returns wrapped error",
			groupID:         "group-1",
			countErr:        errors.New("db error"),
			wantErrContains: "failed to count group members",
		},
		{
			name:        "zero members returns 0",
			groupID:     "group-1",
			countResult: 0,
			want:        0,
		},
		{
			name:        "multiple members returns count",
			groupID:     "group-1",
			countResult: 5,
			want:        5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			groupRepo := &mockGroupRepoForGroup{
				countMembersFn: func(ctx context.Context, groupID string) (int, error) {
					if tt.countErr != nil {
						return 0, tt.countErr
					}
					return tt.countResult, nil
				},
			}
			orgRepo := &mockOrgRepoForGroup{}
			userRepo := &mockUserRepoForGroup{}

			svc := NewGroupService(groupRepo, orgRepo, userRepo)

			got, err := svc.GetMemberCount(context.Background(), tt.groupID)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("GetMemberCount() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if tt.wantErrContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("GetMemberCount() error = %v, wantErrContains %v", err, tt.wantErrContains)
				}
				return
			}
			if err != nil {
				t.Errorf("GetMemberCount() unexpected error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("GetMemberCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ========== Additional Edge Case Tests ==========

// TestGroupServiceImpl_GetDetailByID_MemberConversion tests member response conversion
func TestGroupServiceImpl_GetDetailByID_MemberConversion(t *testing.T) {
	groupRepo := &mockGroupRepoForGroup{
		findByIDFn: func(ctx context.Context, groupID string) (*models.Group, error) {
			return &models.Group{ID: "group-1", OrganizationID: "org-1", Name: "Test Group", InviteCode: "CODE123", CreatedAt: time.Now()}, nil
		},
		findMemberFn: func(ctx context.Context, groupID, userID string) (*models.GroupMember, error) {
			return &models.GroupMember{UserID: "owner-1", Role: "OWNER"}, nil
		},
		findMembersFn: func(ctx context.Context, groupID string) ([]*models.GroupMember, error) {
			return []*models.GroupMember{
				{UserID: "user-1", Role: "OWNER", JoinedAt: time.Now(), User: &models.User{DisplayName: "John", Email: "john@test.com"}},
				{UserID: "user-2", Role: "MEMBER", JoinedAt: time.Now(), User: &models.User{DisplayName: "Jane", Email: "jane@test.com"}},
			}, nil
		},
	}
	orgRepo := &mockOrgRepoForGroup{}
	userRepo := &mockUserRepoForGroup{}

	svc := NewGroupService(groupRepo, orgRepo, userRepo)

	got, err := svc.GetDetailByID(context.Background(), "group-1", "owner-1")
	if err != nil {
		t.Fatalf("GetDetailByID() unexpected error = %v", err)
	}

	if len(got.Members) != 2 {
		t.Fatalf("GetDetailByID() got %d members, want 2", len(got.Members))
	}

	if got.Members[0].DisplayName != "John" {
		t.Errorf("GetDetailByID() first member name = %v, want John", got.Members[0].DisplayName)
	}
	if got.Members[1].Email != "jane@test.com" {
		t.Errorf("GetDetailByID() second member email = %v, want jane@test.com", got.Members[1].Email)
	}
}

// TestGroupServiceImpl_Create_ContextPropagation verifies context is properly passed
func TestGroupServiceImpl_Create_ContextPropagation(t *testing.T) {
	type contextKey string
	const testKey contextKey = "test-key"
	const testValue = "test-value"

	var capturedCtx context.Context

	groupRepo := &mockGroupRepoForGroup{
		createFn: func(ctx context.Context, name string, description *string, orgID string, inviteCode string, createdBy string) (*models.Group, error) {
			capturedCtx = ctx
			return &models.Group{ID: "group-1"}, nil
		},
		generateInviteCodeFn: func() string {
			return "INVITE123"
		},
		addMemberFn: func(ctx context.Context, groupID, userID, role string) (*models.GroupMember, error) {
			return &models.GroupMember{}, nil
		},
	}
	orgRepo := &mockOrgRepoForGroup{
		findByIDFn: func(ctx context.Context, orgID string) (*models.Organization, error) {
			return &models.Organization{ID: "org-1"}, nil
		},
		isAdminFn: func(ctx context.Context, orgID, userID string) (bool, error) {
			return true, nil
		},
	}
	userRepo := &mockUserRepoForGroup{}

	svc := NewGroupService(groupRepo, orgRepo, userRepo)

	ctx := context.WithValue(context.Background(), testKey, testValue)
	_, err := svc.Create(ctx, "user-1", &models.CreateGroupRequest{Name: "Test", OrganizationID: "org-1"})
	if err != nil {
		t.Fatalf("Create() unexpected error = %v", err)
	}

	if capturedCtx == nil {
		t.Fatal("context was not propagated")
	}

	if capturedCtx.Value(testKey) != testValue {
		t.Errorf("context value = %v, want %v", capturedCtx.Value(testKey), testValue)
	}
}
