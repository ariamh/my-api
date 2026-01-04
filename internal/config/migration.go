package config

import (
	"github.com/ariam/my-api/internal/model"
	"github.com/ariam/my-api/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func RunMigration(db *gorm.DB) error {
	logger.Info("Running database migrations...")

	err := db.AutoMigrate(
		&model.User{},
	)

	if err != nil {
		logger.Error("Migration failed", zap.Error(err))
		return err
	}

	logger.Info("Database migrations completed")
	return nil
}