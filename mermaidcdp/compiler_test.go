package mermaidcdp

import (
	"context"
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.abhg.dev/goldmark/mermaid"
	"go.abhg.dev/goldmark/mermaid/internal/svgtest"
	"gopkg.in/yaml.v3"
)

var _regenerate = flag.Bool("regenerate", false, "regenerate testdata")

func TestCompiler_Compile_missingSource(t *testing.T) {
	_, err := New(&Config{})
	assert.ErrorContains(t, err, "source code for MermaidJS must be supplied")
}

func TestCompiler_Compile(t *testing.T) {
	t.Parallel()

	c, err := New(&Config{
		JSSource: loadMermaidJS(t),
	})

	require.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, c.Close())
	})

	testdata, err := os.ReadFile("testdata/render.yaml")
	require.NoError(t, err)

	var tests []struct {
		Name string `yaml:"name"`
		Give string `yaml:"give"`
		Want string `yaml:"want"`
	}
	require.NoError(t, yaml.Unmarshal(testdata, &tests))
	if *_regenerate {
		t.Cleanup(func() {
			out, err := yaml.Marshal(tests)
			require.NoError(t, err)

			require.NoError(t, os.WriteFile("testdata/render.yaml", out, 0o644))
		})
	}

	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.Name, func(t *testing.T) {
			t.Parallel()

			got, err := c.Compile(context.Background(), &mermaid.CompileRequest{
				Source: tt.Give,
			})
			require.NoError(t, err)

			if *_regenerate {
				tt.Want = got.SVG
			} else {
				assert.Equal(t,
					svgtest.Normalize(tt.Want),
					svgtest.Normalize(got.SVG),
				)
			}

			tests[i] = tt
		})
	}
}

func TestCompiler_Compile_closed(t *testing.T) {
	t.Parallel()

	c, err := New(&Config{
		JSSource: loadMermaidJS(t),
	})

	require.NoError(t, err)
	assert.NoError(t, c.Close())

	t.Run("render on closed", func(t *testing.T) {
		t.Parallel()

		assert.Panics(t, func() {
			_, err := c.Compile(context.Background(), &mermaid.CompileRequest{
				Source: "graph TD; A-->B;",
			})
			t.Fatalf("Compile should have panicked, but returned %v", err)
		})
	})

	t.Run("double close", func(t *testing.T) {
		t.Parallel()

		assert.NoError(t, c.Close())
	})
}

func TestCompiler_noChrome(t *testing.T) {
	t.Setenv("PATH", t.TempDir())

	_, err := New(&Config{
		JSSource: loadMermaidJS(t),
	})
	assert.ErrorContains(t, err, "set up headless browser")
}

func loadMermaidJS(t *testing.T) string {
	t.Helper()

	b, err := os.ReadFile("testdata/mermaid.js")
	require.NoError(t, err)

	return string(b)
}
