package handlers

import (
	"database/sql"
)

func getDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "../database/bddforum.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}
