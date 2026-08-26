package services

import (
	"context"
	"strings"

	"menunderfire/internal/models"
	"menunderfire/internal/repositories"
)

// ReportService defines the interface for report-related business logic
type ReportService interface {
	// Create creates a new accountability report
	Create(ctx context.Context, userID string, req *models.CreateReportRequest) (*models.Report, error)

	// List returns all reports for a group
	// The requesterID is used to determine visibility of anonymous reports
	List(ctx context.Context, groupID string, requesterID string) ([]*models.Report, error)

	// GetByID retrieves a report by its ID
	GetByID(ctx context.Context, reportID string) (*models.Report, error)
}

// reportService implements the ReportService interface
type reportService struct {
	reportRepo repositories.ReportRepository
	groupRepo  repositories.GroupRepository
}

// NewReportService creates a new ReportService implementation
func NewReportService(
	reportRepo repositories.ReportRepository,
	groupRepo repositories.GroupRepository,
) ReportService {
	return &reportService{
		reportRepo: reportRepo,
		groupRepo:  groupRepo,
	}
}

// Create creates a new accountability report
// Edge cases:
// - Request is nil -> return ErrValidation
// - UserID is empty or whitespace only -> return ErrValidation
// - GroupID in request is empty or whitespace only -> return ErrValidation
// - Title is empty or whitespace only -> return ErrValidation
// - Content is empty or whitespace only -> return ErrValidation
// - Group not found -> return ErrGroupNotFound
// - User is not a member of the group -> return ErrForbidden
// - Repository error during group lookup -> return ErrGroupNotFound
// - Repository error during membership check -> return wrapped error
// - Repository error during report creation -> return wrapped error
func (s *reportService) Create(ctx context.Context, userID string, req *models.CreateReportRequest) (*models.Report, error) {
	if req == nil {
		return nil, ErrValidation
	}

	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, ErrValidation
	}

	groupID := strings.TrimSpace(req.GroupID)
	if groupID == "" {
		return nil, ErrValidation
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, ErrValidation
	}

	content := strings.TrimSpace(req.Content)
	if content == "" {
		return nil, ErrValidation
	}

	// Check if group exists
	_, err := s.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return nil, ErrGroupNotFound
	}

	// Check if user is a member of the group
	member, err := s.groupRepo.FindMember(ctx, groupID, userID)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, ErrForbidden
	}

	// Create the report
	report, err := s.reportRepo.Create(ctx, groupID, userID, title, content, req.IsAnonymousToGroup)
	if err != nil {
		return nil, err
	}

	return report, nil
}

// List returns all reports for a group
// The requesterID is used to determine visibility of anonymous reports
// Edge cases:
// - GroupID is empty or whitespace only -> return ErrValidation
// - RequesterID is empty or whitespace only -> return ErrValidation
// - Group not found -> return ErrGroupNotFound
// - Requester is not a member of the group -> return ErrForbidden
// - Repository error during group lookup -> return ErrGroupNotFound
// - Repository error during membership check -> return wrapped error
// - Repository error during report lookup -> return wrapped error
// - No reports exist -> return empty slice (not an error)
func (s *reportService) List(ctx context.Context, groupID string, requesterID string) ([]*models.Report, error) {
	groupID = strings.TrimSpace(groupID)
	if groupID == "" {
		return nil, ErrValidation
	}

	requesterID = strings.TrimSpace(requesterID)
	if requesterID == "" {
		return nil, ErrValidation
	}

	// Check if group exists
	_, err := s.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return nil, ErrGroupNotFound
	}

	// Check if requester is a member of the group
	member, err := s.groupRepo.FindMember(ctx, groupID, requesterID)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, ErrForbidden
	}

	// Fetch all reports for the group
	reports, err := s.reportRepo.FindByGroupID(ctx, groupID)
	if err != nil {
		return nil, err
	}

	return reports, nil
}

// GetByID retrieves a report by its ID
// Edge cases:
// - ReportID is empty or whitespace only -> return ErrValidation
// - Report not found -> return ErrReportNotFound
// - Repository error -> return ErrReportNotFound
func (s *reportService) GetByID(ctx context.Context, reportID string) (*models.Report, error) {
	reportID = strings.TrimSpace(reportID)
	if reportID == "" {
		return nil, ErrValidation
	}

	report, err := s.reportRepo.FindByID(ctx, reportID)
	if err != nil {
		return nil, ErrReportNotFound
	}

	return report, nil
}
