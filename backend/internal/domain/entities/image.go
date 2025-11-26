package entities

import (
	"time"
)

type Image struct {
	ID            string
	ExaminationID string
	Filename      string
	FilePath      string
	FileSize      int64
	MimeType      string
	Width         int
	Height        int
	CapturedAt    time.Time
	CreatedAt     time.Time
}

func NewImage(id, examinationID, filename, filepath, mimeType string, fileSize int64, width, height int) *Image {
	now := time.Now()
	return &Image{
		ID:            id,
		ExaminationID: examinationID,
		Filename:      filename,
		FilePath:      filepath,
		FileSize:      fileSize,
		MimeType:      mimeType,
		Width:         width,
		Height:        height,
		CapturedAt:    now,
		CreatedAt:     now,
	}
}
