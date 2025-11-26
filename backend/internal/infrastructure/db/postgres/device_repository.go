package postgres

import (
	"context"
	"database/sql"

	"github.com/project-capillary/backend/internal/domain/entities"
)

type DeviceRepositoryImpl struct {
	db *sql.DB
}

func NewDeviceRepository(db *sql.DB) *DeviceRepositoryImpl {
	return &DeviceRepositoryImpl{db: db}
}

func (r *DeviceRepositoryImpl) Create(ctx context.Context, device *entities.Device) error {
	query := `
		INSERT INTO devices (id, name, device_id, label, status, last_seen, brightness, zoom, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.ExecContext(ctx, query,
		device.ID, device.Name, device.DeviceID, device.Label, device.Status,
		device.LastSeen, device.Brightness, device.Zoom, device.CreatedAt, device.UpdatedAt)
	return err
}

func (r *DeviceRepositoryImpl) GetByID(ctx context.Context, id string) (*entities.Device, error) {
	query := `
		SELECT id, name, device_id, label, status, last_seen, brightness, zoom, created_at, updated_at
		FROM devices WHERE id = $1
	`
	device := &entities.Device{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&device.ID, &device.Name, &device.DeviceID, &device.Label, &device.Status,
		&device.LastSeen, &device.Brightness, &device.Zoom, &device.CreatedAt, &device.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return device, err
}

func (r *DeviceRepositoryImpl) GetByDeviceID(ctx context.Context, deviceID string) (*entities.Device, error) {
	query := `
		SELECT id, name, device_id, label, status, last_seen, brightness, zoom, created_at, updated_at
		FROM devices WHERE device_id = $1
	`
	device := &entities.Device{}
	err := r.db.QueryRowContext(ctx, query, deviceID).Scan(
		&device.ID, &device.Name, &device.DeviceID, &device.Label, &device.Status,
		&device.LastSeen, &device.Brightness, &device.Zoom, &device.CreatedAt, &device.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return device, err
}

func (r *DeviceRepositoryImpl) Update(ctx context.Context, device *entities.Device) error {
	query := `
		UPDATE devices 
		SET name = $2, label = $3, status = $4, last_seen = $5, brightness = $6, zoom = $7, updated_at = $8
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		device.ID, device.Name, device.Label, device.Status, device.LastSeen,
		device.Brightness, device.Zoom, device.UpdatedAt)
	return err
}

func (r *DeviceRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM devices WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *DeviceRepositoryImpl) List(ctx context.Context) ([]*entities.Device, error) {
	query := `
		SELECT id, name, device_id, label, status, last_seen, brightness, zoom, created_at, updated_at
		FROM devices ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []*entities.Device
	for rows.Next() {
		device := &entities.Device{}
		err := rows.Scan(&device.ID, &device.Name, &device.DeviceID, &device.Label, &device.Status,
			&device.LastSeen, &device.Brightness, &device.Zoom, &device.CreatedAt, &device.UpdatedAt)
		if err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}
	return devices, nil
}

func (r *DeviceRepositoryImpl) GetOnlineDevices(ctx context.Context) ([]*entities.Device, error) {
	query := `
		SELECT id, name, device_id, label, status, last_seen, brightness, zoom, created_at, updated_at
		FROM devices WHERE status = 'online' ORDER BY last_seen DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []*entities.Device
	for rows.Next() {
		device := &entities.Device{}
		err := rows.Scan(&device.ID, &device.Name, &device.DeviceID, &device.Label, &device.Status,
			&device.LastSeen, &device.Brightness, &device.Zoom, &device.CreatedAt, &device.UpdatedAt)
		if err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}
	return devices, nil
}
