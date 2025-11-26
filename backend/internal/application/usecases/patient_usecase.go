package usecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/project-capillary/backend/internal/application/dto"
	"github.com/project-capillary/backend/internal/domain/entities"
	"github.com/project-capillary/backend/internal/domain/repositories"
)

type PatientUseCase struct {
	patientRepo repositories.PatientRepository
}

func NewPatientUseCase(patientRepo repositories.PatientRepository) *PatientUseCase {
	return &PatientUseCase{patientRepo: patientRepo}
}

func (uc *PatientUseCase) CreatePatient(ctx context.Context, req dto.CreatePatientRequest) (*dto.PatientResponse, error) {
	patient := entities.NewPatient(
		uuid.New().String(),
		req.FirstName,
		req.LastName,
		req.MiddleName,
		req.DateOfBirth,
		req.Gender,
		req.Phone,
		req.Email,
	)

	err := uc.patientRepo.Create(ctx, patient)
	if err != nil {
		return nil, err
	}

	return &dto.PatientResponse{
		ID:          patient.ID,
		FirstName:   patient.FirstName,
		LastName:    patient.LastName,
		MiddleName:  patient.MiddleName,
		DateOfBirth: patient.DateOfBirth,
		Gender:      patient.Gender,
		Phone:       patient.Phone,
		Email:       patient.Email,
		CreatedAt:   patient.CreatedAt,
		UpdatedAt:   patient.UpdatedAt,
	}, nil
}

func (uc *PatientUseCase) GetPatient(ctx context.Context, id string) (*dto.PatientResponse, error) {
	patient, err := uc.patientRepo.GetByID(ctx, id)
	if err != nil || patient == nil {
		return nil, err
	}

	return &dto.PatientResponse{
		ID:          patient.ID,
		FirstName:   patient.FirstName,
		LastName:    patient.LastName,
		MiddleName:  patient.MiddleName,
		DateOfBirth: patient.DateOfBirth,
		Gender:      patient.Gender,
		Phone:       patient.Phone,
		Email:       patient.Email,
		CreatedAt:   patient.CreatedAt,
		UpdatedAt:   patient.UpdatedAt,
	}, nil
}

func (uc *PatientUseCase) UpdatePatient(ctx context.Context, id string, req dto.UpdatePatientRequest) error {
	patient, err := uc.patientRepo.GetByID(ctx, id)
	if err != nil || patient == nil {
		return err
	}

	patient.Update(req.FirstName, req.LastName, req.MiddleName, req.Phone, req.Email)
	return uc.patientRepo.Update(ctx, patient)
}

func (uc *PatientUseCase) DeletePatient(ctx context.Context, id string) error {
	return uc.patientRepo.Delete(ctx, id)
}

func (uc *PatientUseCase) ListPatients(ctx context.Context, limit, offset int) ([]*dto.PatientResponse, error) {
	patients, err := uc.patientRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	var response []*dto.PatientResponse
	for _, p := range patients {
		response = append(response, &dto.PatientResponse{
			ID:          p.ID,
			FirstName:   p.FirstName,
			LastName:    p.LastName,
			MiddleName:  p.MiddleName,
			DateOfBirth: p.DateOfBirth,
			Gender:      p.Gender,
			Phone:       p.Phone,
			Email:       p.Email,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		})
	}
	return response, nil
}
