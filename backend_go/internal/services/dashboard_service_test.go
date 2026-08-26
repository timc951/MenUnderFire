package services

import (
	"context"
	"errors"
	"strings"
	"testing"

	"menunderfire/internal/models"
)

// mockUserService is a mock implementation of UserService for testing
type mockUserService struct {
	permissions *models.UserPermissionsResponse
	err         error
}

func (m *mockUserService) GetByID(ctx context.Context, userID string) (*models.User, error) {
	return nil, nil // not used by dashboard service
}

func (m *mockUserService) GetByExternalID(ctx context.Context, externalID string) (*models.User, error) {
	return nil, nil // not used by dashboard service
}

func (m *mockUserService) Update(ctx context.Context, userID string, req *models.UpdateUserRequest) (*models.User, error) {
	return nil, nil // not used by dashboard service
}

func (m *mockUserService) GetPermissions(ctx context.Context, userID string) (*models.UserPermissionsResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.permissions, nil
}

func (m *mockUserService) AcceptAgreement(ctx context.Context, userID, version, ipAddress, userAgent string) (*models.AcceptAgreementResponse, error) {
	return nil, nil // not used by dashboard service
}

// mockOrgRepo is a mock implementation of repositories.OrganizationRepository
type mockOrgRepo struct {
	countFn func(ctx context.Context) (int64, error)
}

func (m *mockOrgRepo) FindByID(ctx context.Context, orgID string) (*models.Organization, error) {
	return nil, nil
}
func (m *mockOrgRepo) FindByUserID(ctx context.Context, userID string) ([]*models.Organization, error) {
	return nil, nil
}
func (m *mockOrgRepo) FindAll(ctx context.Context) ([]*models.Organization, error) {
	return nil, nil
}
func (m *mockOrgRepo) Create(ctx context.Context, name string, description *string, createdByID string) (*models.Organization, error) {
	return nil, nil
}
func (m *mockOrgRepo) Update(ctx context.Context, orgID string, name string, description *string) (*models.Organization, error) {
	return nil, nil
}
func (m *mockOrgRepo) FindAdmins(ctx context.Context, orgID string) ([]*models.OrganizationAdmin, error) {
	return nil, nil
}
func (m *mockOrgRepo) FindAdmin(ctx context.Context, orgID, userID string) (*models.OrganizationAdmin, error) {
	return nil, nil
}
func (m *mockOrgRepo) AddAdmin(ctx context.Context, orgID, userID string) error {
	return nil
}
func (m *mockOrgRepo) IsMember(ctx context.Context, orgID, userID string) (bool, error) {
	return false, nil
}
func (m *mockOrgRepo) IsAdmin(ctx context.Context, orgID, userID string) (bool, error) {
	return false, nil
}
func (m *mockOrgRepo) Count(ctx context.Context) (int64, error) {
	if m.countFn != nil {
		return m.countFn(ctx)
	}
	return 0, nil
}

// mockGroupRepo is a mock implementation of repositories.GroupRepository
type mockGroupRepo struct {
	countFn       func(ctx context.Context) (int64, error)
	countByOrgsFn func(ctx context.Context, orgIDs []string) (int64, error)
}

