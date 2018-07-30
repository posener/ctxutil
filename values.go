package ctxutil

import (
	"context"
	"time"
)

// CopyValues values values from src context to dst context.
// It returns a context that exposes the deadline of dst,
// and combined values from dst and src.
// A value in dst overrides a value with the same key in src
func CopyValues(src, dst context.Context) context.Context {
	return &values{src: src, dst: dst}
}

type values struct {
	src, dst context.Context
}

func (c *values) Deadline() (deadline time.Time, ok bool) {
	return c.dst.Deadline()
}

func (c *values) Done() <-chan struct{} {
	return c.dst.Done()
}

func (c *values) Err() error {
	return c.dst.Err()
}

func (c *values) Value(key interface{}) interface{} {
	if v := c.dst.Value(key); v != nil {
		return v
	}
	return c.src.Value(key)
}
