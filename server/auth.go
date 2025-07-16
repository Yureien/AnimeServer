package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/yureien/animeserver/database/models"
	"github.com/yureien/animeserver/server/types"
	"github.com/yureien/animeserver/utils"
	"golang.org/x/crypto/bcrypt"
)

func (s server) LoginPage(w http.ResponseWriter, r *http.Request) {
	data := types.LoginData{
		Title:    "Login - Anime Server",
		Redirect: r.URL.Query().Get("redirect"),
	}

	if error := r.URL.Query().Get("error"); error != "" {
		data.Error = "Invalid username or password"
	}

	err := s.templates.ExecuteTemplate(w, "login", data)
	if err != nil {
		s.logger.Error("Error executing template", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (s server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	redirect := r.URL.Query().Get("redirect")

	// Validate user
	hashedPassword, exists := s.cfg.Users[username]
	if !exists {
		s.logger.Error("User not found", "username", username)
		http.Redirect(w, r, "/login", http.StatusForbidden)
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		s.logger.Error("Error comparing password", "error", err)
		http.Redirect(w, r, "/login", http.StatusForbidden)
		return
	}

	s.logger.Info("User logged in", "username", username)

	// Create session
	token := utils.GenerateSessionToken(username, s.cfg.Session.Secret, hashedPassword)
	sessionValue := fmt.Sprintf("%s:%s", username, token)

	cookie := &http.Cookie{
		Name:     s.cfg.Session.CookieName,
		Value:    sessionValue,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
	}
	http.SetCookie(w, cookie)

	// Redirect
	if redirect != "" {
		http.Redirect(w, r, redirect, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (s server) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     s.cfg.Session.CookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(-1 * time.Hour),
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (s server) GetStreamToken(w http.ResponseWriter, r *http.Request) {
	token, err := models.GenerateToken(s.db)
	if err != nil {
		s.logger.Error("Error generating token", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}
