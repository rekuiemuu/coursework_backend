package usecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/project-capillary/backend/internal/application/dto"
	"github.com/project-capillary/backend/internal/domain/entities"
	"github.com/project-capillary/backend/internal/domain/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	userRepo repositories.UserRepository
}

func NewUserUseCase(userRepo repositories.UserRepository) *UserUseCase {
	return &UserUseCase{userRepo: userRepo}
}

func (uc *UserUseCase) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := entities.NewUser(
		uuid.New().String(),
		req.Username,
		req.Email,
		string(hashedPassword),
		entities.UserRole(req.Role),
		req.FirstName,
		req.LastName,
	)

	err = uc.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      string(user.Role),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (uc *UserUseCase) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := uc.userRepo.GetByUsername(ctx, req.Username)
	if err != nil || user == nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	token := uuid.New().String()

	return &dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      string(user.Role),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}

func (uc *UserUseCase) GetUser(ctx context.Context, id string) (*dto.UserResponse, error) {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil || user == nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      string(user.Role),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
