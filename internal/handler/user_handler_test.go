package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/ariam/my-api/internal/service"
	"github.com/ariam/my-api/pkg/response"
	"github.com/ariam/my-api/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Create(ctx context.Context, input *service.CreateUserInput) (*service.UserResponse, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.UserResponse), args.Error(1)
}

func (m *MockUserService) FindByID(ctx context.Context, id string) (*service.UserResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.UserResponse), args.Error(1)
}

func (m *MockUserService) FindAll(ctx context.Context, page, perPage int) ([]service.UserResponse, int64, error) {
	args := m.Called(ctx, page, perPage)
	return args.Get(0).([]service.UserResponse), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserService) Update(ctx context.Context, id string, input *service.UpdateUserInput) (*service.UserResponse, error) {
	args := m.Called(ctx, id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.UserResponse), args.Error(1)
}

func (m *MockUserService) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupTestApp(handler *UserHandler) *fiber.App {
	validator.Init()
	app := fiber.New()
	app.Post("/users", handler.Create)
	app.Get("/users", handler.FindAll)
	app.Get("/users/:id", handler.FindByID)
	app.Put("/users/:id", handler.Update)
	app.Delete("/users/:id", handler.Delete)
	return app
}

// TestUserHandler_Create implements table-driven tests for the Create endpoint
// Requirements: 3.1, 3.2, 3.3, 3.4, 3.5
func TestUserHandler_Create(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockUserService)
		body           interface{}
		expectedStatus int
		checkResponse  func(*testing.T, response.Response)
	}{
		{
			name: "valid user creation returns 201 with user data",
			setupMock: func(m *MockUserService) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*service.CreateUserInput")).
					Return(&service.UserResponse{
						ID:    "test-uuid",
						Name:  "John Doe",
						Email: "john@example.com",
						Role:  "user",
					}, nil)
			},
			body: map[string]string{
				"name":     "John Doe",
				"email":    "john@example.com",
				"password": "password123",
			},
			expectedStatus: fiber.StatusCreated,
			checkResponse: func(t *testing.T, resp response.Response) {
				assert.True(t, resp.Success)
				data, ok := resp.Data.(map[string]interface{})
				assert.True(t, ok, "Data should be a map")
				assert.Equal(t, "test-uuid", data["id"])
				assert.Equal(t, "John Doe", data["name"])
				assert.Equal(t, "john@example.com", data["email"])
			},
		},
		{
			name:      "invalid JSON body returns 400",
			setupMock: nil,
			body:      "invalid json",
			expectedStatus: fiber.StatusBadRequest,
			checkResponse: func(t *testing.T, resp response.Response) {
				assert.False(t, resp.Success)
				assert.Equal(t, "Invalid request body", resp.Error)
			},
		},
		{
			name:      "validation failure returns 422",
			setupMock: nil,
			body: map[string]string{
				"name":     "",
				"email":    "invalid",
				"password": "123",
			},
			expectedStatus: fiber.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, resp response.Response) {
				assert.False(t, resp.Success)
			},
		},
		{
			name: "duplicate email returns 400",
			setupMock: func(m *MockUserService) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*service.CreateUserInput")).
					Return(nil, service.ErrEmailAlreadyExists)
			},
			body: map[string]string{
				"name":     "John Doe",
				"email":    "existing@example.com",
				"password": "password123",
			},
			expectedStatus: fiber.StatusBadRequest,
			checkResponse: func(t *testing.T, resp response.Response) {
				assert.False(t, resp.Success)
				assert.Equal(t, "email already exists", resp.Error)
			},
		},
		{
			name: "service error returns 500",
			setupMock: func(m *MockUserService) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*service.CreateUserInput")).
					Return(nil, errors.New("database connection failed"))
			},
			body: map[string]string{
				"name":     "John Doe",
				"email":    "john@example.com",
				"password": "password123",
			},
			expectedStatus: fiber.StatusInternalServerError,
			checkResponse: func(t *testing.T, resp response.Response) {
				assert.False(t, resp.Success)
				assert.Equal(t, "Failed to create user", resp.Error)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockUserService)
			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}
			handler := NewUserHandler(mockService)
			app := setupTestApp(handler)

			var body []byte
			switch v := tt.body.(type) {
			case string:
				body = []byte(v)
			default:
				body, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var respBody response.Response
			err = json.NewDecoder(resp.Body).Decode(&respBody)
			assert.NoError(t, err)

			if tt.checkResponse != nil {
				tt.checkResponse(t, respBody)
			}

			mockService.AssertExpectations(t)
		})
	}
}

