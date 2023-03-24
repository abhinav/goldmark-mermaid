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
	js.Global().Set("formatMarkdown", js.FuncOf(func(this js.Value, args []js.Value) any {
		var req request
		req.Decode(args[0])

		return formatMarkdown(&req)
	}))
	select {}
}

type request struct {
	Markdown     string
	ContainerTag string
}

func (r *request) Decode(v js.Value) {
	r.Markdown = v.Get("markdown").String()
	r.ContainerTag = v.Get("containerTag").String()
}

func formatMarkdown(r *request) any {
	input := r.Markdown
	md := goldmark.New(
		goldmark.WithExtensions(
			&mermaid.Extender{
				RenderMode:   mermaid.RenderModeClient,
				NoScript:     true,
				ContainerTag: r.ContainerTag,
			},
		),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(input), &buf); err != nil {
		return err.Error()
	}
	return buf.String()
}
