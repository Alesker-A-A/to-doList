package models

type Task struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Color       string `json:"color"`
	Deadline    string `json:"deadline"`
	Done        bool   `json:"done"`
	CreatedAt   string `json:"created_at"`
	IsPrivate   bool   `json:"is_private"`
}
