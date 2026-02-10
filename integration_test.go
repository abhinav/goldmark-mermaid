package mermaid_test

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
	"go.abhg.dev/goldmark/mermaid"
	"go.abhg.dev/goldmark/mermaid/internal/svgtest"
	"go.abhg.dev/goldmark/mermaid/mermaidcdp"
	"gopkg.in/yaml.v3"
)

var updateIntegration = flag.Bool("update", false, "update integration test fixtures")

func TestIntegration_Client(t *testing.T) {
	testdata, err := os.ReadFile("testdata/client.yaml")
	require.NoError(t, err)

	var tests []struct {
		Desc     string `yaml:"desc"`
		NoScript bool   `yaml:"noscript"`
		Give     string `yaml:"give"`
		Want     string `yaml:"want"`
		Theme    string `yaml:"theme"`

		ContainerTag string `yaml:"containerTag"`
	}
	require.NoError(t, yaml.Unmarshal(testdata, &tests))

	for i, tt := range tests {
		tt := tt
		t.Run(tt.Desc, func(t *testing.T) {
			ext := mermaid.Extender{
				RenderMode:   mermaid.RenderModeClient,
				MermaidURL:   "mermaid.js",
				NoScript:     tt.NoScript,
				ContainerTag: tt.ContainerTag,
				Theme:        tt.Theme,
			}
			md := goldmark.New(goldmark.WithExtensions(&ext))

			var got bytes.Buffer
			require.NoError(t, md.Convert([]byte(tt.Give), &got))

			if *updateIntegration {
				tests[i].Want = got.String()
			} else {
				assert.Equal(t,
					strings.TrimSuffix(tt.Want, "\n"),
					strings.TrimSuffix(got.String(), "\n"),
				)
			}
		})
	}

	if *updateIntegration {
		data, err := yaml.Marshal(tests)
		require.NoError(t, err)
		require.NoError(t, os.WriteFile("testdata/client.yaml", data, 0o644))
	}
}

func TestIntegration_Server_CLI(t *testing.T) {
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

	for i, tt := range tests {
		tt := tt
		t.Run(tt.Desc, func(t *testing.T) {
			ext := mermaid.Extender{
				CLI:          mermaid.MMDC(mmdcPath),
				RenderMode:   mermaid.RenderModeServer,
				ContainerTag: tt.ContainerTag,
			}
			md := goldmark.New(goldmark.WithExtensions(&ext))

			var got bytes.Buffer
			require.NoError(t, md.Convert([]byte(tt.Give), &got))

			if *updateIntegration {
				tests[i].Want = got.String()
			} else {
				assert.Equal(t,
					svgtest.Normalize(tt.Want),
					svgtest.Normalize(got.String()),
				)
			}
		})
	}

	if *updateIntegration {
		data, err := yaml.Marshal(tests)
		require.NoError(t, err)
		require.NoError(t, os.WriteFile("testdata/server_cli.yaml", data, 0o644))
	}
}

func TestIntegration_Server_CDP(t *testing.T) {
	cdpCompiler, err := mermaidcdp.New(&mermaidcdp.Config{
		JSSource:  loadMermaidJS(t),
		NoSandbox: true,
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

	for i, tt := range tests {
		tt := tt
		t.Run(tt.Desc, func(t *testing.T) {
			ext := mermaid.Extender{
				Compiler:   cdpCompiler,
				RenderMode: mermaid.RenderModeServer,
			}

			md := goldmark.New(goldmark.WithExtensions(&ext))

			var got bytes.Buffer
			require.NoError(t, md.Convert([]byte(tt.Give), &got))

			if *updateIntegration {
				tests[i].Want = got.String()
			} else {
				assert.Equal(t,
					svgtest.Normalize(tt.Want),
					svgtest.Normalize(got.String()),
				)
			}
		})
	}

	if *updateIntegration {
		data, err := yaml.Marshal(tests)
		require.NoError(t, err)
		require.NoError(t, os.WriteFile("testdata/server_cdp.yaml", data, 0o644))
	}
}

func loadMermaidJS(t *testing.T) string {
	t.Helper()

	b, err := os.ReadFile("testdata/mermaid.js")
	require.NoError(t, err)

	return string(b)
}
