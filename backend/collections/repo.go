package collections

import (
	"database/sql"
	"time"
)

type CollectionsRepo struct {
	db *sql.DB
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

const deleteItemsFromCollectionSql = `
DELETE FROM item_collection_map
WHERE collection_id = ?;
`

const deleteOrphanedItemsSql = `
DELETE FROM items
WHERE items.id NOT IN (SELECT item_id FROM item_collection_map);
`

func deleteCollection(db *sql.DB, c *Collection) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(deleteCollectionSql, c.Id, c.UserId, c.Name, c.DeletedAt)
	if err != nil {
		return err
	}

	_, err = tx.Exec(deleteItemsFromCollectionSql, c.Id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(deleteOrphanedItemsSql)
	if err != nil {
		return err
	}

	err = tx.Commit()

	return err
}