package db

import "database/sql"

func CreateTables(db *sql.DB) error {

	usersTable := `CREATE TABLE IF NOT EXISTS users (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    username      TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at    TEXT DEFAULT CURRENT_TIMESTAMP
)`
	if _, err := db.Exec(usersTable); err != nil {
		return err
	}

	sessionsTable := `CREATE TABLE IF NOT EXISTS sessions (
    token      TEXT PRIMARY KEY,
    user_id    INTEGER NOT NULL,
    expires_at TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
)`
	if _, err := db.Exec(sessionsTable); err != nil {
		return err
	}

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
    is_private  INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
)`
	if _, err := db.Exec(tasksTable); err != nil {
		return err
	}

	friendships := `CREATE TABLE IF NOT EXISTS friendships (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    requester_id  INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    addressee_id  INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status        TEXT NOT NULL DEFAULT 'pending'
                      CHECK (status IN ('pending', 'accepted', 'declined')),
    created_at    TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (requester_id, addressee_id),
    CHECK (requester_id <> addressee_id)
);
CREATE INDEX IF NOT EXISTS idx_friendships_addressee
    ON friendships(addressee_id)`
	if _, err := db.Exec(friendships); err != nil {
		return err
	}

	calendarAccess := `CREATE TABLE IF NOT EXISTS calendar_access (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    owner_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    grantee_id  INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    level       TEXT NOT NULL CHECK (level IN ('free_busy', 'view')),
    created_at  TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (owner_id, grantee_id),
    CHECK (owner_id <> grantee_id)
);
CREATE INDEX IF NOT EXISTS idx_calendar_access_grantee
    ON calendar_access(grantee_id)`
	if _, err := db.Exec(calendarAccess); err != nil {
		return err
	}

	taskCollaborators := `CREATE TABLE IF NOT EXISTS task_collaborators (
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    task_id   INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    user_id   INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role      TEXT NOT NULL DEFAULT 'viewer'
                  CHECK (role IN ('editor', 'viewer')),
    status    TEXT NOT NULL DEFAULT 'pending'
                  CHECK (status IN ('pending', 'accepted')),
    added_at  TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (task_id, user_id)
);
CREATE INDEX IF NOT EXISTS idx_task_collaborators_user
    ON task_collaborators(user_id, status)`
	if _, err := db.Exec(taskCollaborators); err != nil {
		return err
	}

	if err := addUserIDToTasks(db); err != nil {
		return err
	}
	if err := addIsPrivateToTasks(db); err != nil {
		return err
	}

	return nil
}

func columnExists(db *sql.DB, table, column string) (bool, error) {
	rows, err := db.Query(`PRAGMA table_info(` + table + `)`)
	if err != nil {
		return false, err
	}
	defer rows.Close()

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
			return false, err
		}
		if name == column {
			return true, nil
		}
	}
	if err := rows.Err(); err != nil {
		return false, err
	}

	return false, nil
}

func addUserIDToTasks(db *sql.DB) error {
	exists, err := columnExists(db, "tasks", "user_id")
	if err != nil {
		return err
	}
	if !exists {
		if _, err := db.Exec(`ALTER TABLE tasks ADD COLUMN user_id INTEGER`); err != nil {
			return err
		}
	}
	return nil
}

func addIsPrivateToTasks(db *sql.DB) error {
	exists, err := columnExists(db, "tasks", "is_private")
	if err != nil {
		return err
	}
	if !exists {
		if _, err := db.Exec(`ALTER TABLE tasks ADD COLUMN is_private INTEGER NOT NULL DEFAULT 0`); err != nil {
			return err
		}
	}
	return nil
}
