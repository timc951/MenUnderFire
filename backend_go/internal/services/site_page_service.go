package services

import (
	"context"
	"strings"

	"menunderfire/internal/models"
	"menunderfire/internal/repositories"
)

// SitePageService defines the interface for site page business logic
type SitePageService interface {
	// List returns all site pages (public)
	List(ctx context.Context) ([]*models.SitePage, error)

	// GetBySlug retrieves a site page by its slug (public)
	GetBySlug(ctx context.Context, slug string) (*models.SitePage, error)

	// Create creates a new site page (Site Admin only)
	Create(ctx context.Context, userID string, req *models.CreateSitePageRequest) (*models.SitePage, error)

	// Update updates a site page (Site Admin only)
	Update(ctx context.Context, pageID string, userID string, req *models.UpdateSitePageRequest) (*models.SitePage, error)

	// Delete deletes a site page (Site Admin only)
	Delete(ctx context.Context, pageID string, userID string) error

	// GetByID retrieves a site page by its ID
	GetByID(ctx context.Context, pageID string) (*models.SitePage, error)

	// CreateDraft creates a temporary preview draft (Site Admin only)
	CreateDraft(ctx context.Context, pageID *string, userID string, req *models.CreatePageDraftRequest) (*models.PageDraft, error)

	// GetDraft retrieves a draft by ID (public - for preview)
	GetDraft(ctx context.Context, draftID string) (*models.PageDraft, error)
}

// sitePageService implements the SitePageService interface
type sitePageService struct {
	sitePageRepo  repositories.SitePageRepository
	pageDraftRepo repositories.PageDraftRepository
	userRepo      repositories.UserRepository
}

// NewSitePageService creates a new SitePageService implementation
func NewSitePageService(
	sitePageRepo repositories.SitePageRepository,
	pageDraftRepo repositories.PageDraftRepository,
	userRepo repositories.UserRepository,
) SitePageService {
	return &sitePageService{
		sitePageRepo:  sitePageRepo,
		pageDraftRepo: pageDraftRepo,
		userRepo:      userRepo,
	}
}

// List returns all site pages (public)
// Edge cases:
// - Repository error -> return wrapped error
// - No pages exist -> return empty slice (not an error)
func (s *sitePageService) List(ctx context.Context) ([]*models.SitePage, error) {
	pages, err := s.sitePageRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return pages, nil
}

// GetBySlug retrieves a site page by its slug (public)
// Edge cases:
// - Slug is empty or whitespace only -> return ErrValidation
// - Page not found -> return ErrPageNotFound
// - Repository error -> return ErrPageNotFound
func (s *sitePageService) GetBySlug(ctx context.Context, slug string) (*models.SitePage, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return nil, ErrValidation
	}

	page, err := s.sitePageRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, ErrPageNotFound
	}

	return page, nil
}

// Create creates a new site page (Site Admin only)
// Edge cases:
// - Request is nil -> return ErrValidation
// - UserID is empty or whitespace only -> return ErrValidation
// - Slug is empty or whitespace only -> return ErrValidation
// - Title is empty or whitespace only -> return ErrValidation
// - User is not a site admin -> return ErrForbidden
// - Repository error during isSiteAdmin check -> return wrapped error
// - Repository error during page creation -> return wrapped error
// - Slug already exists -> return ErrConflict (handled by repository)
func (s *sitePageService) Create(ctx context.Context, userID string, req *models.CreateSitePageRequest) (*models.SitePage, error) {
	if req == nil {
		return nil, ErrValidation
	}

	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, ErrValidation
	}

	slug := strings.TrimSpace(req.Slug)
	if slug == "" {
		return nil, ErrValidation
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, ErrValidation
	}

	// Check if user is a site admin
	isSiteAdmin, err := s.userRepo.IsSiteAdmin(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !isSiteAdmin {
		return nil, ErrForbidden
	}

	// Create the page
	page, err := s.sitePageRepo.Create(ctx, slug, title, req.Content, req.IsPublished, userID)
	if err != nil {
		return nil, err
	}

	return page, nil
}

