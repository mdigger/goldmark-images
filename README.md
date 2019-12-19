# goldmark-images

[![GoDoc](https://godoc.org/github.com/mdigger/goldmark-images?status.svg)](https://godoc.org/github.com/mdigger/goldmark-images)

[Goldmark](https://github.com/yuin/goldmark) image replacer extension.

```go
imageURL := func (src string) string {
	return "test-" + src
}

source := []byte(`![alt](image.png "title")`)
gm := goldmark.New(
    images.NewReplacer(imageURL),
    goldmark.WithRendererOptions(html.WithXHTML()),
)
err = gm.Convert(source, os.Stdout)
```

```html
<p><img src="test-image.png" alt="alt" title="title" /></p>
```