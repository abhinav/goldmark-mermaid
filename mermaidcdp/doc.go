// Package mermaidcdp implements a server-side compiler for Mermaid diagrams
// that uses a headless Chromium-based browser to render the diagrams.
//
// It's recommended to use this compiler instead of the CLI compiler
// for processes that render many diagrams.
//
// # Usage
//
// Build the compiler with [New], and be sure to [Close] it when you're done.
//
//	compiler, err := mermaidcdp.New(&mermaidcdp.Config{
//		// ...
//	})
//	if err != nil {
//		// handle error
//	}
//	defer compiler.Close()
//
// Install the compiler into your server-side renderer.
// Do this by setting the 'Compiler' field
// on mermaid.Extender or mermaid.ServerRenderer.
//
//	mermaid.Extender{
//		Compiler: compiler,
//		// ...
//	}
//
// # MermaidJS source code
//
// [Compiler] expects the MermaidJS source code supplied to it.
// This will typically be a minified version of the source code.
// You can download it from the MermaidJS GitHub repository or a CDN.
// For example, https://cdn.jsdelivr.net/npm/mermaid@10.6.0/dist/mermaid.min.js.
//
// It is recommended that you download this once,
// and embed it into your program with go:embed.
//
//	import "embed" // for go:embed
//
//	//go:embed mermaid.min.js
//	var mermaidJSSource string
//
// Then set it on the Config object you pass to [New].
//
//	mermaidcdp.New(&mermaidcdp.Config{
//		JSSource: mermaidJSSource,
//	})
//
// # Downloading MermaidJS source code
//
// As a convenience, you can use [DownloadJSSource] to download
// a minified copy of MermaidJS from a CDN programmatically.
//
//	mermaidcdp.DownloadJSSource(ctx, "10.6.0")
//
// This is useful if you can't embed it into your program.
// For most users, embedding it is recommended.
package mermaidcdp