// TestUserHandler_FindByID implements table-driven tests for the FindByID endpoint
// Requirements: 4.1, 4.2, 4.3
func TestUserHandler_FindByID(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		setupMock      func(*MockUserService)
		expectedStatus int
		checkResponse  func(*testing.T, response.Response)
	}{
		{
			name:   "valid user ID returns 200 with user data",
			userID: "test-uuid",
			setupMock: func(m *MockUserService) {
				m.On("FindByID", mock.Anything, "test-uuid").
					Return(&service.UserResponse{
						ID:    "test-uuid",
						Name:  "John Doe",
						Email: "john@example.com",
						Role:  "user",
					}, nil)
			},
			expectedStatus: fiber.StatusOK,
			checkResponse: func(t *testing.T, resp response.Response) {
				assert.True(t, resp.Success)
				data, ok := resp.Data.(map[string]interface{})
				assert.True(t, ok, "Data should be a map")
				assert.Equal(t, "test-uuid", data["id"])
				assert.Equal(t, "John Doe", data["name"])
				assert.Equal(t, "john@example.com", data["email"])
			},
		},
		{
			name:   "non-existent user ID returns 404",
			userID: "invalid-id",
			setupMock: func(m *MockUserService) {
				m.On("FindByID", mock.Anything, "invalid-id").
					Return(nil, service.ErrUserNotFound)
			},
			expectedStatus: fiber.StatusNotFound,
			checkResponse: func(t *testing.T, resp response.Response) {
				assert.False(t, resp.Success)
				assert.Equal(t, "user not found", resp.Error)
			},
		},
		{
			name:   "service error returns 500",
			userID: "test-uuid",
			setupMock: func(m *MockUserService) {
				m.On("FindByID", mock.Anything, "test-uuid").
					Return(nil, errors.New("database connection failed"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			checkResponse: func(t *testing.T, resp response.Response) {
				assert.False(t, resp.Success)
				assert.Equal(t, "Failed to fetch user", resp.Error)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockUserService)
			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}
			handler := NewUserHandler(mockService)
			app := setupTestApp(handler)

			req := httptest.NewRequest("GET", "/users/"+tt.userID, nil)

			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var respBody response.Response
			err = json.NewDecoder(resp.Body).Decode(&respBody)
			assert.NoError(t, err)

			if tt.checkResponse != nil {
				tt.checkResponse(t, respBody)
			}

			mockService.AssertExpectations(t)
		})
	}
}

