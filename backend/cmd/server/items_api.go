package main

import (
	"database/sql"
	"linked/internal/constants"
	"log"
	"net/http"

	pathrouter "github.com/hehaowen00/path-router"
)

func initItemsApi(db *sql.DB, scope pathrouter.IRoutes) {
	scope.Get("/collections/:collection/items",
		func(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
			userId := r.Context().Value(constants.AuthKey).(string)
			defer r.Body.Close()

			items, err := getItemsByCollection(db, ps.Get("collection"), userId)
			if err != nil {
				log.Println(err)
				writeJson(w, http.StatusInternalServerError, JsonResult{
					Status: "error",
					Error:  err.Error(),
				})
				return
			}

			writeJson(w, http.StatusOK, items)
		})

	scope.Get("/items",
		func(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
			userId := r.Context().Value(constants.AuthKey).(string)
			defer r.Body.Close()

			items, err := getItems(db, userId)
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

	scope.Get("/items/:item",
		func(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
			userId := r.Context().Value(constants.AuthKey).(string)
			defer r.Body.Close()

			item := Item{
				ID:     ps.Get("item"),
				UserId: userId,
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

	scope.Post("/items",
		func(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
			userId := r.Context().Value(constants.AuthKey).(string)
			defer r.Body.Close()

			item := Item{
				UserId: userId,
			}

			err := readJson(r.Body, &item)
			if err != nil {
				log.Println(err)
				writeJson(w, http.StatusBadRequest, JsonResult{
					Status: "error",
					Error:  err.Error(),
				})
				return
			}

			exists, err := getItemByUrl(db, &item)
			if err != nil {
				log.Println(err)
				writeJson(w, http.StatusBadRequest, JsonResult{
					Status: "error",
					Error:  err.Error(),
				})
				return
			}

			if !exists {
				err = createItem(db, &item)
				if err != nil {
					log.Println(err)
					writeJson(w, http.StatusInternalServerError, JsonResult{
						Status: "error",
						Error:  err.Error(),
					})
					return
				}
			} else {
				err = updateItem(db, &item)
				if err != nil {
					log.Println(err)
					writeJson(w, http.StatusInternalServerError, JsonResult{
						Status: "error",
						Error:  err.Error(),
					})
					return
				}

				err = addItemToCollection(db, &item)
				if err != nil {
					log.Println(err)
					writeJson(w, http.StatusInternalServerError, JsonResult{
						Status: "error",
						Error:  err.Error(),
					})
					return
				}
			}

			writeJson(w, http.StatusOK, JsonResult{
				Status: "ok",
				Data:   item,
			})
		})

	scope.Put("/items/:item",
		func(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
			userId := r.Context().Value(constants.AuthKey).(string)
			defer r.Body.Close()

			item := Item{
				ID:     ps.Get("item"),
				UserId: userId,
			}

			err := readJson(r.Body, &item)
			if err != nil {
				log.Println(err)
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

	scope.Delete("/items",
		func(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
		})

	scope.Delete("/collections/:collection/items/:item",
		func(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
			userId := r.Context().Value(constants.AuthKey).(string)
			defer r.Body.Close()

			item := Item{
				ID:           ps.Get("item"),
				CollectionId: ps.Get("collection"),
				UserId:       userId,
			}

			err := deleteItemMapping(db, &item)
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
