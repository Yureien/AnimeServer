package server

import (
	"net/http"
	"strings"

	"github.com/yureien/animeserver/server/types"
)

func (s server) BrowseHandler(w http.ResponseWriter, r *http.Request) {
	requestPath := r.URL.Path
	if requestPath == "/" {
		requestPath = ""
	} else {
		requestPath = strings.TrimPrefix(requestPath, "/")
	}

	files, err := s.fileHandler.GetDirectory(requestPath)
	if err != nil {
		s.logger.Error("Error getting files", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	isLeafPage := true
	for _, file := range files {
		if file.IsDirectory {
			isLeafPage = false
			break
		}
	}

	// Fetch anime data from the first file
	var animeData *types.AnimeData
	if isLeafPage && len(files) > 0 && files[0].DbFile != nil {
		animeData = &types.AnimeData{
			AnimeID:         files[0].DbFile.AnimeID,
			GroupID:         files[0].DbFile.GroupID,
			Year:            files[0].DbFile.Year,
			RomajiName:      files[0].DbFile.RomajiName,
			EnglishName:     files[0].DbFile.EnglishName,
			GroupName:       files[0].DbFile.GroupName,
			Type:            files[0].DbFile.Type,
			Source:          files[0].DbFile.Source,
			Quality:         files[0].DbFile.Quality,
			AudioCodec:      files[0].DbFile.AudioCodec,
			VideoCodec:      files[0].DbFile.VideoCodec,
			VideoResolution: files[0].DbFile.VideoResolution,
		}
	}

	data := types.PageData{
		Title:       "Anime Server - " + requestPath,
		Files:       files,
		CurrentPath: requestPath,
		Anime:       animeData,
		LeafPage:    isLeafPage,
	}

	err = s.templates.ExecuteTemplate(w, "browse", data)
	if err != nil {
		s.logger.Error("Error executing template", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
