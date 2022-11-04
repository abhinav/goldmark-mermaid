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
//	goldmark.New(
//	  // ...
//	  goldmark.WithExtensions(
//	    // ...
//	    &mermaid.Exender{},
//	  ),
//	)
type Extender struct {
	// URL of Mermaid Javascript to be included in the page.
	// Ignored if NoScript is true.
	//
	// Defaults to the latest version available on cdn.jsdelivr.net.
	MermaidJS string

	// If true, don't add a <script> including Mermaid to the end of the
	// page.
	//
	// Use this if the page you're including goldmark-mermaid in
	// already has a MermaidJS script included elsewhere.
	NoScript bool
}

// Extend extends the provided Goldmark parser with support for Mermaid
// diagrams.
func (e *Extender) Extend(md goldmark.Markdown) {
	md.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(&Transformer{
				NoScript: e.NoScript,
			}, 100),
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
