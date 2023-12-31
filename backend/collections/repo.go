package collections

import (
	"database/sql"
	"time"
)

type CollectionsRepo struct {
	db *sql.DB
}

const getCollectionsSql = `
SELECT id, name, created_at, archived, archived_at
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

		err := rows.Scan(&c.Id, &c.Name, &c.CreatedAt, &c.Archived, &c.ArchivedAt)
		if err != nil {
			return nil, err
		}

		collections = append(collections, &c)
	}

	return collections, nil
}

const getCollectionSql = `
SELECT name, created_at, archived, archived_at
FROM collections
WHERE id = ? AND user_id = ?;
`

func getCollection(db *sql.DB, c *Collection) error {
	row := db.QueryRow(getCollectionSql, c.Id, c.UserId)
	err := row.Scan(&c.Name, &c.CreatedAt, &c.Archived, &c.ArchivedAt)
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
SET name = ?, archived = ?, archived_at = ?
WHERE id = ? and user_id = ?;
`

func updateCollection(db *sql.DB, c *Collection) error {
	_, err := db.Exec(updateCollectionSql, c.Name, c.Archived, c.ArchivedAt, c.Id, c.UserId)
	return err
}

const archiveCollectionSql = `
UPDATE collections
SET archived = 1, archived_at = ?
WHERE id = ? and user_id = ? and name = ?;
`

func archiveCollection(db *sql.DB, c *Collection) error {
	c.Archived = true
	c.ArchivedAt = time.Now().UTC().UnixMilli()
	_, err := db.Exec(archiveCollectionSql, c.ArchivedAt, c.Id, c.UserId, c.Name)
	return err
}

const unarchiveCollectionSql = `
UPDATE collections
SET archived = 0, archived_at = 0
WHERE id = ? and user_id = ? and name = ?;
`

func unarchiveCollection(db *sql.DB, c *Collection) error {
	c.Archived = false
	c.ArchivedAt = 0
	_, err := db.Exec(unarchiveCollectionSql, c.Id, c.UserId, c.Name)
	return err
}

const deleteCollectionSql = `
DELETE FROM collections
WHERE id = ? and user_id = ? and name = ? and archived = 1 and archived_at = ?;
`

const deleteCollectionItemsSql = `
DELETE FROM items
WHERE items.id IN (SELECT item_id FROM item_collection_map where collection_id = ?);
`

const deleteItemMappings = `
DELETE FROM item_collection_map
WHERE collection_id = ?;
`

func deleteCollection(db *sql.DB, c *Collection) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(deleteCollectionSql, c.Id, c.UserId, c.Name, c.ArchivedAt)
	if err != nil {
		return err
	}

	_, err = tx.Exec(deleteCollectionItemsSql, c.Id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(deleteItemMappings, c.Id)
	if err != nil {
		return err
	}

	err = tx.Commit()

	return err
}
