package services

import (
	"context"
	"errors"
	"strings"
	"testing"

	"menunderfire/internal/models"
)

// ========== Mock Repositories ==========

// mockGroupMessageRepo mocks GroupMessageRepository
type mockGroupMessageRepo struct {
	findByIDFn                  func(ctx context.Context, messageID string) (*models.GroupMessage, error)
	findByGroupIDFn             func(ctx context.Context, groupID string) ([]*models.GroupMessage, error)
	findFormMessagesByGroupIDFn func(ctx context.Context, groupID string) ([]*models.GroupMessage, error)
	createFn                    func(ctx context.Context, groupID, senderID, content string, notifyMembers bool, formID *string) (*models.GroupMessage, error)
	deleteFn                    func(ctx context.Context, messageID string) error
}

func (m *mockGroupMessageRepo) FindByID(ctx context.Context, messageID string) (*models.GroupMessage, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, messageID)
	}
	return nil, errors.New("findByIDFn not set")
}

func (m *mockGroupMessageRepo) FindByGroupID(ctx context.Context, groupID string) ([]*models.GroupMessage, error) {
	if m.findByGroupIDFn != nil {
		return m.findByGroupIDFn(ctx, groupID)
	}
	return nil, errors.New("findByGroupIDFn not set")
}

func (m *mockGroupMessageRepo) FindFormMessagesByGroupID(ctx context.Context, groupID string) ([]*models.GroupMessage, error) {
	if m.findFormMessagesByGroupIDFn != nil {
		return m.findFormMessagesByGroupIDFn(ctx, groupID)
	}
	return nil, errors.New("findFormMessagesByGroupIDFn not set")
}

func (m *mockGroupMessageRepo) Create(ctx context.Context, groupID, senderID, content string, notifyMembers bool, formID *string) (*models.GroupMessage, error) {
	if m.createFn != nil {
		return m.createFn(ctx, groupID, senderID, content, notifyMembers, formID)
	}
	return nil, errors.New("createFn not set")
}

func (m *mockGroupMessageRepo) Delete(ctx context.Context, messageID string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, messageID)
	}
	return errors.New("deleteFn not set")
}

// mockGroupRepoForMsg mocks GroupRepository (different name to avoid conflicts)
type mockGroupRepoForMsg struct {
	findByIDFn func(ctx context.Context, groupID string) (*models.Group, error)
}

func (m *mockGroupRepoForMsg) FindByID(ctx context.Context, groupID string) (*models.Group, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, groupID)
	}
	return nil, errors.New("findByIDFn not set")
}

func (m *mockGroupRepoForMsg) FindByInviteCode(ctx context.Context, inviteCode string) (*models.Group, error) {
	return nil, errors.New("not implemented")
}

func (m *mockGroupRepoForMsg) FindByUserID(ctx context.Context, userID string) ([]*models.Group, error) {
	return nil, errors.New("not implemented")
}

func (m *mockGroupRepoForMsg) FindByOrganizationID(ctx context.Context, orgID string) ([]*models.Group, error) {
	return nil, errors.New("not implemented")
}

func (m *mockGroupRepoForMsg) Create(ctx context.Context, name string, description *string, orgID string, inviteCode string, createdBy string) (*models.Group, error) {
	return nil, errors.New("not implemented")
}

func (m *mockGroupRepoForMsg) GenerateInviteCode() string {
	return ""
}

func (m *mockGroupRepoForMsg) FindMember(ctx context.Context, groupID, userID string) (*models.GroupMember, error) {
	return nil, errors.New("not implemented")
}

func (m *mockGroupRepoForMsg) FindMembers(ctx context.Context, groupID string) ([]*models.GroupMember, error) {
	return nil, errors.New("not implemented")
}

func (m *mockGroupRepoForMsg) CountMembers(ctx context.Context, groupID string) (int, error) {
	return 0, errors.New("not implemented")
}

func (m *mockGroupRepoForMsg) AddMember(ctx context.Context, groupID, userID, role string) (*models.GroupMember, error) {
	return nil, errors.New("not implemented")
}

func (m *mockGroupRepoForMsg) RemoveMember(ctx context.Context, groupID, userID string) error {
	return errors.New("not implemented")
}

func (m *mockGroupRepoForMsg) Count(ctx context.Context) (int64, error) {
	return 0, errors.New("not implemented")
}

func (m *mockGroupRepoForMsg) CountByOrganizationIDs(ctx context.Context, orgIDs []string) (int64, error) {
	return 0, errors.New("not implemented")
}

func (m *mockGroupRepoForMsg) UpdateMemberRole(ctx context.Context, groupID, userID, role string) error {
	return nil
}

func (m *mockGroupRepoForMsg) UpdateSettings(ctx context.Context, groupID string, requirePostApproval, allowAnonymousPosts bool) error {
	return nil
}

// mockGroupServiceForMsg mocks GroupService (different name to avoid conflicts)
type mockGroupServiceForMsg struct {
	getUserRoleFn func(ctx context.Context, groupID, userID string) (string, error)
}

func (m *mockGroupServiceForMsg) Create(ctx context.Context, userID string, req *models.CreateGroupRequest) (*models.Group, error) {
	return nil, errors.New("not implemented")
}

