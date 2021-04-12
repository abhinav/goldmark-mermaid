package mermaid

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// Extender adds support for Mermaid diagrams to a Goldmark Markdown parser.
//
// Use it by installing it to the goldmark.Markdown object upon creation.
//
//   goldmark.New(
//     // ...
//     goldmark.WithExtensions(
//       // ...
//       &mermaid.Exender{},
//     ),
//   )
type Extender struct {
	// URL of Mermaid Javascript to be included in the page.
	//
	// Defaults to the latest version available on cdn.jsdelivr.net.
	MermaidJS string
}

// Extend extends the provided Goldmark parser with support for Mermaid
// diagrams.
func (e *Extender) Extend(md goldmark.Markdown) {
	md.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(&Transformer{}, 100),
		),
	)
	md.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&Renderer{
				MermaidJS: e.MermaidJS,
			}, 100),
		),
	)
}
