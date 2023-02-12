package mermaid

import "go.abhg.dev/goldmark/mermaid"

// Kind is the kind of a Mermaid block.
var Kind = mermaid.Kind

// Block is a Mermaid block.
//
//	```mermaid
//	graph TD;
//	    A-->B;
//	    A-->C;
//	    B-->D;
//	    C-->D;
//	```
//
// Its raw contents are the plain text of the Mermaid diagram.
type Block = mermaid.Block

// ScriptKind is the kind of a Mermaid Script block.
var ScriptKind = mermaid.ScriptKind

// ScriptBlock marks where the Mermaid Javascript will be included.
//
// This is a placeholder and does not contain anything.
type ScriptBlock = mermaid.ScriptBlock
