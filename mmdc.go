package mermaid

import "go.abhg.dev/goldmark/mermaid"

// MMDC provides access to the MermaidJS CLI.
type MMDC = mermaid.MMDC

// DefaultMMDC is the default [MMDC] implementation
// used for server-side rendering.
//
// It calls out to the 'mmdc' executable to generate SVGs.
var DefaultMMDC = mermaid.DefaultMMDC

// CLI is a basic implementation of [MMDC].
type CLI = mermaid.CLI