func (m *mockGroupServiceForMsg) List(ctx context.Context, userID string) ([]*models.Group, error) {
	return nil, errors.New("not implemented")
}

func (m *mockGroupServiceForMsg) GetByID(ctx context.Context, groupID string) (*models.Group, error) {
	return nil, errors.New("not implemented")
}

func (m *mockGroupServiceForMsg) GetDetailByID(ctx context.Context, groupID string, userID string) (*models.GroupDetailResponse, error) {
	return nil, errors.New("not implemented")
}

func (m *mockGroupServiceForMsg) Join(ctx context.Context, groupID string, userID string, inviteCode string) error {
	return errors.New("not implemented")
}

func (m *mockGroupServiceForMsg) JoinByInviteCode(ctx context.Context, userID string, inviteCode string) error {
	return errors.New("not implemented")
}

func (m *mockGroupServiceForMsg) RemoveMember(ctx context.Context, groupID string, targetUserID string, requestingUserID string) error {
	return errors.New("not implemented")
}

func (m *mockGroupServiceForMsg) GetUserRole(ctx context.Context, groupID string, userID string) (string, error) {
	if m.getUserRoleFn != nil {
		return m.getUserRoleFn(ctx, groupID, userID)
	}
	return "", errors.New("getUserRoleFn not set")
}

func (m *mockGroupServiceForMsg) GetMembers(ctx context.Context, groupID string) ([]*models.GroupMember, error) {
	return nil, errors.New("not implemented")
}

func (m *mockGroupServiceForMsg) GetMemberCount(ctx context.Context, groupID string) (int, error) {
	return 0, errors.New("not implemented")
}

func (m *mockGroupServiceForMsg) UpdateMemberRole(ctx context.Context, groupID string, targetUserID string, newRole string, requestingUserID string) error {
	return nil
}

func (m *mockGroupServiceForMsg) UpdateSettings(ctx context.Context, groupID string, userID string, req *models.UpdateGroupSettingsRequest) error {
	return nil
}

// mockOrgRepoForMsg mocks OrganizationRepository (different name to avoid conflicts)
type mockOrgRepoForMsg struct {
	isAdminFn func(ctx context.Context, orgID, userID string) (bool, error)
}

func (m *mockOrgRepoForMsg) FindByID(ctx context.Context, orgID string) (*models.Organization, error) {
	return nil, errors.New("not implemented")
}

func (m *mockOrgRepoForMsg) FindByUserID(ctx context.Context, userID string) ([]*models.Organization, error) {
	return nil, errors.New("not implemented")
}

func (m *mockOrgRepoForMsg) FindAll(ctx context.Context) ([]*models.Organization, error) {
	return nil, errors.New("not implemented")
}

func (m *mockOrgRepoForMsg) Create(ctx context.Context, name string, description *string, createdByID string) (*models.Organization, error) {
	return nil, errors.New("not implemented")
}

func (m *mockOrgRepoForMsg) Update(ctx context.Context, orgID string, name string, description *string) (*models.Organization, error) {
	return nil, errors.New("not implemented")
}

func (m *mockOrgRepoForMsg) FindAdmins(ctx context.Context, orgID string) ([]*models.OrganizationAdmin, error) {
	return nil, errors.New("not implemented")
}

func (m *mockOrgRepoForMsg) FindAdmin(ctx context.Context, orgID, userID string) (*models.OrganizationAdmin, error) {
	return nil, errors.New("not implemented")
}

func (m *mockOrgRepoForMsg) AddAdmin(ctx context.Context, orgID, userID string) error {
	return errors.New("not implemented")
}

func (m *mockOrgRepoForMsg) IsMember(ctx context.Context, orgID, userID string) (bool, error) {
	return false, errors.New("not implemented")
}

func (m *mockOrgRepoForMsg) IsAdmin(ctx context.Context, orgID, userID string) (bool, error) {
	if m.isAdminFn != nil {
		return m.isAdminFn(ctx, orgID, userID)
	}
	return false, errors.New("isAdminFn not set")
}

func (m *mockOrgRepoForMsg) Count(ctx context.Context) (int64, error) {
	return 0, errors.New("not implemented")
}

// mockUserRepoForMsg mocks UserRepository (different name to avoid conflicts)
type mockUserRepoForMsg struct {
	isSiteAdminFn func(ctx context.Context, userID string) (bool, error)
}

func (m *mockUserRepoForMsg) FindByID(ctx context.Context, userID string) (*models.User, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserRepoForMsg) FindByExternalID(ctx context.Context, externalID string) (*models.User, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserRepoForMsg) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserRepoForMsg) Count(ctx context.Context) (int64, error) {
	return 0, errors.New("not implemented")
}

func (m *mockUserRepoForMsg) Create(ctx context.Context, email, displayName, externalID string) (*models.User, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserRepoForMsg) CreateAsSiteAdmin(ctx context.Context, email, displayName, externalID string) (*models.User, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserRepoForMsg) Update(ctx context.Context, userID, displayName string) (*models.User, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserRepoForMsg) UpdateExternalID(ctx context.Context, userID, externalID string) error {
	return errors.New("not implemented")
}

