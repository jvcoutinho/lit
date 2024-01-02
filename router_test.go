package lit_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jvcoutinho/lit"
	"github.com/stretchr/testify/require"
)

func TestRouter_Handle(t *testing.T) {
	t.Parallel()

	var (
		handler    = func(r *lit.Request) lit.Response { return nil }
		middleware = func(h lit.Handler) lit.Handler { return h }
	)

	tests := []struct {
		description string
		path        string
		method      string
		handler     lit.Handler
		middlewares []lit.Middleware
		function    func(*lit.Router, string, string, lit.Handler, ...lit.Middleware)
		panicValue  any
	}{
		{
			description: "Handle_WhenPathDoesNotContainALeadingSlash_ShouldPanic",
			path:        "users",
			method:      http.MethodGet,
			handler:     handler,
			function: func(r *lit.Router, path string, method string, handler lit.Handler, middlewares ...lit.Middleware) {
				r.Handle(path, method, handler, middlewares...)
			},
			panicValue: `path must begin with '/' in path 'users'`,
		},
		{
			description: "Handle_WhenHandlerIsNil_ShouldPanic",
			path:        "/users",
			method:      http.MethodGet,
			handler:     nil,
			function: func(r *lit.Router, path string, method string, handler lit.Handler, middlewares ...lit.Middleware) {
				r.Handle(path, method, handler, middlewares...)
			},
			panicValue: "handler should not be nil",
		},
		{
			description: "Handle_WhenMethodIsEmpty_ShouldPanic",
			path:        "/users",
			method:      "",
			handler:     handler,
			function: func(r *lit.Router, path string, method string, handler lit.Handler, middlewares ...lit.Middleware) {
				r.Handle(path, method, handler, middlewares...)
			},
			panicValue: "method should not be empty",
		},
		{
			description: "Handle_WhenMiddlewaresContainsANilElement_ShouldPanic",
			path:        "/users",
			method:      http.MethodGet,
			handler:     handler,
			middlewares: []lit.Middleware{nil, middleware},
			function: func(r *lit.Router, path string, method string, handler lit.Handler, middlewares ...lit.Middleware) {
				r.Handle(path, method, handler, middlewares...)
			},
			panicValue: "middlewares should not be nil",
		},
		{
			description: "Handle_ShouldRegisterHandler",
			path:        "/users",
			method:      http.MethodGet,
			handler:     handler,
			middlewares: []lit.Middleware{middleware},
			function: func(r *lit.Router, path string, method string, handler lit.Handler, middlewares ...lit.Middleware) {
				r.Handle(path, method, handler, middlewares...)
			},
			panicValue: `a handle is already registered for path '/users'`,
		},
		{
			description: "GET_WhenPathDoesNotContainALeadingSlash_ShouldPanic",
			path:        "users",
			method:      http.MethodGet,
			handler:     handler,
			function: func(r *lit.Router, path string, _ string, handler lit.Handler, middlewares ...lit.Middleware) {
				r.GET(path, handler, middlewares...)
			},
			panicValue: `path must begin with '/' in path 'users'`,
		},
		{
			description: "GET_WhenHandlerIsNil_ShouldPanic",
			path:        "/users",
			method:      http.MethodGet,
			handler:     nil,
			function: func(r *lit.Router, path string, _ string, handler lit.Handler, middlewares ...lit.Middleware) {
				r.GET(path, handler, middlewares...)
			},
			panicValue: "handler should not be nil",
		},
		{
			description: "GET_WhenMiddlewaresContainsANilElement_ShouldPanic",
			path:        "/users",
			method:      http.MethodGet,
			handler:     handler,
			middlewares: []lit.Middleware{nil, middleware},
			function: func(r *lit.Router, path string, _ string, handler lit.Handler, middlewares ...lit.Middleware) {
				r.GET(path, handler, middlewares...)
			},
			panicValue: "middlewares should not be nil",
		},
		{
			description: "GET_ShouldRegisterHandler",
			path:        "/users",
			method:      http.MethodGet,
			handler:     handler,
			middlewares: []lit.Middleware{middleware},
			function: func(r *lit.Router, path string, _ string, handler lit.Handler, middlewares ...lit.Middleware) {
				r.GET(path, handler, middlewares...)
			},
			panicValue: `a handle is already registered for path '/users'`,
		},
		{
			description: "POST_WhenPathDoesNotContainALeadingSlash_ShouldPanic",
			path:        "users",
			method:      http.MethodPost,
			handler:     handler,
			function: func(r *lit.Router, path string, _ string, handler lit.Handler, middlewares ...lit.Middleware) {
				r.POST(path, handler, middlewares...)
			},
			panicValue: `path must begin with '/' in path 'users'`,
		},
		{
			description: "POST_WhenHandlerIsNil_ShouldPanic",
			path:        "/users",
			method:      http.MethodPost,
			handler:     nil,
			function: func(r *lit.Router, path string, _ string, handler lit.Handler, middlewares ...lit.Middleware) {
				r.POST(path, handler, middlewares...)
			},
			panicValue: "handler should not be nil",
		},
		{
			description: "POST_WhenMiddlewaresContainsANilElement_ShouldPanic",
			path:        "/users",
			method:      http.MethodPost,
			handler:     handler,
			middlewares: []lit.Middleware{nil, middleware},
			function: func(r *lit.Router, path string, _ string, handler lit.Handler, middlewares ...lit.Middleware) {
				r.POST(path, handler, middlewares...)
			},
			panicValue: "middlewares should not be nil",
		},
		{
			description: "POST_ShouldRegisterHandler",
			path:        "/users",
			method:      http.MethodPost,
			handler:     handler,
			middlewares: []lit.Middleware{middleware},
			function: func(r *lit.Router, path string, _ string, handler lit.Handler, middlewares ...lit.Middleware) {
				r.POST(path, handler, middlewares...)
			},
			panicValue: `a handle is already registered for path '/users'`,
		},
		{
			description: "PUT_WhenPathDoesNotContainALeadingSlash_ShouldPanic",
			path:        "users",
			method:      http.MethodPut,
			handler:     handler,
			function: func(r *lit.Router, path string, _ string, handler lit.Handler, middlewares ...lit.Middleware) {
				r.PUT(path, handler, middlewares...)
			},
			panicValue: `path must begin with '/' in path 'users'`,
		},
		{
			description: "PUT_WhenHandlerIsNil_ShouldPanic",
			path:        "/users",
			method:      http.MethodPut,
			handler:     nil,
			function: func(r *lit.Router, path string, _ string, handler lit.Handler, middlewares ...lit.Middleware) {
				r.PUT(path, handler, middlewares...)
			},
			panicValue: "handler should not be nil",
		},
		{
			description: "PUT_WhenMiddlewaresContainsANilElement_ShouldPanic",
			path:        "/users",
			method:      http.MethodPut,
			handler:     handler,
			middlewares: []lit.Middleware{nil, middleware},
			function: func(r *lit.Router, path string, _ string, handler lit.Handler, middlewares ...lit.Middleware) {
				r.PUT(path, handler, middlewares...)
			},
			panicValue: "middlewares should not be nil",
		},
		{
			description: "PUT_ShouldRegisterHandler",
			path:        "/users",
			method:      http.MethodPut,
			handler:     handler,
			middlewares: []lit.Middleware{middleware},
			function: func(r *lit.Router, path string, _ string, handler lit.Handler, middlewares ...lit.Middleware) {
				r.PUT(path, handler, middlewares...)
			},
			panicValue: `a handle is already registered for path '/users'`,
		},
		{
			description: "PATCH_WhenPathDoesNotContainALeadingSlash_ShouldPanic",
			path:        "users",
			method:      http.MethodPatch,
			handler:     handler,
			function: func(r *lit.Router, path string, _ string, handler lit.Handler, middlewares ...lit.Middleware) {
				r.PATCH(path, handler, middlewares...)
			},
			panicValue: `path must begin with '/' in path 'users'`,
		},
		{
			description: "PATCH_WhenHandlerIsNil_ShouldPanic",
			path:        "/users",
			method:      http.MethodPatch,
			handler:     nil,
			function: func(r *lit.Router, path string, _ string, handler lit.Handler, middlewares ...lit.Middleware) {
				r.PATCH(path, handler, middlewares...)
			},
			panicValue: "handler should not be nil",
		},
		{
			description: "PATCH_WhenMiddlewaresContainsANilElement_ShouldPanic",
			path:        "/users",
			method:      http.MethodPatch,
			handler:     handler,
			middlewares: []lit.Middleware{nil, middleware},
			function: func(r *lit.Router, path string, _ string, handler lit.Handler, middlewares ...lit.Middleware) {
				r.PATCH(path, handler, middlewares...)
			},
			panicValue: "middlewares should not be nil",
		},
		{
			description: "PATCH_ShouldRegisterHandler",
			path:        "/users",
			method:      http.MethodPatch,
			handler:     handler,
			middlewares: []lit.Middleware{middleware},
			function: func(r *lit.Router, path string, _ string, handler lit.Handler, middlewares ...lit.Middleware) {
				r.PATCH(path, handler, middlewares...)
			},
			panicValue: `a handle is already registered for path '/users'`,
		},
		{
			description: "DELETE_WhenPathDoesNotContainALeadingSlash_ShouldPanic",
			path:        "users",
			method:      http.MethodDelete,
			handler:     handler,
			function: func(r *lit.Router, path string, _ string, handler lit.Handler, middlewares ...lit.Middleware) {
				r.DELETE(path, handler, middlewares...)
			},
			panicValue: `path must begin with '/' in path 'users'`,
		},
		{
			description: "DELETE_WhenHandlerIsNil_ShouldPanic",
			path:        "/users",
			method:      http.MethodDelete,
			handler:     nil,
			function: func(r *lit.Router, path string, _ string, handler lit.Handler, middlewares ...lit.Middleware) {
				r.DELETE(path, handler, middlewares...)
			},
			panicValue: "handler should not be nil",
		},
		{
			description: "DELETE_WhenMiddlewaresContainsANilElement_ShouldPanic",
			path:        "/users",
			method:      http.MethodDelete,
			handler:     handler,
			middlewares: []lit.Middleware{nil, middleware},
			function: func(r *lit.Router, path string, _ string, handler lit.Handler, middlewares ...lit.Middleware) {
				r.DELETE(path, handler, middlewares...)
			},
			panicValue: "middlewares should not be nil",
		},
		{
			description: "DELETE_ShouldRegisterHandler",
			path:        "/users",
			method:      http.MethodDelete,
			handler:     handler,
			middlewares: []lit.Middleware{middleware},
			function: func(r *lit.Router, path string, _ string, handler lit.Handler, middlewares ...lit.Middleware) {
				r.DELETE(path, handler, middlewares...)
			},
			panicValue: `a handle is already registered for path '/users'`,
		},
		{
			description: "OPTIONS_WhenPathDoesNotContainALeadingSlash_ShouldPanic",
			path:        "users",
			method:      http.MethodOptions,
			handler:     handler,
			function: func(r *lit.Router, path string, _ string, handler lit.Handler, middlewares ...lit.Middleware) {
				r.OPTIONS(path, handler, middlewares...)
			},
			panicValue: `path must begin with '/' in path 'users'`,
		},
		{
			description: "OPTIONS_WhenHandlerIsNil_ShouldPanic",
			path:        "/users",
			method:      http.MethodOptions,
			handler:     nil,
			function: func(r *lit.Router, path string, _ string, handler lit.Handler, middlewares ...lit.Middleware) {
				r.OPTIONS(path, handler, middlewares...)
			},
			panicValue: "handler should not be nil",
		},
		{
			description: "OPTIONS_WhenMiddlewaresContainsANilElement_ShouldPanic",
			path:        "/users",
			method:      http.MethodOptions,
			handler:     handler,
			middlewares: []lit.Middleware{nil, middleware},
			function: func(r *lit.Router, path string, _ string, handler lit.Handler, middlewares ...lit.Middleware) {
				r.OPTIONS(path, handler, middlewares...)
			},
			panicValue: "middlewares should not be nil",
		},
		{
			description: "OPTIONS_ShouldRegisterHandler",
			path:        "/users",
			method:      http.MethodOptions,
			handler:     handler,
			middlewares: []lit.Middleware{middleware},
			function: func(r *lit.Router, path string, _ string, handler lit.Handler, middlewares ...lit.Middleware) {
				r.OPTIONS(path, handler, middlewares...)
			},
			panicValue: `a handle is already registered for path '/users'`,
		},
		{
			description: "HEAD_WhenPathDoesNotContainALeadingSlash_ShouldPanic",
			path:        "users",
			method:      http.MethodHead,
			handler:     handler,
			function: func(r *lit.Router, path string, _ string, handler lit.Handler, middlewares ...lit.Middleware) {
				r.HEAD(path, handler, middlewares...)
			},
			panicValue: `path must begin with '/' in path 'users'`,
		},
		{
			description: "HEAD_WhenHandlerIsNil_ShouldPanic",
			path:        "/users",
			method:      http.MethodHead,
			handler:     nil,
			function: func(r *lit.Router, path string, _ string, handler lit.Handler, middlewares ...lit.Middleware) {
				r.HEAD(path, handler, middlewares...)
			},
			panicValue: "handler should not be nil",
		},
		{
			description: "HEAD_WhenMiddlewaresContainsANilElement_ShouldPanic",
			path:        "/users",
			method:      http.MethodHead,
			handler:     handler,
			middlewares: []lit.Middleware{nil, middleware},
			function: func(r *lit.Router, path string, _ string, handler lit.Handler, middlewares ...lit.Middleware) {
				r.HEAD(path, handler, middlewares...)
			},
			panicValue: "middlewares should not be nil",
		},
		{
			description: "HEAD_ShouldRegisterHandler",
			path:        "/users",
			method:      http.MethodHead,
			handler:     handler,
			middlewares: []lit.Middleware{middleware},
			function: func(r *lit.Router, path string, _ string, handler lit.Handler, middlewares ...lit.Middleware) {
				r.HEAD(path, handler, middlewares...)
			},
			panicValue: `a handle is already registered for path '/users'`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Arrange
			router := lit.NewRouter()

			// Act
			// Assert
			require.PanicsWithValue(t, test.panicValue, func() {
				test.function(router, test.path, test.method, test.handler, test.middlewares...)

				// Attempting to handle again in order to deliberately create a panic. This will only happen
				// if the handler is successfully registered.
				test.function(router, test.path, test.method, test.handler, test.middlewares...)
			})
		})
	}
}

