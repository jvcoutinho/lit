package lit

import (
	"context"
	"net/http"
	"time"

	"github.com/jvcoutinho/lit/internal/maps"
)

// Context is an implementation of context.Context that manages the input and output of incoming requests.
type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request

	arguments map[string]string
	context   context.Context
}

func newContext(writer http.ResponseWriter, request *http.Request) *Context {
	return &Context{
		ResponseWriter: writer,
		Request:        request,
		arguments:      nil,
		context:        context.Background(),
	}
}

func (c *Context) setArguments(arguments map[string]string) {
	c.arguments = arguments
}

// URIArguments returns a copy of the dictionary of pattern parameters and their corresponding substitutions.
//
// For example, if the request URL path is /users/123 and the matching pattern is /users/:id,
// then URIArguments will return { ":id": "123" }.
//
// This is best suited for custom manipulation of the arguments, but one should use bind.URI for
// regular usage.
func (c *Context) URIArguments() map[string]string {
	return maps.Copy(c.arguments)
}

// Deadline returns the time when work done on behalf of this context
// should be canceled. Deadline returns ok==false when no deadline is
// set. Successive calls to Deadline return the same results.
//
// See context.Context.
func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.context.Deadline()
}

// Done returns a channel that's closed when work done on behalf of this
// context should be canceled. Done may return nil if this context can
// never be canceled. Successive calls to Done return the same value.
// The close of the Done channel may happen asynchronously,
// after the cancel function returns.
//
// See context.Context.
func (c *Context) Done() <-chan struct{} {
	return c.context.Done()
}

// Err returns:
//   - context.Canceled if the context was canceled;
//   - context.DeadlineExceeded if the context's deadline passed;
//   - or a nil error if Done is not yet closed.
//
// After Err returns a non-nil error, successive calls to Err return the same error.
func (c *Context) Err() error {
	return c.context.Err()
}

// Value returns the value associated with this context for key, or nil
// if no value is associated with key. Successive calls to Value with
// the same key returns the same result.
//
// See context.Context.
func (c *Context) Value(key any) any {
	return c.context.Value(key)
}