func (m *mockUserRepoForMsg) UpdateInvitationInfo(ctx context.Context, userID, invitedByID, invitationID string) error {
	return errors.New("not implemented")
}

func (m *mockUserRepoForMsg) IsSiteAdmin(ctx context.Context, userID string) (bool, error) {
	if m.isSiteAdminFn != nil {
		return m.isSiteAdminFn(ctx, userID)
	}
	return false, errors.New("isSiteAdminFn not set")
}

func (m *mockUserRepoForMsg) FindAdminOrganizationIDs(ctx context.Context, userID string) ([]string, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserRepoForMsg) FindOwnedGroupIDs(ctx context.Context, userID string) ([]string, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserRepoForMsg) FindMemberGroupIDs(ctx context.Context, userID string) ([]string, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserRepoForMsg) RecordAgreementAcceptance(ctx context.Context, userID, version, signature, ipAddress, userAgent string) error {
	return errors.New("not implemented")
}

// ========== List Tests ==========

func TestGroupMessageServiceImpl_List(t *testing.T) {
	tests := []struct {
		name               string
		groupID            string
		userID             string
		findGroupResult    *models.Group
		findGroupErr       error
		getUserRoleResult  string
		getUserRoleErr     error
		isSiteAdminResult  bool
		isSiteAdminErr     error
		isOrgAdminResult   bool
		isOrgAdminErr      error
		findMessagesResult []*models.GroupMessage
		findMessagesErr    error
		want               []*models.GroupMessage
		wantErr            error
		wantErrContains    string
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
			name:            "get user role error returns wrapped error",
			groupID:         "group-1",
			userID:          "user-1",
			findGroupResult: &models.Group{ID: "group-1", OrganizationID: "org-1"},
			getUserRoleErr:  errors.New("db error"),
			wantErrContains: "failed to get user role",
		},
		{
			name:               "group member can list messages",
			groupID:            "group-1",
			userID:             "member-1",
			findGroupResult:    &models.Group{ID: "group-1", OrganizationID: "org-1"},
			getUserRoleResult:  "MEMBER",
			findMessagesResult: []*models.GroupMessage{{ID: "msg-1", Content: "Hello"}},
			want:               []*models.GroupMessage{{ID: "msg-1", Content: "Hello"}},
		},
		{
			name:               "group owner can list messages",
			groupID:            "group-1",
			userID:             "owner-1",
			findGroupResult:    &models.Group{ID: "group-1", OrganizationID: "org-1"},
			getUserRoleResult:  "OWNER",
			findMessagesResult: []*models.GroupMessage{{ID: "msg-1"}},
			want:               []*models.GroupMessage{{ID: "msg-1"}},
		},
		{
			name:               "group leader can list messages",
			groupID:            "group-1",
			userID:             "leader-1",
			findGroupResult:    &models.Group{ID: "group-1", OrganizationID: "org-1"},
			getUserRoleResult:  "LEADER",
			findMessagesResult: []*models.GroupMessage{{ID: "msg-1"}},
			want:               []*models.GroupMessage{{ID: "msg-1"}},
		},
		{
			name:              "member - find messages error returns wrapped error",
			groupID:           "group-1",
			userID:            "member-1",
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			getUserRoleResult: "MEMBER",
			findMessagesErr:   errors.New("db error"),
			wantErrContains:   "failed to fetch messages",
		},
		{
			name:              "non-member - site admin check error returns wrapped error",
			groupID:           "group-1",
			userID:            "user-1",
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			getUserRoleResult: "",
			isSiteAdminErr:    errors.New("db error"),
			wantErrContains:   "failed to check site admin status",
		},
		{
			name:               "site admin can list messages",
			groupID:            "group-1",
			userID:             "site-admin",
			findGroupResult:    &models.Group{ID: "group-1", OrganizationID: "org-1"},
			getUserRoleResult:  "",
			isSiteAdminResult:  true,
			findMessagesResult: []*models.GroupMessage{{ID: "msg-1"}},
			want:               []*models.GroupMessage{{ID: "msg-1"}},
		},
		{
			name:              "site admin - find messages error returns wrapped error",
			groupID:           "group-1",
			userID:            "site-admin",
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			getUserRoleResult: "",
			isSiteAdminResult: true,
			findMessagesErr:   errors.New("db error"),
			wantErrContains:   "failed to fetch messages",
		},
		{
			name:              "non-member non-site-admin - org admin check error returns wrapped error",
			groupID:           "group-1",
			userID:            "user-1",
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			getUserRoleResult: "",
			isSiteAdminResult: false,
			isOrgAdminErr:     errors.New("db error"),
			wantErrContains:   "failed to check org admin status",
		},
		{
			name:               "org admin can list messages",
			groupID:            "group-1",
			userID:             "org-admin",
			findGroupResult:    &models.Group{ID: "group-1", OrganizationID: "org-1"},
			getUserRoleResult:  "",
			isSiteAdminResult:  false,
			isOrgAdminResult:   true,
			findMessagesResult: []*models.GroupMessage{{ID: "msg-1"}},
			want:               []*models.GroupMessage{{ID: "msg-1"}},
		},
		{
			name:              "org admin - find messages error returns wrapped error",
			groupID:           "group-1",
			userID:            "org-admin",
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			getUserRoleResult: "",
			isSiteAdminResult: false,
			isOrgAdminResult:  true,
			findMessagesErr:   errors.New("db error"),
			wantErrContains:   "failed to fetch messages",
		},
		{
			name:              "non-member non-admin returns ErrForbidden",
			groupID:           "group-1",
			userID:            "random-user",
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			getUserRoleResult: "",
			isSiteAdminResult: false,
			isOrgAdminResult:  false,
			wantErr:           ErrForbidden,
		},
		{
			name:               "empty messages list returns empty slice",
			groupID:            "group-1",
			userID:             "member-1",
			findGroupResult:    &models.Group{ID: "group-1", OrganizationID: "org-1"},
			getUserRoleResult:  "MEMBER",
			findMessagesResult: []*models.GroupMessage{},
			want:               []*models.GroupMessage{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messageRepo := &mockGroupMessageRepo{
				findByGroupIDFn: func(ctx context.Context, groupID string) ([]*models.GroupMessage, error) {
					if tt.findMessagesErr != nil {
						return nil, tt.findMessagesErr
					}
					return tt.findMessagesResult, nil
				},
			}

			groupRepo := &mockGroupRepoForMsg{
				findByIDFn: func(ctx context.Context, groupID string) (*models.Group, error) {
					if tt.findGroupErr != nil {
						return nil, tt.findGroupErr
					}
					return tt.findGroupResult, nil
				},
			}

			groupService := &mockGroupServiceForMsg{
				getUserRoleFn: func(ctx context.Context, groupID, userID string) (string, error) {
					if tt.getUserRoleErr != nil {
						return "", tt.getUserRoleErr
					}
					return tt.getUserRoleResult, nil
				},
			}

			orgRepo := &mockOrgRepoForMsg{
				isAdminFn: func(ctx context.Context, orgID, userID string) (bool, error) {
					if tt.isOrgAdminErr != nil {
						return false, tt.isOrgAdminErr
					}
					return tt.isOrgAdminResult, nil
				},
			}

			userRepo := &mockUserRepoForMsg{
				isSiteAdminFn: func(ctx context.Context, userID string) (bool, error) {
					if tt.isSiteAdminErr != nil {
						return false, tt.isSiteAdminErr
					}
					return tt.isSiteAdminResult, nil
				},
			}

			svc := NewGroupMessageService(messageRepo, groupRepo, groupService, orgRepo, userRepo, nil, nil)

			got, err := svc.List(context.Background(), tt.groupID, tt.userID)

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
				t.Errorf("List() got %d messages, want %d", len(got), len(tt.want))
			}
		})
	}
}

