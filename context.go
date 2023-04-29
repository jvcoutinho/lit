package lit

import (
	"context"
	"net/http"
	"time"
)

// Context is an implementation of context.Context that manages the input and output of incoming requests.
type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request

	context context.Context
}

func newContext(writer http.ResponseWriter, request *http.Request) *Context {
	return &Context{writer, request, context.Background()}
}

// Deadline returns the time when work done on behalf of this context
// should be canceled. Deadline returns ok==false when no deadline is
// set. Successive calls to Deadline return the same results.
//
// See context.Context.
func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return ctx.context.Deadline()
}

// Done returns a channel that's closed when work done on behalf of this
// context should be canceled. Done may return nil if this context can
// never be canceled. Successive calls to Done return the same value.
// The close of the Done channel may happen asynchronously,
// after the cancel function returns.
//
// See context.Context.
func (ctx *Context) Done() <-chan struct{} {
	return ctx.context.Done()
}

// Err returns:
//   - context.Canceled if the context was canceled;
//   - context.DeadlineExceeded if the context's deadline passed;
//   - or a nil error if Done is not yet closed.
//
// After Err returns a non-nil error, successive calls to Err return the same error.
func (ctx *Context) Err() error {
	return ctx.context.Err()
}

// Value returns the value associated with this context for key, or nil
// if no value is associated with key. Successive calls to Value with
// the same key returns the same result.
//
// See context.Context.
func (ctx *Context) Value(key any) any {
	return ctx.context.Value(key)
}
