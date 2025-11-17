package mermaid

import (
	"encoding/json"
	"html/template"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

const _defaultMermaidJS = "https://cdn.jsdelivr.net/npm/mermaid/dist/mermaid.min.js"

// initializationOptions represents the client render mode options
type initializationOptions struct {
	StartOnLoad bool   `json:"startOnLoad"`
	Theme       string `json:"theme,omitempty"`
}

// ClientRenderer renders Mermaid diagrams as HTML,
// to be rendered into images client side.
//
// It operates by installing a <script> tag into the document
// that renders the Mermaid diagrams client-side.
type ClientRenderer struct {
	// URL of Mermaid Javascript to be included in the page.
	//
	// Defaults to the latest version available on cdn.jsdelivr.net.
	MermaidURL string

	// ContainerTag is the name of the HTML tag to use for the container
	// that holds the Mermaid diagram.
	// The name must be without the angle brackets.
	//
	// Defaults to "pre".
	ContainerTag string

	initializeOptions initializationOptions
}

// RegisterFuncs registers the renderer for Mermaid blocks with the provided
// Goldmark Registerer.
func (r *ClientRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(Kind, r.Render)
	reg.Register(ScriptKind, r.RenderScript)
}

// Render renders mermaid.Block nodes.
func (r *ClientRenderer) Render(w util.BufWriter, src []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	tag := r.ContainerTag
	if len(tag) == 0 {
		tag = "pre"
	}

	n := node.(*Block)
	if entering {
		_, _ = w.WriteString("<")
		template.HTMLEscape(w, []byte(tag))
		_, _ = w.WriteString(` class="mermaid">`)

		lines := n.Lines()
		for i := 0; i < lines.Len(); i++ {
			line := lines.At(i)
			template.HTMLEscape(w, line.Value(src))
		}
	} else {
		_, _ = w.WriteString("</")
		template.HTMLEscape(w, []byte(tag))
		_, _ = w.WriteString(">")
	}
	return ast.WalkContinue, nil
}

// RenderScript renders mermaid.ScriptBlock nodes.
func (r *ClientRenderer) RenderScript(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	mermaidJS := r.MermaidURL
	if len(mermaidJS) == 0 {
		mermaidJS = _defaultMermaidJS
	}

	_ = node.(*ScriptBlock) // sanity check
	if entering {
		_, _ = w.WriteString(`<script src="`)
		_, _ = w.WriteString(mermaidJS)
		_, _ = w.WriteString(`"></script>`)
	} else {
		b, err := json.Marshal(r.initializeOptions)
		if err != nil {
			return ast.WalkStop, err
		}

		_, _ = w.WriteString("<script>mermaid.initialize(")
		_, _ = w.Write(b)
		_, _ = w.WriteString(");</script>")
	}

	return ast.WalkContinue, nil
}