func (m *mockGroupRepo) FindByID(ctx context.Context, groupID string) (*models.Group, error) {
	return nil, nil
}
func (m *mockGroupRepo) FindByInviteCode(ctx context.Context, inviteCode string) (*models.Group, error) {
	return nil, nil
}
func (m *mockGroupRepo) FindByUserID(ctx context.Context, userID string) ([]*models.Group, error) {
	return nil, nil
}
func (m *mockGroupRepo) FindByOrganizationID(ctx context.Context, orgID string) ([]*models.Group, error) {
	return nil, nil
}
func (m *mockGroupRepo) Create(ctx context.Context, name string, description *string, orgID string, inviteCode string, createdBy string) (*models.Group, error) {
	return nil, nil
}
func (m *mockGroupRepo) GenerateInviteCode() string {
	return ""
}
func (m *mockGroupRepo) FindMember(ctx context.Context, groupID, userID string) (*models.GroupMember, error) {
	return nil, nil
}
func (m *mockGroupRepo) FindMembers(ctx context.Context, groupID string) ([]*models.GroupMember, error) {
	return nil, nil
}
func (m *mockGroupRepo) CountMembers(ctx context.Context, groupID string) (int, error) {
	return 0, nil
}
func (m *mockGroupRepo) AddMember(ctx context.Context, groupID, userID, role string) (*models.GroupMember, error) {
	return nil, nil
}
func (m *mockGroupRepo) RemoveMember(ctx context.Context, groupID, userID string) error {
	return nil
}
func (m *mockGroupRepo) Count(ctx context.Context) (int64, error) {
	if m.countFn != nil {
		return m.countFn(ctx)
	}
	return 0, nil
}
func (m *mockGroupRepo) CountByOrganizationIDs(ctx context.Context, orgIDs []string) (int64, error) {
	if m.countByOrgsFn != nil {
		return m.countByOrgsFn(ctx, orgIDs)
	}
	return 0, nil
}

func (m *mockGroupRepo) UpdateMemberRole(ctx context.Context, groupID, userID, role string) error {
	return nil
}

func (m *mockGroupRepo) UpdateSettings(ctx context.Context, groupID string, requirePostApproval, allowAnonymousPosts bool) error {
	return nil
}

// Helper function to create pointer to int64
func int64Ptr(v int64) *int64 {
	return &v
}