// ========== Send Tests ==========

func TestGroupMessageServiceImpl_Send(t *testing.T) {
	tests := []struct {
		name                string
		groupID             string
		userID              string
		req                 *models.SendGroupMessageRequest
		findGroupResult     *models.Group
		findGroupErr        error
		getUserRoleResult   string
		getUserRoleErr      error
		createMessageResult *models.GroupMessage
		createMessageErr    error
		want                *models.GroupMessage
		wantErr             error
		wantErrContains     string
	}{
		{
			name:            "find group error returns wrapped error",
			groupID:         "group-1",
			userID:          "user-1",
			req:             &models.SendGroupMessageRequest{Content: "Hello"},
			findGroupErr:    errors.New("db error"),
			wantErrContains: "failed to find group",
		},
		{
			name:            "group not found returns ErrGroupNotFound",
			groupID:         "group-1",
			userID:          "user-1",
			req:             &models.SendGroupMessageRequest{Content: "Hello"},
			findGroupResult: nil,
			wantErr:         ErrGroupNotFound,
		},
		{
			name:            "get user role error returns wrapped error",
			groupID:         "group-1",
			userID:          "user-1",
			req:             &models.SendGroupMessageRequest{Content: "Hello"},
			findGroupResult: &models.Group{ID: "group-1", OrganizationID: "org-1"},
			getUserRoleErr:  errors.New("db error"),
			wantErrContains: "failed to get user role",
		},
		{
			name:              "non-member returns ErrForbidden",
			groupID:           "group-1",
			userID:            "non-member",
			req:               &models.SendGroupMessageRequest{Content: "Hello"},
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			getUserRoleResult: "",
			wantErr:           ErrForbidden,
		},
		{
			name:              "empty content returns ErrValidation",
			groupID:           "group-1",
			userID:            "member-1",
			req:               &models.SendGroupMessageRequest{Content: ""},
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			getUserRoleResult: "MEMBER",
			wantErr:           ErrValidation,
		},
		{
			name:              "whitespace only content returns ErrValidation",
			groupID:           "group-1",
			userID:            "member-1",
			req:               &models.SendGroupMessageRequest{Content: "   "},
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			getUserRoleResult: "MEMBER",
			wantErr:           ErrValidation,
		},
		{
			name:              "create message error returns wrapped error",
			groupID:           "group-1",
			userID:            "member-1",
			req:               &models.SendGroupMessageRequest{Content: "Hello"},
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			getUserRoleResult: "MEMBER",
			createMessageErr:  errors.New("db error"),
			wantErrContains:   "failed to create message",
		},
		{
			name:                "member can send message",
			groupID:             "group-1",
			userID:              "member-1",
			req:                 &models.SendGroupMessageRequest{Content: "Hello", NotifyMembers: false},
			findGroupResult:     &models.Group{ID: "group-1", OrganizationID: "org-1"},
			getUserRoleResult:   "MEMBER",
			createMessageResult: &models.GroupMessage{ID: "msg-1", Content: "Hello", SenderID: "member-1"},
			want:                &models.GroupMessage{ID: "msg-1", Content: "Hello", SenderID: "member-1"},
		},
		{
			name:                "owner can send message",
			groupID:             "group-1",
			userID:              "owner-1",
			req:                 &models.SendGroupMessageRequest{Content: "Announcement", NotifyMembers: true},
			findGroupResult:     &models.Group{ID: "group-1", OrganizationID: "org-1"},
			getUserRoleResult:   "OWNER",
			createMessageResult: &models.GroupMessage{ID: "msg-1", Content: "Announcement", NotifyMembers: true},
			want:                &models.GroupMessage{ID: "msg-1", Content: "Announcement", NotifyMembers: true},
		},
		{
			name:                "leader can send message",
			groupID:             "group-1",
			userID:              "leader-1",
			req:                 &models.SendGroupMessageRequest{Content: "Update"},
			findGroupResult:     &models.Group{ID: "group-1", OrganizationID: "org-1"},
			getUserRoleResult:   "LEADER",
			createMessageResult: &models.GroupMessage{ID: "msg-1", Content: "Update"},
			want:                &models.GroupMessage{ID: "msg-1", Content: "Update"},
		},
		{
			name:                "content with whitespace is trimmed",
			groupID:             "group-1",
			userID:              "member-1",
			req:                 &models.SendGroupMessageRequest{Content: "  Hello World  "},
			findGroupResult:     &models.Group{ID: "group-1", OrganizationID: "org-1"},
			getUserRoleResult:   "MEMBER",
			createMessageResult: &models.GroupMessage{ID: "msg-1", Content: "Hello World"},
			want:                &models.GroupMessage{ID: "msg-1", Content: "Hello World"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messageRepo := &mockGroupMessageRepo{
				createFn: func(ctx context.Context, groupID, senderID, content string, notifyMembers bool, formID *string) (*models.GroupMessage, error) {
					if tt.createMessageErr != nil {
						return nil, tt.createMessageErr
					}
					return tt.createMessageResult, nil
				},
			}

			groupRepo := &mockGroupRepoForMsg{
				findByIDFn: func(ctx context.Context, groupID string) (*models.Group, error) {
					if tt.findGroupErr != nil {
						return nil, tt.findGroupErr
					}
					return tt.findGroupResult, nil
				},
			}

			groupService := &mockGroupServiceForMsg{
				getUserRoleFn: func(ctx context.Context, groupID, userID string) (string, error) {
					if tt.getUserRoleErr != nil {
						return "", tt.getUserRoleErr
					}
					return tt.getUserRoleResult, nil
				},
			}

			orgRepo := &mockOrgRepoForMsg{}
			userRepo := &mockUserRepoForMsg{}

			svc := NewGroupMessageService(messageRepo, groupRepo, groupService, orgRepo, userRepo, nil, nil)

			got, err := svc.Send(context.Background(), tt.groupID, tt.userID, tt.req)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if tt.wantErrContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("Send() error = %v, wantErrContains %v", err, tt.wantErrContains)
				}
				return
			}
			if err != nil {
				t.Errorf("Send() unexpected error = %v", err)
				return
			}
			if got.ID != tt.want.ID {
				t.Errorf("Send() got ID = %v, want %v", got.ID, tt.want.ID)
			}
		})
	}
}

