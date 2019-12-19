package images_test

import (
	"log"
	"os"

	images "github.com/mdigger/goldmark-images"
	"github.com/yuin/goldmark"
)

var source = []byte(`![alt](image.png "title")`)

func imageURL(src string) string {
	return "test-" + src
}

func Example() {
	gm := goldmark.New(
		images.NewReplacer(imageURL),
	)
	if err := gm.Convert(source, os.Stdout); err != nil {
		log.Fatal(err)
	}
	// Output:
	// <p><img src="test-image.png" alt="alt" title="title"></p>
}