func TestDashboardServiceImpl_GetStats(t *testing.T) {
	tests := []struct {
		name               string
		userID             string
		mockPermissions    *models.UserPermissionsResponse
		mockPermissionsErr error
		orgCounterResult   int64
		orgCounterErr      error
		groupCounterResult int64
		groupCounterErr    error
		groupByOrgsResult  int64
		groupByOrgsErr     error
		want               *models.DashboardStatsResponse
		wantErr            error
		wantErrContains    string
	}{
		{
			name:               "user not found - GetPermissions returns ErrUserNotFound",
			userID:             "nonexistent-user",
			mockPermissionsErr: ErrUserNotFound,
			wantErr:            ErrUserNotFound,
		},
		{
			name:   "regular user with no permissions - returns nil org count, 0 group count",
			userID: "regular-user",
			mockPermissions: &models.UserPermissionsResponse{
				IsSiteAdmin:            false,
				AdminOfOrganizationIDs: []string{},
				OwnedGroupIDs:          []string{},
				MemberGroupIDs:         []string{},
			},
			want: &models.DashboardStatsResponse{
				OrganizationCount: nil,
				GroupCount:        0,
			},
		},
		{
			name:   "site admin - returns total org count and total group count",
			userID: "site-admin",
			mockPermissions: &models.UserPermissionsResponse{
				IsSiteAdmin:            true,
				AdminOfOrganizationIDs: []string{},
				OwnedGroupIDs:          []string{},
				MemberGroupIDs:         []string{},
			},
			orgCounterResult:   10,
			groupCounterResult: 50,
			want: &models.DashboardStatsResponse{
				OrganizationCount: int64Ptr(10),
				GroupCount:        50,
			},
		},
		{
			name:   "org admin of one organization - returns nil org count, group count for their org",
			userID: "org-admin-single",
			mockPermissions: &models.UserPermissionsResponse{
				IsSiteAdmin:            false,
				AdminOfOrganizationIDs: []string{"org-1"},
				OwnedGroupIDs:          []string{},
				MemberGroupIDs:         []string{},
			},
			groupByOrgsResult: 5,
			want: &models.DashboardStatsResponse{
				OrganizationCount: nil,
				GroupCount:        5,
			},
		},
		{
			name:   "org admin of multiple organizations - returns nil org count, group count for all their orgs",
			userID: "org-admin-multiple",
			mockPermissions: &models.UserPermissionsResponse{
				IsSiteAdmin:            false,
				AdminOfOrganizationIDs: []string{"org-1", "org-2", "org-3"},
				OwnedGroupIDs:          []string{},
				MemberGroupIDs:         []string{},
			},
			groupByOrgsResult: 15,
			want: &models.DashboardStatsResponse{
				OrganizationCount: nil,
				GroupCount:        15,
			},
		},
		{
			name:   "org admin with no groups in their orgs - returns nil org count, 0 group count",
			userID: "org-admin-no-groups",
			mockPermissions: &models.UserPermissionsResponse{
				IsSiteAdmin:            false,
				AdminOfOrganizationIDs: []string{"org-1"},
				OwnedGroupIDs:          []string{},
				MemberGroupIDs:         []string{},
			},
			groupByOrgsResult: 0,
			want: &models.DashboardStatsResponse{
				OrganizationCount: nil,
				GroupCount:        0,
			},
		},
		{
			name:               "database error when getting permissions - returns wrapped error",
			userID:             "user-db-error",
			mockPermissionsErr: errors.New("database connection failed"),
			wantErrContains:    "failed to get user permissions",
		},
		{
			name:   "database error when counting orgs (site admin) - returns wrapped error",
			userID: "site-admin-org-error",
			mockPermissions: &models.UserPermissionsResponse{
				IsSiteAdmin:            true,
				AdminOfOrganizationIDs: []string{},
				OwnedGroupIDs:          []string{},
				MemberGroupIDs:         []string{},
			},
			orgCounterErr:   errors.New("failed to query organizations table"),
			wantErrContains: "failed to count organizations",
		},
		{
			name:   "database error when counting groups (site admin) - returns wrapped error",
			userID: "site-admin-group-error",
			mockPermissions: &models.UserPermissionsResponse{
				IsSiteAdmin:            true,
				AdminOfOrganizationIDs: []string{},
				OwnedGroupIDs:          []string{},
				MemberGroupIDs:         []string{},
			},
			orgCounterResult: 10,
			groupCounterErr:  errors.New("failed to query groups table"),
			wantErrContains:  "failed to count groups",
		},
		{
			name:   "database error when counting groups by orgs (org admin) - returns wrapped error",
			userID: "org-admin-group-error",
			mockPermissions: &models.UserPermissionsResponse{
				IsSiteAdmin:            false,
				AdminOfOrganizationIDs: []string{"org-1", "org-2"},
				OwnedGroupIDs:          []string{},
				MemberGroupIDs:         []string{},
			},
			groupByOrgsErr:  errors.New("failed to query groups by organization"),
			wantErrContains: "failed to count groups for organizations",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUser := &mockUserService{
				permissions: tt.mockPermissions,
				err:         tt.mockPermissionsErr,
			}

			orgRepo := &mockOrgRepo{
				countFn: func(ctx context.Context) (int64, error) {
					if tt.orgCounterErr != nil {
						return 0, tt.orgCounterErr
					}
					return tt.orgCounterResult, nil
				},
			}

			groupRepo := &mockGroupRepo{
				countFn: func(ctx context.Context) (int64, error) {
					if tt.groupCounterErr != nil {
						return 0, tt.groupCounterErr
					}
					return tt.groupCounterResult, nil
				},
				countByOrgsFn: func(ctx context.Context, orgIDs []string) (int64, error) {
					if tt.groupByOrgsErr != nil {
						return 0, tt.groupByOrgsErr
					}
					return tt.groupByOrgsResult, nil
				},
			}

			service := NewDashboardService(mockUser, orgRepo, groupRepo)

			ctx := context.Background()
			got, err := service.GetStats(ctx, tt.userID)

			// Check error cases
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("GetStats() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("GetStats() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if tt.wantErrContains != "" {
				if err == nil {
					t.Errorf("GetStats() error = nil, wantErrContains %v", tt.wantErrContains)
					return
				}
				if !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("GetStats() error = %v, wantErrContains %v", err, tt.wantErrContains)
				}
				return
			}

			// Check success cases
			if err != nil {
				t.Errorf("GetStats() unexpected error = %v", err)
				return
			}

			if got == nil {
				t.Fatal("GetStats() returned nil response")
			}

			// Compare OrganizationCount (pointer comparison)
			if tt.want.OrganizationCount == nil {
				if got.OrganizationCount != nil {
					t.Errorf("GetStats() OrganizationCount = %v, want nil", *got.OrganizationCount)
				}
			} else {
				if got.OrganizationCount == nil {
					t.Errorf("GetStats() OrganizationCount = nil, want %v", *tt.want.OrganizationCount)
				} else if *got.OrganizationCount != *tt.want.OrganizationCount {
					t.Errorf("GetStats() OrganizationCount = %v, want %v", *got.OrganizationCount, *tt.want.OrganizationCount)
				}
			}

			// Compare GroupCount
			if got.GroupCount != tt.want.GroupCount {
				t.Errorf("GetStats() GroupCount = %v, want %v", got.GroupCount, tt.want.GroupCount)
			}
		})
	}
}

