package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"menunderfire/internal/models"
	"menunderfire/internal/repositories"
)

// PageHitService defines the interface for page hit tracking business logic
type PageHitService interface {
	// RecordHit records a page hit
	RecordHit(ctx context.Context, path string, userID, userEmail, ipAddress, userAgent, referrer, country, region *string) (*models.PageHit, error)

	// GetSummary returns hit summary stats (site admin only)
	GetSummary(ctx context.Context, userID string) (*models.PageHitSummaryResponse, error)

	// GetRecentHits returns recent hit details (site admin only)
	GetRecentHits(ctx context.Context, userID string, limit int) ([]*models.PageHitDetailResponse, error)

	// GetHitsByCountry returns hits grouped by country (site admin only)
	GetHitsByCountry(ctx context.Context, userID string, from, to *time.Time) ([]models.HitsByCountryResponse, error)

	// GetHitsByRegion returns hits grouped by country and region (site admin only)
	GetHitsByRegion(ctx context.Context, userID string, from, to *time.Time) ([]models.HitsByRegionResponse, error)

	// GetHitsHourly returns hits grouped by hour (site admin only)
	GetHitsHourly(ctx context.Context, userID string, from, to *time.Time) ([]models.HitsHourlyResponse, error)

	// GetHitsDaily returns hits grouped by day (site admin only)
	GetHitsDaily(ctx context.Context, userID string, from, to *time.Time) ([]models.HitsDailyResponse, error)

	// GetRecentHitsInRange returns recent hits filtered by date range (site admin only)
	GetRecentHitsInRange(ctx context.Context, userID string, from, to *time.Time, limit int) ([]*models.PageHitDetailResponse, error)
}

type pageHitService struct {
	hitRepo  repositories.PageHitRepository
	userRepo repositories.UserRepository
}

// NewPageHitService creates a new PageHitService implementation
func NewPageHitService(hitRepo repositories.PageHitRepository, userRepo repositories.UserRepository) PageHitService {
	return &pageHitService{
		hitRepo:  hitRepo,
		userRepo: userRepo,
	}
}

func (s *pageHitService) RecordHit(ctx context.Context, path string, userID, userEmail, ipAddress, userAgent, referrer, country, region *string) (*models.PageHit, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, fmt.Errorf("%w: path is required", ErrValidation)
	}

	hit, err := s.hitRepo.Record(ctx, path, userID, userEmail, ipAddress, userAgent, referrer, country, region)
	if err != nil {
		return nil, fmt.Errorf("failed to record hit: %w", err)
	}

	return hit, nil
}

func (s *pageHitService) GetSummary(ctx context.Context, userID string) (*models.PageHitSummaryResponse, error) {
	isSiteAdmin, err := s.userRepo.IsSiteAdmin(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check site admin: %w", err)
	}
	if !isSiteAdmin {
		return nil, ErrForbidden
	}

	totalHits, err := s.hitRepo.GetTotalCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get total hits: %w", err)
	}

	todayHits, err := s.hitRepo.GetTodayCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get today's hits: %w", err)
	}

	pageStats, err := s.hitRepo.GetStatsByPath(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats by path: %w", err)
	}
	if pageStats == nil {
		pageStats = []models.PageHitStatsResponse{}
	}

	return &models.PageHitSummaryResponse{
		TotalHits: totalHits,
		TodayHits: todayHits,
		PageStats: pageStats,
	}, nil
}

func (s *pageHitService) GetRecentHits(ctx context.Context, userID string, limit int) ([]*models.PageHitDetailResponse, error) {
	isSiteAdmin, err := s.userRepo.IsSiteAdmin(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check site admin: %w", err)
	}
	if !isSiteAdmin {
		return nil, ErrForbidden
	}

	if limit <= 0 || limit > 100 {
		limit = 50
	}

	hits, err := s.hitRepo.GetRecentHits(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent hits: %w", err)
	}

	responses := make([]*models.PageHitDetailResponse, len(hits))
	for i, hit := range hits {
		responses[i] = hit.ToDetailResponse()
	}

	return responses, nil
}

func (s *pageHitService) checkSiteAdmin(ctx context.Context, userID string) error {
	isSiteAdmin, err := s.userRepo.IsSiteAdmin(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check site admin: %w", err)
	}
	if !isSiteAdmin {
		return ErrForbidden
	}
	return nil
}

func (s *pageHitService) GetHitsByCountry(ctx context.Context, userID string, from, to *time.Time) ([]models.HitsByCountryResponse, error) {
	if err := s.checkSiteAdmin(ctx, userID); err != nil {
		return nil, err
	}
	results, err := s.hitRepo.GetHitsByCountry(ctx, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to get hits by country: %w", err)
	}
	if results == nil {
		return []models.HitsByCountryResponse{}, nil
	}
	return results, nil
}

func (s *pageHitService) GetHitsByRegion(ctx context.Context, userID string, from, to *time.Time) ([]models.HitsByRegionResponse, error) {
	if err := s.checkSiteAdmin(ctx, userID); err != nil {
		return nil, err
	}
	results, err := s.hitRepo.GetHitsByRegion(ctx, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to get hits by region: %w", err)
	}
	if results == nil {
		return []models.HitsByRegionResponse{}, nil
	}
	return results, nil
}

func (s *pageHitService) GetHitsHourly(ctx context.Context, userID string, from, to *time.Time) ([]models.HitsHourlyResponse, error) {
	if err := s.checkSiteAdmin(ctx, userID); err != nil {
		return nil, err
	}
	results, err := s.hitRepo.GetHitsHourly(ctx, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to get hourly hits: %w", err)
	}
	if results == nil {
		return []models.HitsHourlyResponse{}, nil
	}
	return results, nil
}

func (s *pageHitService) GetHitsDaily(ctx context.Context, userID string, from, to *time.Time) ([]models.HitsDailyResponse, error) {
	if err := s.checkSiteAdmin(ctx, userID); err != nil {
		return nil, err
	}
	results, err := s.hitRepo.GetHitsDaily(ctx, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily hits: %w", err)
	}
	if results == nil {
		return []models.HitsDailyResponse{}, nil
	}
	return results, nil
}

func (s *pageHitService) GetRecentHitsInRange(ctx context.Context, userID string, from, to *time.Time, limit int) ([]*models.PageHitDetailResponse, error) {
	if err := s.checkSiteAdmin(ctx, userID); err != nil {
		return nil, err
	}

	if limit <= 0 || limit > 500 {
		limit = 100
	}

	hits, err := s.hitRepo.GetRecentHitsInRange(ctx, from, to, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent hits in range: %w", err)
	}

	responses := make([]*models.PageHitDetailResponse, len(hits))
	for i, hit := range hits {
		responses[i] = hit.ToDetailResponse()
	}

	return responses, nil
}
