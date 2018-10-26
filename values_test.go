package ctxutil

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithValuesCancel(t *testing.T) {
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

func TestWithValuesDeadline(t *testing.T) {
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
	t.Helper()
	_, deadline := ctx.Deadline()
	assert.False(t, deadline)
	assertNotDone(t, ctx)
	assert.Nil(t, ctx.Err())
}

func assertCancelled(t *testing.T, ctx context.Context) {
	t.Helper()
	_, deadline := ctx.Deadline()
	assert.False(t, deadline)
	assertDone(t, ctx)
	assert.NotNil(t, ctx.Err())
}

func assertDeadlined(t *testing.T, ctx context.Context) {
	t.Helper()
	_, deadline := ctx.Deadline()
	assert.True(t, deadline)
	assertDone(t, ctx)
	assert.NotNil(t, ctx.Err())
}
