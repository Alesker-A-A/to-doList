package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func Connect() (*sql.DB, error) {

	db, err := sql.Open("sqlite", "todo.db?_pragma=foreign_keys(1)")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
