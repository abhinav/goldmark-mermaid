// demo implements a WASM module that can be used to format markdown
// with the goldmark-mermaid extension.
package main

import (
	"bytes"
	"syscall/js"

	"github.com/yuin/goldmark"
	"go.abhg.dev/goldmark/mermaid"
)

func main() {
	js.Global().Set("formatMarkdown", js.FuncOf(formatMarkdown))
	select {}
}

func formatMarkdown(this js.Value, args []js.Value) any {
	input := args[0].String()

	md := goldmark.New(
		goldmark.WithExtensions(
			&mermaid.Extender{
				RenderMode: mermaid.RenderModeClient,
				NoScript:   true,
			},
		),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(input), &buf); err != nil {
		return err.Error()
	}
	return buf.String()
}
