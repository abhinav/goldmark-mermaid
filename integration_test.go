package mermaid_test

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
	"go.abhg.dev/goldmark/mermaid"
	"go.abhg.dev/goldmark/mermaid/mermaidcdp"
	"gopkg.in/yaml.v3"
)

func TestIntegration_Client(t *testing.T) {
	t.Parallel()

	testdata, err := os.ReadFile("testdata/client.yaml")
	require.NoError(t, err)

	var tests []struct {
		Desc     string `yaml:"desc"`
		NoScript bool   `yaml:"noscript"`
		Give     string `yaml:"give"`
		Want     string `yaml:"want"`

		ContainerTag string `yaml:"containerTag"`
	}
	require.NoError(t, yaml.Unmarshal(testdata, &tests))

	for _, tt := range tests {
		tt := tt
		t.Run(tt.Desc, func(t *testing.T) {
			t.Parallel()

			ext := mermaid.Extender{
				RenderMode:   mermaid.RenderModeClient,
				MermaidURL:   "mermaid.js",
				NoScript:     tt.NoScript,
				ContainerTag: tt.ContainerTag,
			}
			md := goldmark.New(goldmark.WithExtensions(&ext))

			var got bytes.Buffer
			require.NoError(t, md.Convert([]byte(tt.Give), &got))
			assert.Equal(t,
				strings.TrimSuffix(tt.Want, "\n"),
				strings.TrimSuffix(got.String(), "\n"),
			)
		})
	}
}

func TestIntegration_Server_CLI(t *testing.T) {
	t.Parallel()

	mmdcPath := filepath.Join("node_modules", ".bin", "mmdc")
	if _, err := os.Stat(mmdcPath); err != nil {
		// 'yarn install' must already have been run.
		t.Fatalf("mmdc not found at %s", mmdcPath)
	}

	testdata, err := os.ReadFile("testdata/server_cli.yaml")
	require.NoError(t, err)

	var tests []struct {
		Desc string `yaml:"desc"`
		Give string `yaml:"give"`
		Want string `yaml:"want"`

		ContainerTag string `yaml:"containerTag"`
	}
	require.NoError(t, yaml.Unmarshal(testdata, &tests))

	// HACK:
	// For some reason,
	// mmdc generates an SVG with specific numbers in the output
	// deterministically on my computer,
	// and for the same diagram, also deterministically,
	// it generates slightly different numbers in CI.
	//
	// This basically 'fixes' those in a string.
	numberRe := regexp.MustCompile(`\d+(\.\d+)?`)
	normalize := func(s string) string {
		s = numberRe.ReplaceAllString(s, `0`)
		return strings.TrimSuffix(s, "\n")
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.Desc, func(t *testing.T) {
			t.Parallel()

			ext := mermaid.Extender{
				CLI:          mermaid.MMDC(mmdcPath),
				RenderMode:   mermaid.RenderModeServer,
				ContainerTag: tt.ContainerTag,
			}
			md := goldmark.New(goldmark.WithExtensions(&ext))

			var got bytes.Buffer
			require.NoError(t, md.Convert([]byte(tt.Give), &got))
			assert.Equal(t,
				normalize(tt.Want),
				normalize(got.String()),
			)
		})
	}
}

func TestIntegration_Server_CDP(t *testing.T) {
	cdpCompiler, err := mermaidcdp.New(&mermaidcdp.Config{
		JSSource: loadMermaidJS(t),
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, cdpCompiler.Close(),
			"unable to stop the CDP compiler")
	})

	testdata, err := os.ReadFile("testdata/server_cdp.yaml")
	require.NoError(t, err)

	var tests []struct {
		Desc string `yaml:"desc"`
		Give string `yaml:"give"`
		Want string `yaml:"want"`
	}
	require.NoError(t, yaml.Unmarshal(testdata, &tests))

	for _, tt := range tests {
		tt := tt
		t.Run(tt.Desc, func(t *testing.T) {
			t.Parallel()

			ext := mermaid.Extender{
				Compiler:   cdpCompiler,
				RenderMode: mermaid.RenderModeServer,
			}

			md := goldmark.New(goldmark.WithExtensions(&ext))

			var got bytes.Buffer
			require.NoError(t, md.Convert([]byte(tt.Give), &got))
			assert.Equal(t, tt.Want, got.String())
		})
	}
}

func loadMermaidJS(t *testing.T) string {
	t.Helper()

	b, err := os.ReadFile("testdata/mermaid.js")
	require.NoError(t, err)

	return string(b)
}
