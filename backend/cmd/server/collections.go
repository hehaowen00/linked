package main

import (
	"database/sql"
	"errors"
	"time"
)

type Collection struct {
	Id        string `json:"id"`
	UserId    string `json:"-"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
	DeletedAt int64  `json:"deleted_at"`
}

func (c *Collection) isValid() error {
	if c.Id == "" {
		return errors.New("missing collection id")
	}
	if c.UserId == "" {
		return errors.New("missing user id")
	}
	if c.Name == "" {
		return errors.New("missing name")
	}
	return nil
}

const getCollectionsSql = `
SELECT id, name, created_at, deleted_at
FROM collections
WHERE user_id = ?
ORDER BY name ASC;
`

func getCollections(db *sql.DB, userId string) ([]*Collection, error) {
	rows, err := db.Query(getCollectionsSql, userId)
	if err != nil {
		return nil, err
	}

	var collections []*Collection

	for rows.Next() {
		c := Collection{
			UserId: userId,
		}

		err := rows.Scan(&c.Id, &c.Name, &c.CreatedAt, &c.DeletedAt)
		if err != nil {
			return nil, err
		}

		collections = append(collections, &c)
	}

	return collections, nil
}

const getCollectionSql = `
SELECT name, created_at, deleted_at
FROM collections
WHERE id = ? AND user_id = ?;
`

func getCollection(db *sql.DB, c *Collection) error {
	row := db.QueryRow(getCollectionSql, c.Id, c.UserId)
	err := row.Scan(&c.Name, &c.CreatedAt, &c.DeletedAt)
	return err
}

const createCollectionSql = `
INSERT INTO collections (id, user_id, name, created_at)
VALUES (?, ?, ?, ?);
`

func createCollection(db *sql.DB, c *Collection) error {
	c.CreatedAt = time.Now().UTC().UnixMilli()
	_, err := db.Exec(createCollectionSql, c.Id, c.UserId, c.Name, c.CreatedAt)
	return err
}

const updateCollectionSql = `
UPDATE collections
SET name = ?, deleted_at = ?
WHERE id = ? and user_id = ?;
`

func updateCollection(db *sql.DB, c *Collection) error {
	_, err := db.Exec(updateCollectionSql, c.Name, c.DeletedAt, c.Id, c.UserId)
	return err
}

const archiveCollectionSql = `
UPDATE collections
SET deleted_at = ?
WHERE id = ? and user_id = ? and name = ? and deleted_at = 0;
`

func archiveCollection(db *sql.DB, c *Collection) error {
	c.DeletedAt = time.Now().UTC().UnixMilli()
	_, err := db.Exec(archiveCollectionSql, c.DeletedAt, c.Id, c.UserId, c.Name)
	return err
}

const deleteCollectionSql = `
DELETE FROM collections
WHERE id = ? and user_id = ? and name = ? and deleted_at = ?;
`

func deleteCollection(db *sql.DB, c *Collection) error {
	_, err := db.Exec(deleteCollectionSql, c.Id, c.UserId, c.Name, c.DeletedAt)
	return err
}
