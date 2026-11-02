package mermaid_test

import (
	"github.com/yuin/goldmark/v2"
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
