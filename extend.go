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

	// URL of Mermaid Javascript to be included in the page.
	//
	// Ignored if NoScript is true
	// or if we're rendering diagrams server-side.
	//
	// Defaults to the latest version available on cdn.jsdelivr.net.
	MermaidJS string

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

	// MMDC provides access to the Mermaid CLI.
	//
	// Ignored if we're rendering diagrams client-side.
	//
	// Uses DefaultMMDC if unset.
	MMDC MMDC

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
	if mode == RenderModeAuto {
		lookPath := exec.LookPath
		if e.execLookPath != nil {
			lookPath = e.execLookPath
		}

		if mmdcPath, err := lookPath("mmdc"); err != nil {
			mode = RenderModeClient
		} else {
			mode = RenderModeServer
			if e.MMDC == nil {
				e.MMDC = &CLI{Path: mmdcPath}
			}
		}
	}

	switch mode {
	case RenderModeClient:
		return RenderModeClient, &ClientRenderer{
			MermaidJS:    e.MermaidJS,
			ContainerTag: e.ContainerTag,
		}
	case RenderModeServer:
		return RenderModeServer, &ServerRenderer{
			MMDC:         e.MMDC,
			Theme:        e.Theme,
			ContainerTag: e.ContainerTag,
		}
	default:
		panic(fmt.Sprintf("unrecognized render mode: %v", mode))
	}
}
