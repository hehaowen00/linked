package opengraph

import (
	"io"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

type Info struct {
	Title       string `json:"title"`
	Type        string `json:"type"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}

func ParseHTML(stream io.Reader) (*Info, error) {
	doc, err := goquery.NewDocumentFromReader(stream)
	if err != nil {
		return nil, err
	}

	info := Info{}

	title := doc.Find("title").Text()
	info.Title = title

	node := doc.Find(`meta[property="og:title"]`)

	v, exists := node.Attr("content")
	if !exists {
		return &info, nil
	}
	info.Title = v

	node = doc.Find(`meta[property="og:type"]`)
	v, exists = node.Attr("content")
	if !exists {
		return &info, nil
	}
	info.Type = v

	node = doc.Find(`meta[property="og:description"]`)
	v, exists = node.Attr("content")
	if !exists {
		return &info, nil
	}
	info.Description = v

	// node = doc.Find(`meta[property="og:image"]`)
	// v, exists = node.Attr("content")
	// if !exists {
	// 	return &info, nil
	// }
	// info.ImageURL = v

	return &info, nil
}

func GetFavicon(requestUrl string) string {
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
