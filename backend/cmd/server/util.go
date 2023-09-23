package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/goccy/go-json"
)

type JsonResult struct {
	Status string      `json:"status"` // can be "ok" or "error"
	Error  string      `json:"error,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

func readJson(r io.ReadCloser, v interface{}) error {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, v)
}

func writeJson(w http.ResponseWriter, statusCode int, v interface{}) error {
	w.WriteHeader(statusCode)
	w.Header().Add("Accept-Encoding", "application/json")
	enc := json.NewEncoder(w)
	return enc.Encode(v)
}

func openSqliteDB(path string) (*sql.DB, error) {
	dbPath := fmt.Sprintf("file:%s?_journal_mode=WAL", path)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return db, nil
}
