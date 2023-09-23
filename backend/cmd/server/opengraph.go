package main

import (
	"linked/internal/opengraph"
	"log"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type openGraphWorker struct {
	httpClient *http.Client
	incoming   chan *openGraphRequest
	stop       chan struct{}
}

type openGraphRequest struct {
	url  string
	recv chan *opengraph.Info
}

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36"

func newOpenGraphWorker() (*openGraphWorker, chan *openGraphRequest) {
	jar, _ := cookiejar.New(nil)
	client := http.Client{
		Jar: jar,
	}

	incoming := make(chan *openGraphRequest)
	ogw := openGraphWorker{
		httpClient: &client,
		incoming:   incoming,
	}

	return &ogw, incoming
}

func (ogw *openGraphWorker) Run() {
	for {
		select {
		case <-ogw.stop:
			break
		case req := <-ogw.incoming:
			resp, err := ogw.httpClient.Get(req.url)
			if err != nil {
				log.Println("error getting webpage", err)
				req.recv <- nil
				continue
			}

			data, err := opengraph.ParseHTML(resp.Body)
			if err != nil {
				log.Println(err)
				req.recv <- nil
				continue
			}
			req.recv <- data

			time.Sleep(1 * time.Second)
		}
	}
}
