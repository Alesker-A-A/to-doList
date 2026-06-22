package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"todo-calendar/auth"
	"todo-calendar/friends"
)

type sendRequestBody struct {
	Username string `json:"username"`
}

type friendActionBody struct {
	RequesterID int `json:"requester_id"`
}

func (h *Handler) SendFriendRequest(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}
	var req sendRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	target, err := auth.GetUserByUsername(h.DB, req.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if target == nil {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	friendship, err := friends.SendRequest(h.DB, userID, target.ID)
	switch {
	case errors.Is(err, friends.ErrCannotFriendSelf):
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	case errors.Is(err, friends.ErrAlreadyFriends),
		errors.Is(err, friends.ErrRequestExists):
		http.Error(w, err.Error(), http.StatusConflict)
		return
	case err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(friendship)
}

func (h *Handler) AcceptFriendRequest(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "не авторизован", http.StatusUnauthorized)
		return
	}

	var req friendActionBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := friends.AcceptRequest(h.DB, req.RequesterID, userID)
	if errors.Is(err, friends.ErrRequestNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeclineFriendRequest(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "не авторизован", http.StatusUnauthorized)
		return
	}

	var req friendActionBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := friends.DeclineRequest(h.DB, req.RequesterID, userID)
	if errors.Is(err, friends.ErrRequestNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) ListFriends(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "не авторизован", http.StatusUnauthorized)
		return
	}

	list, err := friends.ListFriends(h.DB, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *Handler) ListIncomingRequests(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "не авторизован", http.StatusUnauthorized)
		return
	}

	list, err := friends.ListIncomingRequests(h.DB, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}
