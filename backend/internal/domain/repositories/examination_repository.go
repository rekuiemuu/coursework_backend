package repositories

import (
	"context"

	"github.com/project-capillary/backend/internal/domain/entities"
)

type ExaminationRepository interface {
	Create(ctx context.Context, examination *entities.Examination) error
	GetByID(ctx context.Context, id string) (*entities.Examination, error)
	Update(ctx context.Context, examination *entities.Examination) error
	Delete(ctx context.Context, id string) error
	GetByPatientID(ctx context.Context, patientID string) ([]*entities.Examination, error)
	GetByStatus(ctx context.Context, status entities.ExaminationStatus) ([]*entities.Examination, error)
	List(ctx context.Context, limit, offset int) ([]*entities.Examination, error)
	Count(ctx context.Context) (int64, error)
}
