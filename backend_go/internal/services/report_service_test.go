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

type mockReportRepository struct {
	findByID      func(ctx context.Context, reportID string) (*models.Report, error)
	findByGroupID func(ctx context.Context, groupID string) ([]*models.Report, error)
	create        func(ctx context.Context, groupID, reporterID, title, content string, isAnonymous bool) (*models.Report, error)
}

func (m *mockReportRepository) FindByID(ctx context.Context, reportID string) (*models.Report, error) {
	if m.findByID != nil {
		return m.findByID(ctx, reportID)
	}
	return nil, errors.New("FindByID not implemented")
}

func (m *mockReportRepository) FindByGroupID(ctx context.Context, groupID string) ([]*models.Report, error) {
	if m.findByGroupID != nil {
		return m.findByGroupID(ctx, groupID)
	}
	return nil, errors.New("FindByGroupID not implemented")
}

func (m *mockReportRepository) Create(ctx context.Context, groupID, reporterID, title, content string, isAnonymous bool) (*models.Report, error) {
	if m.create != nil {
		return m.create(ctx, groupID, reporterID, title, content, isAnonymous)
	}
	return nil, errors.New("Create not implemented")
}

type mockGroupRepoForReport struct {
	findByID             func(ctx context.Context, groupID string) (*models.Group, error)
	findMember           func(ctx context.Context, groupID, userID string) (*models.GroupMember, error)
	findByOrganizationID func(ctx context.Context, orgID string) ([]*models.Group, error)
}

func (m *mockGroupRepoForReport) FindByID(ctx context.Context, groupID string) (*models.Group, error) {
	if m.findByID != nil {
		return m.findByID(ctx, groupID)
	}
	return nil, errors.New("FindByID not implemented")
}

func (m *mockGroupRepoForReport) FindByInviteCode(ctx context.Context, inviteCode string) (*models.Group, error) {
	return nil, errors.New("FindByInviteCode not implemented")
}

func (m *mockGroupRepoForReport) FindMember(ctx context.Context, groupID, userID string) (*models.GroupMember, error) {
	if m.findMember != nil {
		return m.findMember(ctx, groupID, userID)
	}
	return nil, errors.New("FindMember not implemented")
}

func (m *mockGroupRepoForReport) FindByOrganizationID(ctx context.Context, orgID string) ([]*models.Group, error) {
	if m.findByOrganizationID != nil {
		return m.findByOrganizationID(ctx, orgID)
	}
	return nil, errors.New("FindByOrganizationID not implemented")
}

func (m *mockGroupRepoForReport) FindByUserID(ctx context.Context, userID string) ([]*models.Group, error) {
	return nil, errors.New("FindByUserID not implemented")
}

func (m *mockGroupRepoForReport) Create(ctx context.Context, name string, description *string, orgID string, inviteCode string, createdBy string) (*models.Group, error) {
	return nil, errors.New("Create not implemented")
}

func (m *mockGroupRepoForReport) GenerateInviteCode() string {
	return ""
}

func (m *mockGroupRepoForReport) FindMembers(ctx context.Context, groupID string) ([]*models.GroupMember, error) {
	return nil, errors.New("FindMembers not implemented")
}

func (m *mockGroupRepoForReport) CountMembers(ctx context.Context, groupID string) (int, error) {
	return 0, errors.New("CountMembers not implemented")
}

func (m *mockGroupRepoForReport) AddMember(ctx context.Context, groupID, userID, role string) (*models.GroupMember, error) {
	return nil, errors.New("AddMember not implemented")
}

func (m *mockGroupRepoForReport) RemoveMember(ctx context.Context, groupID, userID string) error {
	return errors.New("RemoveMember not implemented")
}

func (m *mockGroupRepoForReport) Count(ctx context.Context) (int64, error) {
	return 0, errors.New("Count not implemented")
}

func (m *mockGroupRepoForReport) CountByOrganizationIDs(ctx context.Context, orgIDs []string) (int64, error) {
	return 0, errors.New("CountByOrganizationIDs not implemented")
}

func (m *mockGroupRepoForReport) UpdateMemberRole(ctx context.Context, groupID, userID, role string) error {
	return nil
}

