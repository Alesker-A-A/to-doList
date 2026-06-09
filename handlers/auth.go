package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	"todo-calendar/auth"
)

const sessionCookieName = "session_token"

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func setSessionCookie(w http.ResponseWriter, token string, expires time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(req.Username) < 3 {
		http.Error(w, "имя пользователя должно быть не короче 3 символов", http.StatusBadRequest)
		return
	}

	if len(req.Password) < 6 {
		http.Error(w, "пароль должен быть не короче 6 символов", http.StatusBadRequest)
	}

	user, err := auth.CreateUser(h.DB, req.Username, req.Password)
	if err == auth.ErrUserExists {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, expires, err := auth.CreateSession(h.DB, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	setSessionCookie(w, token, expires)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := auth.GetUserByUsername(h.DB, req.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil || !auth.CheckPassword(req.Password, user.PasswordHash) {
		http.Error(w, "неверное имя пользователя или пароль", http.StatusUnauthorized)
		return
	}

	token, expires, err := auth.CreateSession(h.DB, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	setSessionCookie(w, token, expires)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err == nil {
		_ = auth.DeleteSession(h.DB, cookie.Value)
	}
	clearSessionCookie(w)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		http.Error(w, "не авторизирован", http.StatusUnauthorized)
		return
	}

	user, err := auth.GetUserByToken(h.DB, cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "не авторизирован", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
