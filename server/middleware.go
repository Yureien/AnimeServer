package server

import (
	"net/http"
	"strings"

	"github.com/yureien/animeserver/database/models"
	"github.com/yureien/animeserver/utils"
)

func (s server) RequireAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(s.cfg.Session.CookieName)
		if err != nil {
			http.Redirect(w, r, "/login?redirect="+r.URL.Path, http.StatusSeeOther)
			return
		}

		parts := strings.SplitN(cookie.Value, ":", 2)
		if len(parts) != 2 {
			http.Redirect(w, r, "/login?redirect="+r.URL.Path, http.StatusSeeOther)
			return
		}

		username, token := parts[0], parts[1]

		password, ok := s.cfg.Users[username]
		if !ok {
			http.Redirect(w, r, "/login?redirect="+r.URL.Path, http.StatusSeeOther)
			return
		}

		if !utils.ValidateSession(token, username, s.cfg.Session.Secret, password) {
			http.Redirect(w, r, "/login?redirect="+r.URL.Path, http.StatusSeeOther)
			return
		}

		handler(w, r)
	}
}

func (s server) RequireTokenAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			http.Error(w, "Token is required", http.StatusUnauthorized)
			return
		}

		valid, err := models.ValidateToken(s.db, token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		handler(w, r)
	}
}
