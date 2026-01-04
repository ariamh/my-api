package middleware

import (
	"strings"

	"github.com/ariam/my-api/pkg/jwt"
	"github.com/ariam/my-api/pkg/response"
	"github.com/gofiber/fiber/v2"
)

func Auth(jwtManager *jwt.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Unauthorized(c, "Missing authorization header")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return response.Unauthorized(c, "Invalid authorization format")
		}

		claims, err := jwtManager.Validate(parts[1])
		if err != nil {
			return response.Unauthorized(c, err.Error())
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

func RoleRequired(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("role").(string)

		for _, role := range roles {
			if userRole == role {
				return c.Next()
			}
		}

		return response.Forbidden(c, "Insufficient permissions")
	}
}