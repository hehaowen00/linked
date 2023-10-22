package items

import (
	"database/sql"
	"linked/internal/constants"
	"linked/utils"
	"log"
	"net/http"

	pathrouter "github.com/hehaowen00/path-router"
)

type ItemsAPI struct {
	db *sql.DB
}

func NewAPI(db *sql.DB) *ItemsAPI {
	api := ItemsAPI{
		db: db,
	}

	return &api
}

func (api *ItemsAPI) Bind(scope pathrouter.IRoutes) {
	scope.Get("/items/:collection", api.GetItemsByCollection)
	scope.Post("/items/:collection", api.AddItemMapping)
	scope.Delete("/items/:collection", api.RemoveItemMapping)

	scope.Get("/items", api.GetItems)
	scope.Post("/items", api.AddItem)
	scope.Put("/items", api.UpdateItem)
	scope.Delete("/items", api.DeleteItem)
}

func (api *ItemsAPI) GetItemsByCollection(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
	userId := r.Context().Value(constants.AuthKey).(string)

	items, err := getItemsByCollection(api.db, ps.Get("collection"), userId)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.JSON{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, items)
}

func (api *ItemsAPI) AddItemMapping(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
	userId := r.Context().Value(constants.AuthKey).(string)

	item := Item{
		CollectionId: ps.Get("collection"),
		UserId:       userId,
	}

	err := utils.ReadJSON(r.Body, &item)
	if err != nil {
		utils.WriteJSON(w, http.StatusOK, utils.JSON{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	err = item.Validate()
	if err != nil {
		log.Println(err)
		utils.WriteJSON(w, http.StatusOK, utils.JSON{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	err = addItemToCollection(api.db, &item)
	if err != nil {
		utils.WriteJSON(w, http.StatusOK, utils.JSON{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.JSON{
		"status": "ok",
		"data":   item,
	})
}

func (api *ItemsAPI) RemoveItemMapping(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
	userId := r.Context().Value(constants.AuthKey).(string)

	item := Item{
		CollectionId: ps.Get("collection"),
		UserId:       userId,
	}

	err := utils.ReadJSON(r.Body, &item)
	if err != nil {
		utils.WriteJSON(w, http.StatusOK, utils.JSON{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	err = item.Validate()
	if err != nil {
		log.Println(err)
		utils.WriteJSON(w, http.StatusOK, utils.JSON{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	err = deleteItemMapping(api.db, &item)
	if err != nil {
		utils.WriteJSON(w, http.StatusOK, utils.JSON{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.JSON{
		"status": "ok",
	})
}

func (api *ItemsAPI) GetItems(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
	userId := r.Context().Value(constants.AuthKey).(string)
	defer r.Body.Close()

	items, err := getItems(api.db, userId)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.JSON{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, items)
}

func (api *ItemsAPI) AddItem(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
	userId := r.Context().Value(constants.AuthKey).(string)
	defer r.Body.Close()

	item := Item{
		UserId: userId,
	}

	err := utils.ReadJSON(r.Body, &item)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.JSON{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	err = item.Validate()
	if err != nil {
		log.Println(err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.JSON{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	exists, err := getItemByUrl(api.db, &item)
	if err != nil {
		log.Println(err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.JSON{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	if !exists {
		err = createItem(api.db, &item)
		if err != nil {
			log.Println(err)
			utils.WriteJSON(w, http.StatusInternalServerError, utils.JSON{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}

	} else {
		err = updateItem(api.db, &item)
		if err != nil {
			log.Println(err)
			utils.WriteJSON(w, http.StatusInternalServerError, utils.JSON{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}
	}

	if item.CollectionId != "" {
		err = addItemToCollection(api.db, &item)
		if err != nil {
			log.Println(err)
			utils.WriteJSON(w, http.StatusInternalServerError, utils.JSON{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}
	}

	utils.WriteJSON(w, http.StatusOK, utils.JSON{
		"status": "ok",
		"data":   item,
	})
}

func (api *ItemsAPI) UpdateItem(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
	userId := r.Context().Value(constants.AuthKey).(string)
	defer r.Body.Close()

	item := Item{
		ID:     ps.Get("item"),
		UserId: userId,
	}

	err := utils.ReadJSON(r.Body, &item)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.JSON{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	err = item.Validate()
	if err != nil {
		log.Println(err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.JSON{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	err = updateItem(api.db, &item)
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
		"data":   item,
	})
}

func (api *ItemsAPI) DeleteItem(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
	userId := r.Context().Value(constants.AuthKey).(string)
	defer r.Body.Close()

	item := Item{
		UserId: userId,
	}

	err := utils.ReadJSON(r.Body, &item)
	if err != nil {
		utils.WriteJSON(w, http.StatusOK, utils.JSON{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	err = item.Validate()
	if err != nil {
		log.Println(err)
		utils.WriteJSON(w, http.StatusOK, utils.JSON{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	err = deleteItem(api.db, &item)

	if err != nil {
		log.Println(err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.JSON{
			"status": "error",
			"error":  "unable to delete item",
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.JSON{
		"status": "ok",
		"data":   item,
	})
}
