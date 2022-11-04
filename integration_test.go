package mermaid_test

import (
	"testing"

	mermaid "github.com/abhinav/goldmark-mermaid"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/testutil"
)

func TestIntegration(t *testing.T) {
	t.Parallel()

	testutil.DoTestCaseFile(
		goldmark.New(goldmark.WithExtensions(&mermaid.Extender{
			MermaidJS: "mermaid.js",
		})),
		"testdata/tests.txt",
		t,
	)

	t.Run("noscript", func(t *testing.T) {
		testutil.DoTestCaseFile(
			goldmark.New(goldmark.WithExtensions(&mermaid.Extender{
				NoScript: true,
			})),
			"testdata/tests_noscript.txt",
			t,
		)
	})
}
