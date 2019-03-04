package ctxutil

import "context"

// WithValues composes values from multiple contexts.
// It returns a context that exposes the deadline and cancel of `ctx`,
// and combined values from `ctx` and `values`.
// A value in `ctx` context overrides a value with the same key in `values` context.
//
//Consider the following standard HTTP Go server stack:
//
// 1. Middlewares that extract user credentials and request id from the
//    headers and inject them to the `http.Request` context as values.
// 2. The `http.Handler` launches an asynchronous goroutine task which
//    needs those values from the context.
// 3. After launching the asynchronous task the handler returns 202 to
//    the client, the goroutine continues to run in background.
//
// Problem Statement:
// * The async task can't use the request context - it is cancelled
//   as the `http.Handler` returns.
// * There is no way to use the context values in the async task.
// * Specially if those values are used automatically with client
//   `http.Roundtripper` (extract the request id from the context
//   and inject it to http headers in a following request.)
//
// The suggested function `ctx := ctxutil.WithValues(ctx, values)`
// does the following:
//
// 1. When `ctx.Value()` is called, the key is searched in the
//    original `ctx` and if not found it searches in `values`.
// 2. When `Done()`/`Deadline()`/`Err()` are called, it is uses
//    original `ctx`'s state.
//
// ### Example
//
// This is how an `http.Handler` should run a goroutine that need
// values from the context.
//
// 	func handle(w http.ResponseWriter, r *http.Request) {
// 		// [ do something ... ]
//
// 		// Create async task context that enables it run for 1 minute, for example
// 		asyncCtx, asyncCancel = ctxutil.WithTimeout(context.Background(), time.Minute)
// 		// Use values from the request context
// 		asyncCtx = ctxutil.WithValues(asyncCtx, r.Context())
// 		// Run the async task with it's context
// 		go func() {
// 			defer asyncCancel()
// 			asyncTask(asyncCtx)
// 		}()
// 	}
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
