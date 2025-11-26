package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/project-capillary/backend/internal/application/dto"
	"github.com/project-capillary/backend/internal/application/usecases"
)

type ExaminationHandler struct {
	examinationUseCase *usecases.ExaminationUseCase
}

func NewExaminationHandler(examinationUseCase *usecases.ExaminationUseCase) *ExaminationHandler {
	return &ExaminationHandler{examinationUseCase: examinationUseCase}
}

func (h *ExaminationHandler) CreateExamination(c *gin.Context) {
	var req dto.CreateExaminationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	examination, err := h.examinationUseCase.CreateExamination(c.Request.Context(), req)
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
		Data:    examination,
	})
}

func (h *ExaminationHandler) ListExaminations(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	examinations, err := h.examinationUseCase.ListExaminations(c.Request.Context(), limit, offset)
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
		Data:    examinations,
	})
}

func (h *ExaminationHandler) GetExamination(c *gin.Context) {
	id := c.Param("id")
	examination, err := h.examinationUseCase.GetExamination(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "not_found",
			Message: "Examination not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Data:    examination,
	})
}

func (h *ExaminationHandler) StartAnalysis(c *gin.Context) {
	id := c.Param("id")
	err := h.examinationUseCase.StartAnalysis(c.Request.Context(), id)
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
		Message: "Analysis started successfully",
	})
}

func (h *ExaminationHandler) GetPatientExaminations(c *gin.Context) {
	patientID := c.Param("patientId")
	examinations, err := h.examinationUseCase.GetExaminationsByPatient(c.Request.Context(), patientID)
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
		Data:    examinations,
	})
}

func (h *ExaminationHandler) AttachPhotos(c *gin.Context) {
	id := c.Param("id")
	var req dto.AttachPhotosRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	err := h.examinationUseCase.AttachPhotos(c.Request.Context(), id, req.Photos)
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
		Message: "Photos attached successfully",
	})
}
