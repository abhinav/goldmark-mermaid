package mermaid

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
)

// CLI provides access to the MermaidJS CLI.
// Use it with [CLICompiler] to override how the "mmdc" CLI is invoked.
type CLI interface {
	// CommandContext builds an exec.Cmd to run the MermaidJS CLI
	// with the provided arguments.
	//
	// The list of arguments DOES NOT include 'mmdc'.
	CommandContext(context.Context, ...string) *exec.Cmd
}

type mmdcCLI struct{ Path string }

// DefaultCLI is a [CLI] implementation that invokes the "mmdc" CLI
// by searching $PATH for it.
var DefaultCLI = MMDC("")

// MMDC returns a [CLI] implementation that invokes the "mmdc" CLI
// with the provided path.
//
// If path is empty, $PATH will be searched for "mmdc".
func MMDC(path string) CLI {
	return &mmdcCLI{Path: path}
}

func (c *mmdcCLI) CommandContext(ctx context.Context, args ...string) *exec.Cmd {
	path := c.Path
	if path == "" {
		path = "mmdc"
	}
	return exec.CommandContext(ctx, path, args...)
}

// CLICompiler compiles Mermaid diagrams into images
// by shell-executing the "mmdc" command.
//
// Plug it into [ServerRenderer] to use it.
type CLICompiler struct {
	// CLI is the MermaidJS CLI that we'll use
	// to compile Mermaid diagrams into images.
	//
	// If unset, uses DefaultCLI.
	CLI CLI

	// Theme for rendered diagrams.
	//
	// Values include "dark", "default", "forest", and "neutral".
	// See MermaidJS documentation for a full list.
	Theme string
}

var _ Compiler = (*CLICompiler)(nil)

// Compile compiles the provided Mermaid diagram into an SVG.
func (d *CLICompiler) Compile(ctx context.Context, req *CompileRequest) (_ *CompileResponse, err error) {
	mmdc := DefaultCLI
	if d.CLI != nil {
		mmdc = d.CLI
	}

	input, err := os.CreateTemp("", "in.*.mermaid")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = os.Remove(input.Name()) // ignore error
	}()

	_, err = input.WriteString(req.Source)
	if err == nil {
		err = input.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("write input: %w", err)
	}

	output, err := os.CreateTemp("", "out.*.svg")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = os.Remove(output.Name()) // ignore error
	}()
	if err := output.Close(); err != nil {
		return nil, err
	}

	args := []string{
		"--input", input.Name(),
		"--output", output.Name(),
		"--outputFormat", "svg",
		"--quiet",
	}
	if len(d.Theme) > 0 {
		args = append(args, "--theme", d.Theme)
	}

	cmd := mmdc.CommandContext(ctx, args...)
	// If the user-provided MMDC didn't set Stdout/Stderr,
	// capture its output and if anything fails beyond this point,
	// include the output in the error.
	var cmdout bytes.Buffer
	defer func() {
		if err != nil && cmdout.Len() > 0 {
			err = fmt.Errorf("%w\noutput:\n%s", err, cmdout.String())
		}
	}()
	if cmd.Stdout == nil {
		cmd.Stdout = &cmdout
	}
	if cmd.Stderr == nil {
		cmd.Stderr = &cmdout
	}

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("mmdc: %w", err)
	}

	out, err := os.ReadFile(output.Name())
	if err != nil {
		return nil, fmt.Errorf("read svg: %w", err)
	}
	return &CompileResponse{
		SVG: string(out),
	}, nil
}
