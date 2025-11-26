package entities

import (
	"time"
)

type DeviceStatus string

const (
	DeviceStatusOnline  DeviceStatus = "online"
	DeviceStatusOffline DeviceStatus = "offline"
	DeviceStatusError   DeviceStatus = "error"
)

type Device struct {
	ID         string
	Name       string
	DeviceID   string
	Label      string
	Status     DeviceStatus
	LastSeen   time.Time
	Brightness float64
	Zoom       float64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewDevice(id, name, deviceID, label string) *Device {
	now := time.Now()
	return &Device{
		ID:         id,
		Name:       name,
		DeviceID:   deviceID,
		Label:      label,
		Status:     DeviceStatusOffline,
		LastSeen:   now,
		Brightness: 0.5,
		Zoom:       1.0,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func (d *Device) UpdateStatus(status DeviceStatus) {
	d.Status = status
	d.LastSeen = time.Now()
	d.UpdatedAt = time.Now()
}

func (d *Device) SetBrightness(value float64) {
	d.Brightness = value
	d.UpdatedAt = time.Now()
}

func (d *Device) SetZoom(value float64) {
	d.Zoom = value
	d.UpdatedAt = time.Now()
}
