package postgres

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/project-capillary/backend/internal/domain/entities"
)

type AnalysisRepositoryImpl struct {
	db *sql.DB
}

func NewAnalysisRepository(db *sql.DB) *AnalysisRepositoryImpl {
	return &AnalysisRepositoryImpl{db: db}
}

func (r *AnalysisRepositoryImpl) Create(ctx context.Context, analysis *entities.Analysis) error {
	metricsJSON, err := json.Marshal(analysis.Metrics)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO analyses (id, examination_id, image_id, status, metrics, error_message, created_at, updated_at, completed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err = r.db.ExecContext(ctx, query,
		analysis.ID, analysis.ExaminationID, analysis.ImageID, analysis.Status,
		metricsJSON, analysis.ErrorMessage, analysis.CreatedAt, analysis.UpdatedAt, analysis.CompletedAt)
	return err
}

func (r *AnalysisRepositoryImpl) GetByID(ctx context.Context, id string) (*entities.Analysis, error) {
	query := `
		SELECT id, examination_id, image_id, status, metrics, error_message, created_at, updated_at, completed_at
		FROM analyses WHERE id = $1
	`
	analysis := &entities.Analysis{}
	var metricsJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&analysis.ID, &analysis.ExaminationID, &analysis.ImageID, &analysis.Status,
		&metricsJSON, &analysis.ErrorMessage, &analysis.CreatedAt, &analysis.UpdatedAt, &analysis.CompletedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if len(metricsJSON) > 0 {
		if err := json.Unmarshal(metricsJSON, &analysis.Metrics); err != nil {
			return nil, err
		}
	} else {
		analysis.Metrics = make(map[string]interface{})
	}

	return analysis, nil
}

func (r *AnalysisRepositoryImpl) Update(ctx context.Context, analysis *entities.Analysis) error {
	metricsJSON, err := json.Marshal(analysis.Metrics)
	if err != nil {
		return err
	}

	query := `
		UPDATE analyses 
		SET status = $2, metrics = $3, error_message = $4, updated_at = $5, completed_at = $6
		WHERE id = $1
	`
	_, err = r.db.ExecContext(ctx, query,
		analysis.ID, analysis.Status, metricsJSON, analysis.ErrorMessage,
		analysis.UpdatedAt, analysis.CompletedAt)
	return err
}

func (r *AnalysisRepositoryImpl) GetByExaminationID(ctx context.Context, examinationID string) ([]*entities.Analysis, error) {
	query := `
		SELECT id, examination_id, image_id, status, metrics, error_message, created_at, updated_at, completed_at
		FROM analyses WHERE examination_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, examinationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var analyses []*entities.Analysis
	for rows.Next() {
		analysis := &entities.Analysis{}
		var metricsJSON []byte

		err := rows.Scan(&analysis.ID, &analysis.ExaminationID, &analysis.ImageID, &analysis.Status,
			&metricsJSON, &analysis.ErrorMessage, &analysis.CreatedAt, &analysis.UpdatedAt, &analysis.CompletedAt)
		if err != nil {
			return nil, err
		}

		if len(metricsJSON) > 0 {
			if err := json.Unmarshal(metricsJSON, &analysis.Metrics); err != nil {
				return nil, err
			}
		} else {
			analysis.Metrics = make(map[string]interface{})
		}

		analyses = append(analyses, analysis)
	}
	return analyses, nil
}

func (r *AnalysisRepositoryImpl) GetByImageID(ctx context.Context, imageID string) (*entities.Analysis, error) {
	query := `
		SELECT id, examination_id, image_id, status, metrics, error_message, created_at, updated_at, completed_at
		FROM analyses WHERE image_id = $1
	`
	analysis := &entities.Analysis{}
	var metricsJSON []byte

	err := r.db.QueryRowContext(ctx, query, imageID).Scan(
		&analysis.ID, &analysis.ExaminationID, &analysis.ImageID, &analysis.Status,
		&metricsJSON, &analysis.ErrorMessage, &analysis.CreatedAt, &analysis.UpdatedAt, &analysis.CompletedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if len(metricsJSON) > 0 {
		if err := json.Unmarshal(metricsJSON, &analysis.Metrics); err != nil {
			return nil, err
		}
	} else {
		analysis.Metrics = make(map[string]interface{})
	}

	return analysis, nil
}

func (r *AnalysisRepositoryImpl) GetByStatus(ctx context.Context, status entities.AnalysisStatus) ([]*entities.Analysis, error) {
	query := `
		SELECT id, examination_id, image_id, status, metrics, error_message, created_at, updated_at, completed_at
		FROM analyses WHERE status = $1 ORDER BY created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var analyses []*entities.Analysis
	for rows.Next() {
		analysis := &entities.Analysis{}
		var metricsJSON []byte

		err := rows.Scan(&analysis.ID, &analysis.ExaminationID, &analysis.ImageID, &analysis.Status,
			&metricsJSON, &analysis.ErrorMessage, &analysis.CreatedAt, &analysis.UpdatedAt, &analysis.CompletedAt)
		if err != nil {
			return nil, err
		}

		if len(metricsJSON) > 0 {
			if err := json.Unmarshal(metricsJSON, &analysis.Metrics); err != nil {
				return nil, err
			}
		} else {
			analysis.Metrics = make(map[string]interface{})
		}

		analyses = append(analyses, analysis)
	}
	return analyses, nil
}

func (r *AnalysisRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*entities.Analysis, error) {
	query := `
		SELECT id, examination_id, image_id, status, metrics, error_message, created_at, updated_at, completed_at
		FROM analyses ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var analyses []*entities.Analysis
	for rows.Next() {
		analysis := &entities.Analysis{}
		var metricsJSON []byte

		err := rows.Scan(&analysis.ID, &analysis.ExaminationID, &analysis.ImageID, &analysis.Status,
			&metricsJSON, &analysis.ErrorMessage, &analysis.CreatedAt, &analysis.UpdatedAt, &analysis.CompletedAt)
		if err != nil {
			return nil, err
		}

		if len(metricsJSON) > 0 {
			if err := json.Unmarshal(metricsJSON, &analysis.Metrics); err != nil {
				return nil, err
			}
		} else {
			analysis.Metrics = make(map[string]interface{})
		}

		analyses = append(analyses, analysis)
	}
	return analyses, nil
}
