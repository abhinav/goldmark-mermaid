package mermaid

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// Compiler compiles Mermaid diagrams into images.
// It is used with [ServerRenderer] to render Mermaid diagrams server-side.
type Compiler interface {
	Compile(context.Context, *CompileRequest) (*CompileResponse, error)
}

// CompileRequest is a request to compile a Mermaid diagram.
type CompileRequest struct {
	// Source is the raw Mermaid diagram source.
	Source string
}

// CompileResponse is a response from compiling a Mermaid diagram.
type CompileResponse struct {
	// SVG holds the SVG diagram text
	// including the <svg>...</svg> tags.
	SVG string
}

// ServerRenderer renders Mermaid diagrams into images server-side.
//
// By default, it uses [CLICompiler] to compile Mermaid diagrams.
// You can specify a different compiler with [Compiler].
// For long-running processes, you should use the compiler
// provided by the mermaidcdp package.
type ServerRenderer struct {
	// Compiler specifies how to compile Mermaid diagrams into images.
	//
	// If unspecified, this uses CLICompiler.
	Compiler Compiler

	// ContainerTag is the name of the HTML tag to use for the container
	// that holds the Mermaid diagram.
	// The name must be without the angle brackets.
	//
	// Defaults to "div".
	ContainerTag string
}

// RegisterFuncs registers the renderer for Mermaid blocks with the provided
// Goldmark Registerer.
func (r *ServerRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(Kind, r.Render)

	// Normally, we won't hit this
	// because Transformer won't add ScriptBlocks for ServerRenderer.
	//
	// Guard against the possibility that the document used a different
	// transformer.
	reg.Register(ScriptKind, func(util.BufWriter, []byte, ast.Node, bool) (ast.WalkStatus, error) {
		return ast.WalkContinue, nil // no-op
	})
}

// Render renders [Block] nodes.
func (r *ServerRenderer) Render(w util.BufWriter, src []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	compiler := r.Compiler
	if compiler == nil {
		compiler = new(CLICompiler)
	}

	tag := r.ContainerTag
	if len(tag) == 0 {
		tag = "div"
	}

	n := node.(*Block)
	if !entering {
		_, _ = w.WriteString("</")
		template.HTMLEscape(w, []byte(tag))
		_, _ = w.WriteString(">")
		return ast.WalkContinue, nil
	}
	_, _ = w.WriteString("<")
	template.HTMLEscape(w, []byte(tag))
	_, _ = w.WriteString(` class="mermaid">`)

	var buff bytes.Buffer
	lines := n.Lines()
	for i := 0; i < lines.Len(); i++ {
		line := lines.At(i)
		buff.Write(line.Value(src))
	}

	if buff.Len() == 0 {
		return ast.WalkContinue, nil
	}

	res, err := compiler.Compile(context.Background(), &CompileRequest{
		Source: buff.String(),
	})
	if err != nil {
		return ast.WalkContinue, fmt.Errorf("generate svg: %w", err)
	}

	_, err = w.WriteString(res.SVG)
	return ast.WalkContinue, err
}
