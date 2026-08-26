package repositories

import (
	"context"
	"time"

	"menunderfire/internal/models"
)

// PageHitRepository defines the interface for page hit data access
type PageHitRepository interface {
	// Record records a new page hit
	Record(ctx context.Context, path string, userID, userEmail, ipAddress, userAgent, referrer, country, region *string) (*models.PageHit, error)

	// UpdateGeo updates geo information for a hit
	UpdateGeo(ctx context.Context, hitID string, country, city *string, lat, lon *float64) error

	// GetTotalCount returns total hit count
	GetTotalCount(ctx context.Context) (int64, error)

	// GetTodayCount returns today's hit count
	GetTodayCount(ctx context.Context) (int64, error)

	// GetStatsByPath returns hit counts grouped by path
	GetStatsByPath(ctx context.Context) ([]models.PageHitStatsResponse, error)

	// GetRecentHits returns recent hits with details
	GetRecentHits(ctx context.Context, limit int) ([]*models.PageHit, error)

	// GetHitsByCountry returns hit counts grouped by country
	GetHitsByCountry(ctx context.Context, from, to *time.Time) ([]models.HitsByCountryResponse, error)

	// GetHitsByRegion returns hit counts grouped by country and region
	GetHitsByRegion(ctx context.Context, from, to *time.Time) ([]models.HitsByRegionResponse, error)

	// GetHitsHourly returns hit counts grouped by hour
	GetHitsHourly(ctx context.Context, from, to *time.Time) ([]models.HitsHourlyResponse, error)

	// GetHitsDaily returns hit counts grouped by day
	GetHitsDaily(ctx context.Context, from, to *time.Time) ([]models.HitsDailyResponse, error)

	// GetRecentHitsInRange returns recent hits filtered by date range
	GetRecentHitsInRange(ctx context.Context, from, to *time.Time, limit int) ([]*models.PageHit, error)
}
