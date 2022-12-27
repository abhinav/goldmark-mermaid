package mermaid_test

import (
	"github.com/yuin/goldmark"
	"go.abhg.dev/goldmark/mermaid"
)

func ExampleExtender() {
	goldmark.New(
		// ...
		goldmark.WithExtensions(
			&mermaid.Extender{
				RenderMode: mermaid.RenderModeServer,
			},
			// ...
		),
	)
}
