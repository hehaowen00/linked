package main

import (
	"database/sql"
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

const getItemsSql = `
SELECT id, url, title, description, created_at
FROM items
INNER JOIN item_collection_map
ON id = item_id AND items.user_id = item_collection_map.user_id
WHERE collection_id = ? AND items.user_id = $2
ORDER BY title ASC;
`

func getItems(db *sql.DB, collectionId, userId string) ([]*Item, error) {
	var items []*Item

	rows, err := db.Query(getItemsSql, collectionId, userId)

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

const getItemSql = `
SELECT title, url, description, created_at
FROM items
WHERE id = ? and collection_id = ? and user_id = ? and deleted_at = 0;
`

func getItem(db *sql.DB, item *Item) error {
	row := db.QueryRow(getItemSql, item.ID, item.CollectionId, item.UserId)
	err := row.Scan(&item.Title, &item.URL, &item.Description, &item.CreatedAt)
	return err
}

const createItemSql = `
INSERT INTO items (id, collection_id, user_id, title, url, description, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?);
`

func createItem(db *sql.DB, item *Item) error {
	item.ID = uuid.NewString()
	item.CreatedAt = time.Now().UTC().UnixMilli()
	_, err := db.Exec(
		createItemSql,
		item.ID,
		item.CollectionId,
		item.UserId,
		item.Title,
		item.URL,
		item.Description,
		item.CreatedAt,
	)
	return err
}

const updateItemSql = `
UPDATE items
SET title = ?, description = ?
WHERE id = ? and collection_id = ? and user_id = ? and deleted_at = 0;
`

func updateItem(db *sql.DB, item *Item) error {
	_, err := db.Exec(
		updateItemSql,
		item.Title,
		item.Description,
		item.ID,
		item.CollectionId,
		item.UserId,
	)
	return err
}

const deleteItemSql = `
DELETE items
SET deleted_at = ?
WHERE id = ? and collection_id = ? and user_id = ? and deleted_at = 0;
`

func deleteItem(db *sql.DB, item *Item) error {
	item.DeletedAt = time.Now().UTC().UnixMilli()
	_, err := db.Exec(deleteItemSql, item.DeletedAt, item.ID, item.CollectionId, item.UserId)
	return err
}
