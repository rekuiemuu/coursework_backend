package entities

import (
	"time"
)

type Report struct {
	ID              string
	ExaminationID   string
	Title           string
	Content         string
	Summary         string
	Diagnosis       string
	Recommendations string
	GeneratedBy     string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func NewReport(id, examinationID, title, content, summary, diagnosis, recommendations, generatedBy string) *Report {
	now := time.Now()
	return &Report{
		ID:              id,
		ExaminationID:   examinationID,
		Title:           title,
		Content:         content,
		Summary:         summary,
		Diagnosis:       diagnosis,
		Recommendations: recommendations,
		GeneratedBy:     generatedBy,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

func (r *Report) Update(content, summary, diagnosis, recommendations string) {
	r.Content = content
	r.Summary = summary
	r.Diagnosis = diagnosis
	r.Recommendations = recommendations
	r.UpdatedAt = time.Now()
}
