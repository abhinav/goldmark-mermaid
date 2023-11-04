package mermaid

import (
	"fmt"
	"os/exec"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// Extender adds support for Mermaid diagrams to a Goldmark Markdown parser.
//
// Use it by installing it to the goldmark.Markdown object upon creation.
type Extender struct {
	// RenderMode specifies which renderer the Extender should install.
	//
	// Defaults to AutoRenderMode, picking renderers
	// based on the availability of the Mermaid CLI.
	RenderMode RenderMode

	// Compiler specifies how to compile Mermaid diagrams server-side.
	//
	// If specified, and render mode is not set to client-side,
	// this will be used to render diagrams.
	Compiler Compiler

	// CLI specifies how to invoke the Mermaid CLI
	// to compile Mermaid diagrams server-side.
	//
	// If specified, and render mode is not set to client-side,
	// this will be used to render diagrams.
	//
	// If both CLI and Compiler are specified, Compiler takes precedence.
	CLI CLI

	// URL of Mermaid Javascript to be included in the page
	// for client-side rendering.
	//
	// Ignored if NoScript is true or if we're rendering diagrams server-side.
	//
	// Defaults to the latest version available on cdn.jsdelivr.net.
	MermaidURL string

	// HTML tag to use for the container element for diagrams.
	//
	// Defaults to "pre" for client-side rendering,
	// and "div" for server-side rendering.
	ContainerTag string

	// If true, don't add a <script> including Mermaid to the end of the
	// page even if rendering diagrams client-side.
	//
	// Use this if the page you're including goldmark-mermaid in
	// already has a MermaidJS script included elsewhere.
	NoScript bool

	// Theme for mermaid diagrams.
	//
	// Ignored if we're rendering diagrams client-side.
	//
	// Values include "dark", "default", "forest", and "neutral".
	// See MermaidJS documentation for a full list.
	Theme string

	execLookPath func(string) (string, error) // == exec.LookPath
}

// Extend extends the provided Goldmark parser with support for Mermaid
// diagrams.
func (e *Extender) Extend(md goldmark.Markdown) {
	mode, r := e.renderer()

	md.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(&Transformer{
				// If rendering server-side,
				// don't generate <script> tags.
				NoScript: e.NoScript || mode == RenderModeServer,
			}, 100),
		),
	)

	md.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(r, 100),
		),
	)
}

func (e *Extender) renderer() (RenderMode, renderer.NodeRenderer) {
	mode := e.RenderMode
	compiler, ok := e.compiler()
	if mode == RenderModeAuto {
		if ok {
			mode = RenderModeServer
		} else {
			mode = RenderModeClient
		}
	}

	switch mode {
	case RenderModeClient:
		return RenderModeClient, &ClientRenderer{
			MermaidURL:   e.MermaidURL,
			ContainerTag: e.ContainerTag,
		}
	case RenderModeServer:
		return RenderModeServer, &ServerRenderer{
			Compiler:     compiler,
			ContainerTag: e.ContainerTag,
		}
	default:
		panic(fmt.Sprintf("unrecognized render mode: %v", mode))
	}
}

// compiler returns the Compiler to use for server-side rendering
// only if server-side rendering should be used.
//
// The following conditions will cause server-side rendering:
//
//   - Compiler is set
//   - CLI is set (will use CLICompiler)
//   - mmdc is available on $PATH
func (e *Extender) compiler() (c Compiler, ok bool) {
	if e.Compiler != nil {
		return e.Compiler, true
	}

	if e.CLI != nil {
		return &CLICompiler{CLI: e.CLI, Theme: e.Theme}, true
	}

	lookPath := exec.LookPath
	if e.execLookPath != nil {
		lookPath = e.execLookPath
	}

	mmdcPath, err := lookPath("mmdc")
	if err != nil {
		return nil, false
	}

	cli := &mmdcCLI{Path: mmdcPath}
	return &CLICompiler{CLI: cli, Theme: e.Theme}, true
}
