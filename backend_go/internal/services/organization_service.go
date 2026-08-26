package services

import (
	"context"
	"strings"

	"menunderfire/internal/models"
	"menunderfire/internal/repositories"
)

// OrganizationService defines the interface for organization-related business logic
type OrganizationService interface {
	// List returns all organizations the user belongs to
	List(ctx context.Context, userID string) ([]*models.Organization, error)

	// ListAll returns all organizations (Site Admin only)
	ListAll(ctx context.Context) ([]*models.Organization, error)

	// Create creates a new organization (Site Admin only)
	Create(ctx context.Context, userID string, req *models.CreateOrganizationRequest) (*models.Organization, error)

	// GetByID retrieves an organization by its ID
	GetByID(ctx context.Context, orgID string) (*models.Organization, error)

	// Update updates an organization (Admin only)
	Update(ctx context.Context, orgID string, userID string, req *models.UpdateOrganizationRequest) (*models.Organization, error)

	// ListAdmins returns all admins of an organization
	ListAdmins(ctx context.Context, orgID string) ([]*models.OrganizationAdmin, error)

	// ListGroups returns all groups within an organization
	ListGroups(ctx context.Context, orgID string) ([]*models.Group, error)

	// IsAdmin checks if a user is an admin of an organization
	IsAdmin(ctx context.Context, orgID string, userID string) (bool, error)

	// IsMember checks if a user belongs to an organization (via any group)
	IsMember(ctx context.Context, orgID string, userID string) (bool, error)
}

// organizationService implements the OrganizationService interface
type organizationService struct {
	orgRepo   repositories.OrganizationRepository
	groupRepo repositories.GroupRepository
	userRepo  repositories.UserRepository
}

// NewOrganizationService creates a new OrganizationService implementation
func NewOrganizationService(
	orgRepo repositories.OrganizationRepository,
	groupRepo repositories.GroupRepository,
	userRepo repositories.UserRepository,
) OrganizationService {
	return &organizationService{
		orgRepo:   orgRepo,
		groupRepo: groupRepo,
		userRepo:  userRepo,
	}
}

// List returns all organizations the user belongs to
// Edge cases:
// - UserID is empty -> return ErrValidation
// - Repository error -> return wrapped error
// - User belongs to no organizations -> return empty slice (not an error)
func (s *organizationService) List(ctx context.Context, userID string) ([]*models.Organization, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, ErrValidation
	}

	orgs, err := s.orgRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return orgs, nil
}

// ListAll returns all organizations (Site Admin only)
// Edge cases:
// - Repository error -> return wrapped error
// - No organizations exist -> return empty slice (not an error)
// Note: Authorization should be handled at the route/handler level
func (s *organizationService) ListAll(ctx context.Context) ([]*models.Organization, error) {
	orgs, err := s.orgRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return orgs, nil
}

// Create creates a new organization (Site Admin only)
// Edge cases:
// - Request is nil -> return ErrValidation
// - Name is empty or whitespace only -> return ErrValidation
// - UserID is empty -> return ErrValidation
// - User is not a site admin -> return ErrForbidden
// - Repository error during isSiteAdmin check -> return wrapped error
// - Repository error during creation -> return wrapped error
func (s *organizationService) Create(ctx context.Context, userID string, req *models.CreateOrganizationRequest) (*models.Organization, error) {
	if req == nil {
		return nil, ErrValidation
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrValidation
	}

	userID = strings.TrimSpace(userID)
	if userID == "" {
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

	// Create the organization
	org, err := s.orgRepo.Create(ctx, name, req.Description, userID)
	if err != nil {
		return nil, err
	}

	return org, nil
}

// GetByID retrieves an organization by its ID
// Edge cases:
// - OrgID is empty or whitespace only -> return ErrValidation
// - Organization not found -> return ErrOrganizationNotFound
// - Repository error -> return wrapped error
func (s *organizationService) GetByID(ctx context.Context, orgID string) (*models.Organization, error) {
	orgID = strings.TrimSpace(orgID)
	if orgID == "" {
		return nil, ErrValidation
	}

	org, err := s.orgRepo.FindByID(ctx, orgID)
	if err != nil {
		return nil, ErrOrganizationNotFound
	}

	return org, nil
}

// Update updates an organization (Admin only)
// Edge cases:
// - OrgID is empty or whitespace only -> return ErrValidation
// - UserID is empty or whitespace only -> return ErrValidation
// - Request is nil -> return ErrValidation
// - Name is empty or whitespace only -> return ErrValidation
// - Organization not found -> return ErrOrganizationNotFound
// - User is not a site admin and not an org admin -> return ErrForbidden
// - Repository error during isSiteAdmin check -> return wrapped error
// - Repository error during isAdmin check -> return wrapped error
// - Repository error during update -> return wrapped error
func (s *organizationService) Update(ctx context.Context, orgID string, userID string, req *models.UpdateOrganizationRequest) (*models.Organization, error) {
	orgID = strings.TrimSpace(orgID)
	if orgID == "" {
		return nil, ErrValidation
	}

	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, ErrValidation
	}

	if req == nil {
		return nil, ErrValidation
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrValidation
	}

	// Check if organization exists
	_, err := s.orgRepo.FindByID(ctx, orgID)
	if err != nil {
		return nil, ErrOrganizationNotFound
	}

	// Check authorization: must be site admin or org admin
	isSiteAdmin, err := s.userRepo.IsSiteAdmin(ctx, userID)
	if err != nil {
		return nil, err
	}

	if !isSiteAdmin {
		admin, err := s.orgRepo.FindAdmin(ctx, orgID, userID)
		if err != nil {
			return nil, err
		}
		if admin == nil {
			return nil, ErrForbidden
		}
	}

	// Update the organization
	org, err := s.orgRepo.Update(ctx, orgID, name, req.Description)
	if err != nil {
		return nil, err
	}

	return org, nil
}

