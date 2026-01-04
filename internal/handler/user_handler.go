package handler

import (
	"errors"
	"strconv"

	"github.com/ariam/my-api/internal/service"
	"github.com/ariam/my-api/pkg/response"
	"github.com/ariam/my-api/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Create godoc
// @Summary Create new user
// @Description Register a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param request body service.CreateUserInput true "User data"
// @Success 201 {object} response.Response{data=service.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 422 {object} response.Response
// @Router /users [post]
func (h *UserHandler) Create(c *fiber.Ctx) error {
	var input service.CreateUserInput

	if err := c.BodyParser(&input); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if errs := validator.Validate(&input); len(errs) > 0 {
		return response.ValidationError(c, errs)
	}

	user, err := h.userService.Create(c.Context(), &input)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			return response.BadRequest(c, err.Error())
		}
		return response.InternalServerError(c, "Failed to create user")
	}

	return response.Created(c, user)
}

// FindByID godoc
// @Summary Get user by ID
// @Description Get user details by ID
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} response.Response{data=service.UserResponse}
// @Failure 404 {object} response.Response
// @Router /users/{id} [get]
func (h *UserHandler) FindByID(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.userService.FindByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return response.NotFound(c, err.Error())
		}
		return response.InternalServerError(c, "Failed to fetch user")
	}

	return response.Success(c, user)
}

// FindAll godoc
// @Summary Get all users
// @Description Get paginated list of users
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} response.Response{data=response.PaginatedData}
// @Router /users [get]
func (h *UserHandler) FindAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	users, total, err := h.userService.FindAll(c.Context(), page, perPage)
	if err != nil {
		return response.InternalServerError(c, "Failed to fetch users")
	}

	return response.Paginated(c, users, total, page, perPage)
}

// Update godoc
// @Summary Update user
// @Description Update user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body service.UpdateUserInput true "User data"
// @Success 200 {object} response.Response{data=service.UserResponse}
// @Failure 404 {object} response.Response
// @Failure 422 {object} response.Response
// @Router /users/{id} [put]
func (h *UserHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var input service.UpdateUserInput
	if err := c.BodyParser(&input); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if errs := validator.Validate(&input); len(errs) > 0 {
		return response.ValidationError(c, errs)
	}

	user, err := h.userService.Update(c.Context(), id, &input)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return response.NotFound(c, err.Error())
		}
		return response.InternalServerError(c, "Failed to update user")
	}

	return response.Success(c, user)
}

// Delete godoc
// @Summary Delete user
// @Description Delete user by ID (admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 204 "No Content"
// @Failure 404 {object} response.Response
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.userService.Delete(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return response.NotFound(c, err.Error())
		}
		return response.InternalServerError(c, "Failed to delete user")
	}

	return response.NoContent(c)
}