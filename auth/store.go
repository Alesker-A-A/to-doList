package auth

import (
	"database/sql"
	"errors"
	"time"

	"todo-calendar/models"
)

const sessionDuration = 7 * 24 * time.Hour

var ErrUserExists = errors.New("Пользователь с таким именем уже существует")

func CreateUser(db *sql.DB, username, password string) (*models.User, error) {
	hash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	res, err := db.Exec(
		`INSERT INTO users (username, password_hash) VALUES (?,?)`,
		username, hash,
	)
	if err != nil {
		return nil, ErrUserExists
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:           int(id),
		Username:     username,
		PasswordHash: hash,
	}, nil
}

func GetUserByUsername(db *sql.DB, username string) (*models.User, error) {
	var u models.User
	err := db.QueryRow(
		`SELECT id, username, password_hash, created_at FROM users WHERE username = ?`,
		username,
	).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func CreateSession(db *sql.DB, userID int) (string, time.Time, error) {
	token, err := GenerateToken()
	if err != nil {
		return "", time.Time{}, err
	}

	expiresAt := time.Now().Add(sessionDuration)

	_, err = db.Exec(
		`INSERT INTO sessions (token, user_id, expires_at) VALUES (?, ?, ?)`,
		token, userID, expiresAt.Format(time.RFC3339),
	)
	if err != nil {
		return "", time.Time{}, err
	}

	return token, expiresAt, nil
}

func GetUserByToken(db *sql.DB, token string) (*models.User, error) {
	var (
		userID    int
		expiresAt string
	)
	err := db.QueryRow(
		`SELECT user_id, expires_at FROM sessions WHERE token = ?`,
		token,
	).Scan(&userID, &expiresAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	exp, err := time.Parse(time.RFC3339, expiresAt)
	if err != nil {
		return nil, err
	}
	if time.Now().After(exp) {
		_, _ = db.Exec(`DELETE FROM sessions WHERE token = ?`, token)
		return nil, nil
	}

	var u models.User
	err = db.QueryRow(
		`SELECT id, username, password_hash, created_at FROM users WHERE id = ?`,
		userID,
	).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func DeleteSession(db *sql.DB, token string) error {
	_, err := db.Exec(`DELETE FROM sessions WHERE token = ?`, token)
	return err
}
