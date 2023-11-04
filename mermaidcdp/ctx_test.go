package mermaidcdp

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMergeCtxLifetime_cancels(t *testing.T) {
	t.Parallel()

	t.Run("parent canceled", func(t *testing.T) {
		t.Parallel()

		timeCtx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		parentCtx, cancel := context.WithCancel(context.Background())
		cancel()

		ctx, cancel := mergeCtxLifetime(parentCtx, timeCtx)
		defer cancel()

		<-ctx.Done()
		assert.ErrorIs(t, ctx.Err(), context.Canceled)

		assert.NoError(t, context.Cause(timeCtx), "timeCtx should not be canceled yet")
	})

	t.Run("time canceled", func(t *testing.T) {
		t.Parallel()

		parentCtx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		timeCtx, cancel := context.WithCancel(context.Background())
		cancel()

		ctx, cancel := mergeCtxLifetime(parentCtx, timeCtx)
		defer cancel()

		<-ctx.Done()
		assert.ErrorIs(t, ctx.Err(), context.Canceled)

		assert.NoError(t, context.Cause(parentCtx), "parentCtx should not be canceled yet")
	})

	t.Run("child canceled", func(t *testing.T) {
		t.Parallel()

		parentCtx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		timeCtx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		ctx, cancel := mergeCtxLifetime(parentCtx, timeCtx)
		cancel()

		<-ctx.Done()
		assert.ErrorIs(t, ctx.Err(), context.Canceled)

		assert.NoError(t, context.Cause(parentCtx), "parentCtx should not be canceled yet")
		assert.NoError(t, context.Cause(timeCtx), "timeCtx should not be canceled yet")
	})
}

func TestMergeCtxLifetime_Value(t *testing.T) {
	t.Parallel()

	type ContextKey string

	ctxKey := ContextKey("foo")

	t.Run("parent value propagates", func(t *testing.T) {
		t.Parallel()

		parentCtx := context.WithValue(context.Background(), ctxKey, "bar")

		timeCtx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		ctx, cancel := mergeCtxLifetime(parentCtx, timeCtx)
		defer cancel()

		assert.Equal(t, "bar", ctx.Value(ctxKey))
	})

	t.Run("time value ignored", func(t *testing.T) {
		t.Parallel()

		parentCtx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		timeCtx := context.WithValue(context.Background(), ctxKey, "bar")

		ctx, cancel := mergeCtxLifetime(parentCtx, timeCtx)
		defer cancel()

		assert.Nil(t, ctx.Value(ctxKey))
	})
}
