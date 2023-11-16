package collections

import (
	"database/sql"
	"linked/internal/constants"
	"linked/utils"
	"log"
	"net/http"

	"github.com/google/uuid"
	pathrouter "github.com/hehaowen00/path-router"
)

type CollectionAPI struct {
	db *sql.DB
}

func NewAPI(db *sql.DB) *CollectionAPI {
	collectionApi := CollectionAPI{
		db: db,
	}

	return &collectionApi
}

func (api *CollectionAPI) Bind(scope pathrouter.IRoutes) {
	scope.Get("/collections", api.GetCollections)
	scope.Get("/collections/:collection", api.GetCollection)
	scope.Post("/collections", api.AddCollection)
	scope.Put("/collections", api.UpdateCollection)
	scope.Delete("/collections", api.DeleteCollection)
}

func (api *CollectionAPI) GetCollections(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
	userId := r.Context().Value(constants.AuthKey).(string)

	collections, err := getCollections(api.db, userId)
	if err != nil {
		log.Println(err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.JSON{
			"status": "error",
			"error":  "unable to get collections",
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.JSON{
		"status": "ok",
		"data":   collections,
	})
}

func (api *CollectionAPI) GetCollection(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
	userId := r.Context().Value(constants.AuthKey).(string)

	c := Collection{
		Id:     ps.Get("collection"),
		UserId: userId,
	}

	err := getCollection(api.db, &c)
	if err != nil {
		log.Println(err)
		utils.WriteJSON(w, http.StatusOK, utils.JSON{
			"status": "error",
			"error":  "unable to fetch collection",
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.JSON{
		"status": "ok",
		"data":   c,
	})
}

func (api *CollectionAPI) AddCollection(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
	userId := r.Context().Value(constants.AuthKey).(string)

	c := Collection{
		Id:     uuid.NewString(),
		UserId: userId,
	}

	err := utils.ReadJSON(r.Body, &c)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.JSON{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	if err := c.isValid(); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.JSON{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	err = createCollection(api.db, &c)
	if err != nil {
		log.Println(err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.JSON{
			"status": "error",
			"error":  "failed to create collection",
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.JSON{
		"status": "ok",
		"data":   c,
	})
}

func (api *CollectionAPI) UpdateCollection(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
	userId := r.Context().Value(constants.AuthKey).(string)

	c := Collection{
		UserId: userId,
	}

	err := utils.ReadJSON(r.Body, &c)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.JSON{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	if err := c.isValid(); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.JSON{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	err = updateCollection(api.db, &c)
	if err != nil {
		log.Println(err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.JSON{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.JSON{
		"status": "ok",
	})
}

func (api *CollectionAPI) DeleteCollection(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
	userId := r.Context().Value(constants.AuthKey).(string)
	defer r.Body.Close()

	c := Collection{
		UserId: userId,
	}

	err := utils.ReadJSON(r.Body, &c)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.JSON{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	if err := c.isValid(); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.JSON{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	if !c.Archived {
		err = archiveCollection(api.db, &c)
	} else {
		err = deleteCollection(api.db, &c)
	}

	if err != nil {
		log.Println(err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.JSON{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.JSON{
		"status": "ok",
		"data":   c,
	})
}

func (api *CollectionAPI) CreateDefaultCollection(userId string) error {
	err := createCollection(api.db, &Collection{
		Id:     "Unsorted",
		UserId: userId,
		Name:   "Unsorted",
	})
	return err
}
