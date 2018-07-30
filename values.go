package ctxutil

import (
	"context"
)

// WithValues composes values from multiple contexts.
// It returns a context that exposes the deadline and cancel of `ctx`,
// and combined values from `ctx` and `valsCtx`.
// A value in `ctx` context overrides a value with the same key in `valsCtx` context.
func WithValues(ctx, values context.Context) context.Context {
	return &composed{Context: ctx, nextValues: values}
}

type composed struct {
	context.Context
	nextValues context.Context
}

func (c *composed) Value(key interface{}) interface{} {
	if v := c.Context.Value(key); v != nil {
		return v
	}
	return c.nextValues.Value(key)
}
