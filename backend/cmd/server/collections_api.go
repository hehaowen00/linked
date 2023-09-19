package main

import (
	"database/sql"
	"log"
	"net/http"

	pathrouter "github.com/hehaowen00/path-router"
)

func initCollectionsApi(db *sql.DB, router *pathrouter.Group) {
	router.Get("/collections",
		func(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
			collections, err := getCollections(db, "")
			if err != nil {
				log.Println(err)
				writeJson(w, http.StatusInternalServerError, JsonResult{
					Status: "error",
					Error:  "unable to get collections",
				})
			}

			writeJson(w, http.StatusOK, JsonResult{
				Status: "ok",
				Data:   collections,
			})
		})

	router.Get("/collections/:collection",
		func(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
			c := Collection{
				Id: ps.Get("collection"),
			}

			err := getCollection(db, &c)
			if err != nil {
				log.Println(err)
				writeJson(w, http.StatusOK, JsonResult{
					Status: "error",
					Error:  "unable to fetch collection",
				})
				return
			}

			writeJson(w, http.StatusOK, JsonResult{
				Status: "ok",
				Data:   c,
			})
		})

	router.Post("/collections",
		func(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
			c := Collection{}

			err := readJson(r.Body, &c)
			if err != nil {
				log.Println(err)
				writeJson(w, http.StatusBadRequest, JsonResult{
					Status: "error",
					Error:  err.Error(),
				})
				return
			}

			err = createCollection(db, &c)
			if err != nil {
				log.Println(err)
				writeJson(w, http.StatusInternalServerError, JsonResult{
					Status: "error",
					Error:  "failed to create collection",
				})
				return
			}

			writeJson(w, http.StatusOK, JsonResult{
				Status: "ok",
				Data:   c,
			})
		})

	router.Put("/collections/:collection",
		func(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
			c := Collection{}

			err := readJson(r.Body, &c)
			if err != nil {
				log.Println(err)
				writeJson(w, http.StatusBadRequest, JsonResult{
					Status: "error",
					Error:  err.Error(),
				})
				return
			}

			err = updateCollection(db, &c)
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
			})
		})

	router.Delete("/collections/:collection",
		func(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
			c := Collection{
				Id: ps.Get("collection"),
			}

			err := readJson(r.Body, &c)
			if err != nil {
				log.Println(err)
				writeJson(w, http.StatusBadRequest, JsonResult{
					Status: "error",
					Error:  err.Error(),
				})
				return
			}

			err = deleteCollection(db, &c)
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
			})
		})
}
