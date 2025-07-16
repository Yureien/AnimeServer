package filehandler

import (
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/yureien/animeserver/database/models"
	"github.com/zorchenhimer/go-ed2k"
)

func (f FileHandler) ProcessFileAsync(filePath string) {
	f.processingQueue <- filePath
}

func (f FileHandler) processFileWorker() {
	for filePath := range f.processingQueue {
		err := f.processFile(filePath)
		if err != nil {
			f.logger.Error("Error processing file", "error", err)
		}
	}
}

func (f FileHandler) startProcessingQueue() {
	numWorkers := runtime.NumCPU()
	f.logger.Info("Starting file processing workers", "count", numWorkers)
	for range numWorkers {
		go f.processFileWorker()
	}
}

func (f FileHandler) processFile(filePath string) error {
	// Skip hidden files
	if strings.HasPrefix(filePath, ".") {
		return nil
	}

	fullPath := filepath.Join(f.cfg.AnimeDirectory, filePath)

	info, err := os.Stat(fullPath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	if dbFile, err := models.GetFileByPath(f.db, filePath); dbFile != nil || err != nil {
		if err != nil {
			return err
		}

		if dbFile != nil {
			f.logger.Info("Skipping file because it already exists", "path", filePath)
			return nil
		}
	}

	f.logger.Info("Processing file", "path", filePath)

	dbFile := models.File{
		Path: filePath,
		Size: int(info.Size()),
	}

	file, err := os.Open(fullPath)
	if err != nil {
		return err
	}

	ed2kHash := ed2k.New()
	_, err = io.Copy(ed2kHash, file)
	if err != nil {
		return err
	}

	ed2kHashBytes := ed2kHash.Sum(nil)
	dbFile.ED2K = hex.EncodeToString(ed2kHashBytes)

	// Query anihash endpoint to get file info
	aniFile, err := getAniHashFile(dbFile.ED2K, int64(dbFile.Size), f.cfg.AniHashURL)
	if err != nil {
		return err
	}

	f.logger.Debug("AniHash response", "response", aniFile)

	if aniFile != nil {
		dbFile.FileID = uint32(aniFile.FileID)
		dbFile.AnimeID = uint32(aniFile.AnimeID)
		dbFile.EpisodeID = uint32(aniFile.EpisodeID)
		dbFile.GroupID = uint32(aniFile.GroupID)
		dbFile.State = uint16(aniFile.State)
		dbFile.SHA1 = aniFile.SHA1
		dbFile.MD5 = aniFile.MD5
		dbFile.CRC = aniFile.CRC
		dbFile.Quality = aniFile.Quality
		dbFile.Source = aniFile.Source
		dbFile.AudioCodec = aniFile.AudioCodec
		dbFile.AudioBitrate = uint32(aniFile.AudioBitrate)
		dbFile.VideoCodec = aniFile.VideoCodec
		dbFile.VideoBitrate = uint32(aniFile.VideoBitrate)
		dbFile.VideoResolution = aniFile.VideoResolution
		dbFile.Extension = aniFile.Extension
		dbFile.Year = aniFile.Year
		dbFile.Type = aniFile.Type
		dbFile.RomajiName = aniFile.RomajiName
		dbFile.EnglishName = aniFile.EnglishName
		dbFile.EpNum = aniFile.EpNum
		dbFile.EpName = aniFile.EpName
		dbFile.EpRomajiName = aniFile.EpRomajiName
		dbFile.GroupName = aniFile.GroupName
	} else {
		f.logger.Error("AniHash response is nil", "ed2k", dbFile.ED2K, "size", dbFile.Size)
	}

	err = models.UpsertFileByPath(f.db, &dbFile)
	if err != nil {
		return err
	}

	return nil
}