// Update updates a site page (Site Admin only)
// Edge cases:
// - PageID is empty or whitespace only -> return ErrValidation
// - UserID is empty or whitespace only -> return ErrValidation
// - Request is nil -> return ErrValidation
// - Title is empty or whitespace only -> return ErrValidation
// - Page not found -> return ErrPageNotFound
// - User is not a site admin -> return ErrForbidden
// - Repository error during page lookup -> return ErrPageNotFound
// - Repository error during isSiteAdmin check -> return wrapped error
// - Repository error during page update -> return wrapped error
func (s *sitePageService) Update(ctx context.Context, pageID string, userID string, req *models.UpdateSitePageRequest) (*models.SitePage, error) {
	pageID = strings.TrimSpace(pageID)
	if pageID == "" {
		return nil, ErrValidation
	}

	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, ErrValidation
	}

	if req == nil {
		return nil, ErrValidation
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, ErrValidation
	}

	// Check if page exists
	_, err := s.sitePageRepo.FindByID(ctx, pageID)
	if err != nil {
		return nil, ErrPageNotFound
	}

	// Check if user is a site admin
	isSiteAdmin, err := s.userRepo.IsSiteAdmin(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !isSiteAdmin {
		return nil, ErrForbidden
	}

	// Update the page
	page, err := s.sitePageRepo.Update(ctx, pageID, title, req.Content, req.IsPublished, userID)
	if err != nil {
		return nil, err
	}

	return page, nil
}

// Delete deletes a site page (Site Admin only)
// Edge cases:
// - PageID is empty or whitespace only -> return ErrValidation
// - UserID is empty or whitespace only -> return ErrValidation
// - Page not found -> return ErrPageNotFound
// - User is not a site admin -> return ErrForbidden
// - Repository error during page lookup -> return ErrPageNotFound
// - Repository error during isSiteAdmin check -> return wrapped error
// - Repository error during page deletion -> return wrapped error
func (s *sitePageService) Delete(ctx context.Context, pageID string, userID string) error {
	pageID = strings.TrimSpace(pageID)
	if pageID == "" {
		return ErrValidation
	}

	userID = strings.TrimSpace(userID)
	if userID == "" {
		return ErrValidation
	}

	// Check if page exists
	_, err := s.sitePageRepo.FindByID(ctx, pageID)
	if err != nil {
		return ErrPageNotFound
	}

	// Check if user is a site admin
	isSiteAdmin, err := s.userRepo.IsSiteAdmin(ctx, userID)
	if err != nil {
		return err
	}
	if !isSiteAdmin {
		return ErrForbidden
	}

	// Delete the page
	err = s.sitePageRepo.Delete(ctx, pageID)
	if err != nil {
		return err
	}

	return nil
}

// GetByID retrieves a site page by its ID
// Edge cases:
// - PageID is empty or whitespace only -> return ErrValidation
// - Page not found -> return ErrPageNotFound
// - Repository error -> return ErrPageNotFound
func (s *sitePageService) GetByID(ctx context.Context, pageID string) (*models.SitePage, error) {
	pageID = strings.TrimSpace(pageID)
	if pageID == "" {
		return nil, ErrValidation
	}

	page, err := s.sitePageRepo.FindByID(ctx, pageID)
	if err != nil {
		return nil, ErrPageNotFound
	}

	return page, nil
}

// CreateDraft creates a temporary preview draft (Site Admin only)
// Edge cases:
// - UserID is empty -> return ErrValidation
// - Request is nil -> return ErrValidation
// - Title is empty -> return ErrValidation
// - User is not a site admin -> return ErrForbidden
func (s *sitePageService) CreateDraft(ctx context.Context, pageID *string, userID string, req *models.CreatePageDraftRequest) (*models.PageDraft, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, ErrValidation
	}

	if req == nil {
		return nil, ErrValidation
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, ErrValidation
	}

	// Check if user is a site admin
	isSiteAdmin, err := s.userRepo.IsSiteAdmin(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !isSiteAdmin {
		return nil, ErrForbidden
	}

	// Create the draft
	draft, err := s.pageDraftRepo.Create(ctx, pageID, title, req.Content, userID)
	if err != nil {
		return nil, err
	}

	return draft, nil
}

// GetDraft retrieves a draft by ID (public - for preview)
// Edge cases:
// - DraftID is empty -> return ErrValidation
// - Draft not found or expired -> return ErrDraftNotFound
func (s *sitePageService) GetDraft(ctx context.Context, draftID string) (*models.PageDraft, error) {
	draftID = strings.TrimSpace(draftID)
	if draftID == "" {
		return nil, ErrValidation
	}

	draft, err := s.pageDraftRepo.FindByID(ctx, draftID)
	if err != nil {
		return nil, ErrDraftNotFound
	}

	return draft, nil
}
