# ctxutil

[![Build Status](https://travis-ci.org/posener/ctxutil.svg?branch=master)](https://travis-ci.org/posener/ctxutil)
[![codecov](https://codecov.io/gh/posener/ctxutil/branch/master/graph/badge.svg)](https://codecov.io/gh/posener/ctxutil)
[![golangci](https://golangci.com/badges/github.com/posener/ctxutil.svg)](https://golangci.com/r/github.com/posener/ctxutil)
[![GoDoc](https://godoc.org/github.com/posener/ctxutil?status.svg)](http://godoc.org/github.com/posener/ctxutil)
[![goreadme](https://goreadme.herokuapp.com/badge/posener/ctxutil.svg)](https://goreadme.herokuapp.com)

Package ctxutil is a collection of functions for contexts.

## Functions

### Interrupt

`func Interrupt() context.Context`

Interrupt is a convenience function for catching SIGINT on a
background context.

### Example:

```go
func main() {
		ctx := ctxutil.Interrupt()
		// use ctx...
}
```

### WithSignal

`func WithSignal(parent context.Context, sigWhiteList ...os.Signal) context.Context`

WithSignal returns a context which is done when an OS signal is sent.
parent is a parent context to wrap.
sigWhiteList is a list of signals to listen on.
According to the signal.Notify behavior, an empty list will listen
to any OS signal.
If an OS signal closed this context, ErrSignal will be returned in
the Err() method of the returned context.

This method creates the signal channel and invokes a goroutine.

### Examples

```go
func main() {
    // Create a context which will be cancelled on SIGINT.
    ctx := ctxutil.WithSignal(context.Background(), os.Interrupt)
    // use ctx...
}
```

### WithValues

`func WithValues(ctx, values context.Context) context.Context`

WithValues composes values from multiple contexts.
It returns a context that exposes the deadline and cancel of `ctx`,
and combined values from `ctx` and `values`.
A value in `ctx` context overrides a value with the same key in `values` context.

Consider the following standard HTTP Go server stack:

1. Middlewares that extract user credentials and request id from the
   headers and inject them to the `http.Request` context as values.
2. The `http.Handler` launches an asynchronous goroutine task which
   needs those values from the context.
3. After launching the asynchronous task the handler returns 202 to
   the client, the goroutine continues to run in background.

Problem Statement:
* The async task can't use the request context - it is cancelled
  as the `http.Handler` returns.
* There is no way to use the context values in the async task.
* Specially if those values are used automatically with client
  `http.Roundtripper` (extract the request id from the context
  and inject it to http headers in a following request.)

The suggested function `ctx := ctxutil.WithValues(ctx, values)`
does the following:

1. When `ctx.Value()` is called, the key is searched in the
   original `ctx` and if not found it searches in `values`.
2. When `Done()`/`Deadline()`/`Err()` are called, it is uses
   original `ctx`'s state.

### Example

This is how an `http.Handler` should run a goroutine that need
values from the context.

```go
func handle(w http.ResponseWriter, r *http.Request) {
    // [ do something ... ]

    // Create async task context that enables it run for 1 minute, for example
    asyncCtx, asyncCancel = ctxutil.WithTimeout(context.Background(), time.Minute)
    // Use values from the request context
    asyncCtx = ctxutil.WithValues(asyncCtx, r.Context())
    // Run the async task with it's context
    go func() {
        defer asyncCancel()
        asyncTask(asyncCtx)
    }()
}
````


Created by [goreadme](https://github.com/apps/goreadme)
