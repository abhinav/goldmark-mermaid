# Rendering methods

Mermaid diagrams can be rendered
at the time the file is processed ("server-side")
or in-browser when the file is viewed ("client-side").

- With server-side rendering, goldmark-mermaid calls out to the
  [MermaidJS CLI](https://github.com/mermaid-js/mermaid-cli)
  to render SVGs inline into the document.
- With client-side rendering, goldmark-mermaid generates HTML that
  renders diagrams in-browser.

You can pick between the two by setting `RenderMode` on `mermaid.Extender`.

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

By default, goldmark-mermaid will pick between the two,
based on whether it was able to find the `mmdc` executable on your `$PATH`.
