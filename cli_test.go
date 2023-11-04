package mermaid

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.abhg.dev/goldmark/mermaid/internal/exectest"
)

func TestCLI_NoPath(t *testing.T) {
	t.Parallel()

	var cli mmdcCLI
	cmd := cli.CommandContext(context.Background(), "--version")
	assert.Equal(t, []string{"mmdc", "--version"}, cmd.Args)
}

func TestCLI_ExplicitPath(t *testing.T) {
	t.Parallel()

	cli := mmdcCLI{Path: "/bin/false"}
	cmd := cli.CommandContext(context.Background(), "--version")
	assert.Equal(t, []string{"/bin/false", "--version"}, cmd.Args)
}

func TestCLICompiler_Simple(t *testing.T) {
	t.Parallel()

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

	c := CLICompiler{
		CLI:   mmdc,
		Theme: "neutral",
	}
	res, err := c.Compile(context.Background(), &CompileRequest{
		Source: `A -> B`,
	})
	require.NoError(t, err)
	assert.Equal(t, `<svg>A -> B</svg>`, res.SVG)
}

func TestCLICompiler_Error_MermaidRender(t *testing.T) {
	t.Parallel()

	mmdc := exectest.Act(t, func() {
		log.Fatal("great sadness")
	})

	c := CLICompiler{CLI: mmdc}
	_, err := c.Compile(context.Background(), &CompileRequest{
		Source: `A -> B`,
	})
	require.Error(t, err)

	assert.ErrorContains(t, err, "exit status 1")
	assert.ErrorContains(t, err, "great sadness")
}

func TestCLICompiler_Error_NoOutput(t *testing.T) {
	t.Parallel()

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

	c := &CLICompiler{CLI: mmdc}
	_, err := c.Compile(context.Background(), &CompileRequest{
		Source: `A -> B`,
	})
	require.Error(t, err)

	assert.ErrorContains(t, err, "read svg:")
	assert.ErrorContains(t, err, "no such file or directory")
	assert.ErrorContains(t, err, "output:\nintentional no-op")
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
