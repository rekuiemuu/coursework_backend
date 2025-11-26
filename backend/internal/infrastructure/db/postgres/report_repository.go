package postgres

import (
	"context"
	"database/sql"

	"github.com/project-capillary/backend/internal/domain/entities"
)

type ReportRepositoryImpl struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepositoryImpl {
	return &ReportRepositoryImpl{db: db}
}

func (r *ReportRepositoryImpl) Create(ctx context.Context, report *entities.Report) error {
	query := `
		INSERT INTO reports (id, examination_id, title, content, summary, diagnosis, recommendations, generated_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.ExecContext(ctx, query,
		report.ID, report.ExaminationID, report.Title, report.Content, report.Summary,
		report.Diagnosis, report.Recommendations, report.GeneratedBy, report.CreatedAt, report.UpdatedAt)
	return err
}

func (r *ReportRepositoryImpl) GetByID(ctx context.Context, id string) (*entities.Report, error) {
	query := `
		SELECT id, examination_id, title, content, summary, diagnosis, recommendations, generated_by, created_at, updated_at
		FROM reports WHERE id = $1
	`
	report := &entities.Report{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&report.ID, &report.ExaminationID, &report.Title, &report.Content, &report.Summary,
		&report.Diagnosis, &report.Recommendations, &report.GeneratedBy, &report.CreatedAt, &report.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return report, err
}

func (r *ReportRepositoryImpl) Update(ctx context.Context, report *entities.Report) error {
	query := `
		UPDATE reports 
		SET content = $2, summary = $3, diagnosis = $4, recommendations = $5, updated_at = $6
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		report.ID, report.Content, report.Summary, report.Diagnosis,
		report.Recommendations, report.UpdatedAt)
	return err
}

func (r *ReportRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM reports WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *ReportRepositoryImpl) GetByExaminationID(ctx context.Context, examinationID string) (*entities.Report, error) {
	query := `
		SELECT id, examination_id, title, content, summary, diagnosis, recommendations, generated_by, created_at, updated_at
		FROM reports WHERE examination_id = $1 ORDER BY created_at DESC LIMIT 1
	`
	report := &entities.Report{}
	err := r.db.QueryRowContext(ctx, query, examinationID).Scan(
		&report.ID, &report.ExaminationID, &report.Title, &report.Content, &report.Summary,
		&report.Diagnosis, &report.Recommendations, &report.GeneratedBy, &report.CreatedAt, &report.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return report, err
}

func (r *ReportRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*entities.Report, error) {
	query := `
		SELECT id, examination_id, title, content, summary, diagnosis, recommendations, generated_by, created_at, updated_at
		FROM reports ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []*entities.Report
	for rows.Next() {
		report := &entities.Report{}
		err := rows.Scan(&report.ID, &report.ExaminationID, &report.Title, &report.Content,
			&report.Summary, &report.Diagnosis, &report.Recommendations, &report.GeneratedBy,
			&report.CreatedAt, &report.UpdatedAt)
		if err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}
	return reports, nil
}
