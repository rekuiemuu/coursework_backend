package repositories

import (
	"context"

	"github.com/project-capillary/backend/internal/domain/entities"
)

type PatientRepository interface {
	Create(ctx context.Context, patient *entities.Patient) error
	GetByID(ctx context.Context, id string) (*entities.Patient, error)
	Update(ctx context.Context, patient *entities.Patient) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*entities.Patient, error)
	Search(ctx context.Context, query string) ([]*entities.Patient, error)
	Count(ctx context.Context) (int64, error)
}
