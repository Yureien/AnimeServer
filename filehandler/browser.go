package filehandler

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yureien/animeserver/database/models"
)

type FileInfo struct {
	Name        string
	Path        string
	IsDirectory bool
	Size        int64
	ModTime     time.Time
	IsLeafDir   bool

	DbFile *models.File
}

func (f FileHandler) GetDirectory(path string) ([]FileInfo, error) {
	return f.readDirectory(path)
}

func (f FileHandler) readDirectory(requestPath string) ([]FileInfo, error) {
	fullPath := filepath.Join(f.cfg.AnimeDirectory, requestPath)

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	var files []FileInfo

	for _, entry := range entries {
		// Skip hidden files
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		var filePath string
		if requestPath == "" {
			filePath = entry.Name()
		} else {
			filePath = filepath.Join(requestPath, entry.Name())
		}

		dbFile, err := models.GetFileByPath(f.db, filePath)
		if err != nil {
			f.logger.Error("Error getting file by path", "error", err)
		}

		if dbFile == nil {
			f.ProcessFileAsync(filePath)
		}

		fileInfo := FileInfo{
			Name:        entry.Name(),
			Path:        filePath,
			IsDirectory: entry.IsDir(),
			Size:        info.Size(),
			ModTime:     info.ModTime(),
			DbFile:      dbFile,
		}

		// Check if directory is a leaf (no subdirectories)
		if entry.IsDir() {
			fileInfo.IsLeafDir = isLeafDirectory(filepath.Join(fullPath, entry.Name()))
		}

		files = append(files, fileInfo)
	}

	return files, nil
}

func isLeafDirectory(dirPath string) bool {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if entry.IsDir() {
			return false
		}
	}
	return true
}
