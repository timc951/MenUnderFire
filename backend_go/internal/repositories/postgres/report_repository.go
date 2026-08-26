package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"menunderfire/internal/models"
	"menunderfire/internal/repositories"
)

// reportRepository is the PostgreSQL implementation of ReportRepository
type reportRepository struct {
	db *sql.DB
}

// NewReportRepository creates a new instance of ReportRepository
func NewReportRepository(db *sql.DB) repositories.ReportRepository {
	return &reportRepository{db: db}
}

// FindByID retrieves a report by its ID
func (r *reportRepository) FindByID(ctx context.Context, reportID string) (*models.Report, error) {
	query := `
		SELECT id, user_id, group_id, title, content, is_anonymous_to_group, created_at
		FROM reports
		WHERE id = $1
	`

	var report models.Report
	err := r.db.QueryRowContext(ctx, query, reportID).Scan(
		&report.ID,
		&report.ReporterID,
		&report.GroupID,
		&report.Title,
		&report.Content,
		&report.IsAnonymousToGroup,
		&report.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("report not found: %s", reportID)
	}
	if err != nil {
		return nil, fmt.Errorf("error finding report by ID: %w", err)
	}

	return &report, nil
}

// FindByGroupID retrieves all reports for a group with reporter details
func (r *reportRepository) FindByGroupID(ctx context.Context, groupID string) ([]*models.Report, error) {
	query := `
		SELECT
			r.id,
			r.user_id,
			r.group_id,
			r.title,
			r.content,
			r.is_anonymous_to_group,
			r.created_at,
			u.id,
			u.email,
			u.display_name,
			u.external_id,
			u.created_at,
			u.updated_at
		FROM reports r
		INNER JOIN users u ON r.user_id = u.id
		WHERE r.group_id = $1
		ORDER BY r.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, groupID)
	if err != nil {
		return nil, fmt.Errorf("error finding reports by group ID: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var reports []*models.Report
	for rows.Next() {
		var report models.Report
		var reporter models.User

		if err := rows.Scan(
			&report.ID,
			&report.ReporterID,
			&report.GroupID,
			&report.Title,
			&report.Content,
			&report.IsAnonymousToGroup,
			&report.CreatedAt,
			&reporter.ID,
			&reporter.Email,
			&reporter.DisplayName,
			&reporter.ExternalID,
			&reporter.CreatedAt,
			&reporter.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning report: %w", err)
		}

		report.Reporter = &reporter
		reports = append(reports, &report)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating reports: %w", err)
	}

	return reports, nil
}

// Create creates a new accountability report
func (r *reportRepository) Create(ctx context.Context, groupID, reporterID, title, content string, isAnonymous bool) (*models.Report, error) {
	query := `
		INSERT INTO reports (group_id, user_id, title, content, is_anonymous_to_group)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, group_id, title, content, is_anonymous_to_group, created_at
	`

	var report models.Report
	err := r.db.QueryRowContext(ctx, query, groupID, reporterID, title, content, isAnonymous).Scan(
		&report.ID,
		&report.ReporterID,
		&report.GroupID,
		&report.Title,
		&report.Content,
		&report.IsAnonymousToGroup,
		&report.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error creating report: %w", err)
	}

	return &report, nil
}
