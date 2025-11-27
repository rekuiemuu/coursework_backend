package dto

import "time"

type CreateReportRequest struct {
	ExaminationID   string `json:"examination_id" binding:"required"`
	Title           string `json:"title" binding:"required"`
	Content         string `json:"content" binding:"required"`
	Summary         string `json:"summary"`
	Diagnosis       string `json:"diagnosis"`
	Recommendations string `json:"recommendations"`
	GeneratedBy     string `json:"generated_by" binding:"required"`
}

type UpdateReportRequest struct {
	Content         string `json:"content"`
	Summary         string `json:"summary"`
	Diagnosis       string `json:"diagnosis"`
	Recommendations string `json:"recommendations"`
}

type ReportResponse struct {
	ID              string      `json:"id"`
	ExaminationID   string      `json:"examination_id"`
	Title           string      `json:"title"`
	Content         string      `json:"content"`
	Summary         string      `json:"summary"`
	Diagnosis       string      `json:"diagnosis"`
	Recommendations string      `json:"recommendations"`
	GeneratedBy     string      `json:"generated_by"`
	Images          []ImageInfo `json:"images,omitempty"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

type ImageInfo struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	URL      string `json:"url"`
}
