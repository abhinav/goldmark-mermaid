package mermaid_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	mermaid "github.com/abhinav/goldmark-mermaid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
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
	}
	require.NoError(t, yaml.Unmarshal(testdata, &tests))

	for _, tt := range tests {
		tt := tt
		t.Run(tt.Desc, func(t *testing.T) {
			ext := mermaid.Extender{
				MermaidJS: "mermaid.js",
				NoScript:  tt.NoScript,
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
