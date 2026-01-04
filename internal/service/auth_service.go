package service

import (
	"context"

	"github.com/ariam/my-api/internal/repository"
	"github.com/ariam/my-api/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token string        `json:"token"`
	User  *UserResponse `json:"user"`
}

type AuthService interface {
	Login(ctx context.Context, input *LoginInput) (*AuthResponse, error)
}

type authService struct {
	userRepo   repository.UserRepository
	jwtManager *jwt.JWTManager
}

func NewAuthService(userRepo repository.UserRepository, jwtManager *jwt.JWTManager) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

func (s *authService) Login(ctx context.Context, input *LoginInput) (*AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, ErrInvalidCredentials
	}

	token, err := s.jwtManager.Generate(user.ID.String(), user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  toUserResponse(user),
	}, nil
}