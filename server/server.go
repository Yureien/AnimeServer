package server

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/yureien/animeserver/filehandler"
	"goji.io"
	"goji.io/pat"
	"gorm.io/gorm"
)

type server struct {
	db          *gorm.DB
	cfg         *Config
	logger      *slog.Logger
	templates   *template.Template
	fileHandler *filehandler.FileHandler
}

func NewServer(
	logger *slog.Logger,
	db *gorm.DB,
	cfg *Config,
	templates *template.Template,
	fileHandler *filehandler.FileHandler,
) *server {
	return &server{
		db:          db,
		cfg:         cfg,
		logger:      logger,
		templates:   templates,
		fileHandler: fileHandler,
	}
}

func (s server) ListenAndServe() error {
	listenAddress := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)

	mux := goji.NewMux()

	// Static files
	mux.Handle(pat.Get("/static/*"), http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Auth routes (must come before catch-all routes)
	mux.HandleFunc(pat.Get("/login"), s.LoginPage)
	mux.HandleFunc(pat.Post("/login"), s.LoginHandler)
	mux.HandleFunc(pat.Get("/logout"), s.LogoutHandler)

	// Protected routes
	// mux.HandleFunc(pat.Get("/download-dir/*"), s.RequireAuth(handlers.DownloadDirectoryHandler()))
	mux.HandleFunc(pat.Get("/download/*"), s.RequireAuth(s.DownloadHandler))
	mux.HandleFunc(pat.Get("/generate-stream-token"), s.RequireAuth(s.GetStreamToken))
	mux.HandleFunc(pat.Get("/stream/*"), s.RequireTokenAuth(s.StreamHandler))
	mux.HandleFunc(pat.Get("/*"), s.RequireAuth(s.BrowseHandler))

	s.logger.Info("Starting server", "address", listenAddress)

	return http.ListenAndServe(listenAddress, mux)
}