func TestRouter_Use(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		middleware  lit.Middleware
		panicValue  string
	}{
		{
			description: "WhenMiddlewareIsNil_ShouldPanic",
			middleware:  nil,
			panicValue:  "m should not be nil",
		},
		{
			description: "WhenMiddlewareIsNotNil_ShouldNotPanic",
			middleware: func(h lit.Handler) lit.Handler {
				return h
			},
			panicValue: "",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Arrange
			router := lit.NewRouter()

			// Act
			// Assert
			if test.panicValue != "" {
				require.PanicsWithValue(t, test.panicValue, func() {
					router.Use(test.middleware)
				})
			}
		})
	}
}

func TestRouter_HandleNotFound(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		handler     lit.Handler
		panicValue  string
	}{
		{
			description: "WhenHandlerIsNil_ShouldPanic",
			handler:     nil,
			panicValue:  "handler should not be nil",
		},
		{
			description: "WhenHandlerIsNotNil_ShouldNotPanic",
			handler: func(_ *lit.Request) lit.Response {
				return nil
			},
			panicValue: "",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Arrange
			router := lit.NewRouter()

			// Act
			// Assert
			if test.panicValue != "" {
				require.PanicsWithValue(t, test.panicValue, func() {
					router.HandleNotFound(test.handler)
				})
			}
		})
	}
}