// Test that content is properly trimmed before being passed to createMessage
func TestGroupMessageServiceImpl_Send_ContentTrimming(t *testing.T) {
	var capturedContent string

	messageRepo := &mockGroupMessageRepo{
		createFn: func(ctx context.Context, groupID, senderID, content string, notifyMembers bool, formID *string) (*models.GroupMessage, error) {
			capturedContent = content
			return &models.GroupMessage{ID: "msg-1", Content: content}, nil
		},
	}

	groupRepo := &mockGroupRepoForMsg{
		findByIDFn: func(ctx context.Context, groupID string) (*models.Group, error) {
			return &models.Group{ID: "group-1", OrganizationID: "org-1"}, nil
		},
	}

	groupService := &mockGroupServiceForMsg{
		getUserRoleFn: func(ctx context.Context, groupID, userID string) (string, error) {
			return "MEMBER", nil
		},
	}

	orgRepo := &mockOrgRepoForMsg{}
	userRepo := &mockUserRepoForMsg{}

	svc := NewGroupMessageService(messageRepo, groupRepo, groupService, orgRepo, userRepo, nil, nil)

	_, err := svc.Send(context.Background(), "group-1", "member-1", &models.SendGroupMessageRequest{
		Content: "  Hello World  ",
	})
	if err != nil {
		t.Fatalf("Send() unexpected error = %v", err)
	}

	if capturedContent != "Hello World" {
		t.Errorf("Send() content = %q, want %q", capturedContent, "Hello World")
	}
}

