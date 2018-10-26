package ctxutil

import (
	"context"
	"testing"
	"time"
)

const shortDuration = 100 * time.Millisecond

func assertDone(t *testing.T, ctx context.Context) {
	t.Helper()
	select {
	case <-ctx.Done():
	case <-time.After(shortDuration):
		t.Error("context was not done")
	}
}

func assertNotDone(t *testing.T, ctx context.Context) {
	t.Helper()
	select {
	case <-ctx.Done():
		t.Error("context was done")
	case <-time.After(shortDuration):
	}
}
