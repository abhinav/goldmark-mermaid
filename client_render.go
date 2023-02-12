package mermaid

import "go.abhg.dev/goldmark/mermaid"

// Renderer is the client-side renderer for Mermaid diagrams.
//
// Deprecated: Use ClientRenderer.
type Renderer = ClientRenderer

// ClientRenderer renders Mermaid diagrams as HTML,
// to be rendered into images client side.
//
// It operates by installing a <script> tag into the document
// that renders the Mermaid diagrams client-side.
type ClientRenderer = mermaid.ClientRenderer

const _defaultMermaidJS = "https://cdn.jsdelivr.net/npm/mermaid/dist/mermaid.min.js"
