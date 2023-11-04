package mermaid

// RenderMode specifies which renderer the Extender should use.
type RenderMode int

//go:generate stringer -type RenderMode -trimprefix RenderMode

const (
	// RenderModeAuto picks the renderer automatically.
	//
	// If a server-side compiler or CLI is specified,
	// or if the 'mmdc' CLI is available on $PATH,
	// this will generate diagrams server-side.
	//
	// Otherwise, it'll generate them client-side.
	RenderModeAuto RenderMode = iota

	// RenderModeClient renders Mermaid diagrams client-side
	// by adding <script> tags.
	RenderModeClient

	// RenderModeServer renders Mermaid diagrams server-side
	// using the Mermaid CLI.
	//
	// Fails rendering if the Mermaid CLI is absent.
	RenderModeServer
)
