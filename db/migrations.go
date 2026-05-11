package db

import "database/sql"

func CreateTables(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS tasks (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    title       TEXT NOT NULL,
    description TEXT,
    priority    INTEGER DEFAULT 1,
    deadline    TEXT,
    done        INTEGER DEFAULT 0,
    created_at  TEXT DEFAULT CURRENT_TIMESTAMP
)`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
