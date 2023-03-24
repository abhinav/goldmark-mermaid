package mermaid

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/mermaid/internal/exectest"
)

func TestServerRenderer_Simple(t *testing.T) {
	mmdc := exectest.Act(t, func() {
		opts, err := parseMermaidOpts(os.Args[1:])
		if err != nil {
			log.Fatal(err)
		}

		if want, got := "neutral", opts.Theme; want != got {
			log.Fatalf("unexpected theme: want %q, got %q", want, got)
		}

		if want, got := "svg", opts.OutputFormat; want != got {
			log.Fatalf("unexpected output format: want %q, got %q", want, got)
		}

		src, err := os.ReadFile(opts.Input)
		if err != nil {
			log.Fatal(err)
		}

		// This isn't real output, but it's good enough for the test.
		svg := "<svg>" + string(src) + "</svg>"
		if err := os.WriteFile(opts.Output, []byte(svg), 0o644); err != nil {
			log.Fatal(err)
		}
	})

	r := buildNodeRenderer(&ServerRenderer{MMDC: mmdc, Theme: "neutral"})
	reader := text.NewReader([]byte(`A -> B`))
	give := blockFromReader(reader)

	var buff bytes.Buffer
	require.NoError(t, r.Render(&buff, reader.Source(), give), "Render")
	assert.Equal(t, `<div class="mermaid"><svg>A -> B</svg></div>`, buff.String())
}

func TestServerRenderer_ContainerTag(t *testing.T) {
	mmdc := exectest.Act(t, func() {
		opts, err := parseMermaidOpts(os.Args[1:])
		if err != nil {
			log.Fatal(err)
		}

		src, err := os.ReadFile(opts.Input)
		if err != nil {
			log.Fatal(err)
		}

		// This isn't real output, but it's good enough for the test.
		svg := "<svg>" + string(src) + "</svg>"
		if err := os.WriteFile(opts.Output, []byte(svg), 0o644); err != nil {
			log.Fatal(err)
		}
	})

	r := buildNodeRenderer(&ServerRenderer{
		MMDC:         mmdc,
		Theme:        "neutral",
		ContainerTag: "pre",
	})
	reader := text.NewReader([]byte(`A -> B`))
	give := blockFromReader(reader)

	var buff bytes.Buffer
	require.NoError(t, r.Render(&buff, reader.Source(), give), "Render")
	assert.Equal(t, `<pre class="mermaid"><svg>A -> B</svg></pre>`, buff.String())
}

func TestServerRenderer_Error_MermaidRender(t *testing.T) {
	mmdc := exectest.Act(t, func() {
		log.Fatal("great sadness")
	})

	r := buildNodeRenderer(&ServerRenderer{MMDC: mmdc})
	reader := text.NewReader([]byte(`A -> B`))
	give := blockFromReader(reader)

	err := r.Render(io.Discard, reader.Source(), give)

	assert.ErrorContains(t, err, "exit status 1")
	assert.ErrorContains(t, err, "great sadness")
}

func TestServerRenderer_Error_NoOutput(t *testing.T) {
	mmdc := exectest.Act(t, func() {
		opts, err := parseMermaidOpts(os.Args[1:])
		if err != nil {
			log.Fatal(err)
		}
		if err := os.Remove(opts.Output); err != nil {
			log.Fatal(err)
		}
		fmt.Println("intentional no-op")
	})

	r := buildNodeRenderer(&ServerRenderer{MMDC: mmdc})
	reader := text.NewReader([]byte(`A -> B`))
	give := blockFromReader(reader)

	err := r.Render(io.Discard, reader.Source(), give)

	assert.ErrorContains(t, err, "read svg:")
	assert.ErrorContains(t, err, "no such file or directory")
	assert.ErrorContains(t, err, "intentional no-op")
}

func TestServerRenderer_Empty(t *testing.T) {
	mmdc := exectest.Act(t, func() {
		log.Fatal("should not be called")
	})

	r := buildNodeRenderer(&ServerRenderer{MMDC: mmdc})
	reader := text.NewReader([]byte{})
	give := blockFromReader(reader)

	var buff bytes.Buffer
	require.NoError(t, r.Render(&buff, reader.Source(), give), "Render")
	assert.Equal(t, `<div class="mermaid"></div>`, buff.String())
}

func TestServerRenderer_ScriptKindNoop(t *testing.T) {
	mmdc := exectest.Act(t, func() {
		log.Fatal("should not be called")
	})

	r := buildNodeRenderer(&ServerRenderer{MMDC: mmdc})
	reader := text.NewReader([]byte{})
	give := new(ScriptBlock)

	var buff bytes.Buffer
	require.NoError(t, r.Render(&buff, reader.Source(), give), "Render")
	assert.Empty(t, buff.String())
}

type mermaidOpts struct {
	Input        string
	Output       string
	OutputFormat string
	Theme        string
	Quiet        bool
}

// parseMermaidOpts is a helper used in tests pretending to be the mmdc CLI.
func parseMermaidOpts(args []string) (*mermaidOpts, error) {
	var o mermaidOpts
	flag := flag.NewFlagSet("mmdc", flag.ContinueOnError)
	flag.StringVar(&o.Input, "input", "", "")
	flag.StringVar(&o.Output, "output", "", "")
	flag.StringVar(&o.Theme, "theme", "", "")
	flag.StringVar(&o.OutputFormat, "outputFormat", "", "")
	flag.BoolVar(&o.Quiet, "quiet", false, "")
	err := flag.Parse(args)
	return &o, err
}