// ========== Delete Tests ==========

func TestGroupMessageServiceImpl_Delete(t *testing.T) {
	tests := []struct {
		name              string
		groupID           string
		messageID         string
		userID            string
		findMessageResult *models.GroupMessage
		findMessageErr    error
		findGroupResult   *models.Group
		findGroupErr      error
		getUserRoleResult string
		getUserRoleErr    error
		isSiteAdminResult bool
		isSiteAdminErr    error
		isOrgAdminResult  bool
		isOrgAdminErr     error
		deleteMessageErr  error
		wantErr           error
		wantErrContains   string
	}{
		{
			name:            "find message error returns wrapped error",
			groupID:         "group-1",
			messageID:       "msg-1",
			userID:          "user-1",
			findMessageErr:  errors.New("db error"),
			wantErrContains: "failed to find message",
		},
		{
			name:              "message not found returns ErrMessageNotFound",
			groupID:           "group-1",
			messageID:         "msg-1",
			userID:            "user-1",
			findMessageResult: nil,
			wantErr:           ErrMessageNotFound,
		},
		{
			name:              "message belongs to different group returns ErrMessageNotFound",
			groupID:           "group-1",
			messageID:         "msg-1",
			userID:            "user-1",
			findMessageResult: &models.GroupMessage{ID: "msg-1", GroupID: "group-2", SenderID: "other-user"},
			wantErr:           ErrMessageNotFound,
		},
		{
			name:              "find group error returns wrapped error",
			groupID:           "group-1",
			messageID:         "msg-1",
			userID:            "user-1",
			findMessageResult: &models.GroupMessage{ID: "msg-1", GroupID: "group-1", SenderID: "other-user"},
			findGroupErr:      errors.New("db error"),
			wantErrContains:   "failed to find group",
		},
		{
			name:              "group not found returns ErrGroupNotFound",
			groupID:           "group-1",
			messageID:         "msg-1",
			userID:            "user-1",
			findMessageResult: &models.GroupMessage{ID: "msg-1", GroupID: "group-1", SenderID: "other-user"},
			findGroupResult:   nil,
			wantErr:           ErrGroupNotFound,
		},
		{
			name:              "sender can delete own message",
			groupID:           "group-1",
			messageID:         "msg-1",
			userID:            "sender-1",
			findMessageResult: &models.GroupMessage{ID: "msg-1", GroupID: "group-1", SenderID: "sender-1"},
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
		},
		{
			name:              "sender delete error returns wrapped error",
			groupID:           "group-1",
			messageID:         "msg-1",
			userID:            "sender-1",
			findMessageResult: &models.GroupMessage{ID: "msg-1", GroupID: "group-1", SenderID: "sender-1"},
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			deleteMessageErr:  errors.New("db error"),
			wantErrContains:   "failed to delete message",
		},
		{
			name:              "non-sender - get user role error returns wrapped error",
			groupID:           "group-1",
			messageID:         "msg-1",
			userID:            "other-user",
			findMessageResult: &models.GroupMessage{ID: "msg-1", GroupID: "group-1", SenderID: "sender-1"},
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			getUserRoleErr:    errors.New("db error"),
			wantErrContains:   "failed to get user role",
		},
		{
			name:              "group owner can delete any message",
			groupID:           "group-1",
			messageID:         "msg-1",
			userID:            "owner-1",
			findMessageResult: &models.GroupMessage{ID: "msg-1", GroupID: "group-1", SenderID: "sender-1"},
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			getUserRoleResult: "OWNER",
		},
		{
			name:              "owner delete error returns wrapped error",
			groupID:           "group-1",
			messageID:         "msg-1",
			userID:            "owner-1",
			findMessageResult: &models.GroupMessage{ID: "msg-1", GroupID: "group-1", SenderID: "sender-1"},
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			getUserRoleResult: "OWNER",
			deleteMessageErr:  errors.New("db error"),
			wantErrContains:   "failed to delete message",
		},
		{
			name:              "group leader can delete any message",
			groupID:           "group-1",
			messageID:         "msg-1",
			userID:            "leader-1",
			findMessageResult: &models.GroupMessage{ID: "msg-1", GroupID: "group-1", SenderID: "sender-1"},
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			getUserRoleResult: "LEADER",
		},
		{
			name:              "leader delete error returns wrapped error",
			groupID:           "group-1",
			messageID:         "msg-1",
			userID:            "leader-1",
			findMessageResult: &models.GroupMessage{ID: "msg-1", GroupID: "group-1", SenderID: "sender-1"},
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			getUserRoleResult: "LEADER",
			deleteMessageErr:  errors.New("db error"),
			wantErrContains:   "failed to delete message",
		},
		{
			name:              "regular member cannot delete others message - check site admin error",
			groupID:           "group-1",
			messageID:         "msg-1",
			userID:            "member-1",
			findMessageResult: &models.GroupMessage{ID: "msg-1", GroupID: "group-1", SenderID: "sender-1"},
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			getUserRoleResult: "MEMBER",
			isSiteAdminErr:    errors.New("db error"),
			wantErrContains:   "failed to check site admin status",
		},
		{
			name:              "site admin can delete any message",
			groupID:           "group-1",
			messageID:         "msg-1",
			userID:            "site-admin",
			findMessageResult: &models.GroupMessage{ID: "msg-1", GroupID: "group-1", SenderID: "sender-1"},
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			getUserRoleResult: "",
			isSiteAdminResult: true,
		},
		{
			name:              "site admin delete error returns wrapped error",
			groupID:           "group-1",
			messageID:         "msg-1",
			userID:            "site-admin",
			findMessageResult: &models.GroupMessage{ID: "msg-1", GroupID: "group-1", SenderID: "sender-1"},
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			getUserRoleResult: "",
			isSiteAdminResult: true,
			deleteMessageErr:  errors.New("db error"),
			wantErrContains:   "failed to delete message",
		},
		{
			name:              "non-site-admin - check org admin error returns wrapped error",
			groupID:           "group-1",
			messageID:         "msg-1",
			userID:            "user-1",
			findMessageResult: &models.GroupMessage{ID: "msg-1", GroupID: "group-1", SenderID: "sender-1"},
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			getUserRoleResult: "MEMBER",
			isSiteAdminResult: false,
			isOrgAdminErr:     errors.New("db error"),
			wantErrContains:   "failed to check org admin status",
		},
		{
			name:              "org admin can delete any message",
			groupID:           "group-1",
			messageID:         "msg-1",
			userID:            "org-admin",
			findMessageResult: &models.GroupMessage{ID: "msg-1", GroupID: "group-1", SenderID: "sender-1"},
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			getUserRoleResult: "",
			isSiteAdminResult: false,
			isOrgAdminResult:  true,
		},
		{
			name:              "org admin delete error returns wrapped error",
			groupID:           "group-1",
			messageID:         "msg-1",
			userID:            "org-admin",
			findMessageResult: &models.GroupMessage{ID: "msg-1", GroupID: "group-1", SenderID: "sender-1"},
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			getUserRoleResult: "",
			isSiteAdminResult: false,
			isOrgAdminResult:  true,
			deleteMessageErr:  errors.New("db error"),
			wantErrContains:   "failed to delete message",
		},
		{
			name:              "regular member cannot delete others message returns ErrForbidden",
			groupID:           "group-1",
			messageID:         "msg-1",
			userID:            "member-1",
			findMessageResult: &models.GroupMessage{ID: "msg-1", GroupID: "group-1", SenderID: "sender-1"},
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			getUserRoleResult: "MEMBER",
			isSiteAdminResult: false,
			isOrgAdminResult:  false,
			wantErr:           ErrForbidden,
		},
		{
			name:              "non-member non-admin returns ErrForbidden",
			groupID:           "group-1",
			messageID:         "msg-1",
			userID:            "random-user",
			findMessageResult: &models.GroupMessage{ID: "msg-1", GroupID: "group-1", SenderID: "sender-1"},
			findGroupResult:   &models.Group{ID: "group-1", OrganizationID: "org-1"},
			getUserRoleResult: "",
			isSiteAdminResult: false,
			isOrgAdminResult:  false,
			wantErr:           ErrForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messageRepo := &mockGroupMessageRepo{
				findByIDFn: func(ctx context.Context, messageID string) (*models.GroupMessage, error) {
					if tt.findMessageErr != nil {
						return nil, tt.findMessageErr
					}
					return tt.findMessageResult, nil
				},
				deleteFn: func(ctx context.Context, messageID string) error {
					return tt.deleteMessageErr
				},
			}

			groupRepo := &mockGroupRepoForMsg{
				findByIDFn: func(ctx context.Context, groupID string) (*models.Group, error) {
					if tt.findGroupErr != nil {
						return nil, tt.findGroupErr
					}
					return tt.findGroupResult, nil
				},
			}

			groupService := &mockGroupServiceForMsg{
				getUserRoleFn: func(ctx context.Context, groupID, userID string) (string, error) {
					if tt.getUserRoleErr != nil {
						return "", tt.getUserRoleErr
					}
					return tt.getUserRoleResult, nil
				},
			}

			orgRepo := &mockOrgRepoForMsg{
				isAdminFn: func(ctx context.Context, orgID, userID string) (bool, error) {
					if tt.isOrgAdminErr != nil {
						return false, tt.isOrgAdminErr
					}
					return tt.isOrgAdminResult, nil
				},
			}

			userRepo := &mockUserRepoForMsg{
				isSiteAdminFn: func(ctx context.Context, userID string) (bool, error) {
					if tt.isSiteAdminErr != nil {
						return false, tt.isSiteAdminErr
					}
					return tt.isSiteAdminResult, nil
				},
			}

			svc := NewGroupMessageService(messageRepo, groupRepo, groupService, orgRepo, userRepo, nil, nil)

			err := svc.Delete(context.Background(), tt.groupID, tt.messageID, tt.userID)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if tt.wantErrContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("Delete() error = %v, wantErrContains %v", err, tt.wantErrContains)
				}
				return
			}
			if err != nil {
				t.Errorf("Delete() unexpected error = %v", err)
			}
		})
	}
}

