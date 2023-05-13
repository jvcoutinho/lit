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

// Result of an HTTP request.
//
// See the lit/render package.
type Result interface {
	// Render writes this into the HTTP response managed by ctx.
	Render(ctx *Context) error
}

// HandlerFunc is a function that handles requests.
type HandlerFunc func(ctx *Context) Result

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
	defaultNotFoundHandler := func(ctx *Context) Result {
		http.NotFound(ctx.ResponseWriter, ctx.Request)

		return nil
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
func (r *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var (
		context = NewContext(writer, request)
		handler HandlerFunc
	)

	node, arguments := r.trie.Match(request.URL.Path, request.Method)
	if node == nil {
		handler = r.notFoundHandler
	} else {
		context.setArguments(arguments)
		handler = r.handlers[node]
	}

	result := handler(context)
	if result == nil {
		return
	}

	if err := result.Render(context); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

// WithServer sets the server this router uses for listening and serving requests.
func (r *Router) WithServer(server *http.Server) {
	r.server = server
}

// WithNotFoundHandler sets the handler that runs when an incoming request matches no defined routes.
//
// The default handler just runs http.NotFound. If handler is nil, WithNotFoundHandler panics.
func (r *Router) WithNotFoundHandler(handler HandlerFunc) {
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
