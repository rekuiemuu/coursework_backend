package dto

import "time"

type ImageResponse struct {
	ID            string    `json:"id"`
	ExaminationID string    `json:"examination_id"`
	Filename      string    `json:"filename"`
	FilePath      string    `json:"file_path"`
	FileSize      int64     `json:"file_size"`
	MimeType      string    `json:"mime_type"`
	Width         int       `json:"width"`
	Height        int       `json:"height"`
	CapturedAt    time.Time `json:"captured_at"`
	CreatedAt     time.Time `json:"created_at"`
}
