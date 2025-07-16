package filehandler

import (
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/gorm"
)

type FileHandler struct {
	logger *slog.Logger
	cfg    *Config
	db     *gorm.DB

	processingQueue chan string
}

func NewFileHandler(logger *slog.Logger, cfg *Config, db *gorm.DB) *FileHandler {
	fileHandler := &FileHandler{
		logger:          logger,
		cfg:             cfg,
		db:              db,
		processingQueue: make(chan string, 10000),
	}
	fileHandler.startProcessingQueue()
	return fileHandler
}

func (f FileHandler) GetFile(path string) (os.FileInfo, io.ReadSeekCloser, error) {
	fullPath := filepath.Join(f.cfg.AnimeDirectory, path)

	// Security check
	absAnimePath, err := filepath.Abs(f.cfg.AnimeDirectory)
	if err != nil {
		return nil, nil, err
	}

	absRequestPath, err := filepath.Abs(fullPath)
	if err != nil {
		return nil, nil, err
	}

	if !strings.HasPrefix(absRequestPath, absAnimePath) {
		return nil, nil, errors.New("access denied due to security check")
	}

	// Check if file exists
	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, nil, err
	}

	if info.IsDir() {
		return nil, nil, err
	}

	// Serve file
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, nil, err
	}

	return info, file, nil
}
