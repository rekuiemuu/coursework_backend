package repositories

import (
	"context"

	"github.com/project-capillary/backend/internal/domain/entities"
)

type AnalysisRepository interface {
	Create(ctx context.Context, analysis *entities.Analysis) error
	GetByID(ctx context.Context, id string) (*entities.Analysis, error)
	Update(ctx context.Context, analysis *entities.Analysis) error
	GetByExaminationID(ctx context.Context, examinationID string) ([]*entities.Analysis, error)
	GetByImageID(ctx context.Context, imageID string) (*entities.Analysis, error)
	GetByStatus(ctx context.Context, status entities.AnalysisStatus) ([]*entities.Analysis, error)
	List(ctx context.Context, limit, offset int) ([]*entities.Analysis, error)
}
