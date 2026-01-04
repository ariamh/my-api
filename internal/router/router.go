package router

import (
	"github.com/ariam/my-api/internal/handler"
	"github.com/ariam/my-api/internal/middleware"
	"github.com/ariam/my-api/internal/repository"
	"github.com/ariam/my-api/internal/service"
	"github.com/ariam/my-api/pkg/jwt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, db *gorm.DB, jwtManager *jwt.JWTManager) {
	userRepo := repository.NewUserRepository(db)

	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, jwtManager)

	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	auth := v1.Group("/auth")
	auth.Post("/login", authHandler.Login)
	auth.Get("/me", middleware.Auth(jwtManager), authHandler.Me)

	users := v1.Group("/users")
	users.Post("/", userHandler.Create)
	users.Get("/", middleware.Auth(jwtManager), userHandler.FindAll)
	users.Get("/:id", middleware.Auth(jwtManager), userHandler.FindByID)
	users.Put("/:id", middleware.Auth(jwtManager), userHandler.Update)
	users.Delete("/:id", middleware.Auth(jwtManager), middleware.RoleRequired("admin"), userHandler.Delete)
}