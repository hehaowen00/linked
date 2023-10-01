package opengraph

import (
	"io"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Info struct {
	Title       string `json:"title"`
	Type        string `json:"type"`
	Description string `json:"desc"`
	ImageURL    string `json:"image_url"`
}

func parseHtml(stream io.Reader) (*Info, error) {
	doc, err := goquery.NewDocumentFromReader(stream)
	if err != nil {
		return nil, err
	}

	info := Info{}

	title := doc.Find("title").First().Text()
	info.Title = strings.TrimSpace(title)

	node := doc.Find(`meta[property="og:title"]`).First()

	v, exists := node.Attr("content")
	if !exists {
		return &info, nil
	}
	if len(v) >= len(info.Title) {
		info.Title = v
	}

	node = doc.Find(`meta[property="og:type"]`).First()
	v, exists = node.Attr("content")
	if !exists {
		return &info, nil
	}
	info.Type = v

	node = doc.Find(`meta[property="og:description"]`).First()
	v, exists = node.Attr("content")
	if !exists {
		return &info, nil
	}
	info.Description = strings.ReplaceAll(strings.TrimSpace(v), "\n", "")

	return &info, nil
}

func getFavicon(requestUrl string) string {
	req, err := url.Parse(requestUrl)
	if err != nil {
		panic(err)
	}

	favicon := url.URL{
		Scheme: req.Scheme,
		Host:   req.Host,
		Path:   "favicon.ico",
	}

	return favicon.String()
}
