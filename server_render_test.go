package mermaid

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark/text"
)

func TestServerRenderer_Simple(t *testing.T) {
	t.Parallel()

	compiler := compilerStub{
		CompileF: func(ctx context.Context, req *CompileRequest) (*CompileResponse, error) {
			return &CompileResponse{
				SVG: "<svg>" + req.Source + "</svg>",
			}, nil
		},
	}

	r := buildNodeRenderer(&ServerRenderer{
		Compiler: &compiler,
	})
	reader := text.NewReader([]byte(`A -> B`))
	give := blockFromReader(reader)

	var buff bytes.Buffer
	require.NoError(t, r.Render(&buff, reader.Source(), give), "Render")
	assert.Equal(t, `<div class="mermaid"><svg>A -> B</svg></div>`, buff.String())
}

func TestServerRenderer_ContainerTag(t *testing.T) {
	t.Parallel()

	compiler := compilerStub{
		CompileF: func(ctx context.Context, req *CompileRequest) (*CompileResponse, error) {
			return &CompileResponse{
				SVG: "<svg>" + req.Source + "</svg>",
			}, nil
		},
	}

	r := buildNodeRenderer(&ServerRenderer{
		Compiler:     &compiler,
		ContainerTag: "pre",
	})
	reader := text.NewReader([]byte(`A -> B`))
	give := blockFromReader(reader)

	var buff bytes.Buffer
	require.NoError(t, r.Render(&buff, reader.Source(), give), "Render")
	assert.Equal(t, `<pre class="mermaid"><svg>A -> B</svg></pre>`, buff.String())
}

func TestServerRenderer_Empty(t *testing.T) {
	t.Parallel()

	compiler := compilerStub{
		CompileF: func(context.Context, *CompileRequest) (*CompileResponse, error) {
			t.Error("should not be called")
			panic("unreachable")
		},
	}

	r := buildNodeRenderer(&ServerRenderer{
		Compiler: &compiler,
	})
	reader := text.NewReader([]byte{})
	give := blockFromReader(reader)

	var buff bytes.Buffer
	require.NoError(t, r.Render(&buff, reader.Source(), give), "Render")
	assert.Equal(t, `<div class="mermaid"></div>`, buff.String())
}

func TestServerRenderer_ScriptKindNoop(t *testing.T) {
	t.Parallel()

	compiler := compilerStub{
		CompileF: func(context.Context, *CompileRequest) (*CompileResponse, error) {
			t.Error("should not be called")
			panic("unreachable")
		},
	}

	r := buildNodeRenderer(&ServerRenderer{
		Compiler: &compiler,
	})
	reader := text.NewReader([]byte{})
	give := new(ScriptBlock)

	var buff bytes.Buffer
	require.NoError(t, r.Render(&buff, reader.Source(), give), "Render")
	assert.Empty(t, buff.String())
}

type compilerStub struct {
	CompileF func(context.Context, *CompileRequest) (*CompileResponse, error)
}

var _ Compiler = (*compilerStub)(nil)

func (c *compilerStub) Compile(ctx context.Context, req *CompileRequest) (*CompileResponse, error) {
	return c.CompileF(ctx, req)
}
