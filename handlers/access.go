package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"todo-calendar/auth"
	"todo-calendar/friends"
	"todo-calendar/models"
)

type grantAccessBody struct {
	GranteeID int                `json:"grantee_id"`
	Level     models.AccessLevel `json:"level"`
}

type accessActionBody struct {
	GranteeID int `json:"grantee_id"`
}

func (h *Handler) GrantAccess(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	var req grantAccessBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := friends.Grant(h.DB, userID, req.GranteeID, req.Level)
	switch {
	case errors.Is(err, friends.ErrInvalidLevel):
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	case errors.Is(err, friends.ErrNotFriends):
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	case err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) RevokeAccess(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "не авторизован", http.StatusUnauthorized)
		return
	}

	var req accessActionBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := friends.Revoke(h.DB, userID, req.GranteeID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) ListSharedWithMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "не авторизован", http.StatusUnauthorized)
		return
	}

	list, err := friends.ListGrantedToMe(h.DB, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *Handler) ListMyGrants(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "не авторизован", http.StatusUnauthorized)
		return
	}

	list, err := friends.ListMyGrants(h.DB, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *Handler) GetFriendTasks(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "не авторизован", http.StatusUnauthorized)
		return
	}

	ownerID, err := strconv.Atoi(r.PathValue("ownerID"))
	if err != nil {
		http.Error(w, "неверный id пользователя", http.StatusBadRequest)
		return
	}

	level, err := friends.GetLevel(h.DB, ownerID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if level == "" {
		http.Error(w, "нет доступа к этому календарю", http.StatusForbidden)
		return
	}

	rows, err := h.DB.Query(
		`SELECT id, title, description, priority, start_time, end_time, color, deadline, done, created_at
		   FROM tasks
		  WHERE user_id = ? AND is_private = 0`,
		ownerID,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tasks := []models.Task{}
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(
			&t.ID, &t.Title, &t.Description, &t.Priority,
			&t.StartTime, &t.EndTime, &t.Color, &t.Deadline,
			&t.Done, &t.CreatedAt,
		); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if level == models.AccessFreeBusy {
			t.Title = "Занято"
			t.Description = ""
		}

		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}
