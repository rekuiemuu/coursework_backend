package postgres

import (
	"context"
	"database/sql"

	"github.com/project-capillary/backend/internal/domain/entities"
)

type ExaminationRepositoryImpl struct {
	db *sql.DB
}

func NewExaminationRepository(db *sql.DB) *ExaminationRepositoryImpl {
	return &ExaminationRepositoryImpl{db: db}
}

func (r *ExaminationRepositoryImpl) Create(ctx context.Context, examination *entities.Examination) error {
	query := `
		INSERT INTO examinations (id, patient_id, doctor_id, status, description, created_at, updated_at, completed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		examination.ID, examination.PatientID, examination.DoctorID, examination.Status,
		examination.Description, examination.CreatedAt, examination.UpdatedAt, examination.CompletedAt)
	return err
}

func (r *ExaminationRepositoryImpl) GetByID(ctx context.Context, id string) (*entities.Examination, error) {
	query := `
		SELECT id, patient_id, doctor_id, status, description, created_at, updated_at, completed_at
		FROM examinations WHERE id = $1
	`
	examination := &entities.Examination{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&examination.ID, &examination.PatientID, &examination.DoctorID, &examination.Status,
		&examination.Description, &examination.CreatedAt, &examination.UpdatedAt, &examination.CompletedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	imageQuery := `SELECT image_id FROM examination_images WHERE examination_id = $1 ORDER BY image_order`
	rows, err := r.db.QueryContext(ctx, imageQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	examination.Images = []string{}
	for rows.Next() {
		var imageID string
		if err := rows.Scan(&imageID); err != nil {
			return nil, err
		}
		examination.Images = append(examination.Images, imageID)
	}

	return examination, nil
}

func (r *ExaminationRepositoryImpl) Update(ctx context.Context, examination *entities.Examination) error {
	query := `
		UPDATE examinations 
		SET status = $2, description = $3, updated_at = $4, completed_at = $5
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		examination.ID, examination.Status, examination.Description,
		examination.UpdatedAt, examination.CompletedAt)

	if err != nil {
		return err
	}

	deleteQuery := `DELETE FROM examination_images WHERE examination_id = $1`
	if _, err := r.db.ExecContext(ctx, deleteQuery, examination.ID); err != nil {
		return err
	}

	if len(examination.Images) > 0 {
		insertQuery := `INSERT INTO examination_images (examination_id, image_id, image_order) VALUES ($1, $2, $3)`
		for i, imageID := range examination.Images {
			if _, err := r.db.ExecContext(ctx, insertQuery, examination.ID, imageID, i); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *ExaminationRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM examinations WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *ExaminationRepositoryImpl) GetByPatientID(ctx context.Context, patientID string) ([]*entities.Examination, error) {
	query := `
		SELECT id, patient_id, doctor_id, status, description, created_at, updated_at, completed_at
		FROM examinations WHERE patient_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var examinations []*entities.Examination
	for rows.Next() {
		exam := &entities.Examination{}
		err := rows.Scan(&exam.ID, &exam.PatientID, &exam.DoctorID, &exam.Status,
			&exam.Description, &exam.CreatedAt, &exam.UpdatedAt, &exam.CompletedAt)
		if err != nil {
			return nil, err
		}
		exam.Images = []string{}
		examinations = append(examinations, exam)
	}
	return examinations, nil
}

func (r *ExaminationRepositoryImpl) GetByStatus(ctx context.Context, status entities.ExaminationStatus) ([]*entities.Examination, error) {
	query := `
		SELECT id, patient_id, doctor_id, status, description, created_at, updated_at, completed_at
		FROM examinations WHERE status = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var examinations []*entities.Examination
	for rows.Next() {
		exam := &entities.Examination{}
		err := rows.Scan(&exam.ID, &exam.PatientID, &exam.DoctorID, &exam.Status,
			&exam.Description, &exam.CreatedAt, &exam.UpdatedAt, &exam.CompletedAt)
		if err != nil {
			return nil, err
		}
		exam.Images = []string{}
		examinations = append(examinations, exam)
	}
	return examinations, nil
}

func (r *ExaminationRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*entities.Examination, error) {
	query := `
		SELECT id, patient_id, doctor_id, status, description, created_at, updated_at, completed_at
		FROM examinations ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var examinations []*entities.Examination
	for rows.Next() {
		exam := &entities.Examination{}
		err := rows.Scan(&exam.ID, &exam.PatientID, &exam.DoctorID, &exam.Status,
			&exam.Description, &exam.CreatedAt, &exam.UpdatedAt, &exam.CompletedAt)
		if err != nil {
			return nil, err
		}
		exam.Images = []string{}
		examinations = append(examinations, exam)
	}
	return examinations, nil
}

func (r *ExaminationRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM examinations`
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}