// ========== GetByID Tests ==========

func TestGroupMessageServiceImpl_GetByID(t *testing.T) {
	tests := []struct {
		name              string
		messageID         string
		findMessageResult *models.GroupMessage
		findMessageErr    error
		want              *models.GroupMessage
		wantErr           error
		wantErrContains   string
	}{
		{
			name:      "empty messageID returns ErrMessageNotFound",
			messageID: "",
			wantErr:   ErrMessageNotFound,
		},
		{
			name:            "find message error returns wrapped error",
			messageID:       "msg-1",
			findMessageErr:  errors.New("db error"),
			wantErrContains: "failed to get message by ID",
		},
		{
			name:              "message not found returns ErrMessageNotFound",
			messageID:         "msg-1",
			findMessageResult: nil,
			wantErr:           ErrMessageNotFound,
		},
		{
			name:              "message found returns the message",
			messageID:         "msg-1",
			findMessageResult: &models.GroupMessage{ID: "msg-1", Content: "Hello", GroupID: "group-1"},
			want:              &models.GroupMessage{ID: "msg-1", Content: "Hello", GroupID: "group-1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messageRepo := &mockGroupMessageRepo{
				findByIDFn: func(ctx context.Context, messageID string) (*models.GroupMessage, error) {
					if tt.findMessageErr != nil {
						return nil, tt.findMessageErr
					}
					return tt.findMessageResult, nil
				},
			}

			groupRepo := &mockGroupRepoForMsg{}
			groupService := &mockGroupServiceForMsg{}
			orgRepo := &mockOrgRepoForMsg{}
			userRepo := &mockUserRepoForMsg{}

			svc := NewGroupMessageService(messageRepo, groupRepo, groupService, orgRepo, userRepo, nil, nil)

			got, err := svc.GetByID(context.Background(), tt.messageID)

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

// ========== Additional Edge Case Tests ==========

// TestGroupMessageServiceImpl_List_ContextPropagation verifies that context is properly propagated
func TestGroupMessageServiceImpl_List_ContextPropagation(t *testing.T) {
	type contextKey string
	const testKey contextKey = "test-key"
	const testValue = "test-value"

	var capturedCtx context.Context

	messageRepo := &mockGroupMessageRepo{
		findByGroupIDFn: func(ctx context.Context, groupID string) ([]*models.GroupMessage, error) {
			capturedCtx = ctx
			return []*models.GroupMessage{}, nil
		},
	}

	groupRepo := &mockGroupRepoForMsg{
		findByIDFn: func(ctx context.Context, groupID string) (*models.Group, error) {
			return &models.Group{ID: "group-1", OrganizationID: "org-1"}, nil
		},
	}

	groupService := &mockGroupServiceForMsg{
		getUserRoleFn: func(ctx context.Context, groupID, userID string) (string, error) {
			return "MEMBER", nil
		},
	}

	orgRepo := &mockOrgRepoForMsg{}
	userRepo := &mockUserRepoForMsg{}

	svc := NewGroupMessageService(messageRepo, groupRepo, groupService, orgRepo, userRepo, nil, nil)

	ctx := context.WithValue(context.Background(), testKey, testValue)
	_, err := svc.List(ctx, "group-1", "member-1")
	if err != nil {
		t.Fatalf("List() unexpected error = %v", err)
	}

	if capturedCtx == nil {
		t.Fatal("context was not propagated")
	}

	if capturedCtx.Value(testKey) != testValue {
		t.Errorf("context value = %v, want %v", capturedCtx.Value(testKey), testValue)
	}
}

// TestGroupMessageServiceImpl_Delete_OrgIDPropagation verifies the correct org ID is used for admin check
func TestGroupMessageServiceImpl_Delete_OrgIDPropagation(t *testing.T) {
	var capturedOrgID string

	messageRepo := &mockGroupMessageRepo{
		findByIDFn: func(ctx context.Context, messageID string) (*models.GroupMessage, error) {
			return &models.GroupMessage{ID: "msg-1", GroupID: "group-1", SenderID: "other-user"}, nil
		},
		deleteFn: func(ctx context.Context, messageID string) error {
			return nil
		},
	}

	groupRepo := &mockGroupRepoForMsg{
		findByIDFn: func(ctx context.Context, groupID string) (*models.Group, error) {
			return &models.Group{ID: "group-1", OrganizationID: "org-123"}, nil
		},
	}

	groupService := &mockGroupServiceForMsg{
		getUserRoleFn: func(ctx context.Context, groupID, userID string) (string, error) {
			return "", nil
		},
	}

	orgRepo := &mockOrgRepoForMsg{
		isAdminFn: func(ctx context.Context, orgID, userID string) (bool, error) {
			capturedOrgID = orgID
			return true, nil
		},
	}

	userRepo := &mockUserRepoForMsg{
		isSiteAdminFn: func(ctx context.Context, userID string) (bool, error) {
			return false, nil
		},
	}

	svc := NewGroupMessageService(messageRepo, groupRepo, groupService, orgRepo, userRepo, nil, nil)

	err := svc.Delete(context.Background(), "group-1", "msg-1", "org-admin")
	if err != nil {
		t.Fatalf("Delete() unexpected error = %v", err)
	}

	if capturedOrgID != "org-123" {
		t.Errorf("Delete() org ID = %v, want %v", capturedOrgID, "org-123")
	}
}
