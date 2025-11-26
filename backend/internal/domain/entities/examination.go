package entities

import (
	"time"
)

type ExaminationStatus string

const (
	StatusPending    ExaminationStatus = "pending"
	StatusInProgress ExaminationStatus = "in_progress"
	StatusCompleted  ExaminationStatus = "completed"
	StatusFailed     ExaminationStatus = "failed"
)

type Examination struct {
	ID          string
	PatientID   string
	DoctorID    string
	Status      ExaminationStatus
	Description string
	Images      []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CompletedAt *time.Time
}

func NewExamination(id, patientID, doctorID, description string) *Examination {
	now := time.Now()
	return &Examination{
		ID:          id,
		PatientID:   patientID,
		DoctorID:    doctorID,
		Status:      StatusPending,
		Description: description,
		Images:      []string{},
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (e *Examination) AddImage(imageID string) {
	e.Images = append(e.Images, imageID)
	e.UpdatedAt = time.Now()
}

func (e *Examination) StartProcessing() {
	e.Status = StatusInProgress
	e.UpdatedAt = time.Now()
}

func (e *Examination) Complete() {
	e.Status = StatusCompleted
	now := time.Now()
	e.CompletedAt = &now
	e.UpdatedAt = now
}

func (e *Examination) Fail() {
	e.Status = StatusFailed
	e.UpdatedAt = time.Now()
}
