package entities

import (
	"time"
)

type AnalysisStatus string

const (
	AnalysisStatusPending    AnalysisStatus = "pending"
	AnalysisStatusProcessing AnalysisStatus = "processing"
	AnalysisStatusCompleted  AnalysisStatus = "completed"
	AnalysisStatusFailed     AnalysisStatus = "failed"
)

type Analysis struct {
	ID            string
	ExaminationID string
	ImageID       string
	Status        AnalysisStatus
	Metrics       map[string]interface{}
	ErrorMessage  string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	CompletedAt   *time.Time
}

func NewAnalysis(id, examinationID, imageID string) *Analysis {
	now := time.Now()
	return &Analysis{
		ID:            id,
		ExaminationID: examinationID,
		ImageID:       imageID,
		Status:        AnalysisStatusPending,
		Metrics:       make(map[string]interface{}),
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func (a *Analysis) StartProcessing() {
	a.Status = AnalysisStatusProcessing
	a.UpdatedAt = time.Now()
}

func (a *Analysis) Complete(metrics map[string]interface{}) {
	a.Status = AnalysisStatusCompleted
	a.Metrics = metrics
	now := time.Now()
	a.CompletedAt = &now
	a.UpdatedAt = now
}

func (a *Analysis) Fail(errorMsg string) {
	a.Status = AnalysisStatusFailed
	a.ErrorMessage = errorMsg
	a.UpdatedAt = time.Now()
}
