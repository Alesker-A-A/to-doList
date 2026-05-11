package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func Connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "todo.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}
