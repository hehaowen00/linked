package migrations

import (
	"database/sql"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"
)

const createMigrationsTable = `
create table if not exists db_migrations (
	filename text not null,
	time int not null,
	primary key (filename)
);`

const checkIfMigrationExistsSql = `
select time
from db_migrations
where filename = ?;
`

func checkIfMigrationExists(db *sql.DB, filename string) (bool, error) {
	row := db.QueryRow(checkIfMigrationExistsSql, filename)

	var t int64

	err := row.Scan(&t)
	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

const addMigrationSql = `
insert into db_migrations (filename, time)
values (?, ?);
`

func addMigration(db *sql.DB, name string) error {
	_, err := db.Exec(addMigrationSql, name, time.Now().UTC().UnixMilli())
	return err
}

func executeMigration(db *sql.DB, filepath string) error {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(string(bytes))
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

func RunMigrations(db *sql.DB, path string) error {
	log.Println("running sql migrations")

	_, err := db.Exec(createMigrationsTable)
	if err != nil {
		return err
	}

	var migrations []string
	absPath, _ := filepath.Abs(path)
	log.Println("walking", absPath)

	filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(d.Name()) == ".sql" {
			migrations = append(migrations, path)
		}

		return nil
	})

	for _, migration := range migrations {
		name := filepath.Base(migration)
		exists, err := checkIfMigrationExists(db, name)
		if err != nil {
			return err
		}
		if !exists {
			log.Println("executing migration:", migration)
			err = executeMigration(db, migration)
			if err != nil {
				return err
			}

			err = addMigration(db, name)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
