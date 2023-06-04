package lit

import "net/http"

// HandlerFunc is the standard HTTP handler function in Lit ecossystem.
type HandlerFunc func(*Request) Response

// Router manages, listens and serves HTTP requests.
type Router struct{}

// NewRouter creates a new Router instance.
func NewRouter() *Router {
	return &Router{}
}

// Handle registers the handler for the given pattern and method.
func (r *Router) Handle(pattern string, method string, handler HandlerFunc) {
}

// ServeHTTP dispatches the request to the handler whose pattern most closely matches the request URL
// and whose method is the same as the request method.
func (r *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
}
