package opengraph

import (
	"fmt"
	"io"

	"golang.org/x/net/html"
)

type Info struct {
	Title       string `json:"title"`
	Type        string `json:"type"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}

func ParseHTML(stream io.Reader) (*Info, error) {
	tkn := html.NewTokenizer(stream)
	data := Info{}

	headStart := false
	found := false

	for {
		tokenType := tkn.Next()
		switch tokenType {
		case html.ErrorToken:
			return &data, nil
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tkn.Token()
			if token.Data == "head" {
				if headStart {
					if found {
						return &data, nil
					}
					return nil, fmt.Errorf("no opengraph data found")
				}
				headStart = true
			}
			if token.Data == "meta" {
				propertyName := ""
				propertyContent := ""

				for _, attr := range token.Attr {
					if attr.Key == "property" {
						propertyName = attr.Val
					}
					if attr.Key == "content" {
						propertyContent = attr.Val
					}
				}

				switch propertyName {
				case "og:title":
					data.Title = propertyContent
					found = true
				case "og:type":
					data.Type = propertyContent
				case "og:description":
					data.Description = propertyContent
				case "og:image":
					data.ImageURL = propertyContent
				}
			}
		}
	}
}
