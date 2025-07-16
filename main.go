package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/yureien/animeserver/database"
	"github.com/yureien/animeserver/filehandler"
	"github.com/yureien/animeserver/server"
	"github.com/yureien/animeserver/templates"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if len(os.Args) > 1 && os.Args[1] == "create-user" {
		if len(os.Args) != 3 {
			fmt.Println("Usage: animeserver create-user <username>")
			os.Exit(1)
		}
		createUser(os.Args[2], logger)
		return
	}

	cfg, err := LoadConfig("config.yaml")
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	db, err := database.LoadDatabase(logger, &cfg.Database)
	if err != nil {
		logger.Error("Failed to load database", "error", err)
		os.Exit(1)
	}

	fileHandler := filehandler.NewFileHandler(logger, &cfg.FileHandler, db)

	templates, err := templates.LoadTemplates()
	if err != nil {
		logger.Error("Failed to load templates", "error", err)
		os.Exit(1)
	}

	server := server.NewServer(logger, db, &cfg.Server, templates, fileHandler)
	server.ListenAndServe()
}

func createUser(username string, logger *slog.Logger) {
	fmt.Print("Enter password: ")

	// Hide password input
	password, err := readPassword()
	if err != nil {
		logger.Error("Failed to read password", "error", err)
		os.Exit(1)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Failed to hash password", "error", err)
		os.Exit(1)
	}

	hashedPasswordStr := string(hashedPassword)
	fmt.Printf("Hashed password: %s\n", hashedPasswordStr)
}

func readPassword() (string, error) {
	// This is a simple password reader - in production you'd want to use a proper terminal library
	var password string
	fmt.Scanln(&password)
	return password, nil
}
