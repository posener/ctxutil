package ctxutil

import (
	"context"
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignal(t *testing.T) {
	// No parallel - test changes a global variable `signalNotify`
	tests := []struct {
		name      string
		whiteList []os.Signal
		notify    func(n *notifier, cancel context.CancelFunc)
		assert    func(t *testing.T, ctx context.Context)
	}{
		{
			name: "simple signal",
			notify: func(n *notifier, cancel context.CancelFunc) {
				n.notify(os.Interrupt)
			},
			assert: func(t *testing.T, ctx context.Context) {
				assertDone(t, ctx)
				assert.Equal(t, ErrSignal{Signal: os.Interrupt}, ctx.Err())
			},
		},
		{
			name: "parent is cancelled",
			notify: func(n *notifier, cancel context.CancelFunc) {
				cancel()
			},
			assert: func(t *testing.T, ctx context.Context) {
				assertDone(t, ctx)
				assert.Equal(t, context.Canceled, ctx.Err())
			},
		},
		{
			name:   "no error",
			notify: func(n *notifier, cancel context.CancelFunc) {},
			assert: func(t *testing.T, ctx context.Context) {
				assertNotDone(t, ctx)
				assert.Nil(t, ctx.Err())
			},
		},
		{
			name:      "whitelist is passed to notifier",
			whiteList: []os.Signal{syscall.SIGTERM},
			notify: func(n *notifier, cancel context.CancelFunc) {
				n.notify(os.Interrupt)
			},
			assert: func(t *testing.T, ctx context.Context) {
				assertNotDone(t, ctx)
				assert.Nil(t, ctx.Err())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// No parallel - test changes a global variable `signalNotify`
			var n notifier
			signalNotify = func(c chan<- os.Signal, sig ...os.Signal) {
				n = notifier{c: c, whiteList: sig}
			}
			ctx, cancel := context.WithCancel(context.Background())
			ctx = WithSignal(ctx, tt.whiteList...)
			tt.notify(&n, cancel)
			tt.assert(t, ctx)
		})
	}
}

// notifier is a helper for simulating os signals.
type notifier struct {
	c         chan<- os.Signal
	whiteList []os.Signal
}

// notify sends the signal in the notifier channel if it is
// allowed by the notifier whitelist.
func (n *notifier) notify(s os.Signal) {
	if !n.allowed(s) {
		return
	}
	n.c <- s
}

// allowed check if notifier can notify with the given signal
// according to the notifier whitelist.
func (n *notifier) allowed(s os.Signal) bool {
	if len(n.whiteList) == 0 {
		return true
	}
	for _, sig := range n.whiteList {
		if sig == s {
			return true
		}
	}
	return false
}
