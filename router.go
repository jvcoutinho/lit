package lit

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var (
	ErrNilHandler    = errors.New("handler should not be nil")
	ErrMethodIsEmpty = errors.New("method should not be empty")
)

// Handler handles requests.
type Handler func(r *Request) Response

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

// NewRouter creates a new [Router] instance.
func NewRouter() *Router {
	return &Router{
		httprouter.New(),
	}
}

// Handle registers the handler for the given pattern and method.
//
// If handler is nil or method is empty or handler can't be registered, Handle panics.
func (r *Router) Handle(path string, method string, handler Handler) {
	if handler == nil {
		panic(ErrNilHandler)
	}

	if method == "" {
		panic(ErrMethodIsEmpty)
	}

	r.router.Handle(method, path, func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		request := NewRequest(r, getArguments(params))

		response := handler(request)
		if response != nil {
			response.Write(w)
		}
	})
}

// ServeHTTP dispatches the request to the handler whose pattern most closely matches the request URL
// and whose method is the same as the request method.
func (r *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	r.router.ServeHTTP(writer, request)
}
