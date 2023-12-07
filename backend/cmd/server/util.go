package main

import (
	"database/sql"
	"fmt"
	"log"
)

func openSqliteDB(path string) (*sql.DB, error) {
	dbPath := fmt.Sprintf("file:%s?_journal_mode=WAL", path)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return db, nil
}
