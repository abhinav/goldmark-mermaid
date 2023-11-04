package mermaidcdp

import (
	"context"
	_ "embed" // for go:embed
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"sync"

	cdruntime "github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"go.abhg.dev/goldmark/mermaid"
)

//go:embed extras.js
var _extrasJS string

// Config specifies the configuration for a Compiler.
type Config struct {
	// JSSource is JavaScript code for MermaidJS as a string.
	//
	// This will normally be the contents of the mermaid.min.js file
	// acquired from https://mermaid.js.org/intro/#cdn.
	//
	// You can use DownloadJSSource if you don't haev it available.
	JSSource string

	// Theme to use for rendering.
	//
	// Values include "dark", "default", "forest", and "neutral".
	// See MermaidJS documentation for a full list.
	Theme string
}

// Configuration for mermaid.initialize.
// Maps to MermaidConfig.
type mermaidInitializeConfig struct {
	Theme       string `json:"theme,omitempty"`
	StartOnLoad bool   `json:"startOnLoad"`
}

// Compiler compiles Mermaid diagrams into SVGs.
type Compiler struct {
	mu sync.RWMutex // guards ctx

	// While standard practice is to not hold a context in a struct,
	// we do so here because that's where chromedp puts information
	// about the headless browser it's using.
	//
	// ctx is the context scoped to the headless browser.
	ctx context.Context
}

var _ mermaid.Compiler = (*Compiler)(nil)

// New builds a new Compiler with the provided configuration.
//
// The returned Compiler must be closed with [Close] when it is no longer needed.
func New(cfg *Config) (_ *Compiler, err error) {
	if cfg.JSSource == "" {
		return nil, fmt.Errorf(
			"source code for MermaidJS must be supplied: " +
				"use DownloadJSSource if you don't have it")
	}

	// The cdp context should NOT be bound to a context with a limited lifetime
	// because that'll kill the headless browser when the context finishes.
	// Instead, we'll use the background context.
	ctx, cancel := chromedp.NewContext(context.Background())
	defer func(cancel context.CancelFunc) {
		if err != nil {
			cancel() // kill it if this function fails
		}
	}(cancel)

	var ready *cdruntime.RemoteObject
	if err := chromedp.Run(ctx, chromedp.Evaluate(cfg.JSSource, &ready)); err != nil {
		return nil, fmt.Errorf("set up headless browser: %w", err)
	}

	ready = nil
	if err := chromedp.Run(ctx, chromedp.Evaluate(_extrasJS, &ready)); err != nil {
		return nil, fmt.Errorf("inject additional JavaScript: %w", err)
	}

	initConfig := mermaidInitializeConfig{
		Theme:       cfg.Theme,
		StartOnLoad: false,
	}
	var init strings.Builder
	init.WriteString("mermaid.initialize(")
	if err := json.NewEncoder(&init).Encode(initConfig); err != nil {
		return nil, fmt.Errorf("encode mermaid.initialize config: %w", err)
	}
	init.WriteString(")")

	ready = nil
	if err := chromedp.Run(ctx, chromedp.Evaluate(init.String(), &ready)); err != nil {
		return nil, fmt.Errorf("initialize mermaid: %w", err)
	}

	c := &Compiler{ctx: ctx}
	runtime.SetFinalizer(c, func(c *Compiler) {
		// If the engine is garbage collected and not closed, close it.
		_ = c.Close()
	})
	return c, nil
}

// Compile renders a Mermaid diagram into an SVG.
// The context controls how long the rendering is allowed to take.
//
// Panics if the Compiler has already been closed.
func (c *Compiler) Compile(ctx context.Context, req *mermaid.CompileRequest) (*mermaid.CompileResponse, error) {
	var script strings.Builder
	script.WriteString("renderSVG(")
	if err := json.NewEncoder(&script).Encode(req.Source); err != nil {
		return nil, fmt.Errorf("encode source: %w", err)
	}
	script.WriteString(")")

	// TODO: Can we use chromedp.CallFunctionOn instead?
	var result string
	render := chromedp.Evaluate(
		script.String(),
		&result,
		func(p *cdruntime.EvaluateParams) *cdruntime.EvaluateParams {
			return p.WithAwaitPromise(true)
		},
	)

	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.ctx == nil {
		panic("Compiler is closed")
	}
	ctx, cancel := mergeCtxLifetime(c.ctx, ctx)
	defer cancel()

	err := chromedp.Run(ctx, render)
	return &mermaid.CompileResponse{
		SVG: result,
	}, err
}

// Close stops the compiler and releases any resources it holds.
// This method must be called when the compiler is no longer needed.
func (c *Compiler) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ctx := c.ctx; ctx != nil {
		c.ctx = nil
		return chromedp.Cancel(ctx)
	}

	return nil
}
