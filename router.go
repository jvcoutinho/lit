package lit

import (
	"net/http"
	"slices"
	"sync"

	"github.com/julienschmidt/httprouter"
)

// Handler handles requests.
type Handler func(r *Request) Response

// Base returns the equivalent http.HandlerFunc of this handler.
func (h Handler) Base() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		r := NewRequest(req)

		if response := h(r); response != nil {
			response.Write(w)
		}
	}
}

// Middleware transforms a Handler, extending its functionality.
type Middleware func(h Handler) Handler

// Router manages, listens and serves HTTP requests.
type Router struct {
	router      *httprouter.Router
	requestPool sync.Pool
	middlewares []Middleware
}

// NewRouter creates a new [Router] instance.
func NewRouter() *Router {
	return &Router{
		httprouter.New(),
		sync.Pool{New: func() any { return NewEmptyRequest() }},
		make([]Middleware, 0),
	}
}

// Use registers m as a global middleware. They run in every request.
//
// Middlewares transform the request handler. They are applied in reverse order, and local middlewares
// are always applied first. For example, suppose there have been defined global middlewares G1 and G2 in this order and
// local middlewares L1 and L2 in this order. The response for the request r handled by h is
//
//	(G1(G2(L1(L2(h)))))(r)
//
// Global middlewares should be set before any handler is registered.
//
// If m is nil, Use panics.
func (r *Router) Use(m Middleware) {
	if m == nil {
		panic("m should not be nil")
	}

	r.middlewares = append(r.middlewares, m)
}

// Handle registers handler for path and method and optional local middlewares.
//
// Middlewares transform handler. They are applied in reverse order, and local middlewares
// are always applied first. For example, suppose there have been defined global middlewares G1 and G2 in this order and
// local middlewares L1 and L2 in this order. The response for the request r is
//
//	(G1(G2(L1(L2(handler)))))(r)
//
// If path does not contain a leading slash, method is empty, handler is nil or a middleware is nil, Handle panics.
func (r *Router) Handle(path string, method string, handler Handler, middlewares ...Middleware) {
	if handler == nil {
		panic("handler should not be nil")
	}

	if method == "" {
		panic("method should not be empty")
	}

	if slices.ContainsFunc(middlewares, func(m Middleware) bool { return m == nil }) {
		panic("middlewares should not be nil")
	}

	handler = transform(transform(handler, middlewares), r.middlewares)

	r.router.Handle(method, path, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		request := r.requestPool.Get().(*Request).
			WithRequest(req)

		for _, param := range params {
			request.parameters[param.Key] = param.Value
		}

		response := handler(request)

		if response != nil {
			response.Write(w)
		}

		r.requestPool.Put(request)
	})
}

// HandleNotFound registers handler to be called when no matching route is found. By default, Lit uses a
// wrapped http.NotFound.
//
// If handler is nil, HandleNotFound panics.
func (r *Router) HandleNotFound(handler Handler) {
	if handler == nil {
		panic("handler should not be nil")
	}

	r.router.NotFound = handler.Base()
}

// HandleOPTIONS registers handler to be called when the request method is OPTIONS. By default, Lit sets the
// Allow header with supported methods.
//
// Useful to support preflight CORS requests, for instance.
//
// If handler is nil, HandleOPTIONS clears the current set handler. In this case, the behaviour is to call the
// registered handler normally, if there is one.
func (r *Router) HandleOPTIONS(handler Handler) {
	if handler == nil {
		r.router.HandleOPTIONS = false
		r.router.GlobalOPTIONS = nil
	}

	r.router.GlobalOPTIONS = handler.Base()
}

// HandleMethodNotAllowed registers handler to be called when there is a match to a route, but not with that method.
// By default, Lit uses a wrapped http.Error with status code [405 Method Not Allowed].
//
// If handler is nil, HandleMethodNotAllowed clears the current set handler. In this case, the behaviour is to call
// the Not Found handler (either the one defined in HandleNotFound or the default one).
//
// [405 Method Not Allowed]: https://developer.mozilla.org/en-US/docs/web/http/status/405
func (r *Router) HandleMethodNotAllowed(handler Handler) {
	if handler == nil {
		r.router.HandleMethodNotAllowed = false
		return
	}

	r.router.HandleMethodNotAllowed = true
	r.router.MethodNotAllowed = handler.Base()
}

// GET registers handler for path and GET method and optional local middlewares.
//
// It's equivalent to:
//
//	Handle(path, "GET", handler, middlewares)
func (r *Router) GET(path string, handler Handler, middlewares ...Middleware) {
	r.Handle(path, http.MethodGet, handler, middlewares...)
}

// POST registers handler for path and POST method and optional local middlewares.
//
// It's equivalent to:
//
//	Handle(path, "POST", handler, middlewares)
func (r *Router) POST(path string, handler Handler, middlewares ...Middleware) {
	r.Handle(path, http.MethodPost, handler, middlewares...)
}

// PUT registers handler for path and PUT method and optional local middlewares.
//
// It's equivalent to:
//
//	Handle(path, "PUT", handler, middlewares)
func (r *Router) PUT(path string, handler Handler, middlewares ...Middleware) {
	r.Handle(path, http.MethodPut, handler, middlewares...)
}

// PATCH registers handler for path and PATCH method and optional local middlewares.
//
// It's equivalent to:
//
//	Handle(path, "PATCH", handler, middlewares)
func (r *Router) PATCH(path string, handler Handler, middlewares ...Middleware) {
	r.Handle(path, http.MethodPatch, handler, middlewares...)
}

// DELETE registers handler for path and DELETE method and optional local middlewares.
//
// It's equivalent to:
//
//	Handle(path, "DELETE", handler, middlewares)
func (r *Router) DELETE(path string, handler Handler, middlewares ...Middleware) {
	r.Handle(path, http.MethodDelete, handler, middlewares...)
}

// OPTIONS registers handler for path and OPTIONS method and optional local middlewares.
//
// It's equivalent to:
//
//	Handle(path, "OPTIONS", handler, middlewares)
func (r *Router) OPTIONS(path string, handler Handler, middlewares ...Middleware) {
	r.Handle(path, http.MethodOptions, handler, middlewares...)
}

// HEAD registers handler for path and HEAD method and optional local middlewares.
//
// It's equivalent to:
//
//	Handle(path, "HEAD", handler, middlewares)
func (r *Router) HEAD(path string, handler Handler, middlewares ...Middleware) {
	r.Handle(path, http.MethodHead, handler, middlewares...)
}

// ServeHTTP dispatches the request to the handler whose pattern most closely matches the request URL
// and whose method is the same as the request method.
func (r *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	r.router.ServeHTTP(writer, request)
}

func transform(handler Handler, middlewares []Middleware) Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler
}