// TestDashboardServiceImpl_GetStats_ContextPropagation verifies that context is properly propagated
func TestDashboardServiceImpl_GetStats_ContextPropagation(t *testing.T) {
	type contextKey string
	const testKey contextKey = "test-key"
	const testValue = "test-value"

	var capturedCtx context.Context

	mockUser := &mockUserService{
		permissions: &models.UserPermissionsResponse{
			IsSiteAdmin:            true,
			AdminOfOrganizationIDs: []string{},
			OwnedGroupIDs:          []string{},
			MemberGroupIDs:         []string{},
		},
	}

	orgRepo := &mockOrgRepo{
		countFn: func(ctx context.Context) (int64, error) {
			capturedCtx = ctx
			return 5, nil
		},
	}

	groupRepo := &mockGroupRepo{
		countFn: func(ctx context.Context) (int64, error) {
			return 10, nil
		},
	}

	service := NewDashboardService(mockUser, orgRepo, groupRepo)

	ctx := context.WithValue(context.Background(), testKey, testValue)
	_, err := service.GetStats(ctx, "user-id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if capturedCtx == nil {
		t.Fatal("context was not propagated to orgCounter")
	}

	if capturedCtx.Value(testKey) != testValue {
		t.Errorf("context value = %v, want %v", capturedCtx.Value(testKey), testValue)
	}
}

// TestDashboardServiceImpl_GetStats_OrgAdminOrgIDsPropagation verifies org IDs are passed correctly
func TestDashboardServiceImpl_GetStats_OrgAdminOrgIDsPropagation(t *testing.T) {
	var capturedOrgIDs []string

	expectedOrgIDs := []string{"org-1", "org-2", "org-3"}

	mockUser := &mockUserService{
		permissions: &models.UserPermissionsResponse{
			IsSiteAdmin:            false,
			AdminOfOrganizationIDs: expectedOrgIDs,
			OwnedGroupIDs:          []string{},
			MemberGroupIDs:         []string{},
		},
	}

	orgRepo := &mockOrgRepo{}

	groupRepo := &mockGroupRepo{
		countByOrgsFn: func(ctx context.Context, orgIDs []string) (int64, error) {
			capturedOrgIDs = orgIDs
			return 20, nil
		},
	}

	service := NewDashboardService(mockUser, orgRepo, groupRepo)

	ctx := context.Background()
	result, err := service.GetStats(ctx, "org-admin-user")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.GroupCount != 20 {
		t.Errorf("GroupCount = %v, want 20", result.GroupCount)
	}

	if len(capturedOrgIDs) != len(expectedOrgIDs) {
		t.Fatalf("captured org IDs length = %v, want %v", len(capturedOrgIDs), len(expectedOrgIDs))
	}

	for i, orgID := range expectedOrgIDs {
		if capturedOrgIDs[i] != orgID {
			t.Errorf("captured org ID[%d] = %v, want %v", i, capturedOrgIDs[i], orgID)
		}
	}
}
