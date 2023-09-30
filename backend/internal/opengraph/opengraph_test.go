package opengraph

import (
	"bytes"
	"log"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

const data = `
<!DOCTYPE html>
<html>
	<head>
		<title>title</title>
		<meta property="og:title" content="The Rock" />
		<meta property="og:type" content="video.movie" />
		<meta property="og:url" content="https://www.imdb.com/title/tt0117500/" />
		<meta property="og:image" content="https://ia.media-imdb.com/images/rock.jpg" />
	</head>
	<body>
	</body>
</html>
`

func TestOpenGraph(t *testing.T) {
	reader := bytes.NewReader([]byte(data))
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	selection := doc.Find("title").Text()
	log.Println("title", selection)

	node := doc.Find(`meta[property="og:title"]`)

	v, exists := node.Attr("content")
	if !exists {
		log.Println("does not exist")
		t.FailNow()
	}
	log.Println(v)

	node = doc.Find(`meta[property="og:type"]`)
	v, exists = node.Attr("content")
	if !exists {
		log.Println("does not exist")
		t.FailNow()
	}

	log.Println(v)
}
