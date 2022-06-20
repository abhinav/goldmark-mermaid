package mermaid

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

func TestRenderer_Block(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc string
		give string
		want string
	}{
		{
			desc: "empty",
			give: "",
			want: `<div class="mermaid"></div>`,
		},
		{
			desc: "graph",
			give: "graph TD;",
			want: `<div class="mermaid">graph TD;</div>`,
		},
		{
			desc: "newlines",
			give: unlines("foo", "bar"),
			want: `<div class="mermaid">foo` + "\nbar" + "\n</div>",
		},
		{
			desc: "escaping",
			give: "A -> B",
			want: `<div class="mermaid">A -&gt; B</div>`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()

			r := buildNodeRenderer(new(ClientRenderer))

			reader := text.NewReader([]byte(tt.give))
			give := blockFromReader(reader)

			var buff bytes.Buffer
			assert.NoError(t, r.Render(&buff, reader.Source(), give), "Render")
			assert.Equal(t, tt.want, buff.String())
		})
	}
}

func TestRenderer_Script(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc      string
		mermaidJS string
		want      string
	}{
		{
			desc: "default mermaid.js",
			want: fmt.Sprintf("<script src=%q></script><script>mermaid.initialize({startOnLoad: true});</script>", _defaultMermaidJS),
		},
		{
			desc:      "explicit mermaid.js",
			mermaidJS: "mermaid.js",
			want:      `<script src="mermaid.js"></script><script>mermaid.initialize({startOnLoad: true});</script>`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()

			r := buildNodeRenderer(&ClientRenderer{
				MermaidJS: tt.mermaidJS,
			})

			var buff bytes.Buffer
			assert.NoError(t,
				r.Render(&buff, nil /* src */, &ScriptBlock{}))
			assert.Equal(t, tt.want, buff.String())
		})
	}
}

func buildNodeRenderer(r renderer.NodeRenderer) renderer.Renderer {
	return renderer.NewRenderer(
		renderer.WithNodeRenderers(
			util.Prioritized(r, 100),
		),
	)
}
