package usecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/project-capillary/backend/internal/application/dto"
	"github.com/project-capillary/backend/internal/domain/entities"
	"github.com/project-capillary/backend/internal/domain/repositories"
	"github.com/project-capillary/backend/internal/infrastructure/mq"
)

type ExaminationUseCase struct {
	examinationRepo repositories.ExaminationRepository
	analysisRepo    repositories.AnalysisRepository
	imageRepo       repositories.ImageRepository
	mqPublisher     *mq.RabbitMQ
}

func NewExaminationUseCase(
	examinationRepo repositories.ExaminationRepository,
	analysisRepo repositories.AnalysisRepository,
	imageRepo repositories.ImageRepository,
	mqPublisher *mq.RabbitMQ,
) *ExaminationUseCase {
	return &ExaminationUseCase{
		examinationRepo: examinationRepo,
		analysisRepo:    analysisRepo,
		imageRepo:       imageRepo,
		mqPublisher:     mqPublisher,
	}
}

func (uc *ExaminationUseCase) CreateExamination(ctx context.Context, req dto.CreateExaminationRequest) (*dto.ExaminationResponse, error) {
	examination := entities.NewExamination(
		uuid.New().String(),
		req.PatientID,
		req.DoctorID,
		req.Description,
	)

	err := uc.examinationRepo.Create(ctx, examination)
	if err != nil {
		return nil, err
	}

	return &dto.ExaminationResponse{
		ID:          examination.ID,
		PatientID:   examination.PatientID,
		DoctorID:    examination.DoctorID,
		Status:      string(examination.Status),
		Description: examination.Description,
		Images:      examination.Images,
		CreatedAt:   examination.CreatedAt,
		UpdatedAt:   examination.UpdatedAt,
	}, nil
}

func (uc *ExaminationUseCase) GetExamination(ctx context.Context, id string) (*dto.ExaminationResponse, error) {
	examination, err := uc.examinationRepo.GetByID(ctx, id)
	if err != nil || examination == nil {
		return nil, err
	}

	return &dto.ExaminationResponse{
		ID:          examination.ID,
		PatientID:   examination.PatientID,
		DoctorID:    examination.DoctorID,
		Status:      string(examination.Status),
		Description: examination.Description,
		Images:      examination.Images,
		CreatedAt:   examination.CreatedAt,
		UpdatedAt:   examination.UpdatedAt,
		CompletedAt: examination.CompletedAt,
	}, nil
}

func (uc *ExaminationUseCase) StartAnalysis(ctx context.Context, examinationID string) error {
	examination, err := uc.examinationRepo.GetByID(ctx, examinationID)
	if err != nil || examination == nil {
		return err
	}

	images, err := uc.imageRepo.GetByExaminationID(ctx, examinationID)
	if err != nil {
		return err
	}

	examination.StartProcessing()
	err = uc.examinationRepo.Update(ctx, examination)
	if err != nil {
		return err
	}

	for _, image := range images {
		analysis := entities.NewAnalysis(uuid.New().String(), examinationID, image.ID)
		err = uc.analysisRepo.Create(ctx, analysis)
		if err != nil {
			continue
		}

		taskMsg := mq.NewAnalysisTaskMessage(analysis.ID, examinationID, image.ID, image.FilePath)
		uc.mqPublisher.Publish(ctx, taskMsg)
	}

	return nil
}

func (uc *ExaminationUseCase) GetExaminationsByPatient(ctx context.Context, patientID string) ([]*dto.ExaminationResponse, error) {
	examinations, err := uc.examinationRepo.GetByPatientID(ctx, patientID)
	if err != nil {
		return nil, err
	}

	var response []*dto.ExaminationResponse
	for _, e := range examinations {
		response = append(response, &dto.ExaminationResponse{
			ID:          e.ID,
			PatientID:   e.PatientID,
			DoctorID:    e.DoctorID,
			Status:      string(e.Status),
			Description: e.Description,
			Images:      e.Images,
			CreatedAt:   e.CreatedAt,
			UpdatedAt:   e.UpdatedAt,
			CompletedAt: e.CompletedAt,
		})
	}
	return response, nil
}
