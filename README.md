[![Go Reference](https://pkg.go.dev/badge/github.com/abhinav/goldmark-mermaid.svg)](https://pkg.go.dev/github.com/abhinav/goldmark-mermaid)
[![Go](https://github.com/abhinav/goldmark-mermaid/actions/workflows/go.yml/badge.svg)](https://github.com/abhinav/goldmark-mermaid/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/abhinav/goldmark-mermaid/branch/main/graph/badge.svg?token=W98KYF8SPE)](https://codecov.io/gh/abhinav/goldmark-mermaid)

goldmark-mermaid is an extension for the [goldmark] Markdown parser that adds
support for [Mermaid] diagrams.

  [goldmark]: http://github.com/yuin/goldmark
  [Mermaid]: https://mermaid-js.github.io/mermaid/

# Usage

To use goldmark-mermaid, import the `mermaid` package.

```go
import mermaid "github.com/abhinav/goldmark-mermaid"
```

Then include the `mermaid.Extender` in the list of extensions you build your
[`goldmark.Markdown`] with.

  [`goldmark.Markdown`]: https://pkg.go.dev/github.com/yuin/goldmark#Markdown

```go
goldmark.New(
  &mermaid.Extender{}
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

When you render the Markdown as HTML, these will be replaced with HTML blocks.
[Mermaid] will render these into diagrams client-side.