// TestUserHandler_FindAll implements table-driven tests for the FindAll endpoint
// Requirements: 5.1, 5.2, 5.3, 5.4, 5.5
func TestUserHandler_FindAll(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		setupMock      func(*MockUserService)
		expectedStatus int
		checkResponse  func(*testing.T, response.Response)
	}{
		{
			name:        "default pagination (no params) returns 200",
			queryParams: "",
			setupMock: func(m *MockUserService) {
				m.On("FindAll", mock.Anything, 1, 10).
					Return([]service.UserResponse{
						{ID: "user-1", Name: "User One", Email: "user1@example.com", Role: "user"},
						{ID: "user-2", Name: "User Two", Email: "user2@example.com", Role: "user"},
					}, int64(2), nil)
			},
			expectedStatus: fiber.StatusOK,
			checkResponse: func(t *testing.T, resp response.Response) {
				assert.True(t, resp.Success)
				data, ok := resp.Data.(map[string]interface{})
				assert.True(t, ok, "Data should be a map")
				assert.Equal(t, float64(1), data["page"])
				assert.Equal(t, float64(10), data["per_page"])
				assert.Equal(t, float64(2), data["total"])
				items, ok := data["items"].([]interface{})
				assert.True(t, ok, "Items should be an array")
				assert.Len(t, items, 2)
			},
		},
		{
			name:        "custom pagination params returns 200",
			queryParams: "?page=2&per_page=5",
			setupMock: func(m *MockUserService) {
				m.On("FindAll", mock.Anything, 2, 5).
					Return([]service.UserResponse{
						{ID: "user-6", Name: "User Six", Email: "user6@example.com", Role: "user"},
					}, int64(6), nil)
			},
			expectedStatus: fiber.StatusOK,
			checkResponse: func(t *testing.T, resp response.Response) {
				assert.True(t, resp.Success)
				data, ok := resp.Data.(map[string]interface{})
				assert.True(t, ok, "Data should be a map")
				assert.Equal(t, float64(2), data["page"])
				assert.Equal(t, float64(5), data["per_page"])
				assert.Equal(t, float64(6), data["total"])
				assert.Equal(t, float64(2), data["total_pages"])
			},
		},
		{
			name:        "invalid page (< 1) normalized to 1",
			queryParams: "?page=0&per_page=10",
			setupMock: func(m *MockUserService) {
				m.On("FindAll", mock.Anything, 1, 10).
					Return([]service.UserResponse{}, int64(0), nil)
			},
			expectedStatus: fiber.StatusOK,
			checkResponse: func(t *testing.T, resp response.Response) {
				assert.True(t, resp.Success)
				data, ok := resp.Data.(map[string]interface{})
				assert.True(t, ok, "Data should be a map")
				assert.Equal(t, float64(1), data["page"])
			},
		},
		{
			name:        "invalid per_page (< 1) normalized to 10",
			queryParams: "?page=1&per_page=0",
			setupMock: func(m *MockUserService) {
				m.On("FindAll", mock.Anything, 1, 10).
					Return([]service.UserResponse{}, int64(0), nil)
			},
			expectedStatus: fiber.StatusOK,
			checkResponse: func(t *testing.T, resp response.Response) {
				assert.True(t, resp.Success)
				data, ok := resp.Data.(map[string]interface{})
				assert.True(t, ok, "Data should be a map")
				assert.Equal(t, float64(10), data["per_page"])
			},
		},
		{
			name:        "invalid per_page (> 100) normalized to 10",
			queryParams: "?page=1&per_page=150",
			setupMock: func(m *MockUserService) {
				m.On("FindAll", mock.Anything, 1, 10).
					Return([]service.UserResponse{}, int64(0), nil)
			},
			expectedStatus: fiber.StatusOK,
			checkResponse: func(t *testing.T, resp response.Response) {
				assert.True(t, resp.Success)
				data, ok := resp.Data.(map[string]interface{})
				assert.True(t, ok, "Data should be a map")
				assert.Equal(t, float64(10), data["per_page"])
			},
		},
		{
			name:        "service error returns 500",
			queryParams: "",
			setupMock: func(m *MockUserService) {
				m.On("FindAll", mock.Anything, 1, 10).
					Return([]service.UserResponse{}, int64(0), errors.New("database connection failed"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			checkResponse: func(t *testing.T, resp response.Response) {
				assert.False(t, resp.Success)
				assert.Equal(t, "Failed to fetch users", resp.Error)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockUserService)
			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}
			handler := NewUserHandler(mockService)
			app := setupTestApp(handler)

			req := httptest.NewRequest("GET", "/users"+tt.queryParams, nil)

			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var respBody response.Response
			err = json.NewDecoder(resp.Body).Decode(&respBody)
			assert.NoError(t, err)

			if tt.checkResponse != nil {
				tt.checkResponse(t, respBody)
			}

			mockService.AssertExpectations(t)
		})
	}
}

