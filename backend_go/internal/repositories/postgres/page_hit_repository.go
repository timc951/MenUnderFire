package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"menunderfire/internal/models"
	"menunderfire/internal/repositories"
)

type pageHitRepository struct {
	db *sql.DB
}

// NewPageHitRepository creates a new instance of PageHitRepository
func NewPageHitRepository(db *sql.DB) repositories.PageHitRepository {
	return &pageHitRepository{db: db}
}

func (r *pageHitRepository) Record(ctx context.Context, path string, userID, userEmail, ipAddress, userAgent, referrer, country, region *string) (*models.PageHit, error) {
	query := `
		INSERT INTO page_hits (path, user_id, user_email, ip_address, user_agent, referrer, country, region)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, path, user_id, user_email, ip_address, user_agent, referrer, country, region, city, latitude, longitude, created_at
	`

	var hit models.PageHit
	var uid, uemail, ip, ua, ref, ctry, reg, city sql.NullString
	var lat, lon sql.NullFloat64

	if userID != nil {
		uid = sql.NullString{String: *userID, Valid: true}
	}
	if userEmail != nil {
		uemail = sql.NullString{String: *userEmail, Valid: true}
	}
	if ipAddress != nil {
		ip = sql.NullString{String: *ipAddress, Valid: true}
	}
	if userAgent != nil {
		ua = sql.NullString{String: *userAgent, Valid: true}
	}
	if referrer != nil {
		ref = sql.NullString{String: *referrer, Valid: true}
	}
	if country != nil {
		ctry = sql.NullString{String: *country, Valid: true}
	}
	if region != nil {
		reg = sql.NullString{String: *region, Valid: true}
	}

	err := r.db.QueryRowContext(ctx, query, path, uid, uemail, ip, ua, ref, ctry, reg).Scan(
		&hit.ID, &hit.Path, &uid, &uemail, &ip, &ua, &ref, &ctry, &reg, &city, &lat, &lon, &hit.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error recording page hit: %w", err)
	}

	if uid.Valid {
		hit.UserID = &uid.String
	}
	if uemail.Valid {
		hit.UserEmail = &uemail.String
	}
	if ip.Valid {
		hit.IPAddress = &ip.String
	}
	if ua.Valid {
		hit.UserAgent = &ua.String
	}
	if ref.Valid {
		hit.Referrer = &ref.String
	}
	if ctry.Valid {
		hit.Country = &ctry.String
	}
	if reg.Valid {
		hit.Region = &reg.String
	}

	return &hit, nil
}

func (r *pageHitRepository) UpdateGeo(ctx context.Context, hitID string, country, city *string, lat, lon *float64) error {
	query := `
		UPDATE page_hits
		SET country = $2, city = $3, latitude = $4, longitude = $5
		WHERE id = $1
	`

	var c, ci sql.NullString
	var la, lo sql.NullFloat64
	if country != nil {
		c = sql.NullString{String: *country, Valid: true}
	}
	if city != nil {
		ci = sql.NullString{String: *city, Valid: true}
	}
	if lat != nil {
		la = sql.NullFloat64{Float64: *lat, Valid: true}
	}
	if lon != nil {
		lo = sql.NullFloat64{Float64: *lon, Valid: true}
	}

	_, err := r.db.ExecContext(ctx, query, hitID, c, ci, la, lo)
	if err != nil {
		return fmt.Errorf("error updating page hit geo: %w", err)
	}
	return nil
}

func (r *pageHitRepository) GetTotalCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM page_hits`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error getting total hit count: %w", err)
	}
	return count, nil
}

func (r *pageHitRepository) GetTodayCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM page_hits WHERE created_at >= CURRENT_DATE`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error getting today's hit count: %w", err)
	}
	return count, nil
}

