# goldmark-mermaid

[![Go Reference](https://pkg.go.dev/badge/go.abhg.dev/goldmark/mermaid.svg)](https://pkg.go.dev/go.abhg.dev/goldmark/mermaid)
[![Go](https://github.com/abhinav/goldmark-mermaid/actions/workflows/go.yml/badge.svg)](https://github.com/abhinav/goldmark-mermaid/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/abhinav/goldmark-mermaid/branch/main/graph/badge.svg?token=W98KYF8SPE)](https://codecov.io/gh/abhinav/goldmark-mermaid)

goldmark-mermaid is an extension for the [goldmark] Markdown parser that adds
support for [Mermaid] diagrams.

  [goldmark]: http://github.com/yuin/goldmark
  [Mermaid]: https://mermaid-js.github.io/mermaid/

**Demo**:
A web-based demonstration of the extension is available at
<https://abhinav.github.io/goldmark-mermaid/demo/>.

## Installation

```bash
go get go.abhg.dev/goldmark/mermaid@latest
```

## Usage

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

    ```mermaid
    graph TD;
        A-->B;
        A-->C;
        B-->D;
        C-->D;
    ```

When you render the Markdown as HTML, these will be rendered into diagrams.

## Rendering diagrams

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
