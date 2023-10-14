package lit

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

var (
	ErrNilHandler    = errors.New("handler should not be nil")
	ErrMethodIsEmpty = errors.New("method should not be empty")
)

// HandlerFunc is the standard HTTP handler function in Lit ecosystem.
type HandlerFunc func(*Request) Response

func getArguments(params httprouter.Params) map[string]string {
	arguments := make(map[string]string)
	for _, param := range params {
		arguments[param.Key] = param.Value
	}

	return arguments
}

// Router manages, listens and serves HTTP requests.
type Router struct {
	router *httprouter.Router
}

// NewRouter creates a new Router instance.
func NewRouter() *Router {
	return &Router{
		httprouter.New(),
	}
}

// Handle registers the handler for the given pattern and method.
//
// If the route can't be registered, Handle panics.
func (r *Router) Handle(pattern string, method string, handler HandlerFunc) {
	if handler == nil {
		panic(ErrNilHandler)
	}

	if method == "" {
		panic(ErrMethodIsEmpty)
	}

	r.router.Handle(method, pattern, func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		request := newRequest(r, getArguments(params))

		response := handler(request)
		response.Write(w)
	})
}

// ServeHTTP dispatches the request to the handler whose pattern most closely matches the request URL
// and whose method is the same as the request method.
func (r *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	r.router.ServeHTTP(writer, request)
}
