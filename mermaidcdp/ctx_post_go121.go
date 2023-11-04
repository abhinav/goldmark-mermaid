//go:build go1.21

package mermaidcdp

import "context"

func mergeCtxLifetimeInner(parentCtx, timeCtx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancelCause(parentCtx)
	stop := context.AfterFunc(timeCtx, func() {
		cancel(context.Cause(timeCtx))
	})
	return ctx, func() {
		stop()
		cancel(context.Canceled)
	}
}
