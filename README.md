# ctxutil

A collection of functions for contexts.

## WithValues

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
  `http.Toundtripper` (extract the request id from the context
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
```

## WithSignal, Interrupt

Functions for creating a context which is cancelled by OS signals.

### Examples

```go
func main() {
    // Create a context which will be cancelled on SIGINT.
    ctx := ctxutil.WithSignal(context.Background(), os.Interrupt)
    // use ctx...
}
```

Interrupt is a convenience function for the most common use case
of having a background context with whitelist of interrupt signal.

```go
func main() {
    ctx := ctxutil.Interrupt()
    // use ctx...
}
```
