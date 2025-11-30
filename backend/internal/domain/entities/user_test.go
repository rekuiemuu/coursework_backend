package entities

import (
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	user := NewUser(
		"123e4567-e89b-12d3-a456-426614174000",
		"testuser",
		"test@example.com",
		"hashedpassword",
		RoleDoctor,
		"John",
		"Doe",
	)

	if user == nil {
		t.Fatal("NewUser() returned nil")
	}

	if user.ID != "123e4567-e89b-12d3-a456-426614174000" {
		t.Errorf("ID = %v, want %v", user.ID, "123e4567-e89b-12d3-a456-426614174000")
	}

	if user.Username != "testuser" {
		t.Errorf("Username = %v, want %v", user.Username, "testuser")
	}

	if user.Email != "test@example.com" {
		t.Errorf("Email = %v, want %v", user.Email, "test@example.com")
	}

	if user.Role != RoleDoctor {
		t.Errorf("Role = %v, want %v", user.Role, RoleDoctor)
	}

	if !user.IsActive {
		t.Error("IsActive should be true by default")
	}

	if user.FirstName != "John" {
		t.Errorf("FirstName = %v, want %v", user.FirstName, "John")
	}

	if user.LastName != "Doe" {
		t.Errorf("LastName = %v, want %v", user.LastName, "Doe")
	}
}

func TestUser_Deactivate(t *testing.T) {
	user := NewUser("1", "user", "email@test.com", "hash", RoleAdmin, "First", "Last")
	
	before := user.UpdatedAt
	time.Sleep(10 * time.Millisecond)
	
	user.Deactivate()

	if user.IsActive {
		t.Error("User should be deactivated")
	}

	if !user.UpdatedAt.After(before) {
		t.Error("UpdatedAt should be updated after deactivation")
	}
}

func TestUser_Activate(t *testing.T) {
	user := NewUser("1", "user", "email@test.com", "hash", RoleTech, "First", "Last")
	user.Deactivate()

	before := user.UpdatedAt
	time.Sleep(10 * time.Millisecond)
	
	user.Activate()

	if !user.IsActive {
		t.Error("User should be activated")
	}

	if !user.UpdatedAt.After(before) {
		t.Error("UpdatedAt should be updated after activation")
	}
}

func TestUser_UpdatePassword(t *testing.T) {
	user := NewUser("1", "user", "email@test.com", "oldhash", RoleDoctor, "First", "Last")
	
	before := user.UpdatedAt
	time.Sleep(10 * time.Millisecond)
	
	newHash := "newhash"
	user.UpdatePassword(newHash)

	if user.PasswordHash != newHash {
		t.Errorf("PasswordHash = %v, want %v", user.PasswordHash, newHash)
	}

	if !user.UpdatedAt.After(before) {
		t.Error("UpdatedAt should be updated after password change")
	}
}

func TestUser_FullName(t *testing.T) {
	user := NewUser("1", "user", "email@test.com", "hash", RoleAdmin, "John", "Smith")
	
	fullName := user.FullName()
	expected := "John Smith"

	if fullName != expected {
		t.Errorf("FullName() = %v, want %v", fullName, expected)
	}
}