// ListAdmins returns all admins of an organization
// Edge cases:
// - OrgID is empty or whitespace only -> return ErrValidation
// - Organization not found -> return ErrOrganizationNotFound
// - Repository error during org lookup -> return ErrOrganizationNotFound
// - Repository error during admin lookup -> return wrapped error
// - No admins exist -> return empty slice (not an error)
func (s *organizationService) ListAdmins(ctx context.Context, orgID string) ([]*models.OrganizationAdmin, error) {
	orgID = strings.TrimSpace(orgID)
	if orgID == "" {
		return nil, ErrValidation
	}

	// Check if organization exists
	_, err := s.orgRepo.FindByID(ctx, orgID)
	if err != nil {
		return nil, ErrOrganizationNotFound
	}

	admins, err := s.orgRepo.FindAdmins(ctx, orgID)
	if err != nil {
		return nil, err
	}

	return admins, nil
}

// ListGroups returns all groups within an organization
// Edge cases:
// - OrgID is empty or whitespace only -> return ErrValidation
// - Organization not found -> return ErrOrganizationNotFound
// - Repository error during org lookup -> return ErrOrganizationNotFound
// - Repository error during group lookup -> return wrapped error
// - No groups exist -> return empty slice (not an error)
func (s *organizationService) ListGroups(ctx context.Context, orgID string) ([]*models.Group, error) {
	orgID = strings.TrimSpace(orgID)
	if orgID == "" {
		return nil, ErrValidation
	}

	// Check if organization exists
	_, err := s.orgRepo.FindByID(ctx, orgID)
	if err != nil {
		return nil, ErrOrganizationNotFound
	}

	groups, err := s.groupRepo.FindByOrganizationID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

// IsAdmin checks if a user is an admin of an organization
// Edge cases:
// - OrgID is empty or whitespace only -> return ErrValidation
// - UserID is empty or whitespace only -> return ErrValidation
// - Repository error -> return wrapped error
// - User is not an admin -> return false, nil
func (s *organizationService) IsAdmin(ctx context.Context, orgID string, userID string) (bool, error) {
	orgID = strings.TrimSpace(orgID)
	if orgID == "" {
		return false, ErrValidation
	}

	userID = strings.TrimSpace(userID)
	if userID == "" {
		return false, ErrValidation
	}

	admin, err := s.orgRepo.FindAdmin(ctx, orgID, userID)
	if err != nil {
		return false, err
	}

	return admin != nil, nil
}

// IsMember checks if a user belongs to an organization (via any group)
// Edge cases:
// - OrgID is empty or whitespace only -> return ErrValidation
// - UserID is empty or whitespace only -> return ErrValidation
// - Repository error -> return wrapped error
// - User is not a member -> return false, nil
func (s *organizationService) IsMember(ctx context.Context, orgID string, userID string) (bool, error) {
	orgID = strings.TrimSpace(orgID)
	if orgID == "" {
		return false, ErrValidation
	}

	userID = strings.TrimSpace(userID)
	if userID == "" {
		return false, ErrValidation
	}

	isMember, err := s.orgRepo.IsMember(ctx, orgID, userID)
	if err != nil {
		return false, err
	}

	return isMember, nil
}
