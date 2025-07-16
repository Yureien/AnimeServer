package server

import (
	"fmt"
	"net/http"
	"strings"
)

func (s *server) serveFile(w http.ResponseWriter, r *http.Request, asAttachment bool) {
	var requestPath string
	if asAttachment {
		requestPath = strings.TrimPrefix(r.URL.Path, "/download/")
	} else {
		requestPath = strings.TrimPrefix(r.URL.Path, "/stream/")
	}

	fileInfo, file, err := s.fileHandler.GetFile(requestPath)
	if err != nil {
		s.logger.Error("Error getting file", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	if asAttachment {
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileInfo.Name()))
	} else {
		w.Header().Set("Content-Type", "video/mp4")
	}

	http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
}

func (s *server) DownloadHandler(w http.ResponseWriter, r *http.Request) {
	s.serveFile(w, r, true)
}

func (s *server) StreamHandler(w http.ResponseWriter, r *http.Request) {
	s.serveFile(w, r, false)
}
