package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"menunderfire/internal/models"
)

// ===== Mock Repositories =====

type mockSitePageRepository struct {
	findAll    func(ctx context.Context) ([]*models.SitePage, error)
	findByID   func(ctx context.Context, pageID string) (*models.SitePage, error)
	findBySlug func(ctx context.Context, slug string) (*models.SitePage, error)
	create     func(ctx context.Context, slug, title, content string, isPublished bool, createdByID string) (*models.SitePage, error)
	update     func(ctx context.Context, pageID, title, content string, isPublished bool, updatedByID string) (*models.SitePage, error)
	delete     func(ctx context.Context, pageID string) error
}

func (m *mockSitePageRepository) FindAll(ctx context.Context) ([]*models.SitePage, error) {
	if m.findAll != nil {
		return m.findAll(ctx)
	}
	return nil, errors.New("not implemented")
}

func (m *mockSitePageRepository) FindByID(ctx context.Context, pageID string) (*models.SitePage, error) {
	if m.findByID != nil {
		return m.findByID(ctx, pageID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockSitePageRepository) FindBySlug(ctx context.Context, slug string) (*models.SitePage, error) {
	if m.findBySlug != nil {
		return m.findBySlug(ctx, slug)
	}
	return nil, errors.New("not implemented")
}

func (m *mockSitePageRepository) Create(ctx context.Context, slug, title, content string, isPublished bool, createdByID string) (*models.SitePage, error) {
	if m.create != nil {
		return m.create(ctx, slug, title, content, isPublished, createdByID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockSitePageRepository) Update(ctx context.Context, pageID, title, content string, isPublished bool, updatedByID string) (*models.SitePage, error) {
	if m.update != nil {
		return m.update(ctx, pageID, title, content, isPublished, updatedByID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockSitePageRepository) Delete(ctx context.Context, pageID string) error {
	if m.delete != nil {
		return m.delete(ctx, pageID)
	}
	return errors.New("not implemented")
}

// mockPageDraftRepoForSitePage implements repositories.PageDraftRepository
type mockPageDraftRepoForSitePage struct{}

func (m *mockPageDraftRepoForSitePage) Create(ctx context.Context, pageID *string, title, content, createdByID string) (*models.PageDraft, error) {
	return nil, errors.New("not implemented")
}
func (m *mockPageDraftRepoForSitePage) FindByID(ctx context.Context, draftID string) (*models.PageDraft, error) {
	return nil, errors.New("not implemented")
}
func (m *mockPageDraftRepoForSitePage) Delete(ctx context.Context, draftID string) error {
	return errors.New("not implemented")
}
func (m *mockPageDraftRepoForSitePage) DeleteExpired(ctx context.Context) (int64, error) {
	return 0, errors.New("not implemented")
}

// mockUserRepoForSitePage is a partial mock that only implements IsSiteAdmin
type mockUserRepoForSitePage struct {
	isSiteAdmin func(ctx context.Context, userID string) (bool, error)
}

func (m *mockUserRepoForSitePage) FindByID(ctx context.Context, userID string) (*models.User, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserRepoForSitePage) FindByExternalID(ctx context.Context, externalID string) (*models.User, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserRepoForSitePage) Update(ctx context.Context, userID, displayName string) (*models.User, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserRepoForSitePage) UpdateExternalID(ctx context.Context, userID, externalID string) error {
	return errors.New("not implemented")
}

func (m *mockUserRepoForSitePage) UpdateInvitationInfo(ctx context.Context, userID, invitedByID, invitationID string) error {
	return errors.New("not implemented")
}

func (m *mockUserRepoForSitePage) IsSiteAdmin(ctx context.Context, userID string) (bool, error) {
	if m.isSiteAdmin != nil {
		return m.isSiteAdmin(ctx, userID)
	}
	return false, errors.New("not implemented")
}

func (m *mockUserRepoForSitePage) FindAdminOrganizationIDs(ctx context.Context, userID string) ([]string, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserRepoForSitePage) FindOwnedGroupIDs(ctx context.Context, userID string) ([]string, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserRepoForSitePage) FindMemberGroupIDs(ctx context.Context, userID string) ([]string, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserRepoForSitePage) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserRepoForSitePage) Count(ctx context.Context) (int64, error) {
	return 0, errors.New("not implemented")
}

func (m *mockUserRepoForSitePage) Create(ctx context.Context, externalID, email, displayName string) (*models.User, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserRepoForSitePage) CreateAsSiteAdmin(ctx context.Context, externalID, email, displayName string) (*models.User, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserRepoForSitePage) RecordAgreementAcceptance(ctx context.Context, userID, version, signature, ipAddress, userAgent string) error {
	return errors.New("not implemented")
}

// ===== Test Helpers =====

func testPage(id, slug, title string) *models.SitePage {
	return &models.SitePage{
		ID:          id,
		Slug:        slug,
		Title:       title,
		Content:     "Test content",
		IsPublished: true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// ===== List Tests =====

func TestSitePageService_List_Success(t *testing.T) {
	pages := []*models.SitePage{
		testPage("page-1", "about", "About Us"),
		testPage("page-2", "contact", "Contact"),
	}

	pageRepo := &mockSitePageRepository{
		findAll: func(ctx context.Context) ([]*models.SitePage, error) {
			return pages, nil
		},
	}
	svc := NewSitePageService(pageRepo, &mockPageDraftRepoForSitePage{}, &mockUserRepoForSitePage{})

	result, err := svc.List(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 pages, got %d", len(result))
	}
}

func TestSitePageService_List_EmptyResult(t *testing.T) {
	pageRepo := &mockSitePageRepository{
		findAll: func(ctx context.Context) ([]*models.SitePage, error) {
			return []*models.SitePage{}, nil
		},
	}
	svc := NewSitePageService(pageRepo, &mockPageDraftRepoForSitePage{}, &mockUserRepoForSitePage{})

	result, err := svc.List(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected 0 pages, got %d", len(result))
	}
}

func TestSitePageService_List_RepositoryError(t *testing.T) {
	repoErr := errors.New("database error")
	pageRepo := &mockSitePageRepository{
		findAll: func(ctx context.Context) ([]*models.SitePage, error) {
			return nil, repoErr
		},
	}
	svc := NewSitePageService(pageRepo, &mockPageDraftRepoForSitePage{}, &mockUserRepoForSitePage{})

	_, err := svc.List(context.Background())
	if err != repoErr {
		t.Fatalf("expected repo error, got %v", err)
	}
}

// ===== GetBySlug Tests =====

func TestSitePageService_GetBySlug_Success(t *testing.T) {
	page := testPage("page-1", "about", "About Us")

	pageRepo := &mockSitePageRepository{
		findBySlug: func(ctx context.Context, slug string) (*models.SitePage, error) {
			if slug == "about" {
				return page, nil
			}
			return nil, errors.New("not found")
		},
	}
	svc := NewSitePageService(pageRepo, &mockPageDraftRepoForSitePage{}, &mockUserRepoForSitePage{})

	result, err := svc.GetBySlug(context.Background(), "about")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ID != "page-1" {
		t.Fatalf("expected page-1, got %s", result.ID)
	}
}

func TestSitePageService_GetBySlug_EmptySlug(t *testing.T) {
	svc := NewSitePageService(&mockSitePageRepository{}, &mockPageDraftRepoForSitePage{}, &mockUserRepoForSitePage{})

	_, err := svc.GetBySlug(context.Background(), "")
	if err != ErrValidation {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

func TestSitePageService_GetBySlug_WhitespaceOnlySlug(t *testing.T) {
	svc := NewSitePageService(&mockSitePageRepository{}, &mockPageDraftRepoForSitePage{}, &mockUserRepoForSitePage{})

	_, err := svc.GetBySlug(context.Background(), "   ")
	if err != ErrValidation {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

func TestSitePageService_GetBySlug_SlugTrimmed(t *testing.T) {
	page := testPage("page-1", "about", "About Us")

	pageRepo := &mockSitePageRepository{
		findBySlug: func(ctx context.Context, slug string) (*models.SitePage, error) {
			if slug == "about" {
				return page, nil
			}
			return nil, errors.New("not found")
		},
	}
	svc := NewSitePageService(pageRepo, &mockPageDraftRepoForSitePage{}, &mockUserRepoForSitePage{})

	result, err := svc.GetBySlug(context.Background(), "  about  ")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ID != "page-1" {
		t.Fatalf("expected page-1, got %s", result.ID)
	}
}

func TestSitePageService_GetBySlug_PageNotFound(t *testing.T) {
	pageRepo := &mockSitePageRepository{
		findBySlug: func(ctx context.Context, slug string) (*models.SitePage, error) {
			return nil, errors.New("not found")
		},
	}
	svc := NewSitePageService(pageRepo, &mockPageDraftRepoForSitePage{}, &mockUserRepoForSitePage{})

	_, err := svc.GetBySlug(context.Background(), "nonexistent")
	if err != ErrPageNotFound {
		t.Fatalf("expected ErrPageNotFound, got %v", err)
	}
}

// ===== Create Tests =====

func TestSitePageService_Create_Success(t *testing.T) {
	page := testPage("page-1", "about", "About Us")

	pageRepo := &mockSitePageRepository{
		create: func(ctx context.Context, slug, title, content string, isPublished bool, createdByID string) (*models.SitePage, error) {
			return page, nil
		},
	}
	userRepo := &mockUserRepoForSitePage{
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			return true, nil
		},
	}
	svc := NewSitePageService(pageRepo, &mockPageDraftRepoForSitePage{}, userRepo)

	req := &models.CreateSitePageRequest{
		Slug:        "about",
		Title:       "About Us",
		Content:     "Content",
		IsPublished: true,
	}

	result, err := svc.Create(context.Background(), "user-1", req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ID != "page-1" {
		t.Fatalf("expected page-1, got %s", result.ID)
	}
}

func TestSitePageService_Create_NilRequest(t *testing.T) {
	svc := NewSitePageService(&mockSitePageRepository{}, &mockPageDraftRepoForSitePage{}, &mockUserRepoForSitePage{})

	_, err := svc.Create(context.Background(), "user-1", nil)
	if err != ErrValidation {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

func TestSitePageService_Create_EmptyUserID(t *testing.T) {
	svc := NewSitePageService(&mockSitePageRepository{}, &mockPageDraftRepoForSitePage{}, &mockUserRepoForSitePage{})

	req := &models.CreateSitePageRequest{
		Slug:  "about",
		Title: "About Us",
	}

	_, err := svc.Create(context.Background(), "", req)
	if err != ErrValidation {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

func TestSitePageService_Create_WhitespaceOnlyUserID(t *testing.T) {
	svc := NewSitePageService(&mockSitePageRepository{}, &mockPageDraftRepoForSitePage{}, &mockUserRepoForSitePage{})

	req := &models.CreateSitePageRequest{
		Slug:  "about",
		Title: "About Us",
	}

	_, err := svc.Create(context.Background(), "   ", req)
	if err != ErrValidation {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

func TestSitePageService_Create_EmptySlug(t *testing.T) {
	svc := NewSitePageService(&mockSitePageRepository{}, &mockPageDraftRepoForSitePage{}, &mockUserRepoForSitePage{})

	req := &models.CreateSitePageRequest{
		Slug:  "",
		Title: "About Us",
	}

	_, err := svc.Create(context.Background(), "user-1", req)
	if err != ErrValidation {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

func TestSitePageService_Create_WhitespaceOnlySlug(t *testing.T) {
	svc := NewSitePageService(&mockSitePageRepository{}, &mockPageDraftRepoForSitePage{}, &mockUserRepoForSitePage{})

	req := &models.CreateSitePageRequest{
		Slug:  "   ",
		Title: "About Us",
	}

	_, err := svc.Create(context.Background(), "user-1", req)
	if err != ErrValidation {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

func TestSitePageService_Create_EmptyTitle(t *testing.T) {
	svc := NewSitePageService(&mockSitePageRepository{}, &mockPageDraftRepoForSitePage{}, &mockUserRepoForSitePage{})

	req := &models.CreateSitePageRequest{
		Slug:  "about",
		Title: "",
	}

	_, err := svc.Create(context.Background(), "user-1", req)
	if err != ErrValidation {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

func TestSitePageService_Create_WhitespaceOnlyTitle(t *testing.T) {
	svc := NewSitePageService(&mockSitePageRepository{}, &mockPageDraftRepoForSitePage{}, &mockUserRepoForSitePage{})

	req := &models.CreateSitePageRequest{
		Slug:  "about",
		Title: "   ",
	}

	_, err := svc.Create(context.Background(), "user-1", req)
	if err != ErrValidation {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

func TestSitePageService_Create_NotSiteAdmin(t *testing.T) {
	userRepo := &mockUserRepoForSitePage{
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			return false, nil
		},
	}
	svc := NewSitePageService(&mockSitePageRepository{}, &mockPageDraftRepoForSitePage{}, userRepo)

	req := &models.CreateSitePageRequest{
		Slug:  "about",
		Title: "About Us",
	}

	_, err := svc.Create(context.Background(), "user-1", req)
	if err != ErrForbidden {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestSitePageService_Create_IsSiteAdminError(t *testing.T) {
	repoErr := errors.New("database error")
	userRepo := &mockUserRepoForSitePage{
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			return false, repoErr
		},
	}
	svc := NewSitePageService(&mockSitePageRepository{}, &mockPageDraftRepoForSitePage{}, userRepo)

	req := &models.CreateSitePageRequest{
		Slug:  "about",
		Title: "About Us",
	}

	_, err := svc.Create(context.Background(), "user-1", req)
	if err != repoErr {
		t.Fatalf("expected repo error, got %v", err)
	}
}

func TestSitePageService_Create_RepositoryError(t *testing.T) {
	repoErr := errors.New("database error")
	pageRepo := &mockSitePageRepository{
		create: func(ctx context.Context, slug, title, content string, isPublished bool, createdByID string) (*models.SitePage, error) {
			return nil, repoErr
		},
	}
	userRepo := &mockUserRepoForSitePage{
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			return true, nil
		},
	}
	svc := NewSitePageService(pageRepo, &mockPageDraftRepoForSitePage{}, userRepo)

	req := &models.CreateSitePageRequest{
		Slug:  "about",
		Title: "About Us",
	}

	_, err := svc.Create(context.Background(), "user-1", req)
	if err != repoErr {
		t.Fatalf("expected repo error, got %v", err)
	}
}

func TestSitePageService_Create_UserIDTrimmed(t *testing.T) {
	var receivedUserID string
	pageRepo := &mockSitePageRepository{
		create: func(ctx context.Context, slug, title, content string, isPublished bool, createdByID string) (*models.SitePage, error) {
			return testPage("page-1", slug, title), nil
		},
	}
	userRepo := &mockUserRepoForSitePage{
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			receivedUserID = userID
			return true, nil
		},
	}
	svc := NewSitePageService(pageRepo, &mockPageDraftRepoForSitePage{}, userRepo)

	req := &models.CreateSitePageRequest{
		Slug:  "about",
		Title: "About Us",
	}

	_, err := svc.Create(context.Background(), "  user-1  ", req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if receivedUserID != "user-1" {
		t.Fatalf("expected userID to be trimmed to 'user-1', got '%s'", receivedUserID)
	}
}

func TestSitePageService_Create_SlugTrimmed(t *testing.T) {
	var receivedSlug string
	pageRepo := &mockSitePageRepository{
		create: func(ctx context.Context, slug, title, content string, isPublished bool, createdByID string) (*models.SitePage, error) {
			receivedSlug = slug
			return testPage("page-1", slug, title), nil
		},
	}
	userRepo := &mockUserRepoForSitePage{
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			return true, nil
		},
	}
	svc := NewSitePageService(pageRepo, &mockPageDraftRepoForSitePage{}, userRepo)

	req := &models.CreateSitePageRequest{
		Slug:  "  about  ",
		Title: "About Us",
	}

	_, err := svc.Create(context.Background(), "user-1", req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if receivedSlug != "about" {
		t.Fatalf("expected slug to be trimmed to 'about', got '%s'", receivedSlug)
	}
}

func TestSitePageService_Create_TitleTrimmed(t *testing.T) {
	var receivedTitle string
	pageRepo := &mockSitePageRepository{
		create: func(ctx context.Context, slug, title, content string, isPublished bool, createdByID string) (*models.SitePage, error) {
			receivedTitle = title
			return testPage("page-1", slug, title), nil
		},
	}
	userRepo := &mockUserRepoForSitePage{
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			return true, nil
		},
	}
	svc := NewSitePageService(pageRepo, &mockPageDraftRepoForSitePage{}, userRepo)

	req := &models.CreateSitePageRequest{
		Slug:  "about",
		Title: "  About Us  ",
	}

	_, err := svc.Create(context.Background(), "user-1", req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if receivedTitle != "About Us" {
		t.Fatalf("expected title to be trimmed to 'About Us', got '%s'", receivedTitle)
	}
}

// ===== Update Tests =====

func TestSitePageService_Update_Success(t *testing.T) {
	page := testPage("page-1", "about", "About Us")
	updatedPage := testPage("page-1", "about", "Updated Title")

	pageRepo := &mockSitePageRepository{
		findByID: func(ctx context.Context, pageID string) (*models.SitePage, error) {
			if pageID == "page-1" {
				return page, nil
			}
			return nil, errors.New("not found")
		},
		update: func(ctx context.Context, pageID, title, content string, isPublished bool, updatedByID string) (*models.SitePage, error) {
			return updatedPage, nil
		},
	}
	userRepo := &mockUserRepoForSitePage{
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			return true, nil
		},
	}
	svc := NewSitePageService(pageRepo, &mockPageDraftRepoForSitePage{}, userRepo)

	req := &models.UpdateSitePageRequest{
		Title:       "Updated Title",
		Content:     "Updated content",
		IsPublished: true,
	}

	result, err := svc.Update(context.Background(), "page-1", "user-1", req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Title != "Updated Title" {
		t.Fatalf("expected Updated Title, got %s", result.Title)
	}
}

func TestSitePageService_Update_EmptyPageID(t *testing.T) {
	svc := NewSitePageService(&mockSitePageRepository{}, &mockPageDraftRepoForSitePage{}, &mockUserRepoForSitePage{})

	req := &models.UpdateSitePageRequest{
		Title: "Updated Title",
	}

	_, err := svc.Update(context.Background(), "", "user-1", req)
	if err != ErrValidation {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

func TestSitePageService_Update_WhitespaceOnlyPageID(t *testing.T) {
	svc := NewSitePageService(&mockSitePageRepository{}, &mockPageDraftRepoForSitePage{}, &mockUserRepoForSitePage{})

	req := &models.UpdateSitePageRequest{
		Title: "Updated Title",
	}

	_, err := svc.Update(context.Background(), "   ", "user-1", req)
	if err != ErrValidation {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

func TestSitePageService_Update_EmptyUserID(t *testing.T) {
	svc := NewSitePageService(&mockSitePageRepository{}, &mockPageDraftRepoForSitePage{}, &mockUserRepoForSitePage{})

	req := &models.UpdateSitePageRequest{
		Title: "Updated Title",
	}

	_, err := svc.Update(context.Background(), "page-1", "", req)
	if err != ErrValidation {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

func TestSitePageService_Update_WhitespaceOnlyUserID(t *testing.T) {
	svc := NewSitePageService(&mockSitePageRepository{}, &mockPageDraftRepoForSitePage{}, &mockUserRepoForSitePage{})

	req := &models.UpdateSitePageRequest{
		Title: "Updated Title",
	}

	_, err := svc.Update(context.Background(), "page-1", "   ", req)
	if err != ErrValidation {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

func TestSitePageService_Update_NilRequest(t *testing.T) {
	svc := NewSitePageService(&mockSitePageRepository{}, &mockPageDraftRepoForSitePage{}, &mockUserRepoForSitePage{})

	_, err := svc.Update(context.Background(), "page-1", "user-1", nil)
	if err != ErrValidation {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

func TestSitePageService_Update_EmptyTitle(t *testing.T) {
	svc := NewSitePageService(&mockSitePageRepository{}, &mockPageDraftRepoForSitePage{}, &mockUserRepoForSitePage{})

	req := &models.UpdateSitePageRequest{
		Title: "",
	}

	_, err := svc.Update(context.Background(), "page-1", "user-1", req)
	if err != ErrValidation {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

func TestSitePageService_Update_WhitespaceOnlyTitle(t *testing.T) {
	svc := NewSitePageService(&mockSitePageRepository{}, &mockPageDraftRepoForSitePage{}, &mockUserRepoForSitePage{})

	req := &models.UpdateSitePageRequest{
		Title: "   ",
	}

	_, err := svc.Update(context.Background(), "page-1", "user-1", req)
	if err != ErrValidation {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

func TestSitePageService_Update_PageNotFound(t *testing.T) {
	pageRepo := &mockSitePageRepository{
		findByID: func(ctx context.Context, pageID string) (*models.SitePage, error) {
			return nil, errors.New("not found")
		},
	}
	svc := NewSitePageService(pageRepo, &mockPageDraftRepoForSitePage{}, &mockUserRepoForSitePage{})

	req := &models.UpdateSitePageRequest{
		Title: "Updated Title",
	}

	_, err := svc.Update(context.Background(), "nonexistent", "user-1", req)
	if err != ErrPageNotFound {
		t.Fatalf("expected ErrPageNotFound, got %v", err)
	}
}

func TestSitePageService_Update_NotSiteAdmin(t *testing.T) {
	page := testPage("page-1", "about", "About Us")

	pageRepo := &mockSitePageRepository{
		findByID: func(ctx context.Context, pageID string) (*models.SitePage, error) {
			return page, nil
		},
	}
	userRepo := &mockUserRepoForSitePage{
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			return false, nil
		},
	}
	svc := NewSitePageService(pageRepo, &mockPageDraftRepoForSitePage{}, userRepo)

	req := &models.UpdateSitePageRequest{
		Title: "Updated Title",
	}

	_, err := svc.Update(context.Background(), "page-1", "user-1", req)
	if err != ErrForbidden {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestSitePageService_Update_IsSiteAdminError(t *testing.T) {
	page := testPage("page-1", "about", "About Us")
	repoErr := errors.New("database error")

	pageRepo := &mockSitePageRepository{
		findByID: func(ctx context.Context, pageID string) (*models.SitePage, error) {
			return page, nil
		},
	}
	userRepo := &mockUserRepoForSitePage{
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			return false, repoErr
		},
	}
	svc := NewSitePageService(pageRepo, &mockPageDraftRepoForSitePage{}, userRepo)

	req := &models.UpdateSitePageRequest{
		Title: "Updated Title",
	}

	_, err := svc.Update(context.Background(), "page-1", "user-1", req)
	if err != repoErr {
		t.Fatalf("expected repo error, got %v", err)
	}
}

func TestSitePageService_Update_RepositoryError(t *testing.T) {
	page := testPage("page-1", "about", "About Us")
	repoErr := errors.New("database error")

	pageRepo := &mockSitePageRepository{
		findByID: func(ctx context.Context, pageID string) (*models.SitePage, error) {
			return page, nil
		},
		update: func(ctx context.Context, pageID, title, content string, isPublished bool, updatedByID string) (*models.SitePage, error) {
			return nil, repoErr
		},
	}
	userRepo := &mockUserRepoForSitePage{
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			return true, nil
		},
	}
	svc := NewSitePageService(pageRepo, &mockPageDraftRepoForSitePage{}, userRepo)

	req := &models.UpdateSitePageRequest{
		Title: "Updated Title",
	}

	_, err := svc.Update(context.Background(), "page-1", "user-1", req)
	if err != repoErr {
		t.Fatalf("expected repo error, got %v", err)
	}
}

func TestSitePageService_Update_PageIDTrimmed(t *testing.T) {
	page := testPage("page-1", "about", "About Us")
	var receivedPageID string

	pageRepo := &mockSitePageRepository{
		findByID: func(ctx context.Context, pageID string) (*models.SitePage, error) {
			receivedPageID = pageID
			return page, nil
		},
		update: func(ctx context.Context, pageID, title, content string, isPublished bool, updatedByID string) (*models.SitePage, error) {
			return page, nil
		},
	}
	userRepo := &mockUserRepoForSitePage{
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			return true, nil
		},
	}
	svc := NewSitePageService(pageRepo, &mockPageDraftRepoForSitePage{}, userRepo)

	req := &models.UpdateSitePageRequest{
		Title: "Updated Title",
	}

	_, err := svc.Update(context.Background(), "  page-1  ", "user-1", req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if receivedPageID != "page-1" {
		t.Fatalf("expected pageID to be trimmed to 'page-1', got '%s'", receivedPageID)
	}
}

func TestSitePageService_Update_UserIDTrimmed(t *testing.T) {
	page := testPage("page-1", "about", "About Us")
	var receivedUserID string

	pageRepo := &mockSitePageRepository{
		findByID: func(ctx context.Context, pageID string) (*models.SitePage, error) {
			return page, nil
		},
		update: func(ctx context.Context, pageID, title, content string, isPublished bool, updatedByID string) (*models.SitePage, error) {
			return page, nil
		},
	}
	userRepo := &mockUserRepoForSitePage{
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			receivedUserID = userID
			return true, nil
		},
	}
	svc := NewSitePageService(pageRepo, &mockPageDraftRepoForSitePage{}, userRepo)

	req := &models.UpdateSitePageRequest{
		Title: "Updated Title",
	}

	_, err := svc.Update(context.Background(), "page-1", "  user-1  ", req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if receivedUserID != "user-1" {
		t.Fatalf("expected userID to be trimmed to 'user-1', got '%s'", receivedUserID)
	}
}

func TestSitePageService_Update_TitleTrimmed(t *testing.T) {
	page := testPage("page-1", "about", "About Us")
	var receivedTitle string

	pageRepo := &mockSitePageRepository{
		findByID: func(ctx context.Context, pageID string) (*models.SitePage, error) {
			return page, nil
		},
		update: func(ctx context.Context, pageID, title, content string, isPublished bool, updatedByID string) (*models.SitePage, error) {
			receivedTitle = title
			return page, nil
		},
	}
	userRepo := &mockUserRepoForSitePage{
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			return true, nil
		},
	}
	svc := NewSitePageService(pageRepo, &mockPageDraftRepoForSitePage{}, userRepo)

	req := &models.UpdateSitePageRequest{
		Title: "  Updated Title  ",
	}

	_, err := svc.Update(context.Background(), "page-1", "user-1", req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if receivedTitle != "Updated Title" {
		t.Fatalf("expected title to be trimmed to 'Updated Title', got '%s'", receivedTitle)
	}
}

// ===== Delete Tests =====

func TestSitePageService_Delete_Success(t *testing.T) {
	page := testPage("page-1", "about", "About Us")

	pageRepo := &mockSitePageRepository{
		findByID: func(ctx context.Context, pageID string) (*models.SitePage, error) {
			if pageID == "page-1" {
				return page, nil
			}
			return nil, errors.New("not found")
		},
		delete: func(ctx context.Context, pageID string) error {
			return nil
		},
	}
	userRepo := &mockUserRepoForSitePage{
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			return true, nil
		},
	}
	svc := NewSitePageService(pageRepo, &mockPageDraftRepoForSitePage{}, userRepo)

	err := svc.Delete(context.Background(), "page-1", "user-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestSitePageService_Delete_EmptyPageID(t *testing.T) {
	svc := NewSitePageService(&mockSitePageRepository{}, &mockPageDraftRepoForSitePage{}, &mockUserRepoForSitePage{})

	err := svc.Delete(context.Background(), "", "user-1")
	if err != ErrValidation {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

func TestSitePageService_Delete_WhitespaceOnlyPageID(t *testing.T) {
	svc := NewSitePageService(&mockSitePageRepository{}, &mockPageDraftRepoForSitePage{}, &mockUserRepoForSitePage{})

	err := svc.Delete(context.Background(), "   ", "user-1")
	if err != ErrValidation {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

func TestSitePageService_Delete_EmptyUserID(t *testing.T) {
	svc := NewSitePageService(&mockSitePageRepository{}, &mockPageDraftRepoForSitePage{}, &mockUserRepoForSitePage{})

	err := svc.Delete(context.Background(), "page-1", "")
	if err != ErrValidation {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

func TestSitePageService_Delete_WhitespaceOnlyUserID(t *testing.T) {
	svc := NewSitePageService(&mockSitePageRepository{}, &mockPageDraftRepoForSitePage{}, &mockUserRepoForSitePage{})

	err := svc.Delete(context.Background(), "page-1", "   ")
	if err != ErrValidation {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

func TestSitePageService_Delete_PageNotFound(t *testing.T) {
	pageRepo := &mockSitePageRepository{
		findByID: func(ctx context.Context, pageID string) (*models.SitePage, error) {
			return nil, errors.New("not found")
		},
	}
	svc := NewSitePageService(pageRepo, &mockPageDraftRepoForSitePage{}, &mockUserRepoForSitePage{})

	err := svc.Delete(context.Background(), "nonexistent", "user-1")
	if err != ErrPageNotFound {
		t.Fatalf("expected ErrPageNotFound, got %v", err)
	}
}

func TestSitePageService_Delete_NotSiteAdmin(t *testing.T) {
	page := testPage("page-1", "about", "About Us")

	pageRepo := &mockSitePageRepository{
		findByID: func(ctx context.Context, pageID string) (*models.SitePage, error) {
			return page, nil
		},
	}
	userRepo := &mockUserRepoForSitePage{
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			return false, nil
		},
	}
	svc := NewSitePageService(pageRepo, &mockPageDraftRepoForSitePage{}, userRepo)

	err := svc.Delete(context.Background(), "page-1", "user-1")
	if err != ErrForbidden {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestSitePageService_Delete_IsSiteAdminError(t *testing.T) {
	page := testPage("page-1", "about", "About Us")
	repoErr := errors.New("database error")

	pageRepo := &mockSitePageRepository{
		findByID: func(ctx context.Context, pageID string) (*models.SitePage, error) {
			return page, nil
		},
	}
	userRepo := &mockUserRepoForSitePage{
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			return false, repoErr
		},
	}
	svc := NewSitePageService(pageRepo, &mockPageDraftRepoForSitePage{}, userRepo)

	err := svc.Delete(context.Background(), "page-1", "user-1")
	if err != repoErr {
		t.Fatalf("expected repo error, got %v", err)
	}
}

func TestSitePageService_Delete_RepositoryError(t *testing.T) {
	page := testPage("page-1", "about", "About Us")
	repoErr := errors.New("database error")

	pageRepo := &mockSitePageRepository{
		findByID: func(ctx context.Context, pageID string) (*models.SitePage, error) {
			return page, nil
		},
		delete: func(ctx context.Context, pageID string) error {
			return repoErr
		},
	}
	userRepo := &mockUserRepoForSitePage{
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			return true, nil
		},
	}
	svc := NewSitePageService(pageRepo, &mockPageDraftRepoForSitePage{}, userRepo)

	err := svc.Delete(context.Background(), "page-1", "user-1")
	if err != repoErr {
		t.Fatalf("expected repo error, got %v", err)
	}
}

func TestSitePageService_Delete_PageIDTrimmed(t *testing.T) {
	page := testPage("page-1", "about", "About Us")
	var receivedPageID string

	pageRepo := &mockSitePageRepository{
		findByID: func(ctx context.Context, pageID string) (*models.SitePage, error) {
			receivedPageID = pageID
			return page, nil
		},
		delete: func(ctx context.Context, pageID string) error {
			return nil
		},
	}
	userRepo := &mockUserRepoForSitePage{
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			return true, nil
		},
	}
	svc := NewSitePageService(pageRepo, &mockPageDraftRepoForSitePage{}, userRepo)

	err := svc.Delete(context.Background(), "  page-1  ", "user-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if receivedPageID != "page-1" {
		t.Fatalf("expected pageID to be trimmed to 'page-1', got '%s'", receivedPageID)
	}
}

func TestSitePageService_Delete_UserIDTrimmed(t *testing.T) {
	page := testPage("page-1", "about", "About Us")
	var receivedUserID string

	pageRepo := &mockSitePageRepository{
		findByID: func(ctx context.Context, pageID string) (*models.SitePage, error) {
			return page, nil
		},
		delete: func(ctx context.Context, pageID string) error {
			return nil
		},
	}
	userRepo := &mockUserRepoForSitePage{
		isSiteAdmin: func(ctx context.Context, userID string) (bool, error) {
			receivedUserID = userID
			return true, nil
		},
	}
	svc := NewSitePageService(pageRepo, &mockPageDraftRepoForSitePage{}, userRepo)

	err := svc.Delete(context.Background(), "page-1", "  user-1  ")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if receivedUserID != "user-1" {
		t.Fatalf("expected userID to be trimmed to 'user-1', got '%s'", receivedUserID)
	}
}

// ===== GetByID Tests =====

func TestSitePageService_GetByID_Success(t *testing.T) {
	page := testPage("page-1", "about", "About Us")

	pageRepo := &mockSitePageRepository{
		findByID: func(ctx context.Context, pageID string) (*models.SitePage, error) {
			if pageID == "page-1" {
				return page, nil
			}
			return nil, errors.New("not found")
		},
	}
	svc := NewSitePageService(pageRepo, &mockPageDraftRepoForSitePage{}, &mockUserRepoForSitePage{})

	result, err := svc.GetByID(context.Background(), "page-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ID != "page-1" {
		t.Fatalf("expected page-1, got %s", result.ID)
	}
}

func TestSitePageService_GetByID_EmptyPageID(t *testing.T) {
	svc := NewSitePageService(&mockSitePageRepository{}, &mockPageDraftRepoForSitePage{}, &mockUserRepoForSitePage{})

	_, err := svc.GetByID(context.Background(), "")
	if err != ErrValidation {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

func TestSitePageService_GetByID_WhitespaceOnlyPageID(t *testing.T) {
	svc := NewSitePageService(&mockSitePageRepository{}, &mockPageDraftRepoForSitePage{}, &mockUserRepoForSitePage{})

	_, err := svc.GetByID(context.Background(), "   ")
	if err != ErrValidation {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

func TestSitePageService_GetByID_PageNotFound(t *testing.T) {
	pageRepo := &mockSitePageRepository{
		findByID: func(ctx context.Context, pageID string) (*models.SitePage, error) {
			return nil, errors.New("not found")
		},
	}
	svc := NewSitePageService(pageRepo, &mockPageDraftRepoForSitePage{}, &mockUserRepoForSitePage{})

	_, err := svc.GetByID(context.Background(), "nonexistent")
	if err != ErrPageNotFound {
		t.Fatalf("expected ErrPageNotFound, got %v", err)
	}
}

func TestSitePageService_GetByID_PageIDTrimmed(t *testing.T) {
	page := testPage("page-1", "about", "About Us")
	var receivedPageID string

	pageRepo := &mockSitePageRepository{
		findByID: func(ctx context.Context, pageID string) (*models.SitePage, error) {
			receivedPageID = pageID
			return page, nil
		},
	}
	svc := NewSitePageService(pageRepo, &mockPageDraftRepoForSitePage{}, &mockUserRepoForSitePage{})

	result, err := svc.GetByID(context.Background(), "  page-1  ")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ID != "page-1" {
		t.Fatalf("expected page-1, got %s", result.ID)
	}
	if receivedPageID != "page-1" {
		t.Fatalf("expected pageID to be trimmed to 'page-1', got '%s'", receivedPageID)
	}
}
