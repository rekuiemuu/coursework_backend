package usecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/project-capillary/backend/internal/application/dto"
	"github.com/project-capillary/backend/internal/domain/entities"
	"github.com/project-capillary/backend/internal/domain/repositories"
)

type ReportUseCase struct {
	reportRepo      repositories.ReportRepository
	examinationRepo repositories.ExaminationRepository
	analysisRepo    repositories.AnalysisRepository
}

func NewReportUseCase(
	reportRepo repositories.ReportRepository,
	examinationRepo repositories.ExaminationRepository,
	analysisRepo repositories.AnalysisRepository,
) *ReportUseCase {
	return &ReportUseCase{
		reportRepo:      reportRepo,
		examinationRepo: examinationRepo,
		analysisRepo:    analysisRepo,
	}
}

func (uc *ReportUseCase) CreateReport(ctx context.Context, req dto.CreateReportRequest) (*dto.ReportResponse, error) {
	report := entities.NewReport(
		uuid.New().String(),
		req.ExaminationID,
		req.Title,
		req.Content,
		req.Summary,
		req.Diagnosis,
		req.Recommendations,
		req.GeneratedBy,
	)

	err := uc.reportRepo.Create(ctx, report)
	if err != nil {
		return nil, err
	}

	return &dto.ReportResponse{
		ID:              report.ID,
		ExaminationID:   report.ExaminationID,
		Title:           report.Title,
		Content:         report.Content,
		Summary:         report.Summary,
		Diagnosis:       report.Diagnosis,
		Recommendations: report.Recommendations,
		GeneratedBy:     report.GeneratedBy,
		CreatedAt:       report.CreatedAt,
		UpdatedAt:       report.UpdatedAt,
	}, nil
}

func (uc *ReportUseCase) GetReport(ctx context.Context, id string) (*dto.ReportResponse, error) {
	report, err := uc.reportRepo.GetByID(ctx, id)
	if err != nil || report == nil {
		return nil, err
	}

	return &dto.ReportResponse{
		ID:              report.ID,
		ExaminationID:   report.ExaminationID,
		Title:           report.Title,
		Content:         report.Content,
		Summary:         report.Summary,
		Diagnosis:       report.Diagnosis,
		Recommendations: report.Recommendations,
		GeneratedBy:     report.GeneratedBy,
		CreatedAt:       report.CreatedAt,
		UpdatedAt:       report.UpdatedAt,
	}, nil
}

func (uc *ReportUseCase) GetReportByExamination(ctx context.Context, examinationID string) (*dto.ReportResponse, error) {
	report, err := uc.reportRepo.GetByExaminationID(ctx, examinationID)
	if err != nil || report == nil {
		return nil, err
	}

	return &dto.ReportResponse{
		ID:              report.ID,
		ExaminationID:   report.ExaminationID,
		Title:           report.Title,
		Content:         report.Content,
		Summary:         report.Summary,
		Diagnosis:       report.Diagnosis,
		Recommendations: report.Recommendations,
		GeneratedBy:     report.GeneratedBy,
		CreatedAt:       report.CreatedAt,
		UpdatedAt:       report.UpdatedAt,
	}, nil
}

func (uc *ReportUseCase) ListReports(ctx context.Context, limit, offset int) ([]*dto.ReportResponse, error) {
	reports, err := uc.reportRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	var response []*dto.ReportResponse
	for _, report := range reports {
		response = append(response, &dto.ReportResponse{
			ID:              report.ID,
			ExaminationID:   report.ExaminationID,
			Title:           report.Title,
			Content:         report.Content,
			Summary:         report.Summary,
			Diagnosis:       report.Diagnosis,
			Recommendations: report.Recommendations,
			GeneratedBy:     report.GeneratedBy,
			CreatedAt:       report.CreatedAt,
			UpdatedAt:       report.UpdatedAt,
		})
	}

	return response, nil
}