func (m *mockGroupRepoForReport) UpdateSettings(ctx context.Context, groupID string, requirePostApproval, allowAnonymousPosts bool) error {
	return nil
}

// =============================================================================
// Helper Functions
// =============================================================================

// Helper to create a test report
func testReport(id, groupID, reporterID, title, content string, isAnonymous bool) *models.Report {
	return &models.Report{
		ID:                 id,
		GroupID:            groupID,
		ReporterID:         reporterID,
		Title:              title,
		Content:            content,
		IsAnonymousToGroup: isAnonymous,
		CreatedAt:          time.Now(),
	}
}

// Helper to create a test group member
func testMember(id, groupID, userID, role string) *models.GroupMember {
	return &models.GroupMember{
		ID:       id,
		GroupID:  groupID,
		UserID:   userID,
		Role:     role,
		JoinedAt: time.Now(),
	}
}

// =============================================================================
// Create Tests
// =============================================================================

func TestReportService_Create_NilRequest(t *testing.T) {
	svc := NewReportService(&mockReportRepository{}, &mockGroupRepoForReport{})

	_, err := svc.Create(context.Background(), "user-1", nil)
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestReportService_Create_EmptyUserID(t *testing.T) {
	svc := NewReportService(&mockReportRepository{}, &mockGroupRepoForReport{})

	req := &models.CreateReportRequest{
		GroupID: "group-1",
		Title:   "Test Title",
		Content: "Test Content",
	}
	_, err := svc.Create(context.Background(), "", req)
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestReportService_Create_WhitespaceOnlyUserID(t *testing.T) {
	svc := NewReportService(&mockReportRepository{}, &mockGroupRepoForReport{})

	req := &models.CreateReportRequest{
		GroupID: "group-1",
		Title:   "Test Title",
		Content: "Test Content",
	}
	_, err := svc.Create(context.Background(), "   ", req)
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestReportService_Create_EmptyGroupID(t *testing.T) {
	svc := NewReportService(&mockReportRepository{}, &mockGroupRepoForReport{})

	req := &models.CreateReportRequest{
		GroupID: "",
		Title:   "Test Title",
		Content: "Test Content",
	}
	_, err := svc.Create(context.Background(), "user-1", req)
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestReportService_Create_WhitespaceOnlyGroupID(t *testing.T) {
	svc := NewReportService(&mockReportRepository{}, &mockGroupRepoForReport{})

	req := &models.CreateReportRequest{
		GroupID: "   ",
		Title:   "Test Title",
		Content: "Test Content",
	}
	_, err := svc.Create(context.Background(), "user-1", req)
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestReportService_Create_EmptyTitle(t *testing.T) {
	svc := NewReportService(&mockReportRepository{}, &mockGroupRepoForReport{})

	req := &models.CreateReportRequest{
		GroupID: "group-1",
		Title:   "",
		Content: "Test Content",
	}
	_, err := svc.Create(context.Background(), "user-1", req)
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestReportService_Create_WhitespaceOnlyTitle(t *testing.T) {
	svc := NewReportService(&mockReportRepository{}, &mockGroupRepoForReport{})

	req := &models.CreateReportRequest{
		GroupID: "group-1",
		Title:   "   ",
		Content: "Test Content",
	}
	_, err := svc.Create(context.Background(), "user-1", req)
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestReportService_Create_EmptyContent(t *testing.T) {
	svc := NewReportService(&mockReportRepository{}, &mockGroupRepoForReport{})

	req := &models.CreateReportRequest{
		GroupID: "group-1",
		Title:   "Test Title",
		Content: "",
	}
	_, err := svc.Create(context.Background(), "user-1", req)
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestReportService_Create_WhitespaceOnlyContent(t *testing.T) {
	svc := NewReportService(&mockReportRepository{}, &mockGroupRepoForReport{})

	req := &models.CreateReportRequest{
		GroupID: "group-1",
		Title:   "Test Title",
		Content: "   ",
	}
	_, err := svc.Create(context.Background(), "user-1", req)
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestReportService_Create_GroupNotFound(t *testing.T) {
	groupRepo := &mockGroupRepoForReport{
		findByID: func(ctx context.Context, groupID string) (*models.Group, error) {
			return nil, errors.New("not found")
		},
	}
	svc := NewReportService(&mockReportRepository{}, groupRepo)

	req := &models.CreateReportRequest{
		GroupID: "group-1",
		Title:   "Test Title",
		Content: "Test Content",
	}
	_, err := svc.Create(context.Background(), "user-1", req)
	if !errors.Is(err, ErrGroupNotFound) {
		t.Errorf("expected ErrGroupNotFound, got %v", err)
	}
}

func TestReportService_Create_MembershipCheckError(t *testing.T) {
	repoErr := errors.New("database error")
	groupRepo := &mockGroupRepoForReport{
		findByID: func(ctx context.Context, groupID string) (*models.Group, error) {
			return testGroup("group-1", "Test Group", "org-1"), nil
		},
		findMember: func(ctx context.Context, groupID, userID string) (*models.GroupMember, error) {
			return nil, repoErr
		},
	}
	svc := NewReportService(&mockReportRepository{}, groupRepo)

	req := &models.CreateReportRequest{
		GroupID: "group-1",
		Title:   "Test Title",
		Content: "Test Content",
	}
	_, err := svc.Create(context.Background(), "user-1", req)
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repository error, got %v", err)
	}
}

func TestReportService_Create_UserNotMember(t *testing.T) {
	groupRepo := &mockGroupRepoForReport{
		findByID: func(ctx context.Context, groupID string) (*models.Group, error) {
			return testGroup("group-1", "Test Group", "org-1"), nil
		},
		findMember: func(ctx context.Context, groupID, userID string) (*models.GroupMember, error) {
			return nil, nil // User is not a member
		},
	}
	svc := NewReportService(&mockReportRepository{}, groupRepo)

	req := &models.CreateReportRequest{
		GroupID: "group-1",
		Title:   "Test Title",
		Content: "Test Content",
	}
	_, err := svc.Create(context.Background(), "user-1", req)
	if !errors.Is(err, ErrForbidden) {
		t.Errorf("expected ErrForbidden, got %v", err)
	}
}

func TestReportService_Create_CreateReportError(t *testing.T) {
	repoErr := errors.New("create error")
	reportRepo := &mockReportRepository{
		create: func(ctx context.Context, groupID, reporterID, title, content string, isAnonymous bool) (*models.Report, error) {
			return nil, repoErr
		},
	}
	groupRepo := &mockGroupRepoForReport{
		findByID: func(ctx context.Context, groupID string) (*models.Group, error) {
			return testGroup("group-1", "Test Group", "org-1"), nil
		},
		findMember: func(ctx context.Context, groupID, userID string) (*models.GroupMember, error) {
			return testMember("member-1", groupID, userID, "MEMBER"), nil
		},
	}
	svc := NewReportService(reportRepo, groupRepo)

	req := &models.CreateReportRequest{
		GroupID: "group-1",
		Title:   "Test Title",
		Content: "Test Content",
	}
	_, err := svc.Create(context.Background(), "user-1", req)
	if !errors.Is(err, repoErr) {
		t.Errorf("expected create error, got %v", err)
	}
}

func TestReportService_Create_Success(t *testing.T) {
	expectedReport := testReport("report-1", "group-1", "user-1", "Test Title", "Test Content", false)
	reportRepo := &mockReportRepository{
		create: func(ctx context.Context, groupID, reporterID, title, content string, isAnonymous bool) (*models.Report, error) {
			if groupID != "group-1" {
				t.Errorf("expected groupID 'group-1', got '%s'", groupID)
			}
			if reporterID != "user-1" {
				t.Errorf("expected reporterID 'user-1', got '%s'", reporterID)
			}
			if title != "Test Title" {
				t.Errorf("expected title 'Test Title', got '%s'", title)
			}
			if content != "Test Content" {
				t.Errorf("expected content 'Test Content', got '%s'", content)
			}
			return expectedReport, nil
		},
	}
	groupRepo := &mockGroupRepoForReport{
		findByID: func(ctx context.Context, groupID string) (*models.Group, error) {
			return testGroup("group-1", "Test Group", "org-1"), nil
		},
		findMember: func(ctx context.Context, groupID, userID string) (*models.GroupMember, error) {
			return testMember("member-1", groupID, userID, "MEMBER"), nil
		},
	}
	svc := NewReportService(reportRepo, groupRepo)

	req := &models.CreateReportRequest{
		GroupID: "group-1",
		Title:   "Test Title",
		Content: "Test Content",
	}
	report, err := svc.Create(context.Background(), "user-1", req)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if report.ID != "report-1" {
		t.Errorf("expected report ID 'report-1', got '%s'", report.ID)
	}
}

func TestReportService_Create_TrimsInputs(t *testing.T) {
	reportRepo := &mockReportRepository{
		create: func(ctx context.Context, groupID, reporterID, title, content string, isAnonymous bool) (*models.Report, error) {
			if groupID != "group-1" {
				t.Errorf("expected trimmed groupID 'group-1', got '%s'", groupID)
			}
			if reporterID != "user-1" {
				t.Errorf("expected trimmed reporterID 'user-1', got '%s'", reporterID)
			}
			if title != "Test Title" {
				t.Errorf("expected trimmed title 'Test Title', got '%s'", title)
			}
			if content != "Test Content" {
				t.Errorf("expected trimmed content 'Test Content', got '%s'", content)
			}
			return testReport("report-1", groupID, reporterID, title, content, false), nil
		},
	}
	groupRepo := &mockGroupRepoForReport{
		findByID: func(ctx context.Context, groupID string) (*models.Group, error) {
			if groupID != "group-1" {
				t.Errorf("expected trimmed groupID 'group-1', got '%s'", groupID)
			}
			return testGroup("group-1", "Test Group", "org-1"), nil
		},
		findMember: func(ctx context.Context, groupID, userID string) (*models.GroupMember, error) {
			if userID != "user-1" {
				t.Errorf("expected trimmed userID 'user-1', got '%s'", userID)
			}
			return testMember("member-1", groupID, userID, "MEMBER"), nil
		},
	}
	svc := NewReportService(reportRepo, groupRepo)

	req := &models.CreateReportRequest{
		GroupID: "  group-1  ",
		Title:   "  Test Title  ",
		Content: "  Test Content  ",
	}
	_, err := svc.Create(context.Background(), "  user-1  ", req)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestReportService_Create_WithAnonymousFlag(t *testing.T) {
	var capturedAnonymous bool
	reportRepo := &mockReportRepository{
		create: func(ctx context.Context, groupID, reporterID, title, content string, isAnonymous bool) (*models.Report, error) {
			capturedAnonymous = isAnonymous
			return testReport("report-1", groupID, reporterID, title, content, isAnonymous), nil
		},
	}
	groupRepo := &mockGroupRepoForReport{
		findByID: func(ctx context.Context, groupID string) (*models.Group, error) {
			return testGroup("group-1", "Test Group", "org-1"), nil
		},
		findMember: func(ctx context.Context, groupID, userID string) (*models.GroupMember, error) {
			return testMember("member-1", groupID, userID, "MEMBER"), nil
		},
	}
	svc := NewReportService(reportRepo, groupRepo)

	req := &models.CreateReportRequest{
		GroupID:            "group-1",
		Title:              "Test Title",
		Content:            "Test Content",
		IsAnonymousToGroup: true,
	}
	_, err := svc.Create(context.Background(), "user-1", req)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !capturedAnonymous {
		t.Errorf("expected isAnonymous to be true")
	}
}

// =============================================================================
// List Tests
// =============================================================================

func TestReportService_List_EmptyGroupID(t *testing.T) {
	svc := NewReportService(&mockReportRepository{}, &mockGroupRepoForReport{})

	_, err := svc.List(context.Background(), "", "user-1")
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestReportService_List_WhitespaceOnlyGroupID(t *testing.T) {
	svc := NewReportService(&mockReportRepository{}, &mockGroupRepoForReport{})

	_, err := svc.List(context.Background(), "   ", "user-1")
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestReportService_List_EmptyRequesterID(t *testing.T) {
	svc := NewReportService(&mockReportRepository{}, &mockGroupRepoForReport{})

	_, err := svc.List(context.Background(), "group-1", "")
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestReportService_List_WhitespaceOnlyRequesterID(t *testing.T) {
	svc := NewReportService(&mockReportRepository{}, &mockGroupRepoForReport{})

	_, err := svc.List(context.Background(), "group-1", "   ")
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestReportService_List_GroupNotFound(t *testing.T) {
	groupRepo := &mockGroupRepoForReport{
		findByID: func(ctx context.Context, groupID string) (*models.Group, error) {
			return nil, errors.New("not found")
		},
	}
	svc := NewReportService(&mockReportRepository{}, groupRepo)

	_, err := svc.List(context.Background(), "group-1", "user-1")
	if !errors.Is(err, ErrGroupNotFound) {
		t.Errorf("expected ErrGroupNotFound, got %v", err)
	}
}

func TestReportService_List_MembershipCheckError(t *testing.T) {
	repoErr := errors.New("database error")
	groupRepo := &mockGroupRepoForReport{
		findByID: func(ctx context.Context, groupID string) (*models.Group, error) {
			return testGroup("group-1", "Test Group", "org-1"), nil
		},
		findMember: func(ctx context.Context, groupID, userID string) (*models.GroupMember, error) {
			return nil, repoErr
		},
	}
	svc := NewReportService(&mockReportRepository{}, groupRepo)

	_, err := svc.List(context.Background(), "group-1", "user-1")
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repository error, got %v", err)
	}
}

func TestReportService_List_RequesterNotMember(t *testing.T) {
	groupRepo := &mockGroupRepoForReport{
		findByID: func(ctx context.Context, groupID string) (*models.Group, error) {
			return testGroup("group-1", "Test Group", "org-1"), nil
		},
		findMember: func(ctx context.Context, groupID, userID string) (*models.GroupMember, error) {
			return nil, nil // User is not a member
		},
	}
	svc := NewReportService(&mockReportRepository{}, groupRepo)

	_, err := svc.List(context.Background(), "group-1", "user-1")
	if !errors.Is(err, ErrForbidden) {
		t.Errorf("expected ErrForbidden, got %v", err)
	}
}

func TestReportService_List_ReportLookupError(t *testing.T) {
	repoErr := errors.New("database error")
	reportRepo := &mockReportRepository{
		findByGroupID: func(ctx context.Context, groupID string) ([]*models.Report, error) {
			return nil, repoErr
		},
	}
	groupRepo := &mockGroupRepoForReport{
		findByID: func(ctx context.Context, groupID string) (*models.Group, error) {
			return testGroup("group-1", "Test Group", "org-1"), nil
		},
		findMember: func(ctx context.Context, groupID, userID string) (*models.GroupMember, error) {
			return testMember("member-1", groupID, userID, "MEMBER"), nil
		},
	}
	svc := NewReportService(reportRepo, groupRepo)

	_, err := svc.List(context.Background(), "group-1", "user-1")
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repository error, got %v", err)
	}
}

func TestReportService_List_NoReports(t *testing.T) {
	reportRepo := &mockReportRepository{
		findByGroupID: func(ctx context.Context, groupID string) ([]*models.Report, error) {
			return []*models.Report{}, nil
		},
	}
	groupRepo := &mockGroupRepoForReport{
		findByID: func(ctx context.Context, groupID string) (*models.Group, error) {
			return testGroup("group-1", "Test Group", "org-1"), nil
		},
		findMember: func(ctx context.Context, groupID, userID string) (*models.GroupMember, error) {
			return testMember("member-1", groupID, userID, "MEMBER"), nil
		},
	}
	svc := NewReportService(reportRepo, groupRepo)

	reports, err := svc.List(context.Background(), "group-1", "user-1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(reports) != 0 {
		t.Errorf("expected empty slice, got %d reports", len(reports))
	}
}

func TestReportService_List_Success(t *testing.T) {
	expectedReports := []*models.Report{
		testReport("report-1", "group-1", "user-1", "Title 1", "Content 1", false),
		testReport("report-2", "group-1", "user-2", "Title 2", "Content 2", true),
	}
	reportRepo := &mockReportRepository{
		findByGroupID: func(ctx context.Context, groupID string) ([]*models.Report, error) {
			if groupID != "group-1" {
				t.Errorf("expected groupID 'group-1', got '%s'", groupID)
			}
			return expectedReports, nil
		},
	}
	groupRepo := &mockGroupRepoForReport{
		findByID: func(ctx context.Context, groupID string) (*models.Group, error) {
			return testGroup("group-1", "Test Group", "org-1"), nil
		},
		findMember: func(ctx context.Context, groupID, userID string) (*models.GroupMember, error) {
			return testMember("member-1", groupID, userID, "MEMBER"), nil
		},
	}
	svc := NewReportService(reportRepo, groupRepo)

	reports, err := svc.List(context.Background(), "group-1", "user-1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(reports) != 2 {
		t.Errorf("expected 2 reports, got %d", len(reports))
	}
}

func TestReportService_List_TrimsInputs(t *testing.T) {
	reportRepo := &mockReportRepository{
		findByGroupID: func(ctx context.Context, groupID string) ([]*models.Report, error) {
			if groupID != "group-1" {
				t.Errorf("expected trimmed groupID 'group-1', got '%s'", groupID)
			}
			return []*models.Report{}, nil
		},
	}
	groupRepo := &mockGroupRepoForReport{
		findByID: func(ctx context.Context, groupID string) (*models.Group, error) {
			if groupID != "group-1" {
				t.Errorf("expected trimmed groupID 'group-1', got '%s'", groupID)
			}
			return testGroup("group-1", "Test Group", "org-1"), nil
		},
		findMember: func(ctx context.Context, groupID, userID string) (*models.GroupMember, error) {
			if userID != "user-1" {
				t.Errorf("expected trimmed userID 'user-1', got '%s'", userID)
			}
			return testMember("member-1", groupID, userID, "MEMBER"), nil
		},
	}
	svc := NewReportService(reportRepo, groupRepo)

	_, err := svc.List(context.Background(), "  group-1  ", "  user-1  ")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// =============================================================================
// GetByID Tests
// =============================================================================

func TestReportService_GetByID_EmptyReportID(t *testing.T) {
	svc := NewReportService(&mockReportRepository{}, &mockGroupRepoForReport{})

	_, err := svc.GetByID(context.Background(), "")
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestReportService_GetByID_WhitespaceOnlyReportID(t *testing.T) {
	svc := NewReportService(&mockReportRepository{}, &mockGroupRepoForReport{})

	_, err := svc.GetByID(context.Background(), "   ")
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestReportService_GetByID_NotFound(t *testing.T) {
	reportRepo := &mockReportRepository{
		findByID: func(ctx context.Context, reportID string) (*models.Report, error) {
			return nil, errors.New("not found")
		},
	}
	svc := NewReportService(reportRepo, &mockGroupRepoForReport{})

	_, err := svc.GetByID(context.Background(), "report-1")
	if !errors.Is(err, ErrReportNotFound) {
		t.Errorf("expected ErrReportNotFound, got %v", err)
	}
}

func TestReportService_GetByID_Success(t *testing.T) {
	expectedReport := testReport("report-1", "group-1", "user-1", "Test Title", "Test Content", false)
	reportRepo := &mockReportRepository{
		findByID: func(ctx context.Context, reportID string) (*models.Report, error) {
			if reportID != "report-1" {
				t.Errorf("expected reportID 'report-1', got '%s'", reportID)
			}
			return expectedReport, nil
		},
	}
	svc := NewReportService(reportRepo, &mockGroupRepoForReport{})

	report, err := svc.GetByID(context.Background(), "report-1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if report.ID != "report-1" {
		t.Errorf("expected report ID 'report-1', got '%s'", report.ID)
	}
}

func TestReportService_GetByID_TrimsReportID(t *testing.T) {
	reportRepo := &mockReportRepository{
		findByID: func(ctx context.Context, reportID string) (*models.Report, error) {
			if reportID != "report-1" {
				t.Errorf("expected trimmed reportID 'report-1', got '%s'", reportID)
			}
			return testReport("report-1", "group-1", "user-1", "Test Title", "Test Content", false), nil
		},
	}
	svc := NewReportService(reportRepo, &mockGroupRepoForReport{})

	_, err := svc.GetByID(context.Background(), "  report-1  ")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
