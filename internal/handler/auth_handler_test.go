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

// MockAuthService implements service.AuthService interface for testing
type MockAuthService struct {
	mock.Mock
}

// Login implements service.AuthService.Login
func (m *MockAuthService) Login(ctx context.Context, input *service.LoginInput) (*service.AuthResponse, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.AuthResponse), args.Error(1)
}

// setupAuthTestApp creates a Fiber app with auth routes for testing
func setupAuthTestApp(handler *AuthHandler) *fiber.App {
	validator.Init()
	app := fiber.New()

	// Auth routes
	app.Post("/auth/login", handler.Login)
	app.Get("/auth/me", handler.Me)

	return app
}

// TestAuthHandler_Login_Success tests successful login
func TestAuthHandler_Login_Success(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)
	app := setupAuthTestApp(handler)

	input := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}

	expectedResponse := &service.AuthResponse{
		Token: "jwt-token-here",
		User: &service.UserResponse{
			ID:    "user-uuid",
			Name:  "Test User",
			Email: "test@example.com",
			Role:  "user",
		},
	}

	mockService.On("Login", mock.Anything, mock.AnythingOfType("*service.LoginInput")).Return(expectedResponse, nil)

	body, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

// TestAuthHandler_Login_InvalidJSON tests login with invalid JSON body
func TestAuthHandler_Login_InvalidJSON(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)
	app := setupAuthTestApp(handler)

	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

// TestAuthHandler_Login_ValidationError tests login with validation failure
func TestAuthHandler_Login_ValidationError(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)
	app := setupAuthTestApp(handler)

	input := map[string]string{
		"email":    "invalid-email",
		"password": "",
	}

	body, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
}

// TestAuthHandler_Login_InvalidCredentials tests login with wrong credentials
func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)
	app := setupAuthTestApp(handler)

	input := map[string]string{
		"email":    "test@example.com",
		"password": "wrongpassword",
	}

	mockService.On("Login", mock.Anything, mock.AnythingOfType("*service.LoginInput")).Return(nil, service.ErrInvalidCredentials)

	body, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	mockService.AssertExpectations(t)
}

// TestAuthHandler_Login_ServiceError tests login when service returns unexpected error
// Requirements: 1.5
func TestAuthHandler_Login_ServiceError(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)
	app := setupAuthTestApp(handler)

	input := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}

	// Simulate an unexpected service error (e.g., database connection failure)
	mockService.On("Login", mock.Anything, mock.AnythingOfType("*service.LoginInput")).Return(nil, errors.New("database connection failed"))

	body, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	// Verify response body contains error message
	var respBody response.Response
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.False(t, respBody.Success)
	assert.Equal(t, "Login failed", respBody.Error)

	mockService.AssertExpectations(t)
}

// TestAuthHandler_Me implements table-driven tests for the /auth/me endpoint
// Requirements: 2.1, 2.2
func TestAuthHandler_Me(t *testing.T) {
	tests := []struct {
		name           string
		setupContext   func(*fiber.App) *fiber.App
		expectedStatus int
		checkResponse  func(*testing.T, response.Response)
	}{
		{
			name: "returns user context values with 200 status",
			setupContext: func(app *fiber.App) *fiber.App {
				// Create a new app with middleware that sets context values
				mockService := new(MockAuthService)
				handler := NewAuthHandler(mockService)
				validator.Init()
				newApp := fiber.New()

				// Middleware to simulate authenticated user context
				newApp.Use(func(c *fiber.Ctx) error {
					c.Locals("user_id", "test-user-id-123")
					c.Locals("email", "test@example.com")
					c.Locals("role", "user")
					return c.Next()
				})

				newApp.Get("/auth/me", handler.Me)
				return newApp
			},
			expectedStatus: fiber.StatusOK,
			checkResponse: func(t *testing.T, resp response.Response) {
				assert.True(t, resp.Success)
				data, ok := resp.Data.(map[string]interface{})
				assert.True(t, ok, "Data should be a map")
				assert.Equal(t, "test-user-id-123", data["user_id"])
				assert.Equal(t, "test@example.com", data["email"])
				assert.Equal(t, "user", data["role"])
			},
		},
		{
			name: "includes all context fields (user_id, email, role)",
			setupContext: func(app *fiber.App) *fiber.App {
				// Create a new app with middleware that sets all context fields
				mockService := new(MockAuthService)
				handler := NewAuthHandler(mockService)
				validator.Init()
				newApp := fiber.New()

				// Middleware to simulate authenticated admin user context
				newApp.Use(func(c *fiber.Ctx) error {
					c.Locals("user_id", "admin-uuid-456")
					c.Locals("email", "admin@example.com")
					c.Locals("role", "admin")
					return c.Next()
				})

				newApp.Get("/auth/me", handler.Me)
				return newApp
			},
			expectedStatus: fiber.StatusOK,
			checkResponse: func(t *testing.T, resp response.Response) {
				assert.True(t, resp.Success)
				data, ok := resp.Data.(map[string]interface{})
				assert.True(t, ok, "Data should be a map")

				// Verify all three context fields are present
				_, hasUserID := data["user_id"]
				_, hasEmail := data["email"]
				_, hasRole := data["role"]

				assert.True(t, hasUserID, "Response should include user_id")
				assert.True(t, hasEmail, "Response should include email")
				assert.True(t, hasRole, "Response should include role")

				// Verify the values
				assert.Equal(t, "admin-uuid-456", data["user_id"])
				assert.Equal(t, "admin@example.com", data["email"])
				assert.Equal(t, "admin", data["role"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup app with context
			mockService := new(MockAuthService)
			handler := NewAuthHandler(mockService)
			baseApp := setupAuthTestApp(handler)
			app := tt.setupContext(baseApp)

			// Create request
			req := httptest.NewRequest("GET", "/auth/me", nil)
			req.Header.Set("Content-Type", "application/json")

			// Execute request
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			// Parse response
			var respBody response.Response
			err = json.NewDecoder(resp.Body).Decode(&respBody)
			assert.NoError(t, err)

			// Check response
			if tt.checkResponse != nil {
				tt.checkResponse(t, respBody)
			}
		})
	}
}
