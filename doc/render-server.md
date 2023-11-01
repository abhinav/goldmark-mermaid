# Server-side rendering

goldmark-mermaid offers two options for server-side rendering:

- **CLI-based rendering**
  invokes the MermaidJS CLI (`mmdc`) on your system to render diagrams
- **CDP-based rendering**
  uses Chrome DevTools Protocol to drive a headless browser on your system,
  and uses it to render diagrams


## Rendering with the MermaidJS CLI

If server-side rendering is chosen, by default, the CLI-based renderer is used.
You can request it explicitly
by supplying a `CLICompiler` to `mermaid.Extender` or `mermaid.ServerRenderer`.

```go
&mermaid.Extender{
  Compiler: &mermaid.CLICompiler{
    Theme: "neutral",
  },
}
```

By default, the `CLICompiler` will search for `mmdc` on your `$PATH`.
Specify an alternative path with the `CLI` field:

```go
&mermaid.CLICompiler{
  CLI: mermaid.MMDC(pathToMMDC),
}
```

## Rendering with Chrome DevTools Protocol

<a id="render-cdp"></a>

If you have a Chromium-like browser installed on your system
goldmark-mermaid can spin up a long-running headless process of it,
and use that to render MermaidJS diagrams.

To use this, first download a minified copy of the MermaidJS source code.
You can get it from <https://cdn.jsdelivr.net/npm/mermaid@latest/dist/mermaid.min.js>.
Embed this into your program with `go:embed`.

```go
import _ "embed" // needed for go:embed

//go:embed mermaid.min.js
var mermaidJSSource string
```

Next, import `go.abhg.dev/goldmark/mermaid/mermaidcdp`,
and set up a `mermaidcdp.Compiler`.

```go
compiler, err := mermaidcdp.New(&mermaidcdp.Config{
  JSSource: mermaidJSSource,
})
if err != nil {
  panic(err)
}
defer compiler.Close() // Don't forget this!
```

Plug this compiler into the `mermaid.Extender` that
you install into your Goldmark Markdown object,
and use the Markdown object like usual.

```go
md := goldmark.New(
  goldmark.WithExtensions(
    // ...
    &mermaid.Extender{
      Compiler: compiler,
    },
  ),
  // ...
)

md.Convert(...)
```
