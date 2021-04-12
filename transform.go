package mermaid

import (
	"bytes"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// Transformer transforms a Goldmark Markdown AST with support for Mermaid
// diagrams. It makes the following transformations:
//
//  - replace mermaid code blocks with mermaid.Block nodes
//  - add a mermaid.ScriptBlock node if the document uses Mermaid
type Transformer struct {
}

var _mermaid = []byte("mermaid")

// Transform transforms the provided Markdown AST.
func (*Transformer) Transform(doc *ast.Document, reader text.Reader, pctx parser.Context) {
	var (
		used      bool
		hasScript bool
	)
	ast.Walk(doc, func(node ast.Node, enter bool) (ast.WalkStatus, error) {
		// For multiple transforms.
		if _, ok := node.(*ScriptBlock); ok {
			hasScript = true
			return ast.WalkContinue, nil
		}

		cb, ok := node.(*ast.FencedCodeBlock)
		if !ok {
			return ast.WalkContinue, nil
		}

		lang := cb.Language(reader.Source())
		if !bytes.Equal(lang, _mermaid) {
			return ast.WalkContinue, nil
		}

		used = true
		b := new(Block)
		b.SetLines(cb.Lines())

		parent := cb.Parent()
		if parent != nil {
			parent.ReplaceChild(parent, cb, b)
		}

		return ast.WalkSkipChildren, nil
	})

	if !used || hasScript {
		return
	}

	doc.AppendChild(doc, &ScriptBlock{})
}
