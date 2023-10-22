package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

type JSON map[string]interface{}

func ReadJSON(r io.ReadCloser, v interface{}) error {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, v)
}

func WriteJSON(w http.ResponseWriter, statusCode int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(v)
}
