package mermaid

import "go.abhg.dev/goldmark/mermaid"

// Extender adds support for Mermaid diagrams to a Goldmark Markdown parser.
//
// Use it by installing it to the goldmark.Markdown object upon creation.
//
//	goldmark.New(
//		// ...
//		goldmark.WithExtensions(
//			// ...
//			&mermaid.Exender{
//				RenderMode: mermaid.ServerRenderMode,
//			},
//		),
//	)
type Extender = mermaid.Extender
