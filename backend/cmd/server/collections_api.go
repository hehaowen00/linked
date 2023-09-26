package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/google/uuid"
	pathrouter "github.com/hehaowen00/path-router"
)

func initCollectionsApi(db *sql.DB, router *pathrouter.Group) {
	router.Get("/collections",
		func(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
			userId := r.Context().Value("id").(string)

			collections, err := getCollections(db, userId)
			if err != nil {
				log.Println(err)
				writeJson(w, http.StatusInternalServerError, JsonResult{
					Status: "error",
					Error:  "unable to get collections",
				})
				return
			}

			writeJson(w, http.StatusOK, JsonResult{
				Status: "ok",
				Data:   collections,
			})
		})

	router.Get("/collections/:collection",
		func(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
			userId := r.Context().Value("id").(string)

			c := Collection{
				Id:     ps.Get("collection"),
				UserId: userId,
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
			userId := r.Context().Value("id").(string)

			c := Collection{
				Id:     uuid.NewString(),
				UserId: userId,
			}

			if err := c.isValid(); err != nil {
				writeJson(w, http.StatusBadRequest, JsonResult{
					Status: "error",
					Error:  err.Error(),
				})
				return
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
			userId := r.Context().Value("id").(string)

			c := Collection{
				UserId: userId,
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

			if err := c.isValid(); err != nil {
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
			userId := r.Context().Value("id").(string)

			c := Collection{
				Id:     ps.Get("collection"),
				UserId: userId,
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

			if err := c.isValid(); err != nil {
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
