package mermaid

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// ServerRenderer renders Mermaid diagrams into images server-side.
//
// It operates by replacing mermaid code blocks in your document
// with SVGs.
type ServerRenderer struct {
	// MMDC is the MermaidJS CLI that we'll use
	// to render Mermaid diagrams server-side.
	//
	// Uses CLI by default.
	MMDC MMDC

	// ContainerTag is the name of the HTML tag to use for the container
	// that holds the Mermaid diagram.
	// The name must be without the angle brackets.
	//
	// Defaults to "div".
	ContainerTag string

	// Theme for mermaid diagrams.
	//
	// Values include "dark", "default", "forest", and "neutral".
	// See MermaidJS documentation for a full list.
	Theme string
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
	var mmdc MMDC = DefaultMMDC
	if r.MMDC != nil {
		mmdc = r.MMDC
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

	svgout, err := (&mermaidGenerator{
		MMDC:  mmdc,
		Theme: r.Theme,
	}).Generate(buff.Bytes())
	if err != nil {
		return ast.WalkContinue, fmt.Errorf("generate svg: %w", err)
	}

	_, err = w.Write(svgout)
	return ast.WalkContinue, err
}

type mermaidGenerator struct {
	MMDC  MMDC
	Theme string
}

func (d *mermaidGenerator) Generate(src []byte) (_ []byte, err error) {
	input, err := os.CreateTemp("", "in.*.mermaid")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = os.Remove(input.Name()) // ignore error
	}()

	_, err = input.Write(src)
	if err == nil {
		err = input.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("write input: %w", err)
	}

	output, err := os.CreateTemp("", "out.*.svg")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = os.Remove(output.Name()) // ignore error
	}()
	if err := output.Close(); err != nil {
		return nil, err
	}

	args := []string{
		"--input", input.Name(),
		"--output", output.Name(),
		"--outputFormat", "svg",
		"--quiet",
	}
	if len(d.Theme) > 0 {
		args = append(args, "--theme", d.Theme)
	}

	cmd := d.MMDC.Command(args...)
	// If the user-provided MMDC didn't set Stdout/Stderr,
	// capture its output and if anything fails beyond this point,
	// include the output in the error.
	var cmdout bytes.Buffer
	defer func() {
		if err != nil && cmdout.Len() > 0 {
			err = fmt.Errorf("%w\noutput:\n%s", err, cmdout.String())
		}
	}()
	if cmd.Stdout == nil {
		cmd.Stdout = &cmdout
	}
	if cmd.Stderr == nil {
		cmd.Stderr = &cmdout
	}

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("mmdc: %w", err)
	}

	out, err := os.ReadFile(output.Name())
	if err != nil {
		return nil, fmt.Errorf("read svg: %w", err)
	}
	return out, nil
}
