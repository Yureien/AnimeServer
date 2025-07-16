package models

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type File struct {
	gorm.Model

	Path string `gorm:"unique;uniqueIndex:idx_path"`
	Size int
	ED2K string `gorm:"uniqueIndex:idx_ed2k"`

	// AniDB Fields
	FileID          uint32 `gorm:"uniqueIndex:idx_file_id"`
	AnimeID         uint32
	EpisodeID       uint32
	GroupID         uint32
	State           uint16
	SHA1            string
	MD5             string
	CRC             string
	Quality         string
	Source          string
	AudioCodec      string
	AudioBitrate    uint32
	VideoCodec      string
	VideoBitrate    uint32
	VideoResolution string
	Extension       string
	Year            string
	Type            string
	RomajiName      string
	EnglishName     string
	EpNum           string
	EpName          string
	EpRomajiName    string
	GroupName       string
}

func GetFileByPath(db *gorm.DB, path string) (*File, error) {
	var file File
	err := db.Where("path = ?", path).First(&file).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &file, nil
}

func UpsertFileByPath(db *gorm.DB, file *File) error {
	err := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "path"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"path",
			"size",
			"ed2_k",
			"file_id",
			"anime_id",
			"episode_id",
			"group_id",
			"state",
			"sha1",
			"md5",
			"crc",
			"quality",
			"source",
			"audio_codec",
			"audio_bitrate",
			"video_codec",
			"video_bitrate",
			"video_resolution",
			"extension",
			"year",
			"type",
			"romaji_name",
			"english_name",
			"ep_num",
			"ep_name",
			"ep_romaji_name",
			"group_name",
			"updated_at",
		}),
	}).Create(file).Error
	if err != nil {
		return err
	}

	return nil
}
