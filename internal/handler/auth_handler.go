package handler

import (
	"errors"

	"github.com/ariam/my-api/internal/service"
	"github.com/ariam/my-api/pkg/response"
	"github.com/ariam/my-api/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body service.LoginInput true "Login credentials"
// @Success 200 {object} response.Response{data=service.AuthResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var input service.LoginInput

	if err := c.BodyParser(&input); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if errs := validator.Validate(&input); len(errs) > 0 {
		return response.ValidationError(c, errs)
	}

	result, err := h.authService.Login(c.Context(), &input)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return response.Unauthorized(c, "Invalid email or password")
		}
		return response.InternalServerError(c, "Login failed")
	}

	return response.Success(c, result)
}

// Me godoc
// @Summary Get current user
// @Description Get authenticated user info from token
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/me [get]
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	return response.Success(c, fiber.Map{
		"user_id": c.Locals("user_id"),
		"email":   c.Locals("email"),
		"role":    c.Locals("role"),
	})
}