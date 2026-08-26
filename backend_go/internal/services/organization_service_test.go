package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"menunderfire/internal/models"
)

// =============================================================================
// Mock Repositories
// =============================================================================

type mockOrganizationRepository struct {
	findByID     func(ctx context.Context, orgID string) (*models.Organization, error)
	findByUserID func(ctx context.Context, userID string) ([]*models.Organization, error)
	findAll      func(ctx context.Context) ([]*models.Organization, error)
	create       func(ctx context.Context, name string, description *string, createdByID string) (*models.Organization, error)
	update       func(ctx context.Context, orgID string, name string, description *string) (*models.Organization, error)
	findAdmins   func(ctx context.Context, orgID string) ([]*models.OrganizationAdmin, error)
	findAdmin    func(ctx context.Context, orgID, userID string) (*models.OrganizationAdmin, error)
	isMember     func(ctx context.Context, orgID, userID string) (bool, error)
}

func (m *mockOrganizationRepository) FindByID(ctx context.Context, orgID string) (*models.Organization, error) {
	if m.findByID != nil {
		return m.findByID(ctx, orgID)
	}
	return nil, errors.New("FindByID not implemented")
}

func (m *mockOrganizationRepository) FindByUserID(ctx context.Context, userID string) ([]*models.Organization, error) {
	if m.findByUserID != nil {
		return m.findByUserID(ctx, userID)
	}
	return nil, errors.New("FindByUserID not implemented")
}

func (m *mockOrganizationRepository) FindAll(ctx context.Context) ([]*models.Organization, error) {
	if m.findAll != nil {
		return m.findAll(ctx)
	}
	return nil, errors.New("FindAll not implemented")
}

func (m *mockOrganizationRepository) Create(ctx context.Context, name string, description *string, createdByID string) (*models.Organization, error) {
	if m.create != nil {
		return m.create(ctx, name, description, createdByID)
	}
	return nil, errors.New("Create not implemented")
}

func (m *mockOrganizationRepository) Update(ctx context.Context, orgID string, name string, description *string) (*models.Organization, error) {
	if m.update != nil {
		return m.update(ctx, orgID, name, description)
	}
	return nil, errors.New("Update not implemented")
}

func (m *mockOrganizationRepository) FindAdmins(ctx context.Context, orgID string) ([]*models.OrganizationAdmin, error) {
	if m.findAdmins != nil {
		return m.findAdmins(ctx, orgID)
	}
	return nil, errors.New("FindAdmins not implemented")
}

func (m *mockOrganizationRepository) FindAdmin(ctx context.Context, orgID, userID string) (*models.OrganizationAdmin, error) {
	if m.findAdmin != nil {
		return m.findAdmin(ctx, orgID, userID)
	}
	return nil, errors.New("FindAdmin not implemented")
}

func (m *mockOrganizationRepository) IsMember(ctx context.Context, orgID, userID string) (bool, error) {
	if m.isMember != nil {
		return m.isMember(ctx, orgID, userID)
	}
	return false, errors.New("IsMember not implemented")
}

func (m *mockOrganizationRepository) IsAdmin(ctx context.Context, orgID, userID string) (bool, error) {
	return false, errors.New("IsAdmin not implemented")
}

func (m *mockOrganizationRepository) AddAdmin(ctx context.Context, orgID, userID string) error {
	return errors.New("AddAdmin not implemented")
}

func (m *mockOrganizationRepository) Count(ctx context.Context) (int64, error) {
	return 0, errors.New("Count not implemented")
}

type mockGroupRepoForOrg struct {
	findByID             func(ctx context.Context, groupID string) (*models.Group, error)
	findMember           func(ctx context.Context, groupID, userID string) (*models.GroupMember, error)
	findByOrganizationID func(ctx context.Context, orgID string) ([]*models.Group, error)
}

func (m *mockGroupRepoForOrg) FindByID(ctx context.Context, groupID string) (*models.Group, error) {
	if m.findByID != nil {
		return m.findByID(ctx, groupID)
	}
	return nil, errors.New("FindByID not implemented")
}

func (m *mockGroupRepoForOrg) FindByInviteCode(ctx context.Context, inviteCode string) (*models.Group, error) {
	return nil, errors.New("FindByInviteCode not implemented")
}

func (m *mockGroupRepoForOrg) FindMember(ctx context.Context, groupID, userID string) (*models.GroupMember, error) {
	if m.findMember != nil {
		return m.findMember(ctx, groupID, userID)
	}
	return nil, errors.New("FindMember not implemented")
}

func (m *mockGroupRepoForOrg) FindByOrganizationID(ctx context.Context, orgID string) ([]*models.Group, error) {
	if m.findByOrganizationID != nil {
		return m.findByOrganizationID(ctx, orgID)
	}
	return nil, errors.New("FindByOrganizationID not implemented")
}

func (m *mockGroupRepoForOrg) FindByUserID(ctx context.Context, userID string) ([]*models.Group, error) {
	return nil, errors.New("FindByUserID not implemented")
}

func (m *mockGroupRepoForOrg) Create(ctx context.Context, name string, description *string, orgID string, inviteCode string, createdBy string) (*models.Group, error) {
	return nil, errors.New("Create not implemented")
}

func (m *mockGroupRepoForOrg) GenerateInviteCode() string {
	return ""
}

func (m *mockGroupRepoForOrg) FindMembers(ctx context.Context, groupID string) ([]*models.GroupMember, error) {
	return nil, errors.New("FindMembers not implemented")
}

func (m *mockGroupRepoForOrg) CountMembers(ctx context.Context, groupID string) (int, error) {
	return 0, errors.New("CountMembers not implemented")
}

func (m *mockGroupRepoForOrg) AddMember(ctx context.Context, groupID, userID, role string) (*models.GroupMember, error) {
	return nil, errors.New("AddMember not implemented")
}

func (m *mockGroupRepoForOrg) RemoveMember(ctx context.Context, groupID, userID string) error {
	return errors.New("RemoveMember not implemented")
}

func (m *mockGroupRepoForOrg) Count(ctx context.Context) (int64, error) {
	return 0, errors.New("Count not implemented")
}

func (m *mockGroupRepoForOrg) CountByOrganizationIDs(ctx context.Context, orgIDs []string) (int64, error) {
	return 0, errors.New("CountByOrganizationIDs not implemented")
}

func (m *mockGroupRepoForOrg) UpdateMemberRole(ctx context.Context, groupID, userID, role string) error {
	return nil
}

func (m *mockGroupRepoForOrg) UpdateSettings(ctx context.Context, groupID string, requirePostApproval, allowAnonymousPosts bool) error {
	return nil
}

