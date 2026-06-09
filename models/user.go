package models

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	CreatedAt    string `json:"created_at"`
}

type Session struct {
	Token     string `json:"-"`
	UserID    int    `json:"-"`
	ExpiresAt string `json:"-"`
}