func (r *pageHitRepository) GetStatsByPath(ctx context.Context) ([]models.PageHitStatsResponse, error) {
	query := `
		SELECT path, COUNT(*) as hit_count, COUNT(DISTINCT COALESCE(ip_address, id::text)) as unique_hits
		FROM page_hits
		GROUP BY path
		ORDER BY hit_count DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error getting stats by path: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var stats []models.PageHitStatsResponse
	for rows.Next() {
		var s models.PageHitStatsResponse
		if err := rows.Scan(&s.Path, &s.HitCount, &s.UniqueHits); err != nil {
			return nil, fmt.Errorf("error scanning page hit stats: %w", err)
		}
		stats = append(stats, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating page hit stats: %w", err)
	}

	return stats, nil
}

func (r *pageHitRepository) GetRecentHits(ctx context.Context, limit int) ([]*models.PageHit, error) {
	query := `
		SELECT id, path, user_id, user_email, ip_address, user_agent, referrer, country, region, city, latitude, longitude, created_at
		FROM page_hits
		ORDER BY created_at DESC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("error getting recent hits: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var hits []*models.PageHit
	for rows.Next() {
		var hit models.PageHit
		var uid, uemail, ip, ua, ref, country, region, city sql.NullString
		var lat, lon sql.NullFloat64
		if err := rows.Scan(
			&hit.ID, &hit.Path, &uid, &uemail, &ip, &ua, &ref, &country, &region, &city, &lat, &lon, &hit.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning page hit: %w", err)
		}
		if uid.Valid {
			hit.UserID = &uid.String
		}
		if uemail.Valid {
			hit.UserEmail = &uemail.String
		}
		if ip.Valid {
			hit.IPAddress = &ip.String
		}
		if ua.Valid {
			hit.UserAgent = &ua.String
		}
		if ref.Valid {
			hit.Referrer = &ref.String
		}
		if country.Valid {
			hit.Country = &country.String
		}
		if region.Valid {
			hit.Region = &region.String
		}
		if city.Valid {
			hit.City = &city.String
		}
		if lat.Valid {
			hit.Latitude = &lat.Float64
		}
		if lon.Valid {
			hit.Longitude = &lon.Float64
		}
		hits = append(hits, &hit)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating page hits: %w", err)
	}

	return hits, nil
}

// buildDateRangeClause returns a WHERE/AND clause and args for optional date range filtering.
// paramOffset is the starting $N parameter index.
func buildDateRangeClause(from, to *time.Time, paramOffset int, hasExistingWhere bool) (string, []interface{}) {
	clause := ""
	var args []interface{}
	if from != nil {
		if hasExistingWhere || len(args) > 0 {
			clause += fmt.Sprintf(" AND created_at >= $%d", paramOffset)
		} else {
			clause += fmt.Sprintf(" WHERE created_at >= $%d", paramOffset)
		}
		args = append(args, *from)
		paramOffset++
	}
	if to != nil {
		if hasExistingWhere || len(args) > 0 {
			clause += fmt.Sprintf(" AND created_at <= $%d", paramOffset)
		} else {
			clause += fmt.Sprintf(" WHERE created_at <= $%d", paramOffset)
		}
		args = append(args, *to)
	}
	return clause, args
}

func (r *pageHitRepository) GetHitsByCountry(ctx context.Context, from, to *time.Time) ([]models.HitsByCountryResponse, error) {
	dateClause, dateArgs := buildDateRangeClause(from, to, 1, false)
	// #nosec G201 -- only the generated $N date-range clause is interpolated; values are bound
	query := fmt.Sprintf(`
		SELECT COALESCE(country, 'Unknown') as country, COUNT(*) as hit_count
		FROM page_hits
		%s
		GROUP BY COALESCE(country, 'Unknown')
		ORDER BY hit_count DESC
	`, dateClause)

	rows, err := r.db.QueryContext(ctx, query, dateArgs...)
	if err != nil {
		return nil, fmt.Errorf("error getting hits by country: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var results []models.HitsByCountryResponse
	for rows.Next() {
		var r models.HitsByCountryResponse
		if err := rows.Scan(&r.Country, &r.HitCount); err != nil {
			return nil, fmt.Errorf("error scanning hits by country: %w", err)
		}
		results = append(results, r)
	}
	return results, rows.Err()
}

func (r *pageHitRepository) GetHitsByRegion(ctx context.Context, from, to *time.Time) ([]models.HitsByRegionResponse, error) {
	dateClause, dateArgs := buildDateRangeClause(from, to, 1, false)
	// #nosec G201 -- only the generated $N date-range clause is interpolated; values are bound
	query := fmt.Sprintf(`
		SELECT COALESCE(country, 'Unknown') as country, COALESCE(region, 'Unknown') as region, COUNT(*) as hit_count
		FROM page_hits
		%s
		GROUP BY COALESCE(country, 'Unknown'), COALESCE(region, 'Unknown')
		ORDER BY hit_count DESC
	`, dateClause)

	rows, err := r.db.QueryContext(ctx, query, dateArgs...)
	if err != nil {
		return nil, fmt.Errorf("error getting hits by region: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var results []models.HitsByRegionResponse
	for rows.Next() {
		var r models.HitsByRegionResponse
		if err := rows.Scan(&r.Country, &r.Region, &r.HitCount); err != nil {
			return nil, fmt.Errorf("error scanning hits by region: %w", err)
		}
		results = append(results, r)
	}
	return results, rows.Err()
}

func (r *pageHitRepository) GetHitsHourly(ctx context.Context, from, to *time.Time) ([]models.HitsHourlyResponse, error) {
	// Default to last 24 hours if no range specified
	if from == nil && to == nil {
		defaultFrom := time.Now().Add(-24 * time.Hour)
		from = &defaultFrom
	}
	dateClause, dateArgs := buildDateRangeClause(from, to, 1, false)
	// #nosec G201 -- only the generated $N date-range clause is interpolated; values are bound
	query := fmt.Sprintf(`
		SELECT to_char(created_at, 'YYYY-MM-DD HH24:00') as hour, COUNT(*) as hit_count
		FROM page_hits
		%s
		GROUP BY to_char(created_at, 'YYYY-MM-DD HH24:00')
		ORDER BY hour
	`, dateClause)

	rows, err := r.db.QueryContext(ctx, query, dateArgs...)
	if err != nil {
		return nil, fmt.Errorf("error getting hourly hits: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var results []models.HitsHourlyResponse
	for rows.Next() {
		var r models.HitsHourlyResponse
		if err := rows.Scan(&r.Hour, &r.HitCount); err != nil {
			return nil, fmt.Errorf("error scanning hourly hits: %w", err)
		}
		results = append(results, r)
	}
	return results, rows.Err()
}

func (r *pageHitRepository) GetHitsDaily(ctx context.Context, from, to *time.Time) ([]models.HitsDailyResponse, error) {
	// Default to last 30 days if no range specified
	if from == nil && to == nil {
		defaultFrom := time.Now().AddDate(0, 0, -30)
		from = &defaultFrom
	}
	dateClause, dateArgs := buildDateRangeClause(from, to, 1, false)
	// #nosec G201 -- only the generated $N date-range clause is interpolated; values are bound
	query := fmt.Sprintf(`
		SELECT to_char(created_at, 'YYYY-MM-DD') as date, COUNT(*) as hit_count
		FROM page_hits
		%s
		GROUP BY to_char(created_at, 'YYYY-MM-DD')
		ORDER BY date
	`, dateClause)

	rows, err := r.db.QueryContext(ctx, query, dateArgs...)
	if err != nil {
		return nil, fmt.Errorf("error getting daily hits: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var results []models.HitsDailyResponse
	for rows.Next() {
		var r models.HitsDailyResponse
		if err := rows.Scan(&r.Date, &r.HitCount); err != nil {
			return nil, fmt.Errorf("error scanning daily hits: %w", err)
		}
		results = append(results, r)
	}
	return results, rows.Err()
}

func (r *pageHitRepository) GetRecentHitsInRange(ctx context.Context, from, to *time.Time, limit int) ([]*models.PageHit, error) {
	dateClause, dateArgs := buildDateRangeClause(from, to, 1, false)
	limitParam := len(dateArgs) + 1
	// #nosec G201 -- only the generated $N date-range clause is interpolated; values are bound
	query := fmt.Sprintf(`
		SELECT id, path, user_id, user_email, ip_address, user_agent, referrer, country, region, city, latitude, longitude, created_at
		FROM page_hits
		%s
		ORDER BY created_at DESC
		LIMIT $%d
	`, dateClause, limitParam)

	allArgs := append(dateArgs, limit)
	rows, err := r.db.QueryContext(ctx, query, allArgs...)
	if err != nil {
		return nil, fmt.Errorf("error getting recent hits in range: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var hits []*models.PageHit
	for rows.Next() {
		var hit models.PageHit
		var uid, uemail, ip, ua, ref, country, region, city sql.NullString
		var lat, lon sql.NullFloat64
		if err := rows.Scan(
			&hit.ID, &hit.Path, &uid, &uemail, &ip, &ua, &ref, &country, &region, &city, &lat, &lon, &hit.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning page hit: %w", err)
		}
		if uid.Valid {
			hit.UserID = &uid.String
		}
		if uemail.Valid {
			hit.UserEmail = &uemail.String
		}
		if ip.Valid {
			hit.IPAddress = &ip.String
		}
		if ua.Valid {
			hit.UserAgent = &ua.String
		}
		if ref.Valid {
			hit.Referrer = &ref.String
		}
		if country.Valid {
			hit.Country = &country.String
		}
		if region.Valid {
			hit.Region = &region.String
		}
		if city.Valid {
			hit.City = &city.String
		}
		if lat.Valid {
			hit.Latitude = &lat.Float64
		}
		if lon.Valid {
			hit.Longitude = &lon.Float64
		}
		hits = append(hits, &hit)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating page hits: %w", err)
	}

	return hits, nil
}
