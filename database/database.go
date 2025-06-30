package database

import (
	"github.com/hmlylab/common/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDB(dsn string) *gorm.DB {
	logger := logger.NewLogger()

	DB, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{
		PrepareStmt: false,
	})

	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		return nil
	}
	logger.Info("Connected to database successfully")
	return DB
}
