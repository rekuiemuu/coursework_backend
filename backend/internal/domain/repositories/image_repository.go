package repositories

import (
	"context"

	"github.com/project-capillary/backend/internal/domain/entities"
)

type ImageRepository interface {
	Create(ctx context.Context, image *entities.Image) error
	GetByID(ctx context.Context, id string) (*entities.Image, error)
	Delete(ctx context.Context, id string) error
	GetByExaminationID(ctx context.Context, examinationID string) ([]*entities.Image, error)
	List(ctx context.Context, limit, offset int) ([]*entities.Image, error)
	Count(ctx context.Context) (int64, error)
}
