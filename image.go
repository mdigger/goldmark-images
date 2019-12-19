package images

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

// ReplaceFunc is a function for replacing image source link.
type ReplaceFunc = func(link string) string

type withReplacer struct {
	value *Replacer
}

func (o *withReplacer) SetConfig(c *renderer.Config) {
	c.NodeRenderers = append(c.NodeRenderers, util.Prioritized(o.value, 0))
}

// WithReplacer adding src url replacing function to image html render.
func WithReplacer(r ReplaceFunc) goldmark.Option {
	return goldmark.WithRendererOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(New(r), 0),
		),
	)
}

// Replacer render image with replaced source link.
type Replacer struct {
	html.Config
	ReplaceFunc
}

// New return initialized image render with source url replacing support.
func New(r ReplaceFunc, options ...html.Option) *Replacer {
	var config = html.NewConfig()
	for _, opt := range options {
		opt.SetHTMLOption(&config)
	}
	return &Replacer{
		Config:      config,
		ReplaceFunc: r,
	}
}

// RegisterFuncs implements NodeRenderer.RegisterFuncs interface.
func (r *Replacer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindImage, r.renderImage)
}

func (r *Replacer) renderImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.Image)
	// add image link replacing hack
	if r.ReplaceFunc != nil {
		var src = r.ReplaceFunc(util.BytesToReadOnlyString(n.Destination))
		// if src == "" {
		// 	return ast.WalkContinue, nil
		// } else if src == "-" {
		// 	return ast.WalkSkipChildren, nil
		// }
		n.Destination = util.StringToReadOnlyBytes(src)
	}

	w.WriteString("<img src=\"")
	if r.Unsafe || !html.IsDangerousURL(n.Destination) {
		w.Write(util.EscapeHTML(util.URLEscape(n.Destination, true)))
	}
	w.WriteString(`" alt="`)
	w.Write(n.Text(source))
	w.WriteByte('"')
	if n.Title != nil {
		w.WriteString(` title="`)
		r.Writer.Write(w, n.Title)
		w.WriteByte('"')
	}
	if n.Attributes() != nil {
		html.RenderAttributes(w, n, html.ImageAttributeFilter)
	}
	if r.XHTML {
		w.WriteString(" />")
	} else {
		w.WriteString(">")
	}
	return ast.WalkSkipChildren, nil
}

// Extend implement goldmark.Extender interface.
func (r *Replacer) Extend(m goldmark.Markdown) {
	if r.ReplaceFunc == nil {
		return
	}
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(r, 0),
		),
	)
}