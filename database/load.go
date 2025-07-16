package database

import (
	"errors"
	"log/slog"

	slogGorm "github.com/orandin/slog-gorm"
	"github.com/yureien/animeserver/database/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func LoadDatabase(logger *slog.Logger, cfg *DatabaseConfig) (*gorm.DB, error) {
	if cfg.SQLite == nil {
		return nil, errors.New("no database configuration provided")
	}
	db, err := loadSQLiteDatabase(cfg.SQLite, logger)
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.File{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Token{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func loadSQLiteDatabase(cfg *SQLiteConfig, logger *slog.Logger) (*gorm.DB, error) {
	gormLogger := slogGorm.New(
		slogGorm.WithHandler(logger.Handler()),
		slogGorm.WithTraceAll(),
		slogGorm.SetLogLevel(slogGorm.DefaultLogType, slog.LevelDebug),
	)

	db, err := gorm.Open(sqlite.Open(cfg.Path), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}
