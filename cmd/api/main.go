package main

import (
	"os"
	"os/signal"
	"syscall"

	_ "github.com/ariam/my-api/docs"
	"github.com/ariam/my-api/internal/config"
	"github.com/ariam/my-api/internal/middleware"
	"github.com/ariam/my-api/internal/router"
	"github.com/ariam/my-api/pkg/jwt"
	"github.com/ariam/my-api/pkg/logger"
	"github.com/ariam/my-api/pkg/response"
	"github.com/ariam/my-api/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"go.uber.org/zap"
)

// @title My API
// @version 1.0
// @description Production-ready REST API with Go Fiber
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:3000
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter token with Bearer prefix: "Bearer <token>"

func main() {
	cfg := config.Load()

	logger.Init(cfg.App.Env)
	defer logger.Sync()

	validator.Init()

	db, err := config.NewDatabase(&cfg.DB, cfg.App.Env)
	if err != nil {
		logger.Fatal("Database connection failed", zap.Error(err))
	}
	defer config.CloseDatabase(db)

	if err := config.RunMigration(db); err != nil {
		logger.Fatal("Migration failed", zap.Error(err))
	}

	jwtManager := jwt.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpireHours)

	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ErrorHandler: customErrorHandler,
	})

	middleware.SetupSecurity(app, cfg.App.Env)
	app.Use(middleware.RequestLogger())

	app.Get("/health", func(c *fiber.Ctx) error {
		sqlDB, _ := db.DB()
		dbStatus := "ok"
		if err := sqlDB.Ping(); err != nil {
			dbStatus = "error"
		}

		return response.Success(c, fiber.Map{
			"status":   "ok",
			"env":      cfg.App.Env,
			"database": dbStatus,
		})
	})

	app.Get("/swagger/*", swagger.HandlerDefault)

	router.Setup(app, db, jwtManager)

	go func() {
		if err := app.Listen(":" + cfg.App.Port); err != nil {
			logger.Fatal("Server error", zap.Error(err))
		}
	}()

	logger.Info("Server started", zap.String("port", cfg.App.Port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		logger.Error("Server shutdown error", zap.Error(err))
	}
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	logger.Error("Unhandled error",
		zap.Error(err),
		zap.String("path", c.Path()),
		zap.String("method", c.Method()),
	)

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"error":   err.Error(),
	})
}