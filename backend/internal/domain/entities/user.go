package entities

import (
	"time"
)

type UserRole string

const (
	RoleAdmin  UserRole = "admin"
	RoleDoctor UserRole = "doctor"
	RoleTech   UserRole = "technician"
)

type User struct {
	ID           string
	Username     string
	Email        string
	PasswordHash string
	Role         UserRole
	FirstName    string
	LastName     string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(id, username, email, passwordHash string, role UserRole, firstName, lastName string) *User {
	now := time.Now()
	return &User{
		ID:           id,
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
		FirstName:    firstName,
		LastName:     lastName,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

func (u *User) Activate() {
	u.IsActive = true
	u.UpdatedAt = time.Now()
}

func (u *User) UpdatePassword(newPasswordHash string) {
	u.PasswordHash = newPasswordHash
	u.UpdatedAt = time.Now()
}

func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}