// TestUserHandler_Update implements table-driven tests for the Update endpoint
// Requirements: 6.1, 6.2, 6.3, 6.4, 6.5
func TestUserHandler_Update(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		setupMock      func(*MockUserService)
		body           interface{}
		expectedStatus int
		checkResponse  func(*testing.T, response.Response)
	}{
		{
			name:   "valid update returns 200",
			userID: "test-uuid",
			setupMock: func(m *MockUserService) {
				m.On("Update", mock.Anything, "test-uuid", mock.AnythingOfType("*service.UpdateUserInput")).
					Return(&service.UserResponse{
						ID:    "test-uuid",
						Name:  "Updated Name",
						Email: "john@example.com",
						Role:  "user",
					}, nil)
			},
			body: map[string]string{
				"name": "Updated Name",
			},
			expectedStatus: fiber.StatusOK,
			checkResponse: func(t *testing.T, resp response.Response) {
				assert.True(t, resp.Success)
				data, ok := resp.Data.(map[string]interface{})
				assert.True(t, ok, "Data should be a map")
				assert.Equal(t, "test-uuid", data["id"])
				assert.Equal(t, "Updated Name", data["name"])
			},
		},
		{
			name:      "invalid JSON returns 400",
			userID:    "test-uuid",
			setupMock: nil,
			body:      "invalid json",
			expectedStatus: fiber.StatusBadRequest,
			checkResponse: func(t *testing.T, resp response.Response) {
				assert.False(t, resp.Success)
				assert.Equal(t, "Invalid request body", resp.Error)
			},
		},
		{
			name:      "validation failure returns 422",
			userID:    "test-uuid",
			setupMock: nil,
			body: map[string]string{
				"name": "A",
			},
			expectedStatus: fiber.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, resp response.Response) {
				assert.False(t, resp.Success)
			},
		},
		{
			name:   "not found returns 404",
			userID: "non-existent-id",
			setupMock: func(m *MockUserService) {
				m.On("Update", mock.Anything, "non-existent-id", mock.AnythingOfType("*service.UpdateUserInput")).
					Return(nil, service.ErrUserNotFound)
			},
			body: map[string]string{
				"name": "Updated Name",
			},
			expectedStatus: fiber.StatusNotFound,
			checkResponse: func(t *testing.T, resp response.Response) {
				assert.False(t, resp.Success)
				assert.Equal(t, "user not found", resp.Error)
			},
		},
		{
			name:   "service error returns 500",
			userID: "test-uuid",
			setupMock: func(m *MockUserService) {
				m.On("Update", mock.Anything, "test-uuid", mock.AnythingOfType("*service.UpdateUserInput")).
					Return(nil, errors.New("database connection failed"))
			},
			body: map[string]string{
				"name": "Updated Name",
			},
			expectedStatus: fiber.StatusInternalServerError,
			checkResponse: func(t *testing.T, resp response.Response) {
				assert.False(t, resp.Success)
				assert.Equal(t, "Failed to update user", resp.Error)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockUserService)
			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}
			handler := NewUserHandler(mockService)
			app := setupTestApp(handler)

			var body []byte
			switch v := tt.body.(type) {
			case string:
				body = []byte(v)
			default:
				body, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest("PUT", "/users/"+tt.userID, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var respBody response.Response
			err = json.NewDecoder(resp.Body).Decode(&respBody)
			assert.NoError(t, err)

			if tt.checkResponse != nil {
				tt.checkResponse(t, respBody)
			}

			mockService.AssertExpectations(t)
		})
	}
}

// TestUserHandler_Delete implements table-driven tests for the Delete endpoint
// Requirements: 7.1, 7.2, 7.3
func TestUserHandler_Delete(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		setupMock      func(*MockUserService)
		expectedStatus int
		checkResponse  func(*testing.T, *response.Response)
	}{
		{
			name:   "valid delete returns 204",
			userID: "test-uuid",
			setupMock: func(m *MockUserService) {
				m.On("Delete", mock.Anything, "test-uuid").Return(nil)
			},
			expectedStatus: fiber.StatusNoContent,
			checkResponse:  nil,
		},
		{
			name:   "not found returns 404",
			userID: "non-existent-id",
			setupMock: func(m *MockUserService) {
				m.On("Delete", mock.Anything, "non-existent-id").Return(service.ErrUserNotFound)
			},
			expectedStatus: fiber.StatusNotFound,
			checkResponse: func(t *testing.T, resp *response.Response) {
				assert.False(t, resp.Success)
				assert.Equal(t, "user not found", resp.Error)
			},
		},
		{
			name:   "service error returns 500",
			userID: "test-uuid",
			setupMock: func(m *MockUserService) {
				m.On("Delete", mock.Anything, "test-uuid").Return(errors.New("database connection failed"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			checkResponse: func(t *testing.T, resp *response.Response) {
				assert.False(t, resp.Success)
				assert.Equal(t, "Failed to delete user", resp.Error)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockUserService)
			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}
			handler := NewUserHandler(mockService)
			app := setupTestApp(handler)

			req := httptest.NewRequest("DELETE", "/users/"+tt.userID, nil)

			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.checkResponse != nil {
				var respBody response.Response
				err = json.NewDecoder(resp.Body).Decode(&respBody)
				assert.NoError(t, err)
				tt.checkResponse(t, &respBody)
			}

			mockService.AssertExpectations(t)
		})
	}
}