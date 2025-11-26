package postgres

import (
	"context"
	"database/sql"

	"github.com/project-capillary/backend/internal/domain/entities"
)

type ImageRepositoryImpl struct {
	db *sql.DB
}

func NewImageRepository(db *sql.DB) *ImageRepositoryImpl {
	return &ImageRepositoryImpl{db: db}
}

func (r *ImageRepositoryImpl) Create(ctx context.Context, image *entities.Image) error {
	query := `
		INSERT INTO images (id, examination_id, filename, file_path, file_size, mime_type, width, height, captured_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.ExecContext(ctx, query,
		image.ID, image.ExaminationID, image.Filename, image.FilePath, image.FileSize,
		image.MimeType, image.Width, image.Height, image.CapturedAt, image.CreatedAt)
	return err
}

func (r *ImageRepositoryImpl) GetByID(ctx context.Context, id string) (*entities.Image, error) {
	query := `
		SELECT id, examination_id, filename, file_path, file_size, mime_type, width, height, captured_at, created_at
		FROM images WHERE id = $1
	`
	image := &entities.Image{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&image.ID, &image.ExaminationID, &image.Filename, &image.FilePath, &image.FileSize,
		&image.MimeType, &image.Width, &image.Height, &image.CapturedAt, &image.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return image, err
}

func (r *ImageRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM images WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *ImageRepositoryImpl) GetByExaminationID(ctx context.Context, examinationID string) ([]*entities.Image, error) {
	query := `
		SELECT id, examination_id, filename, file_path, file_size, mime_type, width, height, captured_at, created_at
		FROM images WHERE examination_id = $1 ORDER BY captured_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, examinationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []*entities.Image
	for rows.Next() {
		image := &entities.Image{}
		err := rows.Scan(&image.ID, &image.ExaminationID, &image.Filename, &image.FilePath,
			&image.FileSize, &image.MimeType, &image.Width, &image.Height, &image.CapturedAt, &image.CreatedAt)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}
	return images, nil
}

func (r *ImageRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*entities.Image, error) {
	query := `
		SELECT id, examination_id, filename, file_path, file_size, mime_type, width, height, captured_at, created_at
		FROM images ORDER BY captured_at DESC LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []*entities.Image
	for rows.Next() {
		image := &entities.Image{}
		err := rows.Scan(&image.ID, &image.ExaminationID, &image.Filename, &image.FilePath,
			&image.FileSize, &image.MimeType, &image.Width, &image.Height, &image.CapturedAt, &image.CreatedAt)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}
	return images, nil
}

func (r *ImageRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM images`
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}
