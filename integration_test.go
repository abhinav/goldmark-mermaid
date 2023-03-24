package mermaid_test

import (
	"bytes"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
	"go.abhg.dev/goldmark/mermaid"
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
			ext := mermaid.Extender{
				RenderMode:   mermaid.RenderModeClient,
				MermaidJS:    "mermaid.js",
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

func TestIntegration_Server(t *testing.T) {
	t.Parallel()

	testdata, err := os.ReadFile("testdata/server.yaml")
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
			// 'yarn install' must already have been run.
			mmdc := mermaid.CLI{
				Path: "node_modules/.bin/mmdc",
			}

			ext := mermaid.Extender{
				RenderMode:   mermaid.RenderModeServer,
				MMDC:         &mmdc,
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
