package mermaidcdp

import "context"

var _backgroundCtx = context.Background()

// Builds a new child context of parentCtx that will finish
// when either parentCtx or timeCtx finishes.
//
// Values from timeCtx will not be propagated to the child context.
func mergeCtxLifetime(parentCtx, timeCtx context.Context) (context.Context, context.CancelFunc) {
	// Optimization: Avoid the goroutine if either
	// is a background context.
	if parentCtx == _backgroundCtx {
		return context.WithCancel(timeCtx)
	} else if timeCtx == _backgroundCtx {
		return context.WithCancel(parentCtx)
	}

	return mergeCtxLifetimeInner(parentCtx, timeCtx)
}
