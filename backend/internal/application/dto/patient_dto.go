package dto

import "time"

type CreatePatientRequest struct {
	FirstName   string    `json:"first_name" binding:"required"`
	LastName    string    `json:"last_name" binding:"required"`
	MiddleName  string    `json:"middle_name"`
	DateOfBirth time.Time `json:"date_of_birth" binding:"required"`
	Gender      string    `json:"gender" binding:"required"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
}

type UpdatePatientRequest struct {
	FirstName  string `json:"first_name" binding:"required"`
	LastName   string `json:"last_name" binding:"required"`
	MiddleName string `json:"middle_name"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
}

type PatientResponse struct {
	ID          string    `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	MiddleName  string    `json:"middle_name"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Gender      string    `json:"gender"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
