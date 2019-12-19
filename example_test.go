package images_test

import (
	"log"
	"os"

	images "github.com/mdigger/goldmark-images"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

var source = []byte(`![alt](image.png "title")`)

func imageURL(src string) string {
	return "test-" + src
}

func Example() {
	gm := goldmark.New(
		images.WithReplacer(imageURL),
		goldmark.WithRendererOptions(html.WithXHTML()),
	)
	if err := gm.Convert(source, os.Stdout); err != nil {
		log.Fatal(err)
	}
	// Output:
	// <p><img src="test-image.png" alt="alt" title="title" /></p>
}
