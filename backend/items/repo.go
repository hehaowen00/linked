package items

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
)

const getItemsByCollectionSql = `
SELECT id, url, title, description, created_at
FROM items
INNER JOIN item_collection_map
ON id = item_id AND items.user_id = item_collection_map.user_id
WHERE collection_id = ? AND items.user_id = $2
ORDER BY created_at DESC;
`

func getItemsByCollection(db *sql.DB, collectionId, userId string) ([]*Item, error) {
	var items []*Item

	rows, err := db.Query(getItemsByCollectionSql, collectionId, userId)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		item := Item{
			CollectionId: collectionId,
			UserId:       userId,
		}

		err := rows.Scan(&item.ID, &item.URL, &item.Title, &item.Description, &item.CreatedAt)
		if err != nil {
			return nil, err
		}

		items = append(items, &item)
	}

	return items, nil
}

const getItemsSql = `
SELECT id, title, url, description, created_at, deleted_at
FROM items
WHERE items.id NOT IN
	(SELECT item_id FROM item_collection_map WHERE user_id = ?)
ORDER BY created_at DESC;
`

func getItems(db *sql.DB, userId string) ([]*Item, error) {
	var items []*Item

	rows, err := db.Query(getItemsSql, userId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		item := Item{}
		err := rows.Scan(&item.ID, &item.Title, &item.URL, &item.Description, &item.CreatedAt, &item.DeletedAt)
		if err != nil {
			return nil, err
		}

		items = append(items, &item)
	}

	return items, nil
}

const getItemSql = `
SELECT url, title, description, created_at
FROM items
WHERE id = ? AND user_id = ? AND deleted_at = 0;
`

func getItem(db *sql.DB, item *Item) error {
	row := db.QueryRow(getItemSql, item.ID, item.UserId)
	err := row.Scan(&item.URL, &item.Title, &item.Description, &item.CreatedAt)
	return err
}

const getItemByUrlSql = `
SELECT id
FROM items
WHERE user_id = ? AND url = ?;
`

func getItemByUrl(db *sql.DB, item *Item) (bool, error) {
	row := db.QueryRow(getItemByUrlSql, item.UserId, item.URL)

	err := row.Scan(&item.ID)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return true, err
}

const createItemSql = `
INSERT INTO items (id, user_id, url, title, description, created_at)
VALUES (?, ?, ?, ?, ?, ?);
`

func createItem(db *sql.DB, item *Item) error {
	item.ID = uuid.NewString()
	item.CreatedAt = time.Now().UTC().UnixMilli()

	_, err := db.Exec(
		createItemSql,
		item.ID,
		item.UserId,
		item.URL,
		item.Title,
		item.Description,
		item.CreatedAt,
	)
	if err != nil {
		log.Println("user might have already added this item", err)
		return err
	}

	return err
}

const addItemCollectionSql = `
INSERT INTO item_collection_map (user_id, collection_id, item_id)
VALUES (?, ?, ?);
`

func addItemToCollection(db *sql.DB, item *Item) error {
	_, err := db.Exec(
		addItemCollectionSql,
		item.UserId,
		item.CollectionId,
		item.ID,
	)
	return err
}

const updateItemSql = `
UPDATE items
SET title = ?, description = ?
WHERE id = ? AND user_id = ? AND url = ? AND deleted_at = 0;
`

func updateItem(db *sql.DB, item *Item) error {
	_, err := db.Exec(
		updateItemSql,
		item.Title,
		item.Description,
		item.ID,
		item.UserId,
		item.URL,
	)
	return err
}

const deleteItemMappingSql = `
DELETE FROM item_collection_map
WHERE user_id = ? AND item_id = ? AND collection_id = ?;
`

func deleteItemMapping(db *sql.DB, item *Item) error {
	item.DeletedAt = time.Now().UTC().UnixMilli()
	res, err := db.Exec(deleteItemMappingSql, item.UserId, item.ID, item.CollectionId)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.New("item mapping not found")
	}

	return nil
}

const deleteItemSql = `
DELETE FROM items
WHERE user_id = ? AND id = ? AND url = ? AND created_at = ?;
`

const deleteAllItemMappingsSql = `
DELETE FROM item_collection_map
WHERE user_id = ? AND item_id = ?;
`

func deleteItem(db *sql.DB, item *Item) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(deleteAllItemMappingsSql, item.UserId, item.ID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(deleteItemSql, item.UserId, item.ID, item.URL, item.CreatedAt)
	if err != nil {
		return err
	}

	return tx.Commit()
}
