package ctxutil

import (
	"context"
	"fmt"
	"os"
	"os/signal"
)

// ErrSignal is returned when a signal is received.
// The Signal field will contain the received signal.
type ErrSignal struct {
	Signal os.Signal
}

func (e ErrSignal) Error() string {
	return fmt.Sprintf("got signal: %s", e.Signal)
}

// WithSignal returns a context which is done when an OS signal is sent.
// parent is a parent context to wrap.
// sigWhiteList is a list of signals to listen on.
// According to the signal.Notify behavior, an empty list will listen
// to any OS signal.
// If an OS signal closed this context, ErrSignal will be returned in
// the Err() method of the returned context.
//
// This method creates the signal channel and invokes a goroutine.
//
// ### Example
//
// 	func main() {
// 		// Create a context which will be cancelled on SIGINT.
// 		ctx := ctxutil.WithSignal(context.Background(), os.Interrupt)
// 		// use ctx...
// 	}
func WithSignal(parent context.Context, sigWhiteList ...os.Signal) context.Context {
	s := &signalContext{
		Context: parent,
		done:    make(chan struct{}),
	}

	// Register a signal channel before running the goroutine.
	sigCh := make(chan os.Signal, 1)
	signalNotify(sigCh, sigWhiteList...)

	go s.watch(sigCh)
	return s
}

// Interrupt is a convenience function for catching SIGINT on a
// background context.
//
// ### Example:
//
// 	func main() {
//		ctx := ctxutil.Interrupt()
// 		// use ctx...
// 	}
func Interrupt() context.Context {
	return WithSignal(context.Background(), os.Interrupt)
}

// signalContext implements the context interface.
// It is being cancelled with OS signal cancellation.
type signalContext struct {
	// Context is the parent context
	context.Context

	done chan struct{}
	err  error
}

func (s *signalContext) Done() <-chan struct{} {
	return s.done
}

func (s *signalContext) Err() error {
	return s.err
}

// watch checks if the parent context was cancelled
// or an OS signal was received.
// It then sets the appropriate error and closes the
// current context done channel.
func (s *signalContext) watch(sigCh <-chan os.Signal) {
	select {
	case <-s.Context.Done():
		s.err = s.Context.Err()
	case sig := <-sigCh:
		s.err = ErrSignal{Signal: sig}
	}
	close(s.done)
}

// signalNotify is the signal notify function.
// It is used for testing purposes.
var signalNotify = signal.Notify
