package opengraph

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36"

type openGraphWorker struct {
	httpClient *http.Client
	incoming   chan *Request
	stop       chan struct{}
}

type Request struct {
	Url  string
	Recv chan *Info
}

func NewWorker() (*openGraphWorker, chan *Request) {
	jar, _ := cookiejar.New(nil)
	client := http.Client{
		Jar: jar,
	}

	incoming := make(chan *Request)
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
			return
		case req := <-ogw.incoming:
			if strings.Contains(req.Url, "localhost") {
				req.Recv <- nil
				time.Sleep(1 * time.Second)
				continue
			}

			resp, err := ogw.httpClient.Do(makeRequest(req.Url))
			if err != nil {
				log.Println("error getting webpage", err)
				req.Recv <- nil
				time.Sleep(1 * time.Second)
				continue
			}

			if resp.StatusCode != http.StatusOK {
				req.Recv <- nil
				time.Sleep(1 * time.Second)
				return
			}

			data, err := parseHtml(resp.Body)
			if err != nil {
				log.Println(err)
				req.Recv <- nil
				time.Sleep(1 * time.Second)
				continue
			}

			faviconUrl := getFavicon(req.Url)

			resp, err = ogw.httpClient.Do(makeRequest(faviconUrl))
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			}

			if resp.StatusCode == http.StatusOK {
				data.ImageURL = faviconUrl
			}

			req.Recv <- data
			time.Sleep(1 * time.Second)
		}
	}
}

func makeRequest(url string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", userAgent)

	return req
}
