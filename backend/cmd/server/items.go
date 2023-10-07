package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

type Item struct {
	ID           string `json:"id"`
	CollectionId string `json:"collection_id"`
	UserId       string `json:"-"`
	URL          string `json:"url"`
	Title        string `json:"title"`
	Description  string `json:"desc"`
	CreatedAt    int64  `json:"created_at"`
	DeletedAt    int64  `json:"deleted"`
}

const getItemsByCollectionSql = `
SELECT id, url, title, description, created_at
FROM items
INNER JOIN item_collection_map
ON id = item_id AND items.user_id = item_collection_map.user_id
WHERE collection_id = ? AND items.user_id = $2
ORDER BY title ASC;
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

const getItemCountSql = `
SELECT count(id)
FROM items
WHERE user_id = ?;
`

func getItemsCount(db *sql.DB, userId string) (int, error) {
	var count int

	row := db.QueryRow(getItemCountSql, userId)
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

const getItemsSql = `
SELECT id, title, url, description, created_at, deleted_at
FROM items
WHERE user_id = ?
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
WHERE id = ? and user_id = ? and deleted_at = 0;
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
WHERE id = ? and user_id = ? and url = ? and deleted_at = 0;
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
WHERE user_id = ? and item_id = ? and collection_id = ?;
`

func deleteItemMapping(db *sql.DB, item *Item) error {
	item.DeletedAt = time.Now().UTC().UnixMilli()
	_, err := db.Exec(deleteItemMappingSql, item.UserId, item.ID, item.CollectionId)
	return err
}
