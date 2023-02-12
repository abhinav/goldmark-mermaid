package mermaid

import "go.abhg.dev/goldmark/mermaid"

// RenderMode specifies which renderer the Extender should use.
type RenderMode = mermaid.RenderMode

const (
	// RenderModeAuto picks the renderer
	// based on the availability of the Mermaid CLI.
	//
	// If the 'mmdc' CLI is available on $PATH,
	// this will generate diagrams server-side.
	// Otherwise, it'll generate them client-side.
	RenderModeAuto = mermaid.RenderModeAuto

	// RenderModeClient renders Mermaid diagrams client-side
	// by adding <script> tags.
	RenderModeClient = mermaid.RenderModeClient

	// RenderModeServer renders Mermaid diagrams server-side
	// using the Mermaid CLI.
	//
	// Fails rendering if the Mermaid CLI is absent.
	RenderModeServer = mermaid.RenderModeServer
)
