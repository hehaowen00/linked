package main

import (
	"database/sql"
	"linked/internal/constants"
	"log"
	"net/http"

	"github.com/google/uuid"
	pathrouter "github.com/hehaowen00/path-router"
)

func initCollectionsApi(db *sql.DB, scope pathrouter.IRoutes) {
	scope.Get("/collections",
		func(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
			userId := r.Context().Value(constants.AuthKey).(string)
			defer r.Body.Close()

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

	scope.Get("/collections/:collection",
		func(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
			userId := r.Context().Value(constants.AuthKey).(string)
			defer r.Body.Close()

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

	scope.Post("/collections",
		func(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
			userId := r.Context().Value(constants.AuthKey).(string)
			defer r.Body.Close()

			c := Collection{
				Id:     uuid.NewString(),
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

	scope.Put("/collections/:collection",
		func(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
			userId := r.Context().Value(constants.AuthKey).(string)
			defer r.Body.Close()

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

	scope.Delete("/collections/:collection",
		func(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
			userId := r.Context().Value(constants.AuthKey).(string)
			defer r.Body.Close()

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

			if c.DeletedAt == 0 {
				log.Println("archive collection")
				err = archiveCollection(db, &c)
			} else {
				log.Println("delete collection")
				err = deleteCollection(db, &c)
			}

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
