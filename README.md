# goldmark-mermaid

- [Introduction](#introduction)
- [Installation](#installation)
- [Usage](#usage)
- [Rendering](#rendering)
  - [Rendering modes](#rendering-methods)
  - [Server-side rendering](#server-side-rendering)
- [License](#license)

## Introduction

[![Go Reference](https://pkg.go.dev/badge/go.abhg.dev/goldmark/mermaid.svg)](https://pkg.go.dev/go.abhg.dev/goldmark/mermaid)
[![CI](https://github.com/abhinav/goldmark-mermaid/actions/workflows/ci.yml/badge.svg)](https://github.com/abhinav/goldmark-mermaid/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/abhinav/goldmark-mermaid/branch/main/graph/badge.svg?token=W98KYF8SPE)](https://codecov.io/gh/abhinav/goldmark-mermaid)

goldmark-mermaid is an extension for the [goldmark](http://github.com/yuin/goldmark) Markdown parser that adds
support for [Mermaid](https://mermaid-js.github.io/mermaid/) diagrams.

**Demo**:
A web-based demonstration of the extension is available at
https://abhinav.github.io/goldmark-mermaid/demo/.

### Features

- Pluggable components
- Supports client-side rendering by injecting JavaScript
- Supports server-side rendering with the MermaidJS CLI or with your browser

## Installation

Install the latest version of the library with Go modules.

```bash
go get go.abhg.dev/goldmark/mermaid@latest
```

## Usage

To use goldmark-mermaid, import the `mermaid` package.

```go
import "go.abhg.dev/goldmark/mermaid"
```

Then include the `mermaid.Extender` in the list of extensions you build your
[`goldmark.Markdown`](https://pkg.go.dev/github.com/yuin/goldmark#Markdown) with.

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
installed. See [Rendering with CDP](#render-cdp) for details.

## Rendering

### Rendering methods

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

### Server-side rendering

goldmark-mermaid offers two options for server-side rendering:

- **CLI-based rendering**
  invokes the MermaidJS CLI (`mmdc`) on your system to render diagrams
- **CDP-based rendering**
  uses Chrome DevTools Protocol to drive a headless browser on your system,
  and uses it to render diagrams

#### Rendering with the MermaidJS CLI

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

#### Rendering with Chrome DevTools Protocol

<a id="render-cdp"></a>

If you have a Chromium-like browser installed on your system
goldmark-mermaid can spin up a long-running headless process of it,
and use that to render MermaidJS diagrams.

To use this, first download a minified copy of the MermaidJS source code.
You can get it from https://cdn.jsdelivr.net/npm/mermaid@latest/dist/mermaid.min.js.
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

## License

This software is made available under the MIT license.
