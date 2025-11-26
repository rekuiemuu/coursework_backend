package postgres

import (
	"context"
	"database/sql"

	"github.com/project-capillary/backend/internal/domain/entities"
)

type PatientRepositoryImpl struct {
	db *sql.DB
}

func NewPatientRepository(db *sql.DB) *PatientRepositoryImpl {
	return &PatientRepositoryImpl{db: db}
}

func (r *PatientRepositoryImpl) Create(ctx context.Context, patient *entities.Patient) error {
	query := `
		INSERT INTO patients (id, first_name, last_name, middle_name, date_of_birth, gender, phone, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.ExecContext(ctx, query,
		patient.ID, patient.FirstName, patient.LastName, patient.MiddleName,
		patient.DateOfBirth, patient.Gender, patient.Phone, patient.Email,
		patient.CreatedAt, patient.UpdatedAt)
	return err
}

func (r *PatientRepositoryImpl) GetByID(ctx context.Context, id string) (*entities.Patient, error) {
	query := `
		SELECT id, first_name, last_name, middle_name, date_of_birth, gender, phone, email, created_at, updated_at
		FROM patients WHERE id = $1
	`
	patient := &entities.Patient{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&patient.ID, &patient.FirstName, &patient.LastName, &patient.MiddleName,
		&patient.DateOfBirth, &patient.Gender, &patient.Phone, &patient.Email,
		&patient.CreatedAt, &patient.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return patient, nil
}

func (r *PatientRepositoryImpl) Update(ctx context.Context, patient *entities.Patient) error {
	query := `
		UPDATE patients 
		SET first_name = $2, last_name = $3, middle_name = $4, phone = $5, email = $6, updated_at = $7
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		patient.ID, patient.FirstName, patient.LastName, patient.MiddleName,
		patient.Phone, patient.Email, patient.UpdatedAt)
	return err
}

func (r *PatientRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM patients WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PatientRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*entities.Patient, error) {
	query := `
		SELECT id, first_name, last_name, middle_name, date_of_birth, gender, phone, email, created_at, updated_at
		FROM patients ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []*entities.Patient
	for rows.Next() {
		patient := &entities.Patient{}
		err := rows.Scan(&patient.ID, &patient.FirstName, &patient.LastName, &patient.MiddleName,
			&patient.DateOfBirth, &patient.Gender, &patient.Phone, &patient.Email,
			&patient.CreatedAt, &patient.UpdatedAt)
		if err != nil {
			return nil, err
		}
		patients = append(patients, patient)
	}
	return patients, nil
}

func (r *PatientRepositoryImpl) Search(ctx context.Context, query string) ([]*entities.Patient, error) {
	searchQuery := `
		SELECT id, first_name, last_name, middle_name, date_of_birth, gender, phone, email, created_at, updated_at
		FROM patients 
		WHERE first_name ILIKE $1 OR last_name ILIKE $1 OR middle_name ILIKE $1 OR email ILIKE $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, searchQuery, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []*entities.Patient
	for rows.Next() {
		patient := &entities.Patient{}
		err := rows.Scan(&patient.ID, &patient.FirstName, &patient.LastName, &patient.MiddleName,
			&patient.DateOfBirth, &patient.Gender, &patient.Phone, &patient.Email,
			&patient.CreatedAt, &patient.UpdatedAt)
		if err != nil {
			return nil, err
		}
		patients = append(patients, patient)
	}
	return patients, nil
}

func (r *PatientRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM patients`
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}
