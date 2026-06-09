package db

import "database/sql"

func CreateTables(db *sql.DB) error {
	// Таблица пользователей.
	usersTable := `CREATE TABLE IF NOT EXISTS users (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    username      TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at    TEXT DEFAULT CURRENT_TIMESTAMP
)`
	if _, err := db.Exec(usersTable); err != nil {
		return err
	}

	// Таблица сессий.
	sessionsTable := `CREATE TABLE IF NOT EXISTS sessions (
    token      TEXT PRIMARY KEY,
    user_id    INTEGER NOT NULL,
    expires_at TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
)`
	if _, err := db.Exec(sessionsTable); err != nil {
		return err
	}

	// Таблица задач.
	tasksTable := `CREATE TABLE IF NOT EXISTS tasks (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id     INTEGER,
    title       TEXT NOT NULL,
    description TEXT,
    priority    INTEGER DEFAULT 1,
    deadline    TEXT,
    start_time  TEXT,
    end_time    TEXT,
    color       TEXT,
    done        INTEGER DEFAULT 0,
    created_at  TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
)`
	if _, err := db.Exec(tasksTable); err != nil {
		return err
	}

	if err := addUserIDToTasks(db); err != nil {
		return err
	}

	return nil
}

func addUserIDToTasks(db *sql.DB) error {
	rows, err := db.Query(`PRAGMA table_info(tasks)`)
	if err != nil {
		return err
	}
	defer rows.Close()

	hasUserID := false
	for rows.Next() {
		var (
			cid       int
			name      string
			colType   string
			notNull   int
			dfltValue sql.NullString
			pk        int
		)
		if err := rows.Scan(&cid, &name, &colType, &notNull, &dfltValue, &pk); err != nil {
			return err
		}
		if name == "user_id" {
			hasUserID = true
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}

	if !hasUserID {
		if _, err := db.Exec(`ALTER TABLE tasks ADD COLUMN user_id INTEGER`); err != nil {
			return err
		}
	}

	return nil
}
