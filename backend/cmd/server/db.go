package main

import (
	"database/sql"
	"linked/internal/migrations"
)

func loadDB(dbPath, migrationPath string) (*sql.DB, error) {
	db, err := openSqliteDB(dbPath)
	if err != nil {
		return nil, err
	}

	err = migrations.RunMigrations(db, migrationPath)
	if err != nil {
		return nil, err
	}

	return db, nil
}
