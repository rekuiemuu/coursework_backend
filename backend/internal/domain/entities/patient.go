package entities

import (
	"time"
)

type Patient struct {
	ID          string
	FirstName   string
	LastName    string
	MiddleName  string
	DateOfBirth time.Time
	Gender      string
	Phone       string
	Email       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewPatient(id, firstName, lastName, middleName string, dateOfBirth time.Time, gender, phone, email string) *Patient {
	now := time.Now()
	return &Patient{
		ID:          id,
		FirstName:   firstName,
		LastName:    lastName,
		MiddleName:  middleName,
		DateOfBirth: dateOfBirth,
		Gender:      gender,
		Phone:       phone,
		Email:       email,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (p *Patient) Update(firstName, lastName, middleName, phone, email string) {
	p.FirstName = firstName
	p.LastName = lastName
	p.MiddleName = middleName
	p.Phone = phone
	p.Email = email
	p.UpdatedAt = time.Now()
}

func (p *Patient) FullName() string {
	if p.MiddleName != "" {
		return p.LastName + " " + p.FirstName + " " + p.MiddleName
	}
	return p.LastName + " " + p.FirstName
}
