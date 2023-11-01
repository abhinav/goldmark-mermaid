# Rendering methods

goldmark-mermaid supports two rendering modes:

- **Client-side**:
  Diagrams are rendered in-browser by injecting MermaidJS
  into the generated HTML.
- **Server-side**:
  Diagrams are rendered at the time the HTML is generated.
  The browser receives only the final SVGs.

You can pick between these by setting the `RenderMode` field.

```go
goldmark.New(
  goldmark.WithExtensions(
    &mermaid.Extender{
      RenderMode: mermaid.RenderModeServer, // or RenderModeClient
    },
  ),
  // ...
).Convert(src, out)
```

A third automatic mode is provided as a convenience.
It automatically picks between client-side and server-side rendering
based on other configurations and system functionality.
This mode is the default.
