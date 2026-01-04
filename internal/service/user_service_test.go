package service

import (
	"context"
	// "errors"
	"testing"

	"github.com/ariam/my-api/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindAll(ctx context.Context, page, perPage int) ([]model.User, int64, error) {
	args := m.Called(ctx, page, perPage)
	return args.Get(0).([]model.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) Update(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestUserService_Create_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	input := &CreateUserInput{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	mockRepo.On("FindByEmail", ctx, input.Email).Return(nil, gorm.ErrRecordNotFound)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*model.User")).Return(nil)

	result, err := service.Create(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, input.Name, result.Name)
	assert.Equal(t, input.Email, result.Email)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Create_EmailExists(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	input := &CreateUserInput{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	existingUser := &model.User{
		Base:  model.Base{ID: uuid.New()},
		Email: input.Email,
	}

	mockRepo.On("FindByEmail", ctx, input.Email).Return(existingUser, nil)

	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrEmailAlreadyExists, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_FindByID_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	user := &model.User{
		Base:  model.Base{ID: userID},
		Name:  "John Doe",
		Email: "john@example.com",
		Role:  "user",
	}

	mockRepo.On("FindByID", ctx, userID.String()).Return(user, nil)

	result, err := service.FindByID(ctx, userID.String())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.Name, result.Name)
	mockRepo.AssertExpectations(t)
}

func TestUserService_FindByID_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	mockRepo.On("FindByID", ctx, "invalid-id").Return(nil, gorm.ErrRecordNotFound)

	result, err := service.FindByID(ctx, "invalid-id")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrUserNotFound, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Delete_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	userID := uuid.New()
	user := &model.User{
		Base: model.Base{ID: userID},
	}

	mockRepo.On("FindByID", ctx, userID.String()).Return(user, nil)
	mockRepo.On("Delete", ctx, userID.String()).Return(nil)

	err := service.Delete(ctx, userID.String())

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Delete_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	mockRepo.On("FindByID", ctx, "invalid-id").Return(nil, gorm.ErrRecordNotFound)

	err := service.Delete(ctx, "invalid-id")

	assert.Error(t, err)
	assert.Equal(t, ErrUserNotFound, err)
	mockRepo.AssertExpectations(t)
}