type mockUserRepoForOrg struct {
	findByID                 func(ctx context.Context, userID string) (*models.User, error)
	findByExternalID         func(ctx context.Context, externalID string) (*models.User, error)
	update                   func(ctx context.Context, userID, displayName string) (*models.User, error)
	isSiteAdmin              func(ctx context.Context, userID string) (bool, error)
	findAdminOrganizationIDs func(ctx context.Context, userID string) ([]string, error)
	findOwnedGroupIDs        func(ctx context.Context, userID string) ([]string, error)
	findMemberGroupIDs       func(ctx context.Context, userID string) ([]string, error)
}

func (m *mockUserRepoForOrg) FindByID(ctx context.Context, userID string) (*models.User, error) {
	if m.findByID != nil {
		return m.findByID(ctx, userID)
	}
	return nil, errors.New("FindByID not implemented")
}

func (m *mockUserRepoForOrg) FindByExternalID(ctx context.Context, externalID string) (*models.User, error) {
	if m.findByExternalID != nil {
		return m.findByExternalID(ctx, externalID)
	}
	return nil, errors.New("FindByExternalID not implemented")
}

func (m *mockUserRepoForOrg) Update(ctx context.Context, userID, displayName string) (*models.User, error) {
	if m.update != nil {
		return m.update(ctx, userID, displayName)
	}
	return nil, errors.New("Update not implemented")
}

func (m *mockUserRepoForOrg) UpdateExternalID(ctx context.Context, userID, externalID string) error {
	return errors.New("UpdateExternalID not implemented")
}

func (m *mockUserRepoForOrg) UpdateInvitationInfo(ctx context.Context, userID, invitedByID, invitationID string) error {
	return errors.New("UpdateInvitationInfo not implemented")
}

func (m *mockUserRepoForOrg) IsSiteAdmin(ctx context.Context, userID string) (bool, error) {
	if m.isSiteAdmin != nil {
		return m.isSiteAdmin(ctx, userID)
	}
	return false, errors.New("IsSiteAdmin not implemented")
}

func (m *mockUserRepoForOrg) FindAdminOrganizationIDs(ctx context.Context, userID string) ([]string, error) {
	if m.findAdminOrganizationIDs != nil {
		return m.findAdminOrganizationIDs(ctx, userID)
	}
	return nil, errors.New("FindAdminOrganizationIDs not implemented")
}

func (m *mockUserRepoForOrg) FindOwnedGroupIDs(ctx context.Context, userID string) ([]string, error) {
	if m.findOwnedGroupIDs != nil {
		return m.findOwnedGroupIDs(ctx, userID)
	}
	return nil, errors.New("FindOwnedGroupIDs not implemented")
}

func (m *mockUserRepoForOrg) FindMemberGroupIDs(ctx context.Context, userID string) ([]string, error) {
	if m.findMemberGroupIDs != nil {
		return m.findMemberGroupIDs(ctx, userID)
	}
	return nil, errors.New("FindMemberGroupIDs not implemented")
}

func (m *mockUserRepoForOrg) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	return nil, errors.New("FindByEmail not implemented")
}

func (m *mockUserRepoForOrg) Count(ctx context.Context) (int64, error) {
	return 0, errors.New("Count not implemented")
}

func (m *mockUserRepoForOrg) Create(ctx context.Context, email, displayName, externalID string) (*models.User, error) {
	return nil, errors.New("Create not implemented")
}

func (m *mockUserRepoForOrg) CreateAsSiteAdmin(ctx context.Context, email, displayName, externalID string) (*models.User, error) {
	return nil, errors.New("CreateAsSiteAdmin not implemented")
}

func (m *mockUserRepoForOrg) RecordAgreementAcceptance(ctx context.Context, userID, version, signature, ipAddress, userAgent string) error {
	return errors.New("not implemented")
}

// =============================================================================
// Helper Functions
// =============================================================================

