package models

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
	Deadline    string `json:"deadline"`
	Done        bool   `json:"done"`
	CreatedAt   string `json:"created_at"`
}
