package usecases

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/project-capillary/backend/internal/application/dto"
	"github.com/project-capillary/backend/internal/domain/entities"
	"github.com/project-capillary/backend/internal/domain/repositories"
)

type ReportUseCase struct {
	reportRepo      repositories.ReportRepository
	examinationRepo repositories.ExaminationRepository
	analysisRepo    repositories.AnalysisRepository
	imageRepo       repositories.ImageRepository
}

func NewReportUseCase(
	reportRepo repositories.ReportRepository,
	examinationRepo repositories.ExaminationRepository,
	analysisRepo repositories.AnalysisRepository,
	imageRepo repositories.ImageRepository,
) *ReportUseCase {
	return &ReportUseCase{
		reportRepo:      reportRepo,
		examinationRepo: examinationRepo,
		analysisRepo:    analysisRepo,
		imageRepo:       imageRepo,
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

	images, err := uc.imageRepo.GetByExaminationID(ctx, report.ExaminationID)
	if err != nil {
		images = []*entities.Image{}
	}

	imageInfos := make([]dto.ImageInfo, 0, len(images))
	for _, img := range images {
		imageInfos = append(imageInfos, dto.ImageInfo{
			ID:       img.ID,
			Filename: img.Filename,
			URL:      "/api/photos/" + img.Filename,
		})
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
		Images:          imageInfos,
		CreatedAt:       report.CreatedAt,
		UpdatedAt:       report.UpdatedAt,
	}, nil
}

func (uc *ReportUseCase) GetReportByExamination(ctx context.Context, examinationID string) (*dto.ReportResponse, error) {
	report, err := uc.reportRepo.GetByExaminationID(ctx, examinationID)
	if err != nil || report == nil {
		return nil, err
	}

	images, err := uc.imageRepo.GetByExaminationID(ctx, examinationID)
	if err != nil {
		images = []*entities.Image{}
	}

	imageInfos := make([]dto.ImageInfo, 0, len(images))
	for _, img := range images {
		imageInfos = append(imageInfos, dto.ImageInfo{
			ID:       img.ID,
			Filename: img.Filename,
			URL:      "/api/photos/" + img.Filename,
		})
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
		Images:          imageInfos,
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

func (uc *ReportUseCase) UpdateReport(ctx context.Context, id string, req dto.UpdateReportRequest) (*dto.ReportResponse, error) {
	report, err := uc.reportRepo.GetByID(ctx, id)
	if err != nil || report == nil {
		return nil, err
	}

	if req.Content != "" {
		report.Content = req.Content
	}
	if req.Summary != "" {
		report.Summary = req.Summary
	}
	if req.Diagnosis != "" {
		report.Diagnosis = req.Diagnosis
	}
	if req.Recommendations != "" {
		report.Recommendations = req.Recommendations
	}

	report.UpdatedAt = time.Now()

	err = uc.reportRepo.Update(ctx, report)
	if err != nil {
		return nil, err
	}

	return uc.GetReport(ctx, id)
}
