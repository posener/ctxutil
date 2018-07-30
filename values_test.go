package ctxutil

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const shortDuration = 100 * time.Millisecond

func TestCopyValuesCancel(t *testing.T) {
	t.Parallel()

	// Create main context with key value and cancel
	ctxOrig, cancel := context.WithCancel(context.Background())
	ctxOrig = context.WithValue(ctxOrig, "key1", "value1")
	ctxOrig = context.WithValue(ctxOrig, "key2", "value2")

	// Create a copy of the context
	ctxCopy := WithValues(context.Background(), ctxOrig)
	ctxCopy = context.WithValue(ctxCopy, "key2", "value2-2")
	ctxCopy = context.WithValue(ctxCopy, "key3", "value3")

	// Test copy of key and values of the copied context and the original context
	assert.Equal(t, "value1", ctxCopy.Value("key1").(string))
	assert.Equal(t, "value2-2", ctxCopy.Value("key2").(string))
	assert.Equal(t, "value3", ctxCopy.Value("key3").(string))

	assert.Equal(t, "value1", ctxOrig.Value("key1").(string))
	assert.Equal(t, "value2", ctxOrig.Value("key2").(string))
	assert.Nil(t, ctxOrig.Value("key3"))

	// Cancel the original context
	cancel()

	assertCancelled(t, ctxOrig)
	assertValid(t, ctxCopy)
}

func TestCopyValuesDeadline(t *testing.T) {
	t.Parallel()

	// Create main context with timeout
	ctxOrig, cancel := context.WithTimeout(context.Background(), shortDuration)
	defer cancel()

	// Create a copy of the context
	ctxCopy := WithValues(context.Background(), ctxOrig)

	// Wait for deadline
	time.Sleep(2 * shortDuration)

	assertDeadlined(t, ctxOrig)
	assertValid(t, ctxCopy)
}

func assertValid(t *testing.T, ctx context.Context) {
	_, deadline := ctx.Deadline()
	assert.False(t, deadline)
	assert.Nil(t, ctx.Err())
	select {
	case <-ctx.Done():
		t.Error("context was done")
	case <-time.After(shortDuration):
	}
}

func assertCancelled(t *testing.T, ctx context.Context) {
	_, deadline := ctx.Deadline()
	assert.False(t, deadline)
	assert.NotNil(t, ctx.Err())
	select {
	case <-ctx.Done():
	case <-time.After(shortDuration):
		t.Error("context was not done")
	}
}

func assertDeadlined(t *testing.T, ctx context.Context) {
	_, deadline := ctx.Deadline()
	assert.True(t, deadline)
	assert.NotNil(t, ctx.Err())
	select {
	case <-ctx.Done():
	case <-time.After(shortDuration):
		t.Error("context was not done")
	}
}
