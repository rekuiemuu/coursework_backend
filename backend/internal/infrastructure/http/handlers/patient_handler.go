package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/project-capillary/backend/internal/application/dto"
	"github.com/project-capillary/backend/internal/application/usecases"
)

type PatientHandler struct {
	patientUseCase *usecases.PatientUseCase
}

func NewPatientHandler(patientUseCase *usecases.PatientUseCase) *PatientHandler {
	return &PatientHandler{patientUseCase: patientUseCase}
}

func (h *PatientHandler) CreatePatient(c *gin.Context) {
	var req dto.CreatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	patient, err := h.patientUseCase.CreatePatient(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse{
		Success: true,
		Data:    patient,
	})
}

func (h *PatientHandler) GetPatient(c *gin.Context) {
	id := c.Param("id")
	patient, err := h.patientUseCase.GetPatient(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "not_found",
			Message: "Patient not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Data:    patient,
	})
}

func (h *PatientHandler) UpdatePatient(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	err := h.patientUseCase.UpdatePatient(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Patient updated successfully",
	})
}

func (h *PatientHandler) DeletePatient(c *gin.Context) {
	id := c.Param("id")
	err := h.patientUseCase.DeletePatient(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Patient deleted successfully",
	})
}

func (h *PatientHandler) ListPatients(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	patients, err := h.patientUseCase.ListPatients(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Data:    patients,
	})
}
