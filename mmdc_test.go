package mermaid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCLI_NoPath(t *testing.T) {
	var cli CLI
	cmd := cli.Command("--version")
	assert.Equal(t, []string{"mmdc", "--version"}, cmd.Args)
}

func TestCLI_ExplicitPath(t *testing.T) {
	cli := CLI{Path: "/bin/false"}
	cmd := cli.Command("--version")
	assert.Equal(t, []string{"/bin/false", "--version"}, cmd.Args)
}
