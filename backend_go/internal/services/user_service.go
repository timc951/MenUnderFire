package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"menunderfire/internal/models"
	"menunderfire/internal/repositories"
)

// UserService defines the interface for user-related business logic
type UserService interface {
	// GetByID retrieves a user by their ID
	GetByID(ctx context.Context, userID string) (*models.User, error)

	// GetByExternalID retrieves a user by their external authentication provider ID
	GetByExternalID(ctx context.Context, externalID string) (*models.User, error)

	// Update updates a user's profile information
	Update(ctx context.Context, userID string, req *models.UpdateUserRequest) (*models.User, error)

	// GetPermissions retrieves the user's roles and permissions across the system
	GetPermissions(ctx context.Context, userID string) (*models.UserPermissionsResponse, error)

	// AcceptAgreement records the user's acceptance of the testing agreement with HMAC signature
	AcceptAgreement(ctx context.Context, userID, version, ipAddress, userAgent string) (*models.AcceptAgreementResponse, error)
}

// userService implements the UserService interface
type userService struct {
	userRepo        repositories.UserRepository
	agreementSecret string
}

// NewUserService creates a new UserService implementation
func NewUserService(userRepo repositories.UserRepository, agreementSecret string) UserService {
	return &userService{
		userRepo:        userRepo,
		agreementSecret: agreementSecret,
	}
}

// GetByID retrieves a user by their ID
// Edge cases:
// - UserID is empty or whitespace only -> return ErrValidation
// - User not found -> return ErrUserNotFound
// - Repository error -> return ErrUserNotFound
func (s *userService) GetByID(ctx context.Context, userID string) (*models.User, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, ErrValidation
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// GetByExternalID retrieves a user by their external authentication provider ID
// Edge cases:
// - ExternalID is empty or whitespace only -> return ErrValidation
// - User not found -> return ErrUserNotFound
// - Repository error -> return ErrUserNotFound
func (s *userService) GetByExternalID(ctx context.Context, externalID string) (*models.User, error) {
	externalID = strings.TrimSpace(externalID)
	if externalID == "" {
		return nil, ErrValidation
	}

	user, err := s.userRepo.FindByExternalID(ctx, externalID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// Update updates a user's profile information
// Edge cases:
// - UserID is empty or whitespace only -> return ErrValidation
// - Request is nil -> return ErrValidation
// - DisplayName is empty or whitespace only -> return ErrValidation
// - User not found -> return ErrUserNotFound
// - Repository error during user lookup -> return ErrUserNotFound
// - Repository error during update -> return wrapped error
func (s *userService) Update(ctx context.Context, userID string, req *models.UpdateUserRequest) (*models.User, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, ErrValidation
	}

	if req == nil {
		return nil, ErrValidation
	}

	displayName := strings.TrimSpace(req.DisplayName)
	if displayName == "" {
		return nil, ErrValidation
	}

	// Check if user exists
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Update the user
	user, err := s.userRepo.Update(ctx, userID, displayName)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetPermissions retrieves the user's roles and permissions across the system
// Edge cases:
// - UserID is empty or whitespace only -> return ErrValidation
// - User not found -> return ErrUserNotFound
// - Repository error during user lookup -> return ErrUserNotFound
// - Repository error during isSiteAdmin check -> return wrapped error
// - Repository error during findAdminOrganizationsByUser -> return wrapped error
// - Repository error during findOwnedGroupIDsByUser -> return wrapped error
// - Repository error during findMemberGroupIDsByUser -> return wrapped error
// - User has no permissions -> return response with all empty/false values
func (s *userService) GetPermissions(ctx context.Context, userID string) (*models.UserPermissionsResponse, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, ErrValidation
	}

	// Check if user exists
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Check if user is a site admin
	isSiteAdmin, err := s.userRepo.IsSiteAdmin(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get organization IDs where user is an admin
	adminOrgIDs, err := s.userRepo.FindAdminOrganizationIDs(ctx, userID)
	if err != nil {
		return nil, err
	}
	if adminOrgIDs == nil {
		adminOrgIDs = []string{}
	}

	// Get group IDs where user is an owner
	ownedGroupIDs, err := s.userRepo.FindOwnedGroupIDs(ctx, userID)
	if err != nil {
		return nil, err
	}
	if ownedGroupIDs == nil {
		ownedGroupIDs = []string{}
	}

	// Get group IDs where user is a member
	memberGroupIDs, err := s.userRepo.FindMemberGroupIDs(ctx, userID)
	if err != nil {
		return nil, err
	}
	if memberGroupIDs == nil {
		memberGroupIDs = []string{}
	}

	return &models.UserPermissionsResponse{
		IsSiteAdmin:            isSiteAdmin,
		AdminOfOrganizationIDs: adminOrgIDs,
		OwnedGroupIDs:          ownedGroupIDs,
		MemberGroupIDs:         memberGroupIDs,
	}, nil
}

// AcceptAgreement records the user's acceptance of the testing agreement
// Generates an HMAC-SHA256 signature binding user ID, version, and timestamp
func (s *userService) AcceptAgreement(ctx context.Context, userID, version, ipAddress, userAgent string) (*models.AcceptAgreementResponse, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, ErrValidation
	}

	version = strings.TrimSpace(version)
	if version == "" {
		return nil, ErrValidation
	}

	// Verify user exists
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Check if already accepted this version
	if user.AgreementAcceptedAt != nil && user.AgreementVersion != nil && *user.AgreementVersion == version {
		return nil, fmt.Errorf("%w: agreement already accepted", ErrConflict)
	}

	// Generate HMAC signature: HMAC-SHA256(secret, "userID|version|timestamp")
	now := time.Now().UTC()
	message := fmt.Sprintf("%s|%s|%s", userID, version, now.Format(time.RFC3339))
	mac := hmac.New(sha256.New, []byte(s.agreementSecret))
	mac.Write([]byte(message))
	signature := hex.EncodeToString(mac.Sum(nil))

	// Store in database
	if err := s.userRepo.RecordAgreementAcceptance(ctx, userID, version, signature, ipAddress, userAgent); err != nil {
		return nil, fmt.Errorf("failed to record agreement acceptance: %w", err)
	}

	return &models.AcceptAgreementResponse{
		AcceptedAt: now.Format(time.RFC3339),
		Version:    version,
		Signature:  signature,
	}, nil
}
