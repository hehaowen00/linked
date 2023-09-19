package main

import (
	"database/sql"
	"log"
	"net/http"

	pathrouter "github.com/hehaowen00/path-router"
)

func initItemsApi(db *sql.DB, router *pathrouter.Group) {
	router.Get("/collections/:collection/items",
		func(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
			items, err := getItems(db, ps.Get("collection"), "")
			if err != nil {
				log.Println(err)
				writeJson(w, http.StatusInternalServerError, JsonResult{
					Status: "error",
					Error:  err.Error(),
				})
				return
			}

			writeJson(w, http.StatusOK, JsonResult{
				Status: "ok",
				Data:   items,
			})
		})

	router.Get("/collections/:collection/items/:item",
		func(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
			item := Item{
				ID:           ps.Get("item"),
				CollectionId: ps.Get("collection"),
			}

			err := getItem(db, &item)
			if err != nil {
				log.Println(err)
				writeJson(w, http.StatusInternalServerError, JsonResult{
					Status: "error",
					Error:  err.Error(),
				})
				return
			}

			writeJson(w, http.StatusOK, JsonResult{
				Status: "ok",
				Data:   item,
			})
		})

	router.Post("/collections/:collection/items",
		func(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
			item := Item{
				CollectionId: ps.Get("collection"),
			}

			err := readJson(r.Body, &item)
			if err != nil {
				writeJson(w, http.StatusBadRequest, JsonResult{
					Status: "error",
					Error:  err.Error(),
				})
				return
			}

			err = createItem(db, &item)
			if err != nil {
				log.Println(err)
				writeJson(w, http.StatusInternalServerError, JsonResult{
					Status: "error",
					Error:  err.Error(),
				})
				return
			}

			writeJson(w, http.StatusOK, JsonResult{
				Status: "ok",
				Data:   item,
			})
		})

	router.Put("/collections/:collection/items/:item",
		func(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
			item := Item{
				ID:           ps.Get("item"),
				CollectionId: ps.Get("collection"),
			}

			err := readJson(r.Body, &item)
			if err != nil {
				writeJson(w, http.StatusBadRequest, JsonResult{
					Status: "error",
					Error:  err.Error(),
				})
				return
			}

			err = updateItem(db, &item)
			if err != nil {
				log.Println(err)
				writeJson(w, http.StatusInternalServerError, JsonResult{
					Status: "error",
					Error:  err.Error(),
				})
				return
			}

			writeJson(w, http.StatusOK, JsonResult{
				Status: "ok",
				Data:   item,
			})
		})

	router.Delete("/collections/:collection/items/:item",
		func(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
			item := Item{
				ID:           ps.Get("item"),
				CollectionId: ps.Get("collection"),
			}

			err := deleteItem(db, &item)
			if err != nil {
				log.Println(err)
				writeJson(w, http.StatusInternalServerError, JsonResult{
					Status: "error",
					Error:  err.Error(),
				})
				return
			}

			writeJson(w, http.StatusOK, JsonResult{
				Status: "ok",
				Data:   item,
			})
		})
}
