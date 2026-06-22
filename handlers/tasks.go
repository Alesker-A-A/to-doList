package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"todo-calendar/auth"
	"todo-calendar/models"
)

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Color       string `json:"color"`
	Deadline    string `json:"deadline"`
	IsPrivate   bool   `json:"is_private"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Color       string `json:"color"`
	Deadline    string `json:"deadline"`
	Done        bool   `json:"done"`
	IsPrivate   bool   `json:"is_private"`
}

type Handler struct {
	DB *sql.DB
}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "не авторизирован", http.StatusUnauthorized)
		return
	}
	rows, err := h.DB.Query(`SELECT id, title, description, priority, start_time, end_time, color, deadline, done, is_private, created_at FROM tasks WHERE user_id = ?`, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Priority, &t.StartTime, &t.EndTime, &t.Color, &t.Deadline, &t.Done, &t.IsPrivate, &t.CreatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, t)
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(tasks)
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "не авторизирован", http.StatusUnauthorized)
		return
	}

	var req CreateTaskRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = h.DB.Exec(`INSERT INTO tasks (user_id, title, description, priority, start_time, end_time, color, deadline, is_private) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		userID, req.Title, req.Description, req.Priority, req.StartTime, req.EndTime, req.Color, req.Deadline, req.IsPrivate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "не авторизован", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")
	taskID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "неверный номер задачи", http.StatusBadRequest)
		return
	}

	var req UpdateTaskRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = h.DB.Exec(`UPDATE tasks SET title=?, description=?, priority=?, start_time=?, end_time=?, color=?, deadline=?, done=?, is_private=? WHERE id=? AND user_id=?`,
		req.Title, req.Description, req.Priority, req.StartTime, req.EndTime, req.Color, req.Deadline, req.Done, req.IsPrivate, taskID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "не авторизован", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")
	taskID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "неверный номер задачи", http.StatusBadRequest)
		return
	}
	_, err = h.DB.Exec(`DELETE FROM tasks WHERE id =? AND user_id=?`, taskID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
