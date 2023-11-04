//go:build !go1.21

package mermaidcdp

import "context"

// Pre Go 1.21, we don't have AfterFunc.
// Spawn a goroutine and use a select.

func mergeCtxLifetimeInner(parentCtx, timeCtx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancelCause(parentCtx)
	go func() {
		select {
		case <-timeCtx.Done():
			cancel(context.Cause(timeCtx))
		case <-ctx.Done():
		}
	}()
	return ctx, func() {
		cancel(context.Canceled)
	}
}
