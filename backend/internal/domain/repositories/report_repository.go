package repositories

import (
	"context"

	"github.com/project-capillary/backend/internal/domain/entities"
)

type ReportRepository interface {
	Create(ctx context.Context, report *entities.Report) error
	GetByID(ctx context.Context, id string) (*entities.Report, error)
	Update(ctx context.Context, report *entities.Report) error
	Delete(ctx context.Context, id string) error
	GetByExaminationID(ctx context.Context, examinationID string) (*entities.Report, error)
	List(ctx context.Context, limit, offset int) ([]*entities.Report, error)
}
