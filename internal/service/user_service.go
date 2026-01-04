package service

import (
	"context"
	"errors"

	"github.com/ariam/my-api/internal/model"
	"github.com/ariam/my-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type CreateUserInput struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UpdateUserInput struct {
	Name string `json:"name" validate:"omitempty,min=2,max=100"`
}

type UserResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}

type UserService interface {
	Create(ctx context.Context, input *CreateUserInput) (*UserResponse, error)
	FindByID(ctx context.Context, id string) (*UserResponse, error)
	FindAll(ctx context.Context, page, perPage int) ([]UserResponse, int64, error)
	Update(ctx context.Context, id string, input *UpdateUserInput) (*UserResponse, error)
	Delete(ctx context.Context, id string) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Create(ctx context.Context, input *CreateUserInput) (*UserResponse, error) {
	existing, _ := s.userRepo.FindByEmail(ctx, input.Email)
	if existing != nil {
		return nil, ErrEmailAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     "user",
		IsActive: true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return toUserResponse(user), nil
}

func (s *userService) FindByID(ctx context.Context, id string) (*UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return toUserResponse(user), nil
}

func (s *userService) FindAll(ctx context.Context, page, perPage int) ([]UserResponse, int64, error) {
	users, total, err := s.userRepo.FindAll(ctx, page, perPage)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]UserResponse, len(users))
	for i, user := range users {
		responses[i] = *toUserResponse(&user)
	}

	return responses, total, nil
}

func (s *userService) Update(ctx context.Context, id string, input *UpdateUserInput) (*UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if input.Name != "" {
		user.Name = input.Name
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return toUserResponse(user), nil
}

func (s *userService) Delete(ctx context.Context, id string) error {
	_, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	return s.userRepo.Delete(ctx, id)
}

func toUserResponse(user *model.User) *UserResponse {
	return &UserResponse{
		ID:       user.ID.String(),
		Name:     user.Name,
		Email:    user.Email,
		Role:     user.Role,
		IsActive: user.IsActive,
	}
}