# Usage

To use goldmark-mermaid, import the `mermaid` package.

```go
import "go.abhg.dev/goldmark/mermaid"
```

Then include the `mermaid.Extender` in the list of extensions you build your
[`goldmark.Markdown`] with.

  [`goldmark.Markdown`]: https://pkg.go.dev/github.com/yuin/goldmark#Markdown

```go
goldmark.New(
  goldmark.WithExtensions(
    // ...
    &mermaid.Extender{},
  ),
  // ...
).Convert(src, out)
```

The package supports Mermaid diagrams inside fenced code blocks with the language `mermaid`. For example,

<pre>
```mermaid
graph TD;
    A-->B;
    A-->C;
    B-->D;
    C-->D;
```
</pre>

When you render the Markdown as HTML, these will be rendered into diagrams.

You can also render diagrams server-side if you have a Chromium-like browser
installed. See [Rendering with CDP](render-server.md#render-cdp) for details.
