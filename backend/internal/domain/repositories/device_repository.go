package repositories

import (
	"context"

	"github.com/project-capillary/backend/internal/domain/entities"
)

type DeviceRepository interface {
	Create(ctx context.Context, device *entities.Device) error
	GetByID(ctx context.Context, id string) (*entities.Device, error)
	GetByDeviceID(ctx context.Context, deviceID string) (*entities.Device, error)
	Update(ctx context.Context, device *entities.Device) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*entities.Device, error)
	GetOnlineDevices(ctx context.Context) ([]*entities.Device, error)
}
