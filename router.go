package lit

import (
	"errors"
	"net/http"
	"time"

	"github.com/jvcoutinho/lit/internal/routes"
)

// ErrNilHandler is returned when a handler function is passed as parameter and is nil.
var ErrNilHandler = errors.New("handler should not be nil")

const defaultReadHeaderTimeout = 3 * time.Second

// HandlerFunc is a function that handles requests.
type HandlerFunc func(req *Request) Response

// Router manages API routes.
//
// It is the entrypoint of a Lit-based application.
type Router struct {
	trie     *routes.Trie
	handlers map[*routes.Node]HandlerFunc

	server          *http.Server
	notFoundHandler HandlerFunc
}

// NewRouter creates a new Router instance.
func NewRouter() *Router {
	defaultNotFoundHandler := func(req *Request) Response {
		return CustomResponse(func(writer http.ResponseWriter) error {
			http.NotFound(writer, req.httpRequest)

			return nil
		})
	}

	defaultServer := &http.Server{ReadHeaderTimeout: defaultReadHeaderTimeout}

	return &Router{
		trie:            routes.NewTrie(),
		handlers:        make(map[*routes.Node]HandlerFunc),
		server:          defaultServer,
		notFoundHandler: defaultNotFoundHandler,
	}
}

// Handle registers the handler for the given pattern and method.
// If a handler already exists for pattern, Handle panics.
func (r *Router) Handle(pattern string, method string, handler HandlerFunc) {
	if handler == nil {
		panic(ErrNilHandler)
	}

	node, err := r.trie.Insert(pattern, method)
	if err != nil {
		panic(err)
	}

	r.handlers[node] = handler
}

// ServeHTTP dispatches the request to the handler whose pattern and method most closely matches one previously defined.
func (r *Router) ServeHTTP(writer http.ResponseWriter, httpRequest *http.Request) {
	var (
		request = NewRequest(httpRequest)
		handler HandlerFunc
	)

	node, arguments := r.trie.Match(httpRequest.URL.Path, httpRequest.Method)
	if node == nil {
		handler = r.notFoundHandler
	} else {
		request.setURIArguments(arguments)
		handler = r.handlers[node]
	}

	response := handler(request)
	if response == nil {
		return
	}

	if err := response.Write(writer); err != nil {
		panic(err)
	}
}

// WithServer sets the server this router uses for listening and serving requests.
func (r *Router) WithServer(server *http.Server) {
	r.server = server
}

// HandleNotFound registers the handler that runs when an incoming request matches no registered routes.
// If not explicitly called, Lit runs http.NotFound.
//
// If handler is nil, HandleNotFound panics.
func (r *Router) HandleNotFound(handler HandlerFunc) {
	if handler == nil {
		panic(ErrNilHandler)
	}

	r.notFoundHandler = handler
}

// ListenAndServe listens on the TCP network address addr and then handle requests on incoming connections.
func (r *Router) ListenAndServe(addr string) error {
	r.server.Addr = addr
	r.server.Handler = r

	return r.server.ListenAndServe()
}
