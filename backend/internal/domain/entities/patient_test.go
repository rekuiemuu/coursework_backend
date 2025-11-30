package entities

import (
	"testing"
	"time"
)

func TestNewPatient(t *testing.T) {
	dob := time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC)
	patient := NewPatient(
		"patient-123",
		"Ivan",
		"Ivanov",
		"Ivanovich",
		dob,
		"male",
		"+79991234567",
		"ivan@example.com",
	)

	if patient == nil {
		t.Fatal("NewPatient() returned nil")
	}

	if patient.ID != "patient-123" {
		t.Errorf("ID = %v, want %v", patient.ID, "patient-123")
	}

	if patient.FirstName != "Ivan" {
		t.Errorf("FirstName = %v, want %v", patient.FirstName, "Ivan")
	}

	if patient.LastName != "Ivanov" {
		t.Errorf("LastName = %v, want %v", patient.LastName, "Ivanov")
	}

	if patient.MiddleName != "Ivanovich" {
		t.Errorf("MiddleName = %v, want %v", patient.MiddleName, "Ivanovich")
	}

	if !patient.DateOfBirth.Equal(dob) {
		t.Errorf("DateOfBirth = %v, want %v", patient.DateOfBirth, dob)
	}

	if patient.Gender != "male" {
		t.Errorf("Gender = %v, want %v", patient.Gender, "male")
	}
}

func TestPatient_Update(t *testing.T) {
	dob := time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC)
	patient := NewPatient("1", "Ivan", "Ivanov", "", dob, "male", "111", "old@test.com")

	before := patient.UpdatedAt
	time.Sleep(10 * time.Millisecond)

	patient.Update("Petr", "Petrov", "Petrovich", "+79999999999", "new@test.com")

	if patient.FirstName != "Petr" {
		t.Errorf("FirstName = %v, want %v", patient.FirstName, "Petr")
	}

	if patient.LastName != "Petrov" {
		t.Errorf("LastName = %v, want %v", patient.LastName, "Petrov")
	}

	if patient.MiddleName != "Petrovich" {
		t.Errorf("MiddleName = %v, want %v", patient.MiddleName, "Petrovich")
	}

	if patient.Phone != "+79999999999" {
		t.Errorf("Phone = %v, want %v", patient.Phone, "+79999999999")
	}

	if patient.Email != "new@test.com" {
		t.Errorf("Email = %v, want %v", patient.Email, "new@test.com")
	}

	if !patient.UpdatedAt.After(before) {
		t.Error("UpdatedAt should be updated after Update()")
	}
}

func TestPatient_FullName(t *testing.T) {
	dob := time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC)

	t.Run("with middle name", func(t *testing.T) {
		patient := NewPatient("1", "Ivan", "Ivanov", "Ivanovich", dob, "male", "", "")
		expected := "Ivanov Ivan Ivanovich"
		if patient.FullName() != expected {
			t.Errorf("FullName() = %v, want %v", patient.FullName(), expected)
		}
	})

	t.Run("without middle name", func(t *testing.T) {
		patient := NewPatient("1", "Ivan", "Ivanov", "", dob, "male", "", "")
		expected := "Ivanov Ivan"
		if patient.FullName() != expected {
			t.Errorf("FullName() = %v, want %v", patient.FullName(), expected)
		}
	})
}
