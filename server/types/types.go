package types

import "github.com/yureien/animeserver/filehandler"

type PageData struct {
	Title       string
	Files       []filehandler.FileInfo
	CurrentPath string
	Error       string
	Anime       *AnimeData
	LeafPage    bool
}

type LoginData struct {
	Title    string
	Error    string
	Redirect string
}

type AnimeData struct {
	AnimeID         uint32
	GroupID         uint32
	Year            string
	RomajiName      string
	EnglishName     string
	GroupName       string
	Type            string
	Source          string
	Quality         string
	AudioCodec      string
	VideoCodec      string
	VideoResolution string
}
