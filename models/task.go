package models

type Task struct {
	ID          int
	Title       string
	Description string
	Priority    int
	Deadline    string
	Done        bool
	CreatedAt   string
}
