package mermaid

import (
	"os/exec"
)

// MMDC provides access to the MermaidJS CLI.
type MMDC interface {
	// Command builds an exec.Cmd to run the Mermaid CLI
	// with the provided arguments.
	//
	// The list of arguments does not include 'mmdc' itself.
	Command(args ...string) *exec.Cmd
}

// DefaultMMDC is the default [MMDC] implementation
// used for server-side rendering.
//
// It calls out to the 'mmdc' executable to generate SVGs.
var DefaultMMDC = new(defaultMMDC)

// defaultMMDC is a transparent wrapper around CLI.
//
// Keeping the type of DefaultMMDC private
// gives us freedom to swap out the default implementation
// for something else in the future.
type defaultMMDC struct{ CLI }

var _ MMDC = (*defaultMMDC)(nil)

// CLI is a basic implementation of [MMDC].
type CLI struct {
	// Path to the "mmdc" executable.
	//
	// If unspecified, search $PATH for it.
	Path string
}

var _ MMDC = (*CLI)(nil)

// Command builds an exec.Cmd to run 'mmdc' with the given arguments.
func (c *CLI) Command(args ...string) *exec.Cmd {
	path := "mmdc"
	if len(c.Path) != 0 {
		path = c.Path
	}

	return exec.Command(path, args...)
}
