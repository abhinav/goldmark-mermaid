package mermaid

import "go.abhg.dev/goldmark/mermaid"

// Transformer transforms a Goldmark Markdown AST with support for Mermaid
// diagrams. It makes the following transformations:
//
//   - replace mermaid code blocks with mermaid.Block nodes
//   - add a mermaid.ScriptBlock node if the document uses Mermaid
//     and one does not already exist
type Transformer = mermaid.Transformer
