package main

import (
	"database/sql"
	"linked/internal/opengraph"
	"log"
	"net/http"

	pathrouter "github.com/hehaowen00/path-router"
)

type urlRequest struct {
	Url string `json:"url"`
}

func initOpenGraphApi(db *sql.DB, router *pathrouter.Group) {
	ogw, queue := opengraph.NewWorker()
	go ogw.Run()

	router.Post(
		"/opengraph/info",
		func(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
			urlReq := urlRequest{}

			err := readJson(r.Body, &urlReq)
			if err != nil {
				log.Println(err)
				writeJson(w, http.StatusBadRequest, JsonResult{
					Status: "error",
					Error:  err.Error(),
				})
				return
			}

			recv := make(chan *opengraph.Info)
			queue <- &opengraph.Request{
				Url:  urlReq.Url,
				Recv: recv,
			}

			info := <-recv
			if info == nil {
				writeJson(w, http.StatusNotFound, JsonResult{
					Status: "error",
					Error:  "no opengraph metadata found",
				})
				return
			}

			writeJson(w, http.StatusOK, info)
		},
	)
}
