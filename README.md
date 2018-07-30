# ctxutil

## CopyValues

Consider the following standard HTTP Go server stack:

1. Middlewares that extract user credentials and request id from the headers and inject them to the `http.Request` context as values.
2. The `http.Handler` launches an asynchronous goroutin task which needs those values from the context.
3. After launching the asynchronous task the handler returns 202 to the client, the gorouin continues to run in background.

Problem Statement:

* The async  task can't use the request context - it is cancelled as the `http.Handler` returns.
* There is no way to use the context values in the async task.
* Specially if those values are used automatically with client `http.Toundtripper`
  (extract the request id from the context and inject it to http headers in a following request.)

The suggested function `ctx := ctxutil.CopyValues(src, dst)` does the following:
1. When `ctx.Value()` is called, the key is searched in `dst` and if not found it searches in `src`.
2. When `Done()`/`Deadline()`/`Err()` are called, it is uses `dst`'s state.
