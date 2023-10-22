package opengraph

import (
	"database/sql"
	"linked/utils"
	"net/http"

	pathrouter "github.com/hehaowen00/path-router"
)

type urlRequest struct {
	Url string `json:"url"`
}

func InitOpenGraphApi(db *sql.DB, scope pathrouter.IRoutes) {
	ogw, queue := NewWorker()
	go ogw.Run()

	scope.Post(
		"/opengraph/info",
		func(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
			urlReq := urlRequest{}

			err := utils.ReadJSON(r.Body, &urlReq)
			if err != nil {
				utils.WriteJSON(w, http.StatusBadRequest, utils.JSON{
					"status": "error",
					"error":  err.Error(),
				})
				return
			}

			recv := make(chan *Info)
			queue <- &Request{
				Url:  urlReq.Url,
				Recv: recv,
			}

			info := <-recv
			if info == nil {
				utils.WriteJSON(w, http.StatusNotFound, utils.JSON{
					"status": "error",
					"error":  "no opengraph metadata found",
				})
				return
			}

			utils.WriteJSON(w, http.StatusOK, info)
		},
	)
}