// Helper to create a test organization
func testOrg(id, name string) *models.Organization {
	return &models.Organization{
		ID:          id,
		Name:        name,
		Description: nil,
		CreatedByID: "creator-1",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// Helper to create a test organization admin
func testAdmin(id, userID, orgID string) *models.OrganizationAdmin {
	return &models.OrganizationAdmin{
		ID:             id,
		UserID:         userID,
		OrganizationID: orgID,
		JoinedAt:       time.Now(),
	}
}

// Helper to create a test group
func testGroup(id, name, orgID string) *models.Group {
	return &models.Group{
		ID:             id,
		Name:           name,
		OrganizationID: orgID,
		InviteCode:     "ABC123",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

// =============================================================================
// List Tests
// =============================================================================

func TestOrganizationService_List_EmptyUserID(t *testing.T) {
	svc := NewOrganizationService(
		&mockOrganizationRepository{},
		&mockGroupRepoForOrg{},
		&mockUserRepoForOrg{},
	)

	_, err := svc.List(context.Background(), "")
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestOrganizationService_List_WhitespaceOnlyUserID(t *testing.T) {
	svc := NewOrganizationService(
		&mockOrganizationRepository{},
		&mockGroupRepoForOrg{},
		&mockUserRepoForOrg{},
	)

	_, err := svc.List(context.Background(), "   ")
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestOrganizationService_List_RepositoryError(t *testing.T) {
	repoErr := errors.New("database error")
	orgRepo := &mockOrganizationRepository{
		findByUserID: func(ctx context.Context, userID string) ([]*models.Organization, error) {
			return nil, repoErr
		},
	}
	svc := NewOrganizationService(orgRepo, &mockGroupRepoForOrg{}, &mockUserRepoForOrg{})

	_, err := svc.List(context.Background(), "user-1")
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repository error, got %v", err)
	}
}

func TestOrganizationService_List_NoOrganizations(t *testing.T) {
	orgRepo := &mockOrganizationRepository{
		findByUserID: func(ctx context.Context, userID string) ([]*models.Organization, error) {
			return []*models.Organization{}, nil
		},
	}
	svc := NewOrganizationService(orgRepo, &mockGroupRepoForOrg{}, &mockUserRepoForOrg{})

	orgs, err := svc.List(context.Background(), "user-1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(orgs) != 0 {
		t.Errorf("expected empty slice, got %d organizations", len(orgs))
	}
}

func TestOrganizationService_List_Success(t *testing.T) {
	expectedOrgs := []*models.Organization{
		testOrg("org-1", "Org 1"),
		testOrg("org-2", "Org 2"),
	}
	orgRepo := &mockOrganizationRepository{
		findByUserID: func(ctx context.Context, userID string) ([]*models.Organization, error) {
			if userID != "user-1" {
				t.Errorf("expected userID 'user-1', got '%s'", userID)
			}
			return expectedOrgs, nil
		},
	}
	svc := NewOrganizationService(orgRepo, &mockGroupRepoForOrg{}, &mockUserRepoForOrg{})

	orgs, err := svc.List(context.Background(), "user-1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(orgs) != 2 {
		t.Errorf("expected 2 organizations, got %d", len(orgs))
	}
}

func TestOrganizationService_List_TrimsUserID(t *testing.T) {
	orgRepo := &mockOrganizationRepository{
		findByUserID: func(ctx context.Context, userID string) ([]*models.Organization, error) {
			if userID != "user-1" {
				t.Errorf("expected trimmed userID 'user-1', got '%s'", userID)
			}
			return []*models.Organization{}, nil
		},
	}
	svc := NewOrganizationService(orgRepo, &mockGroupRepoForOrg{}, &mockUserRepoForOrg{})

	_, err := svc.List(context.Background(), "  user-1  ")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// =============================================================================
// ListAll Tests
// =============================================================================

func TestOrganizationService_ListAll_RepositoryError(t *testing.T) {
	repoErr := errors.New("database error")
	orgRepo := &mockOrganizationRepository{
		findAll: func(ctx context.Context) ([]*models.Organization, error) {
			return nil, repoErr
		},
	}
	svc := NewOrganizationService(orgRepo, &mockGroupRepoForOrg{}, &mockUserRepoForOrg{})

	_, err := svc.ListAll(context.Background())
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repository error, got %v", err)
	}
}

func TestOrganizationService_ListAll_NoOrganizations(t *testing.T) {
	orgRepo := &mockOrganizationRepository{
		findAll: func(ctx context.Context) ([]*models.Organization, error) {
			return []*models.Organization{}, nil
		},
	}
	svc := NewOrganizationService(orgRepo, &mockGroupRepoForOrg{}, &mockUserRepoForOrg{})

	orgs, err := svc.ListAll(context.Background())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(orgs) != 0 {
		t.Errorf("expected empty slice, got %d organizations", len(orgs))
	}
}

func TestOrganizationService_ListAll_Success(t *testing.T) {
	expectedOrgs := []*models.Organization{
		testOrg("org-1", "Org 1"),
		testOrg("org-2", "Org 2"),
		testOrg("org-3", "Org 3"),
	}
	orgRepo := &mockOrganizationRepository{
		findAll: func(ctx context.Context) ([]*models.Organization, error) {
			return expectedOrgs, nil
		},
	}
	svc := NewOrganizationService(orgRepo, &mockGroupRepoForOrg{}, &mockUserRepoForOrg{})

	orgs, err := svc.ListAll(context.Background())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(orgs) != 3 {
		t.Errorf("expected 3 organizations, got %d", len(orgs))
	}
}

// =============================================================================
// Create Tests
// =============================================================================

func TestOrganizationService_Create_NilRequest(t *testing.T) {
	svc := NewOrganizationService(
		&mockOrganizationRepository{},
		&mockGroupRepoForOrg{},
		&mockUserRepoForOrg{},
	)

	_, err := svc.Create(context.Background(), "user-1", nil)
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestOrganizationService_Create_EmptyName(t *testing.T) {
	svc := NewOrganizationService(
		&mockOrganizationRepository{},
		&mockGroupRepoForOrg{},
		&mockUserRepoForOrg{},
	)

	req := &models.CreateOrganizationRequest{Name: ""}
	_, err := svc.Create(context.Background(), "user-1", req)
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestOrganizationService_Create_WhitespaceOnlyName(t *testing.T) {
	svc := NewOrganizationService(
		&mockOrganizationRepository{},
		&mockGroupRepoForOrg{},
		&mockUserRepoForOrg{},
	)

	req := &models.CreateOrganizationRequest{Name: "   "}
	_, err := svc.Create(context.Background(), "user-1", req)
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestOrganizationService_Create_EmptyUserID(t *testing.T) {
	svc := NewOrganizationService(
		&mockOrganizationRepository{},
		&mockGroupRepoForOrg{},
		&mockUserRepoForOrg{},
	)

	req := &models.CreateOrganizationRequest{Name: "My Org"}
	_, err := svc.Create(context.Background(), "", req)
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestOrganizationService_Create_WhitespaceOnlyUserID(t *testing.T) {
	svc := NewOrganizationService(
		&mockOrganizationRepository{},
		&mockGroupRepoForOrg{},
		&mockUserRepoForOrg{},
	)

	req := &models.CreateOrganizationRequest{Name: "My Org"}
	_, err := svc.Create(context.Background(), "   ", req)
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestOrganizationService_Create_IsSiteAdminError(t *testing.T) {
	repoErr := errors.New("database error")
	userRepo := &mockUserRepoForOrg{
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			return false, repoErr
		},
	}
	svc := NewOrganizationService(&mockOrganizationRepository{}, &mockGroupRepoForOrg{}, userRepo)

	req := &models.CreateOrganizationRequest{Name: "My Org"}
	_, err := svc.Create(context.Background(), "user-1", req)
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repository error, got %v", err)
	}
}

func TestOrganizationService_Create_NotSiteAdmin(t *testing.T) {
	userRepo := &mockUserRepoForOrg{
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			return false, nil
		},
	}
	svc := NewOrganizationService(&mockOrganizationRepository{}, &mockGroupRepoForOrg{}, userRepo)

	req := &models.CreateOrganizationRequest{Name: "My Org"}
	_, err := svc.Create(context.Background(), "user-1", req)
	if !errors.Is(err, ErrForbidden) {
		t.Errorf("expected ErrForbidden, got %v", err)
	}
}

func TestOrganizationService_Create_CreateOrgError(t *testing.T) {
	repoErr := errors.New("create error")
	orgRepo := &mockOrganizationRepository{
		create: func(ctx context.Context, name string, description *string, createdByID string) (*models.Organization, error) {
			return nil, repoErr
		},
	}
	userRepo := &mockUserRepoForOrg{
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			return true, nil
		},
	}
	svc := NewOrganizationService(orgRepo, &mockGroupRepoForOrg{}, userRepo)

	req := &models.CreateOrganizationRequest{Name: "My Org"}
	_, err := svc.Create(context.Background(), "user-1", req)
	if !errors.Is(err, repoErr) {
		t.Errorf("expected create error, got %v", err)
	}
}

func TestOrganizationService_Create_Success(t *testing.T) {
	expectedOrg := testOrg("org-1", "My Org")
	orgRepo := &mockOrganizationRepository{
		create: func(ctx context.Context, name string, description *string, createdByID string) (*models.Organization, error) {
			if name != "My Org" {
				t.Errorf("expected name 'My Org', got '%s'", name)
			}
			if createdByID != "user-1" {
				t.Errorf("expected createdByID 'user-1', got '%s'", createdByID)
			}
			return expectedOrg, nil
		},
	}
	userRepo := &mockUserRepoForOrg{
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			return true, nil
		},
	}
	svc := NewOrganizationService(orgRepo, &mockGroupRepoForOrg{}, userRepo)

	req := &models.CreateOrganizationRequest{Name: "My Org"}
	org, err := svc.Create(context.Background(), "user-1", req)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if org.ID != "org-1" {
		t.Errorf("expected org ID 'org-1', got '%s'", org.ID)
	}
}

func TestOrganizationService_Create_TrimsNameAndUserID(t *testing.T) {
	orgRepo := &mockOrganizationRepository{
		create: func(ctx context.Context, name string, description *string, createdByID string) (*models.Organization, error) {
			if name != "My Org" {
				t.Errorf("expected trimmed name 'My Org', got '%s'", name)
			}
			if createdByID != "user-1" {
				t.Errorf("expected trimmed createdByID 'user-1', got '%s'", createdByID)
			}
			return testOrg("org-1", name), nil
		},
	}
	userRepo := &mockUserRepoForOrg{
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			if userID != "user-1" {
				t.Errorf("expected trimmed userID 'user-1', got '%s'", userID)
			}
			return true, nil
		},
	}
	svc := NewOrganizationService(orgRepo, &mockGroupRepoForOrg{}, userRepo)

	req := &models.CreateOrganizationRequest{Name: "  My Org  "}
	_, err := svc.Create(context.Background(), "  user-1  ", req)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestOrganizationService_Create_WithDescription(t *testing.T) {
	desc := "A test organization"
	var capturedDesc *string
	orgRepo := &mockOrganizationRepository{
		create: func(ctx context.Context, name string, description *string, createdByID string) (*models.Organization, error) {
			capturedDesc = description
			return testOrg("org-1", name), nil
		},
	}
	userRepo := &mockUserRepoForOrg{
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			return true, nil
		},
	}
	svc := NewOrganizationService(orgRepo, &mockGroupRepoForOrg{}, userRepo)

	req := &models.CreateOrganizationRequest{Name: "My Org", Description: &desc}
	_, err := svc.Create(context.Background(), "user-1", req)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if capturedDesc == nil || *capturedDesc != desc {
		t.Errorf("expected description '%s', got %v", desc, capturedDesc)
	}
}

// =============================================================================
// GetByID Tests
// =============================================================================

func TestOrganizationService_GetByID_EmptyOrgID(t *testing.T) {
	svc := NewOrganizationService(
		&mockOrganizationRepository{},
		&mockGroupRepoForOrg{},
		&mockUserRepoForOrg{},
	)

	_, err := svc.GetByID(context.Background(), "")
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestOrganizationService_GetByID_WhitespaceOnlyOrgID(t *testing.T) {
	svc := NewOrganizationService(
		&mockOrganizationRepository{},
		&mockGroupRepoForOrg{},
		&mockUserRepoForOrg{},
	)

	_, err := svc.GetByID(context.Background(), "   ")
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestOrganizationService_GetByID_NotFound(t *testing.T) {
	orgRepo := &mockOrganizationRepository{
		findByID: func(ctx context.Context, orgID string) (*models.Organization, error) {
			return nil, errors.New("not found")
		},
	}
	svc := NewOrganizationService(orgRepo, &mockGroupRepoForOrg{}, &mockUserRepoForOrg{})

	_, err := svc.GetByID(context.Background(), "org-1")
	if !errors.Is(err, ErrOrganizationNotFound) {
		t.Errorf("expected ErrOrganizationNotFound, got %v", err)
	}
}

func TestOrganizationService_GetByID_Success(t *testing.T) {
	expectedOrg := testOrg("org-1", "My Org")
	orgRepo := &mockOrganizationRepository{
		findByID: func(ctx context.Context, orgID string) (*models.Organization, error) {
			if orgID != "org-1" {
				t.Errorf("expected orgID 'org-1', got '%s'", orgID)
			}
			return expectedOrg, nil
		},
	}
	svc := NewOrganizationService(orgRepo, &mockGroupRepoForOrg{}, &mockUserRepoForOrg{})

	org, err := svc.GetByID(context.Background(), "org-1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if org.ID != "org-1" {
		t.Errorf("expected org ID 'org-1', got '%s'", org.ID)
	}
}

func TestOrganizationService_GetByID_TrimsOrgID(t *testing.T) {
	orgRepo := &mockOrganizationRepository{
		findByID: func(ctx context.Context, orgID string) (*models.Organization, error) {
			if orgID != "org-1" {
				t.Errorf("expected trimmed orgID 'org-1', got '%s'", orgID)
			}
			return testOrg("org-1", "My Org"), nil
		},
	}
	svc := NewOrganizationService(orgRepo, &mockGroupRepoForOrg{}, &mockUserRepoForOrg{})

	_, err := svc.GetByID(context.Background(), "  org-1  ")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// =============================================================================
// Update Tests
// =============================================================================

func TestOrganizationService_Update_EmptyOrgID(t *testing.T) {
	svc := NewOrganizationService(
		&mockOrganizationRepository{},
		&mockGroupRepoForOrg{},
		&mockUserRepoForOrg{},
	)

	req := &models.UpdateOrganizationRequest{Name: "Updated Org"}
	_, err := svc.Update(context.Background(), "", "user-1", req)
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestOrganizationService_Update_WhitespaceOnlyOrgID(t *testing.T) {
	svc := NewOrganizationService(
		&mockOrganizationRepository{},
		&mockGroupRepoForOrg{},
		&mockUserRepoForOrg{},
	)

	req := &models.UpdateOrganizationRequest{Name: "Updated Org"}
	_, err := svc.Update(context.Background(), "   ", "user-1", req)
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestOrganizationService_Update_EmptyUserID(t *testing.T) {
	svc := NewOrganizationService(
		&mockOrganizationRepository{},
		&mockGroupRepoForOrg{},
		&mockUserRepoForOrg{},
	)

	req := &models.UpdateOrganizationRequest{Name: "Updated Org"}
	_, err := svc.Update(context.Background(), "org-1", "", req)
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestOrganizationService_Update_WhitespaceOnlyUserID(t *testing.T) {
	svc := NewOrganizationService(
		&mockOrganizationRepository{},
		&mockGroupRepoForOrg{},
		&mockUserRepoForOrg{},
	)

	req := &models.UpdateOrganizationRequest{Name: "Updated Org"}
	_, err := svc.Update(context.Background(), "org-1", "   ", req)
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestOrganizationService_Update_NilRequest(t *testing.T) {
	svc := NewOrganizationService(
		&mockOrganizationRepository{},
		&mockGroupRepoForOrg{},
		&mockUserRepoForOrg{},
	)

	_, err := svc.Update(context.Background(), "org-1", "user-1", nil)
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestOrganizationService_Update_EmptyName(t *testing.T) {
	svc := NewOrganizationService(
		&mockOrganizationRepository{},
		&mockGroupRepoForOrg{},
		&mockUserRepoForOrg{},
	)

	req := &models.UpdateOrganizationRequest{Name: ""}
	_, err := svc.Update(context.Background(), "org-1", "user-1", req)
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestOrganizationService_Update_WhitespaceOnlyName(t *testing.T) {
	svc := NewOrganizationService(
		&mockOrganizationRepository{},
		&mockGroupRepoForOrg{},
		&mockUserRepoForOrg{},
	)

	req := &models.UpdateOrganizationRequest{Name: "   "}
	_, err := svc.Update(context.Background(), "org-1", "user-1", req)
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestOrganizationService_Update_OrgNotFound(t *testing.T) {
	orgRepo := &mockOrganizationRepository{
		findByID: func(ctx context.Context, orgID string) (*models.Organization, error) {
			return nil, errors.New("not found")
		},
	}
	svc := NewOrganizationService(orgRepo, &mockGroupRepoForOrg{}, &mockUserRepoForOrg{})

	req := &models.UpdateOrganizationRequest{Name: "Updated Org"}
	_, err := svc.Update(context.Background(), "org-1", "user-1", req)
	if !errors.Is(err, ErrOrganizationNotFound) {
		t.Errorf("expected ErrOrganizationNotFound, got %v", err)
	}
}

func TestOrganizationService_Update_IsSiteAdminError(t *testing.T) {
	repoErr := errors.New("database error")
	orgRepo := &mockOrganizationRepository{
		findByID: func(ctx context.Context, orgID string) (*models.Organization, error) {
			return testOrg("org-1", "My Org"), nil
		},
	}
	userRepo := &mockUserRepoForOrg{
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			return false, repoErr
		},
	}
	svc := NewOrganizationService(orgRepo, &mockGroupRepoForOrg{}, userRepo)

	req := &models.UpdateOrganizationRequest{Name: "Updated Org"}
	_, err := svc.Update(context.Background(), "org-1", "user-1", req)
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repository error, got %v", err)
	}
}

func TestOrganizationService_Update_NotSiteAdminAndAdminCheckError(t *testing.T) {
	repoErr := errors.New("admin check error")
	orgRepo := &mockOrganizationRepository{
		findByID: func(ctx context.Context, orgID string) (*models.Organization, error) {
			return testOrg("org-1", "My Org"), nil
		},
		findAdmin: func(ctx context.Context, orgID, userID string) (*models.OrganizationAdmin, error) {
			return nil, repoErr
		},
	}
	userRepo := &mockUserRepoForOrg{
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			return false, nil
		},
	}
	svc := NewOrganizationService(orgRepo, &mockGroupRepoForOrg{}, userRepo)

	req := &models.UpdateOrganizationRequest{Name: "Updated Org"}
	_, err := svc.Update(context.Background(), "org-1", "user-1", req)
	if !errors.Is(err, repoErr) {
		t.Errorf("expected admin check error, got %v", err)
	}
}

func TestOrganizationService_Update_NotSiteAdminAndNotOrgAdmin(t *testing.T) {
	orgRepo := &mockOrganizationRepository{
		findByID: func(ctx context.Context, orgID string) (*models.Organization, error) {
			return testOrg("org-1", "My Org"), nil
		},
		findAdmin: func(ctx context.Context, orgID, userID string) (*models.OrganizationAdmin, error) {
			return nil, nil // User is not an admin
		},
	}
	userRepo := &mockUserRepoForOrg{
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			return false, nil
		},
	}
	svc := NewOrganizationService(orgRepo, &mockGroupRepoForOrg{}, userRepo)

	req := &models.UpdateOrganizationRequest{Name: "Updated Org"}
	_, err := svc.Update(context.Background(), "org-1", "user-1", req)
	if !errors.Is(err, ErrForbidden) {
		t.Errorf("expected ErrForbidden, got %v", err)
	}
}

func TestOrganizationService_Update_UpdateOrgError(t *testing.T) {
	repoErr := errors.New("update error")
	orgRepo := &mockOrganizationRepository{
		findByID: func(ctx context.Context, orgID string) (*models.Organization, error) {
			return testOrg("org-1", "My Org"), nil
		},
		update: func(ctx context.Context, orgID string, name string, description *string) (*models.Organization, error) {
			return nil, repoErr
		},
	}
	userRepo := &mockUserRepoForOrg{
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			return true, nil // Site admin
		},
	}
	svc := NewOrganizationService(orgRepo, &mockGroupRepoForOrg{}, userRepo)

	req := &models.UpdateOrganizationRequest{Name: "Updated Org"}
	_, err := svc.Update(context.Background(), "org-1", "user-1", req)
	if !errors.Is(err, repoErr) {
		t.Errorf("expected update error, got %v", err)
	}
}

func TestOrganizationService_Update_SuccessAsSiteAdmin(t *testing.T) {
	expectedOrg := testOrg("org-1", "Updated Org")
	orgRepo := &mockOrganizationRepository{
		findByID: func(ctx context.Context, orgID string) (*models.Organization, error) {
			return testOrg("org-1", "My Org"), nil
		},
		update: func(ctx context.Context, orgID string, name string, description *string) (*models.Organization, error) {
			if name != "Updated Org" {
				t.Errorf("expected name 'Updated Org', got '%s'", name)
			}
			return expectedOrg, nil
		},
	}
	userRepo := &mockUserRepoForOrg{
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			return true, nil
		},
	}
	svc := NewOrganizationService(orgRepo, &mockGroupRepoForOrg{}, userRepo)

	req := &models.UpdateOrganizationRequest{Name: "Updated Org"}
	org, err := svc.Update(context.Background(), "org-1", "user-1", req)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if org.Name != "Updated Org" {
		t.Errorf("expected org name 'Updated Org', got '%s'", org.Name)
	}
}

func TestOrganizationService_Update_SuccessAsOrgAdmin(t *testing.T) {
	expectedOrg := testOrg("org-1", "Updated Org")
	orgRepo := &mockOrganizationRepository{
		findByID: func(ctx context.Context, orgID string) (*models.Organization, error) {
			return testOrg("org-1", "My Org"), nil
		},
		findAdmin: func(ctx context.Context, orgID, userID string) (*models.OrganizationAdmin, error) {
			return testAdmin("admin-1", userID, orgID), nil // User is an admin
		},
		update: func(ctx context.Context, orgID string, name string, description *string) (*models.Organization, error) {
			return expectedOrg, nil
		},
	}
	userRepo := &mockUserRepoForOrg{
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			return false, nil // Not a site admin
		},
	}
	svc := NewOrganizationService(orgRepo, &mockGroupRepoForOrg{}, userRepo)

	req := &models.UpdateOrganizationRequest{Name: "Updated Org"}
	org, err := svc.Update(context.Background(), "org-1", "user-1", req)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if org.Name != "Updated Org" {
		t.Errorf("expected org name 'Updated Org', got '%s'", org.Name)
	}
}

func TestOrganizationService_Update_TrimsInputs(t *testing.T) {
	orgRepo := &mockOrganizationRepository{
		findByID: func(ctx context.Context, orgID string) (*models.Organization, error) {
			if orgID != "org-1" {
				t.Errorf("expected trimmed orgID 'org-1', got '%s'", orgID)
			}
			return testOrg("org-1", "My Org"), nil
		},
		update: func(ctx context.Context, orgID string, name string, description *string) (*models.Organization, error) {
			if name != "Updated Org" {
				t.Errorf("expected trimmed name 'Updated Org', got '%s'", name)
			}
			return testOrg("org-1", name), nil
		},
	}
	userRepo := &mockUserRepoForOrg{
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			if userID != "user-1" {
				t.Errorf("expected trimmed userID 'user-1', got '%s'", userID)
			}
			return true, nil
		},
	}
	svc := NewOrganizationService(orgRepo, &mockGroupRepoForOrg{}, userRepo)

	req := &models.UpdateOrganizationRequest{Name: "  Updated Org  "}
	_, err := svc.Update(context.Background(), "  org-1  ", "  user-1  ", req)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestOrganizationService_Update_WithDescription(t *testing.T) {
	desc := "Updated description"
	var capturedDesc *string
	orgRepo := &mockOrganizationRepository{
		findByID: func(ctx context.Context, orgID string) (*models.Organization, error) {
			return testOrg("org-1", "My Org"), nil
		},
		update: func(ctx context.Context, orgID string, name string, description *string) (*models.Organization, error) {
			capturedDesc = description
			return testOrg("org-1", name), nil
		},
	}
	userRepo := &mockUserRepoForOrg{
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			return true, nil
		},
	}
	svc := NewOrganizationService(orgRepo, &mockGroupRepoForOrg{}, userRepo)

	req := &models.UpdateOrganizationRequest{Name: "Updated Org", Description: &desc}
	_, err := svc.Update(context.Background(), "org-1", "user-1", req)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if capturedDesc == nil || *capturedDesc != desc {
		t.Errorf("expected description '%s', got %v", desc, capturedDesc)
	}
}

// =============================================================================
// ListAdmins Tests
// =============================================================================

func TestOrganizationService_ListAdmins_EmptyOrgID(t *testing.T) {
	svc := NewOrganizationService(
		&mockOrganizationRepository{},
		&mockGroupRepoForOrg{},
		&mockUserRepoForOrg{},
	)

	_, err := svc.ListAdmins(context.Background(), "")
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestOrganizationService_ListAdmins_WhitespaceOnlyOrgID(t *testing.T) {
	svc := NewOrganizationService(
		&mockOrganizationRepository{},
		&mockGroupRepoForOrg{},
		&mockUserRepoForOrg{},
	)

	_, err := svc.ListAdmins(context.Background(), "   ")
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestOrganizationService_ListAdmins_OrgNotFound(t *testing.T) {
	orgRepo := &mockOrganizationRepository{
		findByID: func(ctx context.Context, orgID string) (*models.Organization, error) {
			return nil, errors.New("not found")
		},
	}
	svc := NewOrganizationService(orgRepo, &mockGroupRepoForOrg{}, &mockUserRepoForOrg{})

	_, err := svc.ListAdmins(context.Background(), "org-1")
	if !errors.Is(err, ErrOrganizationNotFound) {
		t.Errorf("expected ErrOrganizationNotFound, got %v", err)
	}
}

func TestOrganizationService_ListAdmins_FindAdminsError(t *testing.T) {
	repoErr := errors.New("database error")
	orgRepo := &mockOrganizationRepository{
		findByID: func(ctx context.Context, orgID string) (*models.Organization, error) {
			return testOrg("org-1", "My Org"), nil
		},
		findAdmins: func(ctx context.Context, orgID string) ([]*models.OrganizationAdmin, error) {
			return nil, repoErr
		},
	}
	svc := NewOrganizationService(orgRepo, &mockGroupRepoForOrg{}, &mockUserRepoForOrg{})

	_, err := svc.ListAdmins(context.Background(), "org-1")
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repository error, got %v", err)
	}
}

func TestOrganizationService_ListAdmins_NoAdmins(t *testing.T) {
	orgRepo := &mockOrganizationRepository{
		findByID: func(ctx context.Context, orgID string) (*models.Organization, error) {
			return testOrg("org-1", "My Org"), nil
		},
		findAdmins: func(ctx context.Context, orgID string) ([]*models.OrganizationAdmin, error) {
			return []*models.OrganizationAdmin{}, nil
		},
	}
	svc := NewOrganizationService(orgRepo, &mockGroupRepoForOrg{}, &mockUserRepoForOrg{})

	admins, err := svc.ListAdmins(context.Background(), "org-1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(admins) != 0 {
		t.Errorf("expected empty slice, got %d admins", len(admins))
	}
}

func TestOrganizationService_ListAdmins_Success(t *testing.T) {
	expectedAdmins := []*models.OrganizationAdmin{
		testAdmin("admin-1", "user-1", "org-1"),
		testAdmin("admin-2", "user-2", "org-1"),
	}
	orgRepo := &mockOrganizationRepository{
		findByID: func(ctx context.Context, orgID string) (*models.Organization, error) {
			return testOrg("org-1", "My Org"), nil
		},
		findAdmins: func(ctx context.Context, orgID string) ([]*models.OrganizationAdmin, error) {
			if orgID != "org-1" {
				t.Errorf("expected orgID 'org-1', got '%s'", orgID)
			}
			return expectedAdmins, nil
		},
	}
	svc := NewOrganizationService(orgRepo, &mockGroupRepoForOrg{}, &mockUserRepoForOrg{})

	admins, err := svc.ListAdmins(context.Background(), "org-1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(admins) != 2 {
		t.Errorf("expected 2 admins, got %d", len(admins))
	}
}

func TestOrganizationService_ListAdmins_TrimsOrgID(t *testing.T) {
	orgRepo := &mockOrganizationRepository{
		findByID: func(ctx context.Context, orgID string) (*models.Organization, error) {
			if orgID != "org-1" {
				t.Errorf("expected trimmed orgID 'org-1', got '%s'", orgID)
			}
			return testOrg("org-1", "My Org"), nil
		},
		findAdmins: func(ctx context.Context, orgID string) ([]*models.OrganizationAdmin, error) {
			return []*models.OrganizationAdmin{}, nil
		},
	}
	svc := NewOrganizationService(orgRepo, &mockGroupRepoForOrg{}, &mockUserRepoForOrg{})

	_, err := svc.ListAdmins(context.Background(), "  org-1  ")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// =============================================================================
// ListGroups Tests
// =============================================================================

func TestOrganizationService_ListGroups_EmptyOrgID(t *testing.T) {
	svc := NewOrganizationService(
		&mockOrganizationRepository{},
		&mockGroupRepoForOrg{},
		&mockUserRepoForOrg{},
	)

	_, err := svc.ListGroups(context.Background(), "")
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestOrganizationService_ListGroups_WhitespaceOnlyOrgID(t *testing.T) {
	svc := NewOrganizationService(
		&mockOrganizationRepository{},
		&mockGroupRepoForOrg{},
		&mockUserRepoForOrg{},
	)

	_, err := svc.ListGroups(context.Background(), "   ")
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestOrganizationService_ListGroups_OrgNotFound(t *testing.T) {
	orgRepo := &mockOrganizationRepository{
		findByID: func(ctx context.Context, orgID string) (*models.Organization, error) {
			return nil, errors.New("not found")
		},
	}
	svc := NewOrganizationService(orgRepo, &mockGroupRepoForOrg{}, &mockUserRepoForOrg{})

	_, err := svc.ListGroups(context.Background(), "org-1")
	if !errors.Is(err, ErrOrganizationNotFound) {
		t.Errorf("expected ErrOrganizationNotFound, got %v", err)
	}
}

func TestOrganizationService_ListGroups_FindGroupsError(t *testing.T) {
	repoErr := errors.New("database error")
	orgRepo := &mockOrganizationRepository{
		findByID: func(ctx context.Context, orgID string) (*models.Organization, error) {
			return testOrg("org-1", "My Org"), nil
		},
	}
	groupRepo := &mockGroupRepoForOrg{
		findByOrganizationID: func(ctx context.Context, orgID string) ([]*models.Group, error) {
			return nil, repoErr
		},
	}
	svc := NewOrganizationService(orgRepo, groupRepo, &mockUserRepoForOrg{})

	_, err := svc.ListGroups(context.Background(), "org-1")
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repository error, got %v", err)
	}
}

func TestOrganizationService_ListGroups_NoGroups(t *testing.T) {
	orgRepo := &mockOrganizationRepository{
		findByID: func(ctx context.Context, orgID string) (*models.Organization, error) {
			return testOrg("org-1", "My Org"), nil
		},
	}
	groupRepo := &mockGroupRepoForOrg{
		findByOrganizationID: func(ctx context.Context, orgID string) ([]*models.Group, error) {
			return []*models.Group{}, nil
		},
	}
	svc := NewOrganizationService(orgRepo, groupRepo, &mockUserRepoForOrg{})

	groups, err := svc.ListGroups(context.Background(), "org-1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(groups) != 0 {
		t.Errorf("expected empty slice, got %d groups", len(groups))
	}
}

func TestOrganizationService_ListGroups_Success(t *testing.T) {
	expectedGroups := []*models.Group{
		testGroup("group-1", "Group 1", "org-1"),
		testGroup("group-2", "Group 2", "org-1"),
		testGroup("group-3", "Group 3", "org-1"),
	}
	orgRepo := &mockOrganizationRepository{
		findByID: func(ctx context.Context, orgID string) (*models.Organization, error) {
			return testOrg("org-1", "My Org"), nil
		},
	}
	groupRepo := &mockGroupRepoForOrg{
		findByOrganizationID: func(ctx context.Context, orgID string) ([]*models.Group, error) {
			if orgID != "org-1" {
				t.Errorf("expected orgID 'org-1', got '%s'", orgID)
			}
			return expectedGroups, nil
		},
	}
	svc := NewOrganizationService(orgRepo, groupRepo, &mockUserRepoForOrg{})

	groups, err := svc.ListGroups(context.Background(), "org-1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(groups) != 3 {
		t.Errorf("expected 3 groups, got %d", len(groups))
	}
}

func TestOrganizationService_ListGroups_TrimsOrgID(t *testing.T) {
	orgRepo := &mockOrganizationRepository{
		findByID: func(ctx context.Context, orgID string) (*models.Organization, error) {
			if orgID != "org-1" {
				t.Errorf("expected trimmed orgID 'org-1', got '%s'", orgID)
			}
			return testOrg("org-1", "My Org"), nil
		},
	}
	groupRepo := &mockGroupRepoForOrg{
		findByOrganizationID: func(ctx context.Context, orgID string) ([]*models.Group, error) {
			return []*models.Group{}, nil
		},
	}
	svc := NewOrganizationService(orgRepo, groupRepo, &mockUserRepoForOrg{})

	_, err := svc.ListGroups(context.Background(), "  org-1  ")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// =============================================================================
// IsAdmin Tests
// =============================================================================

func TestOrganizationService_IsAdmin_EmptyOrgID(t *testing.T) {
	svc := NewOrganizationService(
		&mockOrganizationRepository{},
		&mockGroupRepoForOrg{},
		&mockUserRepoForOrg{},
	)

	_, err := svc.IsAdmin(context.Background(), "", "user-1")
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestOrganizationService_IsAdmin_WhitespaceOnlyOrgID(t *testing.T) {
	svc := NewOrganizationService(
		&mockOrganizationRepository{},
		&mockGroupRepoForOrg{},
		&mockUserRepoForOrg{},
	)

	_, err := svc.IsAdmin(context.Background(), "   ", "user-1")
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestOrganizationService_IsAdmin_EmptyUserID(t *testing.T) {
	svc := NewOrganizationService(
		&mockOrganizationRepository{},
		&mockGroupRepoForOrg{},
		&mockUserRepoForOrg{},
	)

	_, err := svc.IsAdmin(context.Background(), "org-1", "")
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestOrganizationService_IsAdmin_WhitespaceOnlyUserID(t *testing.T) {
	svc := NewOrganizationService(
		&mockOrganizationRepository{},
		&mockGroupRepoForOrg{},
		&mockUserRepoForOrg{},
	)

	_, err := svc.IsAdmin(context.Background(), "org-1", "   ")
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestOrganizationService_IsAdmin_RepositoryError(t *testing.T) {
	repoErr := errors.New("database error")
	orgRepo := &mockOrganizationRepository{
		findAdmin: func(ctx context.Context, orgID, userID string) (*models.OrganizationAdmin, error) {
			return nil, repoErr
		},
	}
	svc := NewOrganizationService(orgRepo, &mockGroupRepoForOrg{}, &mockUserRepoForOrg{})

	_, err := svc.IsAdmin(context.Background(), "org-1", "user-1")
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repository error, got %v", err)
	}
}

func TestOrganizationService_IsAdmin_NotAdmin(t *testing.T) {
	orgRepo := &mockOrganizationRepository{
		findAdmin: func(ctx context.Context, orgID, userID string) (*models.OrganizationAdmin, error) {
			return nil, nil
		},
	}
	svc := NewOrganizationService(orgRepo, &mockGroupRepoForOrg{}, &mockUserRepoForOrg{})

	isAdmin, err := svc.IsAdmin(context.Background(), "org-1", "user-1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if isAdmin {
		t.Errorf("expected isAdmin to be false")
	}
}

func TestOrganizationService_IsAdmin_IsAdmin(t *testing.T) {
	orgRepo := &mockOrganizationRepository{
		findAdmin: func(ctx context.Context, orgID, userID string) (*models.OrganizationAdmin, error) {
			return testAdmin("admin-1", userID, orgID), nil
		},
	}
	svc := NewOrganizationService(orgRepo, &mockGroupRepoForOrg{}, &mockUserRepoForOrg{})

	isAdmin, err := svc.IsAdmin(context.Background(), "org-1", "user-1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !isAdmin {
		t.Errorf("expected isAdmin to be true")
	}
}

func TestOrganizationService_IsAdmin_TrimsInputs(t *testing.T) {
	orgRepo := &mockOrganizationRepository{
		findAdmin: func(ctx context.Context, orgID, userID string) (*models.OrganizationAdmin, error) {
			if orgID != "org-1" {
				t.Errorf("expected trimmed orgID 'org-1', got '%s'", orgID)
			}
			if userID != "user-1" {
				t.Errorf("expected trimmed userID 'user-1', got '%s'", userID)
			}
			return nil, nil
		},
	}
	svc := NewOrganizationService(orgRepo, &mockGroupRepoForOrg{}, &mockUserRepoForOrg{})

	_, err := svc.IsAdmin(context.Background(), "  org-1  ", "  user-1  ")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// =============================================================================
// IsMember Tests
// =============================================================================

func TestOrganizationService_IsMember_EmptyOrgID(t *testing.T) {
	svc := NewOrganizationService(
		&mockOrganizationRepository{},
		&mockGroupRepoForOrg{},
		&mockUserRepoForOrg{},
	)

	_, err := svc.IsMember(context.Background(), "", "user-1")
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestOrganizationService_IsMember_WhitespaceOnlyOrgID(t *testing.T) {
	svc := NewOrganizationService(
		&mockOrganizationRepository{},
		&mockGroupRepoForOrg{},
		&mockUserRepoForOrg{},
	)

	_, err := svc.IsMember(context.Background(), "   ", "user-1")
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestOrganizationService_IsMember_EmptyUserID(t *testing.T) {
	svc := NewOrganizationService(
		&mockOrganizationRepository{},
		&mockGroupRepoForOrg{},
		&mockUserRepoForOrg{},
	)

	_, err := svc.IsMember(context.Background(), "org-1", "")
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestOrganizationService_IsMember_WhitespaceOnlyUserID(t *testing.T) {
	svc := NewOrganizationService(
		&mockOrganizationRepository{},
		&mockGroupRepoForOrg{},
		&mockUserRepoForOrg{},
	)

	_, err := svc.IsMember(context.Background(), "org-1", "   ")
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestOrganizationService_IsMember_RepositoryError(t *testing.T) {
	repoErr := errors.New("database error")
	orgRepo := &mockOrganizationRepository{
		isMember: func(ctx context.Context, orgID, userID string) (bool, error) {
			return false, repoErr
		},
	}
	svc := NewOrganizationService(orgRepo, &mockGroupRepoForOrg{}, &mockUserRepoForOrg{})

	_, err := svc.IsMember(context.Background(), "org-1", "user-1")
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repository error, got %v", err)
	}
}

func TestOrganizationService_IsMember_NotMember(t *testing.T) {
	orgRepo := &mockOrganizationRepository{
		isMember: func(ctx context.Context, orgID, userID string) (bool, error) {
			return false, nil
		},
	}
	svc := NewOrganizationService(orgRepo, &mockGroupRepoForOrg{}, &mockUserRepoForOrg{})

	isMember, err := svc.IsMember(context.Background(), "org-1", "user-1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if isMember {
		t.Errorf("expected isMember to be false")
	}
}

func TestOrganizationService_IsMember_IsMember(t *testing.T) {
	orgRepo := &mockOrganizationRepository{
		isMember: func(ctx context.Context, orgID, userID string) (bool, error) {
			return true, nil
		},
	}
	svc := NewOrganizationService(orgRepo, &mockGroupRepoForOrg{}, &mockUserRepoForOrg{})

	isMember, err := svc.IsMember(context.Background(), "org-1", "user-1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !isMember {
		t.Errorf("expected isMember to be true")
	}
}

func TestOrganizationService_IsMember_TrimsInputs(t *testing.T) {
	orgRepo := &mockOrganizationRepository{
		isMember: func(ctx context.Context, orgID, userID string) (bool, error) {
			if orgID != "org-1" {
				t.Errorf("expected trimmed orgID 'org-1', got '%s'", orgID)
			}
			if userID != "user-1" {
				t.Errorf("expected trimmed userID 'user-1', got '%s'", userID)
			}
			return false, nil
		},
	}
	svc := NewOrganizationService(orgRepo, &mockGroupRepoForOrg{}, &mockUserRepoForOrg{})

	_, err := svc.IsMember(context.Background(), "  org-1  ", "  user-1  ")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