func TestRouter_ServeHTTP(t *testing.T) {
	t.Parallel()

	var (
		notFoundHandler = func(r *lit.Request) lit.Response {
			return lit.ResponseFunc(func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusNotFound)
			})
		}

		notContentHandler = func(r *lit.Request) lit.Response {
			return lit.ResponseFunc(func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusNoContent)
			})
		}

		methodNotAllowedHandler = func(r *lit.Request) lit.Response {
			return lit.ResponseFunc(func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusMethodNotAllowed)
			})
		}

		printUsersHandler = func(r *lit.Request) lit.Response {
			return lit.ResponseFunc(func(w http.ResponseWriter) {
				w.Write([]byte("users"))
			})
		}

		getUserBookHandler = func(r *lit.Request) lit.Response {
			return lit.ResponseFunc(func(w http.ResponseWriter) {
				uri := r.URIParameters()

				w.Write([]byte(uri["user_id"] + "\n"))
				w.Write([]byte(uri["book_id"]))
			})
		}

		helloWorldMiddleware = func(h lit.Handler) lit.Handler {
			return func(r *lit.Request) lit.Response {
				res := h(r)

				return lit.ResponseFunc(func(w http.ResponseWriter) {
					w.WriteHeader(http.StatusAccepted) // this assures the order matters in tests

					w.Write([]byte("Hello, World!\n"))
					res.Write(w)
				})
			}
		}

		byeWorldMiddleware = func(h lit.Handler) lit.Handler {
			return func(r *lit.Request) lit.Response {
				res := h(r)

				return lit.ResponseFunc(func(w http.ResponseWriter) {
					w.WriteHeader(http.StatusOK) // this assures the order matters in tests

					res.Write(w)
					w.Write([]byte("\nBye, World!"))
				})
			}
		}
	)

	tests := []struct {
		description        string
		setupRouter        func(*lit.Router)
		request            *http.Request
		expectedBody       string
		expectedStatusCode int
		expectedHeader     http.Header
	}{
		{
			description:        "GivenHandlerIsNotRegistered_ShouldRespondNotFound",
			request:            httptest.NewRequest(http.MethodGet, "/users", nil),
			expectedBody:       "404 page not found\n",
			expectedStatusCode: http.StatusNotFound,
			expectedHeader: http.Header{
				"Content-Type":           {"text/plain; charset=utf-8"},
				"X-Content-Type-Options": {"nosniff"},
			},
		},
		{
			description: "GivenHandlerIsNotRegistered_AndNotFoundHandlerIsSet_ShouldRespondHandlerResponse",
			setupRouter: func(r *lit.Router) {
				r.HandleNotFound(notFoundHandler)
			},
			request:            httptest.NewRequest(http.MethodGet, "/users", nil),
			expectedBody:       "",
			expectedStatusCode: http.StatusNotFound,
			expectedHeader:     http.Header{},
		},
		{
			description: "GivenHandlerIsRegisteredForAnotherMethod_ShouldRespondMethodNotAllowed",
			setupRouter: func(r *lit.Router) {
				r.Handle("/users", http.MethodPost, printUsersHandler)
			},
			request:            httptest.NewRequest(http.MethodGet, "/users", nil),
			expectedBody:       "Method Not Allowed\n",
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectedHeader: http.Header{
				"Allow":                  {"OPTIONS, POST"},
				"Content-Type":           {"text/plain; charset=utf-8"},
				"X-Content-Type-Options": {"nosniff"}},
		},
		{
			description: "GivenHandlerIsRegisteredForAnotherMethod_AndMethodNotAllowedHandlerIsSet_ShouldRespondHandlerResponse",
			setupRouter: func(r *lit.Router) {
				r.HandleMethodNotAllowed(methodNotAllowedHandler)
				r.Handle("/users", http.MethodPost, printUsersHandler)
			},
			request:            httptest.NewRequest(http.MethodGet, "/users", nil),
			expectedBody:       "",
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectedHeader: http.Header{
				"Allow": {"OPTIONS, POST"},
			},
		},
		{
			description: "GivenHandlerIsRegisteredForAnotherMethod_AndMethodNotAllowedHandlerIsNotSet_ShouldRespondNotFound",
			setupRouter: func(r *lit.Router) {
				r.HandleMethodNotAllowed(nil)
			},
			request:            httptest.NewRequest(http.MethodGet, "/users", nil),
			expectedBody:       "404 page not found\n",
			expectedStatusCode: http.StatusNotFound,
			expectedHeader: http.Header{
				"Content-Type":           {"text/plain; charset=utf-8"},
				"X-Content-Type-Options": {"nosniff"},
			},
		},
		{
			description: "GivenHandleOPTIONSIsSet_AndMethodIsOPTIONS_ShouldRespondAllowHeaders",
			setupRouter: func(r *lit.Router) {
				r.HandleOPTIONS(notContentHandler)
				r.Handle("/users", http.MethodGet, printUsersHandler)
			},
			request:            httptest.NewRequest(http.MethodOptions, "/users", nil),
			expectedBody:       "",
			expectedStatusCode: http.StatusNoContent,
			expectedHeader: http.Header{
				"Allow": {"GET, OPTIONS"},
			},
		},
		{
			description: "GivenHandleOPTIONSIsNotSet_AndMethodIsOPTIONS_ShouldRespondMethodNotAllowed",
			setupRouter: func(r *lit.Router) {
				r.HandleOPTIONS(nil)
				r.Handle("/users", http.MethodGet, printUsersHandler)
			},
			request:            httptest.NewRequest(http.MethodOptions, "/users", nil),
			expectedBody:       "Method Not Allowed\n",
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectedHeader: http.Header{
				"Allow":                  {"GET, OPTIONS"},
				"Content-Type":           {"text/plain; charset=utf-8"},
				"X-Content-Type-Options": {"nosniff"},
			},
		},
		{
			description: "GivenHandlerIsRegistered_ShouldReturnHandlerResponse",
			setupRouter: func(r *lit.Router) {
				r.Handle("/users", http.MethodGet, printUsersHandler)
			},
			request:            httptest.NewRequest(http.MethodGet, "/users", nil),
			expectedBody:       "users",
			expectedStatusCode: http.StatusOK,
			expectedHeader: http.Header{
				"Content-Type": {"text/plain; charset=utf-8"},
			},
		},
		{
			description: "WhenPathHasParameters_ShouldParseThem",
			setupRouter: func(r *lit.Router) {
				r.Handle("/users/:user_id/books/:book_id", http.MethodGet, getUserBookHandler)
			},
			request:            httptest.NewRequest(http.MethodGet, "/users/1/books/2", nil),
			expectedBody:       "1\n2",
			expectedStatusCode: http.StatusOK,
			expectedHeader: http.Header{
				"Content-Type": {"text/plain; charset=utf-8"},
			},
		},
		{
			description: "WhenHandleHasLocalMiddlewares_ShouldUseThem",
			setupRouter: func(r *lit.Router) {
				r.Handle("/users", http.MethodGet, printUsersHandler,
					helloWorldMiddleware,
					byeWorldMiddleware,
				)
			},
			request:            httptest.NewRequest(http.MethodGet, "/users", nil),
			expectedBody:       "Hello, World!\nusers\nBye, World!",
			expectedStatusCode: http.StatusAccepted,
			expectedHeader:     http.Header{},
		},
		{
			description: "GivenRouterHasGlobalMiddlewares_ShouldUseThem",
			setupRouter: func(r *lit.Router) {
				r.Use(helloWorldMiddleware)
				r.Use(byeWorldMiddleware)
				r.Handle("/users", http.MethodGet, printUsersHandler)
			},
			request:            httptest.NewRequest(http.MethodGet, "/users", nil),
			expectedBody:       "Hello, World!\nusers\nBye, World!",
			expectedStatusCode: http.StatusAccepted,
			expectedHeader:     http.Header{},
		},
		{
			description: "GivenRouterHasGlobalMiddlewares_AndHandleHasLocalMiddlewares_ShouldUseThem",
			setupRouter: func(r *lit.Router) {
				r.Use(helloWorldMiddleware)
				r.Handle("/users", http.MethodGet, printUsersHandler,
					byeWorldMiddleware,
				)
			},
			request:            httptest.NewRequest(http.MethodGet, "/users", nil),
			expectedBody:       "Hello, World!\nusers\nBye, World!",
			expectedStatusCode: http.StatusAccepted,
			expectedHeader:     http.Header{},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Arrange
			router := lit.NewRouter()
			if test.setupRouter != nil {
				test.setupRouter(router)
			}

			recorder := httptest.NewRecorder()

			// Act
			router.ServeHTTP(recorder, test.request)

			// Assert
			require.Equal(t, test.expectedBody, recorder.Body.String())
			require.Equal(t, test.expectedStatusCode, recorder.Code)
			require.Equal(t, test.expectedHeader, recorder.Header())
		})
	}
}
