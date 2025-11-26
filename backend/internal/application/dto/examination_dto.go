package dto

import "time"

type CreateExaminationRequest struct {
	PatientID   string `json:"patient_id" binding:"required"`
	DoctorID    string `json:"doctor_id" binding:"required"`
	Description string `json:"description"`
}

type UpdateExaminationRequest struct {
	Description string `json:"description"`
	Status      string `json:"status"`
}

type ExaminationResponse struct {
	ID          string     `json:"id"`
	PatientID   string     `json:"patient_id"`
	DoctorID    string     `json:"doctor_id"`
	Status      string     `json:"status"`
	Description string     `json:"description"`
	Images      []string   `json:"images"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

type StartAnalysisRequest struct {
	ExaminationID string `json:"examination_id" binding:"required"`
}

type AnalysisResponse struct {
	ID            string                 `json:"id"`
	ExaminationID string                 `json:"examination_id"`
	ImageID       string                 `json:"image_id"`
	Status        string                 `json:"status"`
	Metrics       map[string]interface{} `json:"metrics"`
	ErrorMessage  string                 `json:"error_message,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
	CompletedAt   *time.Time             `json:"completed_at"`
}
