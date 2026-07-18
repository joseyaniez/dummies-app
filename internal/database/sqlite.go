package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func OpenSQLiteDatabase(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	return db, nil
}
